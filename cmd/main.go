package main

import (
	"log"
	"os"

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
    // Загрузка .env файла
    if err := godotenv.Load(); err != nil {
        log.Fatal("Error loading .env file")
    }

    logger.Init()

    dbConn, err := sqlx.Connect("postgres", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

    db.ExecMigrations(dbConn)

    router := setupRouter(dbConn)
    router.Run(":8080")
}