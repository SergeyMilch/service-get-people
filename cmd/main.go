package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SergeyMilch/service-get-people/handlers"
	"github.com/SergeyMilch/service-get-people/utils/db"
	"github.com/SergeyMilch/service-get-people/utils/logger"
	"github.com/jmoiron/sqlx"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)


func setupRouter(db *sqlx.DB) *gin.Engine {
    router := gin.Default()

    router.GET("/people", handlers.GetPeople(db))
    router.GET("/people/:id", handlers.GetPersonByID(db))
    router.POST("/people", handlers.AddPerson(db))
    router.DELETE("/people/:id", handlers.DeletePerson(db))
    router.PUT("/people/:id", handlers.UpdatePerson(db))

    return router
}

func main() {
    // Логирование начала выполнения
    logger.Init()
    logger.Info("Запуск приложения", "")

    // Загрузка .env файла
    if err := godotenv.Load(); err != nil {
        logger.Error("Ошибка при загрузке .env файла", err.Error())
        log.Fatal("Error loading .env file")
    }

    // Подключение к базе данных
    dbConn, err := sqlx.Connect("postgres", os.Getenv("DB_URL"))
    if err != nil {
        logger.Error("Ошибка подключения к базе данных", err.Error())
        log.Fatal("Error connecting to the database: ", err.Error())
    }
    defer dbConn.Close()

    // Логирование успешного подключения к БД
    logger.Info("Успешное подключение к базе данных", "")

    // Инициализация и выполнение миграций БД
    db.ExecMigrations(dbConn)

    // Настройка и запуск HTTP-сервера
    server := &http.Server{
        Addr:    ":8080",
        Handler: setupRouter(dbConn),
    }

    go func() {
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            logger.Error("Ошибка запуска HTTP сервера", err.Error())
            log.Fatalf("Error starting HTTP server: %v", err)
        }
    }()

    // Обработка сигналов для корректного завершения работы
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    // Начало процесса корректного завершения работы
    logger.Info("Завершение работы сервера", "")

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := server.Shutdown(ctx); err != nil {
        logger.Error("Ошибка при завершении работы сервера", err.Error())
        log.Fatalf("Server forced to shutdown: %v", err)
    }

    logger.Info("Сервер корректно завершил работу", "")
}