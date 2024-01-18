package handlers

import (
	"fmt"
	"net/http"

	"github.com/SergeyMilch/service-get-people/models"
	"github.com/SergeyMilch/service-get-people/utils/logger"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)


func UpdatePerson(db *sqlx.DB) gin.HandlerFunc {
    return func(c *gin.Context) {

        id := c.Param("id")

        logger.Debug(fmt.Sprintf("Начало обработки запроса на обновление пользователя с ID: %s", id), fmt.Sprintf("id = %s", id))

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

        var person models.Person
        if err := c.ShouldBindJSON(&person); err != nil {
            logger.Warn("Ошибка при разборе JSON: ", err.Error())
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
            return
        }

        // Запросы к API
        ageChan := make(chan int)
        genderChan := make(chan string)
        nationalityChan := make(chan string)
        errChan := make(chan error, 3)

        go func() {
            age, err := GetAge(person.UserName)
            if err != nil {
                errChan <- err
                return
            }
            ageChan <- int(age)
        }()

        go func() {
            gender, err := GetGender(person.UserName)
            if err != nil {
                errChan <- err
                return
            }
            genderChan <- gender
        }()

        go func() {
            nationality, err := GetNationality(person.UserName)
            if err != nil {
                errChan <- err
                return
            }
            nationalityChan <- nationality
        }()

        // Обработка результатов запросов
        for i := 0; i < 3; i++ {
            select {
            case age := <-ageChan:
                person.Age = age
            case gender := <-genderChan:
                person.Gender = gender
            case nationality := <-nationalityChan:
                person.Nationality = nationality
            case err := <-errChan:
                logger.Warn("Ошибка при вызове API: ", err.Error())
            }
        }

        // Обновление данных пользователя
        _, err = db.Exec("UPDATE people SET name=$1, surname=$2, patronymic=$3, age=$4, gender=$5, nationality=$6 WHERE id=$7",
            person.UserName, person.Surname, person.Patronymic, person.Age, person.Gender, person.Nationality, id)
        if err != nil {
            logger.Error("Ошибка при обновлении данных пользователя: ", err.Error())
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        logger.Debug(fmt.Sprintf("Завершение обработки запроса на обновление пользователя с ID: %s", id), fmt.Sprintf("id = %s", id))

        // Возвращаем обновлённые данные пользователя, включая его ID
        c.JSON(http.StatusOK, gin.H{
        "id":          id,
        "name":        person.UserName,
        "surname":     person.Surname,
        "patronymic":  person.Patronymic,
        "age":         person.Age,
        "gender":      person.Gender,
        "nationality": person.Nationality,
    })
    }
}
