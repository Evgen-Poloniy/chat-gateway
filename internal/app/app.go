package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/Evgen-Poloniy/chat-gateway/internal/config"
	"github.com/Evgen-Poloniy/chat-gateway/internal/middleware"
	kfk "github.com/Evgen-Poloniy/chat-gateway/internal/repository/kafka"
	pg "github.com/Evgen-Poloniy/chat-gateway/internal/repository/postgres"
	rds "github.com/Evgen-Poloniy/chat-gateway/internal/repository/redis"
	httpserver "github.com/Evgen-Poloniy/chat-gateway/internal/server/http"
	"github.com/Evgen-Poloniy/chat-gateway/internal/service/messenger"
	"github.com/Evgen-Poloniy/chat-gateway/internal/service/resolver"
	router "github.com/Evgen-Poloniy/chat-gateway/internal/transport/http"
	v1 "github.com/Evgen-Poloniy/chat-gateway/internal/transport/http/v1"
	"github.com/Evgen-Poloniy/chat-gateway/internal/transport/ws"
	kfpc "github.com/Evgen-Poloniy/chat-gateway/pkg/database/kafka"
	"github.com/Evgen-Poloniy/chat-gateway/pkg/database/postgres"
	"github.com/Evgen-Poloniy/chat-gateway/pkg/database/redis"
	"github.com/Evgen-Poloniy/chat-gateway/pkg/logs"
	"github.com/sirupsen/logrus"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Run() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		logrus.Fatalf("error when loading env CONFIG_PATH")
	}

	config, err := config.LoadConfig(configPath)
	if err != nil {
		logrus.Fatalf("error when loading config: %v", err)
	}

	logger := logs.NewLogrusLogger(
		logs.WithLevel(config.Logger.Level),
		logs.WithFormat(config.Logger.Format),
	)

	jwks, err := middleware.InitKeyfunc(config.Auth.JWKSURL, logger)
	if err != nil {
		logger.Fatalf("keyfunc initialization error: %v", err)
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Postgres.Host,
		config.Postgres.Port,
		config.Postgres.Username,
		config.Postgres.Password,
		config.Postgres.DBName,
		config.Postgres.SSLMode,
	)

	db, err := postgres.NewPostgreSQL(
		dsn,
		postgres.WithMaxOpenConns(config.Postgres.MaxOpenConns),
		postgres.WithMaxIdleConns(config.Postgres.MaxIdleConns),
		postgres.WithConnMaxLifetime(config.Postgres.ConnMaxLifetime),
		postgres.WithConnMaxIdleLifetime(config.Postgres.ConnMaxIdleLifetime),
	)
	if err != nil {
		logger.Fatalf("database error: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Errorf("database error: %v", err)
		}
	}()

	kafkaCfg := &kfpc.Config{
		BootstrapServers: strings.Split(config.Kafka.BootstrapServers, ","),
		Acks:             config.Kafka.Acks,
		Retries:          config.Kafka.Retries,
		LingerMs:         config.Kafka.LingerMs,
		BatchNumMessages: config.Kafka.BatchNumMessages,
		CompressionType:  config.Kafka.CompressionType,
		MessageTimeoutMs: config.Kafka.MessageTimeout,
	}

	producer, err := kfpc.NewKafkaProducer(
		kafkaCfg,
		kfpc.WithSecurityProtocol(config.Kafka.SecurityProtocol),
		kfpc.WithSASLMechanism(config.Kafka.SASLMechanism),
		kfpc.WithSASLUsername(config.Kafka.User),
		kfpc.WithSASLPassword(config.Kafka.Password),
	)
	if err != nil {
		logger.Fatalf("message broker error: %v", err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			logger.Errorf("error when closing message broker: %v", err)
		}
	}()

	rdb, err := redis.NewRedisCache(
		redis.WithAddr(config.Redis.Host, config.Redis.Port),
		redis.WithPassword(config.Redis.Password),
		redis.WithDialTimeout(config.Redis.DialTimeout),
		redis.WithReadTimeout(config.Redis.ReadTimeout),
		redis.WithWriteTimeout(config.Redis.WriteTimeout),
	)
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
	messageBroker := kfk.NewKafkaRepository(producer)

	cache := rds.NewRedisCache(rdb, &config.Redis, &config.Resolver)

	messenger := messenger.NewMessengerService(messengerRepository, messageBroker, cache)
	resolver := resolver.NewResolverService(cache, cache)

	address := fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)
	wsHub := ws.NewHub(resolver, &config.Dispatch, address, logger)
	v1Handler := v1.NewHandler(messenger, wsHub, jwks, &config.Auth, logger)
	wsHandler := ws.NewHandler(wsHub, jwks, &config.Auth, logger)
	router := router.NewRouter(&config.CORS, logger)
	v1.NewRouter(router, v1Handler)
	ws.NewRouter(router, wsHandler)

	logger.Info("starting dispatch messages")
	if err := wsHub.StartDispatchMessage(); err != nil {
		logger.Fatalf("start dispatch error: %v", err)
	}

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

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	ctx, cancel := context.WithTimeout(
		context.Background(),
		config.Server.TimeForGracefulShutdown*time.Second,
	)
	defer cancel()

	logger.Info("shutting down dispatch messages workers")
	if err := wsHub.ShutdownDispatchMessage(ctx); err != nil {
		logger.Errorf("dispatch messages workers forced to shutdown: %v", err)
	}

	logger.Info("shutting down server")
	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Errorf("server forced to shutdown: %v", err)
	}

	wg.Wait()

	logger.Info("server stop")
}
