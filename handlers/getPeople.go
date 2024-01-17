package handlers

import (
	"net/http"

	"github.com/SergeyMilch/service-get-people/models"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)


func GetPeople(db *sqlx.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var people []models.Person
        err := db.Select(&people, "SELECT * FROM people")
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, people)
    }
}