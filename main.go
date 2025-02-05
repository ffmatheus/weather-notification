package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"weather-notification/internal/domain/service"
	middleware "weather-notification/internal/middlewares"

	_ "weather-notification/docs"
	"weather-notification/internal/infrastructure/adapter/api/handler"
	"weather-notification/internal/infrastructure/adapter/cptec"
	postgres "weather-notification/internal/infrastructure/adapter/persistence/postgres"
	"weather-notification/internal/infrastructure/adapter/queue"
	"weather-notification/internal/infrastructure/worker"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title API de Notificação de Previsão do Tempo
// @version 1.0
// @description Serviço de notificações de previsão do tempo
// @host localhost:8080
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @BasePath /api
func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Arquivo .env não encontrado, usando variáveis do sistema")
	}

	//DATABASE CONFIG
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Erro ao pingar banco: %v", err)
	}

	// REPOSITORIES
	userRepo := postgres.NewUserRepository(db)
	locationRepo := postgres.NewLocationRepository(db)
	notificationRepo := postgres.NewNotificationRepository(db)
	globalNotificationRepo := postgres.NewGlobalNotificationRepository(db)

	// ADAPTERS
	cptecClient := cptec.NewClient()

	queueService, err := queue.NewRabbitMQService(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		log.Fatalf("Erro ao conectar ao RabbitMQ: %v", err)
	}
	defer queueService.Close()

	// SERVICES
	weatherService := service.NewWeatherService(cptecClient, locationRepo)
	notificationService := service.NewNotificationService(
		notificationRepo,
		userRepo,
		weatherService,
		queueService,
	)
	globalNotificationService := service.NewGlobalNotificationService(
		globalNotificationRepo,
		userRepo,
		queueService,
		weatherService,
		notificationRepo,
	)
	userService := service.NewUserService(userRepo)

	// WORKERS
	notificationWorker := worker.NewNotificationWorker(
		context.Background(),
		notificationService,
		weatherService,
		queueService,
	)
	globalWorker := worker.NewGlobalNotificationWorker(context.Background(), globalNotificationService)

	go func() {
		if err := notificationWorker.Start(); err != nil {
			log.Printf("Erro no worker: %v", err)
		}
	}()

	go func() {
		if err := globalWorker.Start(); err != nil {
			log.Printf("Erro no worker de notificações globais: %v", err)
		}
	}()

	// API
	forecastHandler := handler.NewWeatherHandler(weatherService)
	notificationHandler := handler.NewNotificationHandler(notificationService)
	globalNotificationHandler := handler.NewGlobalNotificationHandler(globalNotificationService)
	userHandler := handler.NewUserHandler(userService, weatherService)
	webhookHandler := handler.NewWebhookHandler()

	gin.SetMode(os.Getenv("GIN_MODE"))
	router := gin.Default()

	router.Use(cors.Default())

	// SWAGGO
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api", middleware.AuthMiddleware())
	{
		forecastHandler.SetupRoutes(api)
		notificationHandler.SetupRoutes(api)
		globalNotificationHandler.SetupRoutes(api)
		userHandler.SetupRoutes(api)
		webhookHandler.SetupRoutes(api)
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	})

	srv := &http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: router,
	}

	go func() {
		log.Printf("Servidor iniciado na porta %s", os.Getenv("PORT"))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erro ao iniciar servidor: %v", err)
		}
	}()

	// SHUTDOWN
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Desligando servidor...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forçado a fechar:", err)
	}

	log.Println("Servidor encerrado")
}
