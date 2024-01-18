package handlers

import (
	"fmt"
	"net/http"

	"github.com/SergeyMilch/service-get-people/models"
	"github.com/SergeyMilch/service-get-people/utils/logger"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)


func GetPersonByID(db *sqlx.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        id := c.Param("id")

        logger.Debug(fmt.Sprintf("Запрос к базе данных для получения пользователя с ID: %s", id), fmt.Sprintf("id = %s", id))

        var person models.Person
        err := db.Get(&person, "SELECT * FROM people WHERE id = $1", id)
        if err != nil {
            logger.Warn(fmt.Sprintf("Не найден пользователь с ID: %s", id), fmt.Sprintf("id = %s", id))
            c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
            return
        }

        logger.Debug(fmt.Sprintf("Получены данные пользователя с ID: %s", id), fmt.Sprintf("id = %s", id))

        c.JSON(http.StatusOK, person)
    }
}