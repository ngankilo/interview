package router

import (
	"fmt"
	"go-streaming/model"
	"go-streaming/request_api"
	"go-streaming/utils"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func SseHandler(config *model.Config, jwtToken model.JWT, channelManager *model.ChannelManager, rdb *model.PubsubManager) gin.HandlerFunc {
	const maxConnections = 1000 // Configurable max connections
	// semaphore := limit.NewSemaphore(maxConnections)
	pricePrefix := "price" // Configurable via config if needed
	maxMessageSize := 1024 * 1024 // 1MB
	idleTimeout := 30 * time.Second

	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		prefix, keys, arrayChannels := utils.GetPrefixStreamingByGinContext(c)
		path := c.Request.URL.Path

		// Acquire connection slot
		// if !semaphore.Acquire() {
		// 	log.Warn().Str("ip", clientIP).Msg("Connection limit reached")
		// 	c.JSON(http.StatusServiceUnavailable, gin.H{
		// 		"code":    http.StatusServiceUnavailable,
		// 		"message": "Connection limit reached",
		// 	})
		// 	return
		// }
		// defer semaphore.Release()

		// Validate JWT
		var token jwt.MapClaims
		if config.CheckJwt {
			var err error
			token, _, err = jwtToken.Validate(c)
			if err != nil {
				utils.AddMetricsStatus(fmt.Sprintf("%v", token["iss"]), clientIP, strconv.Itoa(http.StatusUnauthorized), path, "sse")
				log.Error().Err(err).Str("user", fmt.Sprintf("%v", token["iss"])).Str("path", path).Msg("JWT validation failed")
				c.JSON(http.StatusUnauthorized, gin.H{
					"code":    http.StatusUnauthorized,
					"message": err.Error(),
				})
				return
			}
		}

		// Start SSE
		deviceId := token["device_id"]
		userId := fmt.Sprintf("%v", token["iss"])
		subId := fmt.Sprintf("%v", token["sub"])
		sseId := uuid.New()

		if userId != config.UserAccess && strings.Contains(prefix, pricePrefix) {
			if !request_api.CheckRealtimeMarketData(subId, config.RedisTTL) {
				log.Warn().Str("user_id", userId).Msg("Access denied for price channel")
				c.JSON(http.StatusForbidden, gin.H{
					"code":    http.StatusForbidden,
					"message": fmt.Sprintf("Access denied for user: %s", userId),
				})
				return
			}

			if deviceIdStr, ok := deviceId.(string); !ok || !request_api.CheckLatestDeviceId(subId, deviceIdStr, config.RedisTTL) {
				log.Warn().Str("user_id", userId).Msg("Invalid session")
				c.JSON(http.StatusForbidden, gin.H{
					"code":    http.StatusForbidden,
					"message": fmt.Sprintf("Invalid session for user: %s", userId),
				})
				return
			}
		}

		startTime := time.Now()
		utils.AddMetricsStatus(userId, clientIP, "200", path, "sse")
		utils.AddMetricsConnection(userId, clientIP, path, "sse", true)
		log.Info().
			Str("ip", clientIP).
			Str("user", userId).
			Str("device", fmt.Sprintf("%v", deviceId)).
			Str("sse_id", sseId.String()).
			Str("path", path).
			Str("keys", keys).
			Msg("SSE connection established")

		listener := channelManager.OpenListener(prefix, keys)
		if !config.SubAll {
			rdb.RegisterRedisPubsub(arrayChannels)
		}

		defer func() {
			channelManager.CloseListener(prefix, keys, listener)
			if !config.SubAll {
				rdb.UnregisterRedisPubsub(arrayChannels)
			}
			utils.AddMetricsConnection(userId, clientIP, path, "sse", false)
			// metrics.Observe("sse.connection_duration", time.Since(startTime).Seconds())
			log.Info().
				Str("ip", clientIP).
				Str("user", userId).
				Str("device", fmt.Sprintf("%v", deviceId)).
				Str("sse_id", sseId.String()).
				Str("path", path).
				Str("keys", keys).
				Msg("SSE connection closed")
		}()

		clientGone := c.Request.Context().Done()
		idleTimer := time.NewTimer(idleTimeout)
		defer idleTimer.Stop()

		c.Stream(func(w io.Writer) bool {
			select {
			case <-clientGone:
				return false
			case <-idleTimer.C:
				log.Info().Str("sse_id", sseId.String()).Msg("SSE connection timed out")
				return false
			case message := <-listener:
				idleTimer.Reset(idleTimeout)
				if msgStr, ok := message.(string); ok && len(msgStr) > maxMessageSize {
					log.Warn().Str("sse_id", sseId.String()).Msg("Message size exceeds limit")
					return true
				}
				utils.AddMetricsMesage(userId, clientIP, path, "sse")
				c.SSEvent("", message)
				return true
			}
		})
	}
}