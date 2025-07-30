package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/company/microservice-template/internal/auth"
	"github.com/company/microservice-template/internal/config"
	"github.com/company/microservice-template/internal/domain"
	"github.com/company/microservice-template/internal/handlers"
	"github.com/company/microservice-template/internal/middleware"
	"github.com/company/microservice-template/internal/repositories"
	"github.com/company/microservice-template/internal/services"
	"github.com/company/microservice-template/pkg/logger"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

// @title Microservice Template API
// @version 1.0
// @description Template para microservicios en Go
// @host localhost:8080
// @BasePath /api/v1
func main() {
	// Cargar configuración
	cfg := config.Load()

	// Inicializar logger
	logger := logger.NewLogger(cfg.LogLevel)

	// Inicializar base de datos (no fatal si falla)
	var db *sql.DB
	var err error

	// Intentar conectar a la base de datos, pero no fallar si no puede
	if cfg.Database.Host != "" && cfg.Database.Password != "" {
		db, err = initDatabase(&cfg.Database, logger)
		if err != nil {
			logger.Error("Failed to initialize database, continuing without DB", err)
		} else {
			defer db.Close()
		}
	} else {
		logger.Info("Database configuration not complete, running without DB")
	}

	// Inicializar Redis (opcional)
	var redisClient *redis.Client
	if cfg.Redis.Enabled {
		redisClient = initRedis(&cfg.Redis, logger)
		defer redisClient.Close()
	}

	// Inicializar JWT manager
	jwtManager := auth.NewJWTManager(cfg.JWT.SecretKey, cfg.JWT.Issuer)

	// Inicializar repositorios (con manejo de DB nula)
	var conversationRepo domain.ConversationRepository
	var messageRepo domain.MessageRepository
	var attachmentRepo domain.AttachmentRepository

	if db != nil {
		conversationRepo = repositories.NewPostgresConversationRepository(db, logger)
		messageRepo = repositories.NewPostgresMessageRepository(db, logger)
		attachmentRepo = repositories.NewPostgresAttachmentRepository(db, logger)
	} else {
		// Usar repositorios mock/no-op cuando no hay DB
		conversationRepo = repositories.NewNoOpConversationRepository()
		messageRepo = repositories.NewNoOpMessageRepository()
		attachmentRepo = repositories.NewNoOpAttachmentRepository()
	}

	// Inicializar servicios auxiliares
	var cacheService services.CacheService
	if redisClient != nil {
		cacheService = services.NewRedisCacheService(redisClient, logger)
	} else {
		cacheService = services.NewNoOpCacheService()
	}

	var eventPublisher services.EventPublisher
	if redisClient != nil && cfg.Events.Provider == "redis" {
		eventPublisher = services.NewRedisEventPublisher(redisClient, cfg.Events.Topic, logger)
	} else {
		eventPublisher = services.NewNoOpEventPublisher()
	}

	fileService := services.NewLocalFileService(&cfg.FileStorage, logger)

	// Inicializar servicios principales
	healthService := services.NewHealthService()
	messagingService := services.NewMessagingService(
		conversationRepo,
		messageRepo,
		attachmentRepo,
		eventPublisher,
		cacheService,
		logger,
	)

	// Configurar Gin
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.Logger(logger))
	router.Use(middleware.CORS())
	router.Use(middleware.Metrics())

	// Rutas
	handlers.SetupRoutes(router, healthService, messagingService, fileService, jwtManager, logger)

	// Servidor HTTP
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	// Iniciar servidor en goroutine
	go func() {
		logger.Info("Starting server on port " + cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", err)
	}

	logger.Info("Server exited")
}

func initDatabase(dbCfg *config.Database, logger logger.Logger) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbCfg.Host,
		dbCfg.Port,
		dbCfg.User,
		dbCfg.Password,
		dbCfg.Name,
		dbCfg.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Create a context with a timeout for the ping to avoid long waits on startup
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Test connection
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	logger.Info("Database connection established successfully")
	return db, nil
}

func initRedis(redisCfg *config.Redis, logger logger.Logger) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisCfg.Host, redisCfg.Port),
		Password: redisCfg.Password,
		DB:       redisCfg.DB,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		logger.Error("Failed to connect to Redis", err)
		return nil
	}

	logger.Info("Redis connection established successfully")
	return client
}
