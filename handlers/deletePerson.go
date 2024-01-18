package handlers

import (
	"fmt"
	"net/http"

	"github.com/SergeyMilch/service-get-people/utils/logger"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)


func DeletePerson(db *sqlx.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        id := c.Param("id")

        logger.Info(fmt.Sprintf("Попытка удаления пользователя с ID: %s", id), fmt.Sprintf("id = %s", id))

        // Проверка наличия пользователя
        var exists bool
        err := db.Get(&exists, "SELECT EXISTS(SELECT 1 FROM people WHERE id = $1)", id)
        if err != nil {
            logger.Error("Ошибка при проверке существования пользователя: ", err.Error())
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        if !exists {
            logger.Info(fmt.Sprintf("Пользователь с ID: %s не найден", id), fmt.Sprintf("id = %s", id))
            c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
            return
        }

        // Удаление пользователя
        _, err = db.Exec("DELETE FROM people WHERE id = $1", id)
        if err != nil {
            logger.Error("Ошибка при удалении пользователя: ", err.Error())
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        logger.Info(fmt.Sprintf("Успешное удаление пользователя с ID: %s", id), fmt.Sprintf("id = %s", id))
        c.JSON(http.StatusOK, gin.H{"message": "Пользователь успешно удален", "id": id})
    }
}