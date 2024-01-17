package handlers

import (
	"net/http"

	"github.com/SergeyMilch/service-get-people/models"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)


func GetPersonByID(db *sqlx.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        id := c.Param("id")
        var person models.Person
        err := db.Get(&person, "SELECT * FROM people WHERE id = $1", id)
        if err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
            return
        }
        c.JSON(http.StatusOK, person)
    }
}