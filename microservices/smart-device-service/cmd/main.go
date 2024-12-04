package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	migratepgx "github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file" // Для работы с локальными файлами миграций
	_ "github.com/jackc/pgx/v5/stdlib"                   // Драйвер pgx для database/sql
	httpUpdateStatus "smart-device-service/internal/http/update_status"

	"smart-device-service/internal/usecases/publish_status_changed"
	"smart-device-service/internal/usecases/update_status"
)

func main() {
	dbURL := os.Getenv("POSTGRESQL_URL")
	kafkaBrokerURL := os.Getenv("KAFKA_BROKER_URL")

	var err error
	db, err := sql.Open("pgx", dbURL) // Используем драйвер pgx
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	fmt.Println("ping", db.Ping())
	runMigrations(db)

	router := gin.Default()

	statusChangePublisher := publish_status_changed.NewStatusChangedPublisher(kafkaBrokerURL)

	updateStatusUsecase := update_status.NewUpdateStatusUsecase(statusChangePublisher)
	updateStatusHandler := httpUpdateStatus.NewHandler(updateStatusUsecase)

	router.PUT("/devices/:deviceId/status", updateStatusHandler.Handle)

	log.Println("Starting server on :8088")
	if err := router.Run(":8088"); err != nil {
		log.Fatalf("Error starting server: %v", err)
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
