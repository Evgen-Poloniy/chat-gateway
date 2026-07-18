package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Evgen-Poloniy/chat-gateway/internal/config"
	kf "github.com/Evgen-Poloniy/chat-gateway/internal/repository/kafka"
	pg "github.com/Evgen-Poloniy/chat-gateway/internal/repository/postgres"
	rds "github.com/Evgen-Poloniy/chat-gateway/internal/repository/redis"
	httpserver "github.com/Evgen-Poloniy/chat-gateway/internal/server/http"
	"github.com/Evgen-Poloniy/chat-gateway/internal/service/messenger"
	"github.com/Evgen-Poloniy/chat-gateway/internal/service/resolver"
	router "github.com/Evgen-Poloniy/chat-gateway/internal/transport/http"
	v1 "github.com/Evgen-Poloniy/chat-gateway/internal/transport/http/v1"
	"github.com/Evgen-Poloniy/chat-gateway/internal/transport/http/ws"

	"github.com/Evgen-Poloniy/chat-gateway/pkg/database"
	logs "github.com/Evgen-Poloniy/chat-gateway/pkg/logger"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func Run() {
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		logrus.Fatal("required API_KEY")
	}

	apiKeyHash, err := bcrypt.GenerateFromPassword([]byte(apiKey), bcrypt.DefaultCost)
	if err != nil {
		logrus.Fatalf("error when generation API_KEY hash: %v", err)
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		logrus.Fatalf("error when loading env CONFIG_PATH")
	}

	config, err := config.LoadConfig(configPath)
	if err != nil {
		logrus.Fatalf("error when loading config: %v", err)
	}

	logger, err := logs.NewLogrusLogger(&config.Logger)
	if err != nil {
		logrus.Fatalf("error when loading the logger: %v", err)
	}
	defer func() {
		if err := logs.Close(); err != nil {
			logger.Errorf("error when closing the logger: %v", err)
		}
	}()

	db, err := database.NewPostgreSQL(&config.Postgres)
	if err != nil {
		logger.Errorf("database error: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Errorf("database error: %v", err)
		}
	}()

	producer, err := database.NewKafkaProducer(&config.Kafka, logger)
	if err != nil {
		logger.Errorf("message broker error: %v", err)
	}
	defer func() {
		unflushedCount := producer.Flush(config.Kafka.FlashTimeout * 1000)
		if unflushedCount > 0 {
			logger.Warnf("warning: %d messages were not flushed and might be lost", unflushedCount)
		}

		producer.Close()
	}()

	rdb, err := database.NewRedisClient(&config.Redis)
	if err != nil {
		logger.Fatalf("error: %v", err)
	}
	defer func() {
		err := rdb.Close()
		if err != nil {
			logger.Errorf("error when closing redis: %v", err)
		}

		logger.Info("redis successful closed")
	}()

	messengerRepository := pg.NewPostgresRepository(db)
	messageBroker := kf.NewKafkaRepository(producer)

	cache := rds.NewRedisCache(rdb, &config.Chat.Redis, &config.PubSub.Redis)

	messenger := messenger.NewMessengerService(messengerRepository, messageBroker, cache)
	resolver := resolver.NewResolverService(cache, cache)

	address := fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)
	wsHub := ws.NewHub(resolver, address, logger)
	v1Handler := v1.NewHandler(messenger, wsHub)
	wsHandler := ws.NewHandler(wsHub, logger)
	router := router.NewRouter(&config.CORS, logger)
	v1.NewRouter(router, v1Handler, apiKeyHash)
	ws.NewRouter(router, wsHandler)

	httpServer := httpserver.NewServer(&config.Server, router)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		logger.Infof("http server is running on %s", address)
		if err := httpServer.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf("http server error: %v", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		logger.Info("starting dispatch messages worker")
		wsHub.StartDispatchMessage()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	ctx, cancel := context.WithTimeout(
		context.Background(),
		config.Server.TimeForGracefulShutdown*time.Second,
	)
	defer cancel()

	logger.Info("shutting down dispatch messages worker")
	wsHub.ShutdownDispatchMessage(ctx)

	logger.Info("shutting down server")
	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Errorf("server forced to shutdown: %v", err)
	}

	wg.Wait()

	logger.Info("server stop")
}
