// @title Models API
// @version 1.0
// @description AI 模型数据聚合服务，从 models.dev 同步数据并提供 RESTful API
// @host api.model.false.ltd
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/false-ltd/model/api/internal/cache"
	"github.com/false-ltd/model/api/internal/config"
	"github.com/false-ltd/model/api/internal/handler"
	"github.com/false-ltd/model/api/internal/repository"
	"github.com/false-ltd/model/api/internal/router"
	"github.com/false-ltd/model/api/internal/service"
)

var version = "dev"

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("invalid configuration: %v", err)
	}
	log.Printf("models %s starting", version)

	if len(cfg.Auth.APIKeys) == 0 {
		log.Printf("WARNING: MODEL_API_KEYS is empty; POST /api/v1/sync is disabled (fail closed)")
	}

	db, err := gorm.Open(mysql.Open(cfg.Database.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get underlying sql.DB: %v", err)
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Printf("connected to MySQL")

	appCache := cache.New()

	providerRepo := repository.NewProviderRepo(db)
	modelRepo := repository.NewModelRepo(db)

	providerService := service.NewProviderService(providerRepo, appCache)
	modelService := service.NewModelService(modelRepo)
	statsService := service.NewStatsService(modelRepo, providerRepo, appCache)
	atlasService := service.NewAtlasService(modelRepo, providerRepo, appCache)
	syncService := service.NewSyncService(providerRepo, modelRepo, db, &cfg.Sync, appCache)

	handlers := &router.Handlers{
		Model:    handler.NewModelHandler(modelService),
		Provider: handler.NewProviderHandler(providerService),
		Stats:    handler.NewStatsHandler(statsService),
		Atlas:    handler.NewAtlasHandler(atlasService),
		Compare:  handler.NewCompareHandler(modelService),
		Sync:     handler.NewSyncHandler(syncService),
		Sitemap:  handler.NewSitemapHandler(modelRepo, cfg.Server.SiteURL, appCache),
	}

	engine := router.Setup(cfg, handlers)
	router.SetupSPA(engine)

	// Start background sync cron
	if cfg.Sync.CronMinutes > 0 {
		interval := time.Duration(cfg.Sync.CronMinutes) * time.Minute
		go func() {
			log.Printf("sync cron started, interval: %dm", cfg.Sync.CronMinutes)
			// Initial sync on startup
			if result, err := syncService.Trigger(); err != nil {
				log.Printf("initial sync failed: %v", err)
			} else if result.Skipped {
				log.Printf("initial sync skipped, last sync: %s", result.SyncedAt)
			} else {
				log.Printf("initial sync completed, %d providers, %d models", result.Providers, result.Models)
			}
			// Periodic sync
			ticker := time.NewTicker(interval)
			defer ticker.Stop()
			for range ticker.C {
				if result, err := syncService.Trigger(); err != nil {
					log.Printf("cron sync failed: %v", err)
				} else if result.Skipped {
					log.Printf("cron sync skipped, last sync: %s", result.SyncedAt)
				} else {
					log.Printf("cron sync completed, %d providers, %d models", result.Providers, result.Models)
				}
			}
		}()
	}

	srv := &http.Server{
		Addr:              ":" + cfg.Server.Port,
		Handler:           engine,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	go func() {
		log.Printf("server starting on :%s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Printf("shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("forced shutdown: %v", err)
	}
	log.Printf("server stopped")
}
