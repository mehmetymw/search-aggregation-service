package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/mehmetymw/search-aggregation-service/backend/application/usecase"
	"github.com/mehmetymw/search-aggregation-service/backend/domain/entity"
	"github.com/mehmetymw/search-aggregation-service/backend/domain/service"
	"github.com/mehmetymw/search-aggregation-service/backend/infrastructure/cache"
	"github.com/mehmetymw/search-aggregation-service/backend/infrastructure/config"
	"github.com/mehmetymw/search-aggregation-service/backend/infrastructure/db"
	loggerPkg "github.com/mehmetymw/search-aggregation-service/backend/infrastructure/logger"
	"github.com/mehmetymw/search-aggregation-service/backend/infrastructure/providers"
	"github.com/mehmetymw/search-aggregation-service/backend/infrastructure/repositories"
	"github.com/mehmetymw/search-aggregation-service/backend/infrastructure/resilience"
	contentpb "github.com/mehmetymw/search-aggregation-service/backend/proto/gen"
	grpcTransport "github.com/mehmetymw/search-aggregation-service/backend/transport/grpc"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger, err := loggerPkg.NewZapLogger()
	if err != nil {
		fmt.Printf("failed to create logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	configPath := getEnv("CONFIG_PATH", "config.yaml")
	configProvider, err := config.NewViperConfig(configPath)
	if err != nil {
		logger.Error("failed to load config", loggerPkg.Error(err))
		os.Exit(1)
	}

	appConfig := configProvider.GetAppConfig()
	logger.Info("configuration loaded", loggerPkg.String("config_path", configPath))

	database, err := db.NewPostgresConnection(appConfig.Database.DSN)
	if err != nil {
		logger.Error("failed to connect to database", loggerPkg.Error(err))
		os.Exit(1)
	}
	defer database.Close()
	logger.Info("connected to database")

	logger.Info("running database migrations")
	if err := db.AutoMigrate(database); err != nil {
		logger.Warn("migration failed (tables may already exist)", loggerPkg.Error(err))
	} else {
		logger.Info("migrations completed successfully")
	}

	cacheClient, err := cache.NewRedisCache(appConfig.Redis.Addr)
	if err != nil {
		logger.Error("failed to connect to redis", loggerPkg.Error(err))
		os.Exit(1)
	}
	logger.Info("connected to redis")

	contentRepo := repositories.NewContentRepository(database)
	contentStatsRepo := repositories.NewContentStatsRepository(database)
	providerRepo := repositories.NewProviderRepository(database)
	tagRepo := repositories.NewTagRepository(database)
	scoringRepo := repositories.NewScoringRepository(database)

	dbConfigProvider := config.NewDatabaseConfigProvider(configProvider, scoringRepo)

	scoringConfig := dbConfigProvider.GetScoringConfig()
	timeProvider := func() time.Time {
		return time.Now()
	}
	scoringService := service.NewScoringService(scoringConfig, timeProvider)

	searchUseCase := usecase.NewSearchContentsUseCase(
		contentRepo,
		contentStatsRepo,
		cacheClient,
		scoringService,
		logger,
		appConfig.Cache.GetTTL(),
	)

	getByIDUseCase := usecase.NewGetContentByIDUseCase(
		contentRepo,
		contentStatsRepo,
		scoringService,
	)

	jsonProviderClient := providers.NewJsonProviderClient()
	xmlProviderClient := providers.NewXmlProviderClient()

	// Wrap provider clients with Circuit Breaker
	cbConfig := appConfig.CircuitBreaker
	jsonProviderClientWithCB := resilience.NewCircuitBreakerProviderClient(jsonProviderClient, cbConfig)
	xmlProviderClientWithCB := resilience.NewCircuitBreakerProviderClient(xmlProviderClient, cbConfig)

	tagNormalizer := service.NewTagNormalizer()

	syncUseCase := usecase.NewSyncProviderContentsUseCase(
		providerRepo,
		contentRepo,
		contentStatsRepo,
		tagRepo,
		jsonProviderClientWithCB,
		xmlProviderClientWithCB,
		tagNormalizer,
		logger,
	)

	go startSyncWorker(ctx, syncUseCase, appConfig, logger)

	metadataRepo := repositories.NewMetadataRepository(database)

	// Initialize Rate Limiter
	rateLimitInterceptor := grpcTransport.NewRateLimitInterceptor(appConfig.RateLimit)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(rateLimitInterceptor.Unary()),
	)
	contentServer := grpcTransport.NewContentServiceServer(
		searchUseCase,
		getByIDUseCase,
		metadataRepo,
		*appConfig,
		logger,
	)
	contentpb.RegisterContentServiceServer(grpcServer, contentServer)

	grpcAddr := fmt.Sprintf(":%d", appConfig.Server.GRPCPort)
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		logger.Error("failed to listen on grpc port", loggerPkg.Error(err))
		os.Exit(1)
	}

	go func() {
		logger.Info("starting gRPC server", loggerPkg.Int("port", appConfig.Server.GRPCPort))
		if err := grpcServer.Serve(lis); err != nil {
			logger.Error("grpc server error", loggerPkg.Error(err))
		}
	}()

	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames:   true,
				EmitUnpopulated: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
	)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err = contentpb.RegisterContentServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%d", appConfig.Server.GRPCPort), opts)
	if err != nil {
		logger.Error("failed to register grpc-gateway", loggerPkg.Error(err))
		os.Exit(1)
	}

	httpMux := http.NewServeMux()
	httpMux.Handle("/api/", mux)
	httpMux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", appConfig.Server.HTTPPort),
		Handler: corsMiddleware(httpMux),
	}

	go func() {
		logger.Info("starting HTTP server", loggerPkg.Int("port", appConfig.Server.HTTPPort))
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("http server error", loggerPkg.Error(err))
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	logger.Info("shutting down servers...")
	cancel()

	httpServer.Close()
	grpcServer.GracefulStop()

	logger.Info("servers stopped")
}

func startSyncWorker(ctx context.Context, syncUseCase *usecase.SyncProviderContentsUseCase, config *entity.AppConfig, logger *loggerPkg.ZapLogger) {
	interval := config.Sync.GetInterval()
	logger.Info("starting sync worker", loggerPkg.String("interval", interval.String()))

	syncUseCase.ExecuteAll(ctx)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("sync worker stopped")
			return
		case <-ticker.C:
			if err := syncUseCase.ExecuteAll(ctx); err != nil {
				logger.Error("sync all failed", loggerPkg.Error(err))
			}
		}
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
