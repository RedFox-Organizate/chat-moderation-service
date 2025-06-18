package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"chat-moderation-service/internal/app"
	"chat-moderation-service/internal/db"
	"chat-moderation-service/internal/utils"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	config := zap.NewProductionConfig()

	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	appLogger, err := config.Build()
	if err != nil {
		log.Fatalf("Zap logger başlatılamadı: %v", err)
	}
	defer func() {
		syncErr := appLogger.Sync()
		if syncErr != nil {
			log.Printf("Zap logger senkronizasyon hatası: %v", syncErr)
		}
	}()

	cfg, err := utils.LoadConfig("config/config.yaml")
	if err != nil {
		appLogger.Fatal("Config yüklenemedi", zap.Error(err))
	}

	mongoClient, err := db.NewMongoClient(cfg.MongoURI)
	if err != nil {
		appLogger.Fatal("MongoDB bağlantısı kurulamadı", zap.Error(err))
	}
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			appLogger.Error("MongoDB bağlantısı kapatılamadı", zap.Error(err))
		}
	}()

	badwordManager, err := app.NewBadWordManager(cfg.BadWordsFile, appLogger)
	if err != nil {
		appLogger.Fatal("BadWords yöneticisi oluşturulamadı", zap.Error(err))
	}
	go badwordManager.WatchFileChanges()

	moderator, err := app.NewModerator(mongoClient, cfg.DatabaseName, cfg.CollectionName, badwordManager, appLogger, cfg.AllowedPlayers)
	if err != nil {
		appLogger.Fatal("Moderator oluşturulamadı", zap.Error(err))
	}
	defer moderator.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go moderator.StartMonitoring()

	<-stop
	appLogger.Info("Program kapatılıyor...")
}