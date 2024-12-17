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
	_ "github.com/jackc/pgx/v5/stdlib"

	getDeviceTelemetryHandler "smart-telemetry-service/internal/http/get_device_telemetry"
	"smart-telemetry-service/internal/subscribers"
	"smart-telemetry-service/internal/usecases/get_device_telemetry"
	getEventsStorage "smart-telemetry-service/internal/usecases/get_device_telemetry/storage"
	"smart-telemetry-service/internal/usecases/handle_sensor_event"
	handleSensorEventStorage "smart-telemetry-service/internal/usecases/handle_sensor_event/storage"
)

const (
	readHeaderTimeout           = 5 * time.Second
	shutdownTimeout             = 2 * time.Second
	signalChannelBufferCapacity = 1
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
	defer func() { _ = db.Close() }()
	runMigrations(db)

	// run subscribers
	telemetryStotage := handleSensorEventStorage.NewStorage(db)
	sensorDataHandler := handle_sensor_event.NewHandleSensorEventUsecase(telemetryStotage)
	sensorSub := subscribers.NewSensorDataSubscriber(kafkaBrokerURL, sensorDataHandler)
	defer func() { _ = sensorSub.Stop() }()

	go func() {
		if err := sensorSub.Run(ctx); err != nil {
			fmt.Printf("sensor subscriber error: %v", err)
		}
	}()

	// Usecases.
	getTelemetryStorage := getEventsStorage.New(db)
	getDeviceTelemetryUsecase := get_device_telemetry.NewGetEventsUsecase(getTelemetryStorage)
	getDeviceTelemetryApi := getDeviceTelemetryHandler.NewHandler(getDeviceTelemetryUsecase)

	router := gin.Default()
	router.GET("/telemetry/devices/:deviceId", getDeviceTelemetryApi.Handle)
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to Smart Telemetry Service!",
		})
	})
	srv := &http.Server{
		Addr:              ":8089",
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
