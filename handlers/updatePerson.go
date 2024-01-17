package handlers

import (
	"net/http"

	"github.com/SergeyMilch/service-get-people/models"
	"github.com/SergeyMilch/service-get-people/utils/logger"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)


func UpdatePerson(db *sqlx.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        id := c.Param("id")
        var person models.Person
        if err := c.ShouldBindJSON(&person); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
            return
        }

		// Обогащение данных, при ошибке используем значения по умолчанию
        if age, err := GetAge(person.UserName); err == nil {
            person.Age = int(age)
        } else {
            logger.Warn("Не удалось обогатить возраст: ", err.Error())
        }

        if gender, err := GetGender(person.UserName); err == nil {
            person.Gender = gender
        } else {
            logger.Warn("Не удалось обогатить пол: ", err.Error())
        }

        if nationality, err := GetNationality(person.UserName); err == nil {
            person.Nationality = nationality
        } else {
            logger.Warn("Не удалось обогатить национальность: ", err.Error())
        }
        _, err := db.Exec("UPDATE people SET name=$1, surname=$2, patronymic=$3, age=$4, gender=$5, nationality=$6 WHERE id=$7",
            person.UserName, person.Surname, person.Patronymic, person.Age, person.Gender, person.Nationality, id)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, person)
    }
}