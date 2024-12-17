package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	migratepgx "github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file" // Для работы с локальными файлами миграций
	_ "github.com/jackc/pgx/v5/stdlib"                   // Драйвер pgx для database/sql

	httpGetDevice "smart-device-service/internal/http/get_device"
	httpUpdateStatus "smart-device-service/internal/http/update_status"
	"smart-device-service/internal/subscribers"
	"smart-device-service/internal/usecases/get_device"
	getDeviceStorage "smart-device-service/internal/usecases/get_device/storage"
	"smart-device-service/internal/usecases/send_command"
	"smart-device-service/internal/usecases/update_status"
	"smart-device-service/internal/usecases/update_status/events_sender"
	updateStatusStorage "smart-device-service/internal/usecases/update_status/storage"
	"smart-device-service/internal/usecases/update_status/wrappers"
)

const (
	readHeaderTimeout           = 5 * time.Second
	shutdownTimeout             = 2 * time.Second
	signalChannelBufferCapacity = 1
	outboxTimeout               = 2 * time.Second
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbURL := os.Getenv("POSTGRESQL_URL")
	kafkaBrokerURL := os.Getenv("KAFKA_BROKER_URL")

	var err error
	db, err := sql.Open("pgx", dbURL) // Используем драйвер pgx
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		cancel()
		return
	}
	defer db.Close()
	runMigrations(db)

	// run subscribers
	commandsHandler := send_command.NewSendCommandUsecase()
	commandSub := subscribers.NewCommandSubscriber(kafkaBrokerURL, commandsHandler)
	defer func() { _ = commandSub.Stop() }()

	go func() {
		if err := commandSub.Run(ctx); err != nil {
			fmt.Printf("command subscriber error: %v", err)
		}
	}()

	router := gin.Default()

	statusChangePublisher := wrappers.NewStatusChangedPublisher(kafkaBrokerURL)
	updateStatusStorage := updateStatusStorage.New(db)

	// Запуск transactional outbox
	outboxSender := events_sender.NewSender(updateStatusStorage, statusChangePublisher)
	outboxSender.StartProcessEvents(ctx, outboxTimeout)

	updateStatusUsecase := update_status.NewUpdateStatusUsecase(updateStatusStorage)
	updateStatusHandler := httpUpdateStatus.NewHandler(updateStatusUsecase)

	getDeviceStorage := getDeviceStorage.New(db)
	getDeviceUsecase := get_device.NewGetDeviceUsecase(getDeviceStorage)
	getDeviceHandler := httpGetDevice.NewHandler(getDeviceUsecase)

	router.PUT("/devices/:deviceId/status", updateStatusHandler.Handle)
	router.GET("/devices/:deviceId", getDeviceHandler.Handle)

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to Smart Device Service!",
		})
	})
	srv := &http.Server{
		Addr:              ":8088",
		Handler:           router,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("Error starting server: %v\n", err)
		}
	}()

	sigchan := make(chan os.Signal, signalChannelBufferCapacity)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// Ожидание сигнала от системы (например, при завершении работы)
	<-sigchan
	fmt.Println("Shutdown signal received, initiating graceful shutdown...")

	// Отмена контекста при получении сигнала
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, shutdownTimeout)
	defer shutdownCancel()

	// Ожидаем завершения работы сервера с таймаутом
	if err := srv.Shutdown(shutdownCtx); err != nil {
		fmt.Printf("Server forced to shutdown: %v\n", err)
	} else {
		fmt.Println("Server gracefully shut down")
	}
}

func runMigrations(db *sql.DB) {
	// Создаём экземпляр драйвера базы данных для миграций
	driver, err := migratepgx.WithInstance(db, &migratepgx.Config{})
	if err != nil {
		log.Fatalf("Failed to create migration driver: %v", err)
	}

	// Настраиваем миграции
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations", // Путь к папке с миграциями
		"master",            // Имя базы данных
		driver,
	)
	if err != nil {
		log.Fatalf("Failed to create migration instance: %v", err)
	}

	// Применяем миграции
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Migration failed: %v", err)
	}
}
