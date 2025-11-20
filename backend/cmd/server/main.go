package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/yourusername/shorty/internal/config"
	"github.com/yourusername/shorty/internal/database"
	"github.com/yourusername/shorty/internal/handlers"
	"github.com/yourusername/shorty/internal/repositories"
	"github.com/yourusername/shorty/internal/services"
)

func main() {
    // load env file if present (dev)
    _ = godotenv.Load()

    // setup logger
    zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
    log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

    cfg, err := config.LoadConfigFromEnv()
    if err != nil {
        log.Fatal().Err(err).Msg("load config")
    }

    dbpool, err := database.NewPgxPool(cfg)
    if err != nil {
        log.Fatal().Err(err).Msg("connect db")
    }
    defer dbpool.Close()

    repo := repositories.NewPgUrlRepository(dbpool)
    svc := services.NewURLService(repo, cfg)
    h := handlers.NewUrlHandler(svc, cfg)

    router := gin.New()
    router.Use(gin.Recovery())
    router.Use(handlers.ZerologMiddleware())
    router.Use(handlers.CORSMiddleware(cfg))
    

    api := router.Group("/api/v1")
    {
        api.POST("/shorten", h.Shorten)
        api.GET("/stats/:code", h.Stats)
        api.GET("/urls", h.GetAll)
    }
    // redirect route
    router.GET("/:code", h.Redirect)
    
    srv := &http.Server{
        Addr:    fmt.Sprintf(":%s", cfg.Port),
        Handler: router,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  30 * time.Second,
    }

    // graceful shutdown
    go func() {
        log.Info().Msgf("server starting on %s", srv.Addr)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatal().Err(err).Msg("listen")
        }
    }()

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    log.Info().Msg("shutting down server...")

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    if err := srv.Shutdown(ctx); err != nil {
        log.Fatal().Err(err).Msg("server shutdown")
    }
    log.Info().Msg("server stopped")
}