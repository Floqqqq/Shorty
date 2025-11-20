package handlers

import (
    "time"

    "github.com/gin-gonic/gin"
    "github.com/rs/zerolog/log"

    "github.com/yourusername/shorty/internal/config"
)

func CORSMiddleware(cfg *config.Config) gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        c.Next()
    }
}

func ZerologMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        c.Next()
        dur := time.Since(start)
        status := c.Writer.Status()
        log.Info().
            Str("method", c.Request.Method).
            Str("path", c.Request.URL.Path).
            Int("status", status).
            Dur("duration", dur).
            Msg("")
    }
}