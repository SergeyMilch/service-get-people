package handlers

import (
	"net/http"

	"github.com/SergeyMilch/service-get-people/models"
	"github.com/SergeyMilch/service-get-people/utils/logger"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func AddPerson(db *sqlx.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
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

        // Проверяем, существует ли уже пользователь с таким ID
        var existingUser models.Person
        err := db.Get(&existingUser, "SELECT id FROM people WHERE id = $1", person.ID)
        if err == nil {
            logger.Warn("Пользователь уже существует", err.Error())
            c.JSON(http.StatusBadRequest, gin.H{"error": "Пользователь уже существует"})
            return
        }

        // Если пользователь не найден, добавляем его в базу данных
        _, err = db.Exec("INSERT INTO people (name, surname, patronymic, age, gender, nationality) VALUES ($1, $2, $3, $4, $5, $6)",
            person.UserName, person.Surname, person.Patronymic, person.Age, person.Gender, person.Nationality)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, person)
    }
}