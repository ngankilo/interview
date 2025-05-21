package router

import (
	"go-streaming/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func Options(config *model.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Info().Str("ip", c.ClientIP()).Str("path", c.Request.URL.Path).Msg("Handling OPTIONS request")
		c.Header("Access-Control-Allow-Origin", config.CorsOrigin)
		c.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
		c.Writer.WriteHeader(http.StatusNoContent)
	}
}