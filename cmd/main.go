package main

import (
	"github.com/ZnNr/songs-library/internal/repository/config"
	"github.com/ZnNr/songs-library/migration"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
)

var (
	cfgPath = "config/config.yaml"
)

func main() {
	logger := initializeLogger()
	defer func() {
		_ = logger.Sync()
	}()

	cfg := loadConfig(cfgPath, logger)

	ordersRepo := initializeRepository(cfg, logger)
	defer closeRepository(ordersRepo, logger)
	migration.InitializeDatabaseSchema(songRepo.DB, logger)
	appCache := initializeCache(songRepo, logger)

	server := initializeController(cfgPath, appCache, logger)
	startServer(server, logger)

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)

	logger.Info("Application shutting down")
}

func initializeLogger() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	return logger
}

func loadConfig(cfgPath string, logger *zap.Logger) *config.Config {
	cfg, err := config.Load(cfgPath)
	if err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
	}
	logger.Info("Configuration loaded successfully")
	return cfg
}

func initializeRepository(cfg *config.Config, logger *zap.Logger) *repository.OrdersRepo {
	ordersRepo, err := repository.New(cfg)
	if err != nil {
		logger.Fatal("Failed to initialize repository", zap.Error(err))
	}
	logger.Info("Repository initialized successfully",
		zap.String("host", cfg.DB.Host),
		zap.String("port", cfg.DB.Port),
		zap.String("db", cfg.DB.Name),
		zap.String("user", cfg.DB.User),
	)
	return ordersRepo
}

func closeRepository(repo *repository.OrdersRepo, logger *zap.Logger) {
	if err := repo.DB.Close(); err != nil {
		logger.Error("Error closing repository", zap.Error(err))
	} else {
		logger.Info("Repository closed successfully")
	}
}

func initializeCache(ordersRepo *repository.OrdersRepo, logger *zap.Logger) *cache.Cache {
	appCache := cache.New(100)

	orders, err := ordersRepo.GetOrders()
	if err != nil {
		logger.Fatal("Orders Load error", zap.Error(err))
	}

	logger.Info("Cache initialized successfully")
	for _, order := range orders {
		appCache.SaveOrder(order)
		logger.Info("order cached successfully", zap.String("order_uid", order.OrderUID))
	}

	return appCache
}

func initializeController(cfgPath string, cache *cache.Cache, logger *zap.Logger) *server.Server {
	server, err := server.New(cfgPath, cache)
	if err != nil {
		logger.Fatal("Controller initialization error", zap.Error(err))
	}
	logger.Info("Controller initialized successfully")
	return server
}

func startServer(server *server.Server, logger *zap.Logger) {

	go func() {
		if err := server.Launch(); err != nil {
			logger.Fatal("Server error", zap.Error(err))
		}
	}()

	logger.Info("Server started successfully")
}
