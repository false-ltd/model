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
	"log"
	"time"

	"github.com/false-ltd/model/api/internal/config"
	"github.com/false-ltd/model/api/internal/handler"
	"github.com/false-ltd/model/api/internal/repository"
	"github.com/false-ltd/model/api/internal/router"
	"github.com/false-ltd/model/api/internal/service"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var version = "dev"

func main() {
	cfg := config.Load()
	log.Printf("models %s starting", version)

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

	log.Printf("connected to MySQL")

	providerRepo := repository.NewProviderRepo(db)
	modelRepo := repository.NewModelRepo(db)

	providerService := service.NewProviderService(providerRepo)
	modelService := service.NewModelService(modelRepo)
	statsService := service.NewStatsService(modelRepo, providerRepo)
	syncService := service.NewSyncService(providerRepo, modelRepo, &cfg.Sync)

	handlers := &router.Handlers{
		Model:    handler.NewModelHandler(modelService),
		Provider: handler.NewProviderHandler(providerService),
		Stats:    handler.NewStatsHandler(statsService),
		Compare:  handler.NewCompareHandler(modelService),
		Sync:     handler.NewSyncHandler(syncService),
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

	log.Printf("server starting on :%s", cfg.Server.Port)
	if err := engine.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
