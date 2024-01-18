package handlers

import (
	"net/http"

	"github.com/SergeyMilch/service-get-people/models"
	"github.com/SergeyMilch/service-get-people/utils/logger"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)


func GetPeople(db *sqlx.DB) gin.HandlerFunc {
    return func(c *gin.Context) {

        logger.Debug("Начало обработки запроса на получение всех пользователей", "func GetPeople")

        var people []models.Person
        err := db.Select(&people, "SELECT * FROM people")
        if err != nil {
            logger.Error("Ошибка при получении пользователей из базы", err.Error())
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        logger.Debug("Завершение обработки запроса на получение всех пользователей", "func GetPeople")

        c.JSON(http.StatusOK, people)
    }
}