package main

import (
	"context"
	"go-streaming/model"
	"go-streaming/request_api"
	"go-streaming/router"
	"go-streaming/utils"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/penglongli/gin-metrics/ginmetrics"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	jwtToken       model.JWT
	sseInstanceId  = uuid.New().String()
	channelManager *model.ChannelManager
	pubsubManager  *model.PubsubManager
	rdb            *redis.Client
)

func main() {
	// Initialize logging
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	gin.SetMode(gin.ReleaseMode)

	// Load configuration
	config, err := model.LoadConfig("config-go-streaming-api")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Initialize JWT
	jwtToken, err = model.NewJWT()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize JWT")
	}

	// Initialize Redis
	rdb = redis.NewClient(&redis.Options{
		Addr:     config.Redis.Host + ":" + config.Redis.Port,
		Password: config.Redis.Password,
		DB:       config.Redis.Database,
	})
	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to Redis")
	}
	defer rdb.Close()

	// Initialize Redis cache
	cacheConfig := &model.ConfigRedisCache{
		Host:     config.Redis.Host,
		Port:     config.Redis.Port,
		Password: config.Redis.Password,
		Database: config.Redis.Database,
	}
	if err := request_api.InitRedisCache(cacheConfig); err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize Redis cache")
	}

	// Initialize channel manager
	channelManager = model.NewChannelManager(config.BufferSize, time.Millisecond*time.Duration(config.BroadcasterTimeout), config.EnablePing)

	// Initialize Redis PubSub
	if config.SubAll {
		rdbMsg := make(chan *redis.Message, config.BufferSize)
		go func() {
			channel := config.RedisPattern // Configurable pattern
			log.Info().Str("channel", channel).Msg("Subscribing to Redis pattern")
			subscribeRedis := rdb.PSubscribe(context.Background(), channel)
			defer subscribeRedis.Close()
			for msg := range subscribeRedis.Channel() {
				rdbMsg <- msg
			}
		}()
		go func() {
			for msg := range rdbMsg {
				utils.SendData(channelManager, msg)
			}
		}()
	} else {
		subscribeRedis := rdb.Subscribe(context.Background())
		pubsubManager = model.NewPubsubManager(subscribeRedis, config.BufferSize)
		rdbMsg := make(chan *redis.Message, config.BufferSize)
		go func() {
			for {
				msg, err := subscribeRedis.ReceiveMessage(context.Background())
				if err != nil {
					log.Error().Err(err).Msg("Failed to receive Redis message")
					continue
				}
				if msg != nil {
					rdbMsg <- msg
				}
			}
		}()
		go func() {
			for msg := range rdbMsg {
				utils.SendData(channelManager, msg)
			}
		}()
	}

	// Heartbeat
	go func() {
		ticker := time.Tick(10 * time.Second)
		for range ticker {
			utils.SendPing(channelManager, sseInstanceId)
		}
	}()

	// Initialize Gin router
	router := gin.New()
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		log.Info().
			Int("status", param.StatusCode).
			Str("ip", param.ClientIP).
			Str("method", param.Method).
			Str("path", param.Path).
			Dur("latency", param.Latency).
			Str("user_agent", param.Request.UserAgent()).
			Str("error", param.ErrorMessage).
			Msg("HTTP request")
		return ""
	}))

	// Metrics
	metrics := ginmetrics.GetMonitor()
	metrics.SetMetricPath("/metrics")
	metrics.Use(router)
	metrics.AddMetric(&ginmetrics.Metric{
		Type:        ginmetrics.Counter,
		Name:        "go_gin_sse_status_total",
		Description: "Metric status response",
		Labels:      []string{"userId", "ip", "status", "path", "streamType"},
	})
	metrics.AddMetric(&ginmetrics.Metric{
		Type:        ginmetrics.Gauge,
		Name:        "go_gin_sse_total",
		Description: "Metric connections",
		Labels:      []string{"userId", "ip", "path", "streamType"},
	})
	metrics.AddMetric(&ginmetrics.Metric{
		Type:        ginmetrics.Gauge,
		Name:        "go_gin_sse_message_total",
		Description: "Metric message total",
		Labels:      []string{"userId", "ip", "path", "streamType"},
	})

	// Routes
	config.RedisTTL = 45 // Example TTL
	config.CorsOrigin = "*" // Example CORS origin
	config.RedisPattern = "*" // Example pattern
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})
	router.GET("/:p1/:p2/:p3/:p4/:p5/:p6", router.SseHandler(config, jwtToken, channelManager, pubsubManager))
	router.GET("/:p1/:p2/:p3/:p4/:p5", router.SseHandler(config, jwtToken, channelManager, pubsubManager))
	router.GET("/:p1/:p2/:p3/:p4", router.SseHandler(config, jwtToken, channelManager, pubsubManager))
	router.GET("/:p1/:p2/:p3", router.SseHandler(config, jwtToken, channelManager, pubsubManager))
	router.GET("/:p1/:p2", router.SseHandler(config, jwtToken, channelManager, pubsubManager))
	router.OPTIONS("/:p1/:p2/:p3/:p4/:p5/:p6/p7", router.Options(config))
	router.OPTIONS("/:p1/:p2/:p3/:p4/:p5/:p6", router.Options(config))
	router.OPTIONS("/:p1/:p2/:p3/:p4/:p5", router.Options(config))
	router.OPTIONS("/:p1/:p2/:p3/:p4", router.Options(config))
	router.OPTIONS("/:p1/:p2/:p3", router.Options(config))
	router.OPTIONS("/:p1/:p2", router.Options(config))

	// Start server with graceful shutdown
	srv := &http.Server{
		Addr:    ":" + config.Port,
		Handler: router,
	}

	go func() {
		log.Info().Str("port", config.Port).Msg("Starting server")
		// metrics.Increment("server.start")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Server failed")
		}
	}()

	// Handle shutdown
	ctx, cancel := context.WithCancel(context.Background())
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Info().Msg("Shutting down server")
		// metrics.Increment("server.shutdown")

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Error().Err(err).Msg("Server shutdown failed")
		}
		if err := channelManager.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close channel manager")
		}
		if pubsubManager != nil {
			if err := pubsubManager.Close(); err != nil {
				log.Error().Err(err).Msg("Failed to close pubsub manager")
			}
		}
		if err := rdb.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close Redis")
		}
		cancel()
	}()

	<-ctx.Done()
	log.Info().Msg("Server stopped")
}