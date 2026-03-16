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
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/yourusername/shorty/internal/config"
	"github.com/yourusername/shorty/internal/database"
	"github.com/yourusername/shorty/internal/handlers"
	"github.com/yourusername/shorty/internal/repositories"
	"github.com/yourusername/shorty/internal/services"
)

type app struct {
	cfg     *config.Config
	handler *handlers.URLHandler
	dbpool  interface{ Close() }
}

func main() {
	setupLogger()

	application := mustBuildApp()
	defer application.dbpool.Close()

	router := buildRouter(application.handler, application.cfg)
	server := newHTTPServer(application.cfg, router)

	runServer(server)
	gracefulShutdown(server)
}
func setupLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
}

func mustBuildApp() *app {
	cfg, err := config.LoadConfigFromEnv()
	if err != nil {
		log.Fatal().Err(err).Msg("load config")
	}

	dbpool, err := database.NewPgxPool(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("connect db")
	}

	repo := repositories.NewPgURLRepository(dbpool)
	svc := services.NewURLService(repo, cfg)
	handler := handlers.NewURLHandler(svc, cfg)

	return &app{
		cfg:     cfg,
		handler: handler,
		dbpool:  dbpool,
	}
}

func buildRouter(handler *handlers.URLHandler, cfg *config.Config) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(handlers.ZerologMiddleware())
	router.Use(handlers.CORSMiddleware(cfg))

	api := router.Group("/api/v1")
	{
		api.POST("/shorten", handler.Shorten)
		api.GET("/stats/:code", handler.Stats)
		api.GET("/urls", handler.GetAll)
	}

	router.GET("/:code", handler.Redirect)

	return router
}

func newHTTPServer(cfg *config.Config, router *gin.Engine) *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}
}

func runServer(srv *http.Server) {
	go func() {
		log.Info().Msgf("server starting on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("listen")
		}
	}()
}
func gracefulShutdown(srv *http.Server) {
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
