package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jvzito/airball/internal/cache"
	"github.com/jvzito/airball/internal/config"
	"github.com/jvzito/airball/internal/handlers"
	"github.com/jvzito/airball/internal/httpclient"
	"github.com/jvzito/airball/internal/middleware"
	"github.com/jvzito/airball/internal/repository"
	"github.com/jvzito/airball/internal/service"
	"github.com/jvzito/airball/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()
	logger.Init(cfg.AppEnv)
	defer logger.Log.Sync()

	redisCache, err := cache.New(cfg)
	if err != nil {
		logger.Fatal("redis", zap.Error(err))
	}
	defer redisCache.Close()

	db, err := repository.NewPostgres(cfg)
	if err != nil {
		logger.Fatal("postgres", zap.Error(err))
	}
	defer db.Close()

	nbaClient := httpclient.NewNBAClient(cfg)
	userRepo := repository.NewUserRepo(db)
	playerSvc := service.NewPlayerService(nbaClient, redisCache)
	authSvc := service.NewAuthService(userRepo, cfg.JWTSecret)

	playerH := handlers.NewPlayerHandler(playerSvc)
	searchH := handlers.NewSearchHandler()
	authH := handlers.NewAuthHandler(authSvc)

	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(middleware.RequestLogger(), middleware.CORS(), gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "env": cfg.AppEnv})
	})

	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		auth.POST("/register", authH.Register)
		auth.POST("/login", authH.Login)

		v1.GET("/leaders/:category", playerH.GetLeaders)
		v1.GET("/players/search", searchH.Search)
		v1.GET("/players/:id/shotchart", playerH.GetShotChart)
	}

	srv := &http.Server{
		Addr:         ":" + cfg.AppPort,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
	go func() {
		logger.Info("Airball API iniciada", zap.String("port", cfg.AppPort))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("encerrando...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
	logger.Info("encerrado")
}
