package handlers

import (
	"database/sql"
	"net/http"

	"github.com/SergeyMilch/service-get-people/models"
	"github.com/SergeyMilch/service-get-people/utils/logger"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func AddPerson(db *sqlx.DB) gin.HandlerFunc {
    return func(c *gin.Context) {

        logger.Debug("Начало обработки запроса на добавление пользователя", "AddPerson")

        var person models.Person
        if err := c.ShouldBindJSON(&person); err != nil {
            logger.Warn("Ошибка при разборе JSON: ", err.Error())
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
            return
        }

		// Обогащение данных
        // Каналы для результатов API и ошибок
        ageChan := make(chan int)
        genderChan := make(chan string)
        nationalityChan := make(chan string)
        errChan := make(chan error, 3) // Буферизированный канал для 3 потенциальных ошибок

        // Асинхронный вызов GetAge
        go func() {
            age, err := GetAge(person.UserName)
            if err != nil {
                errChan <- err
                return
            }
            ageChan <- int(age)
        }()

        // Асинхронный вызов GetGender
        go func() {
            gender, err := GetGender(person.UserName)
            if err != nil {
                errChan <- err
                return
            }
            genderChan <- gender
        }()

        // Асинхронный вызов GetNationality
        go func() {
            nationality, err := GetNationality(person.UserName)
            if err != nil {
                errChan <- err
                return
            }
            nationalityChan <- nationality
        }()

        // Ждем результаты или ошибки от API
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

        // Проверяем, существует ли уже пользователь с таким ID (т.к. уникальных полей (email и т.п.), при создании пользователя нет, то проверим по ФИО)
        var existingUser models.Person
        err := db.Get(&existingUser, "SELECT * FROM people WHERE name = $1 AND surname = $2 AND patronymic = $3", person.UserName, person.Surname, person.Patronymic)
        if err == nil {
            // Пользователь существует, проверяем и обновляем недостающие данные
            update := false
            if existingUser.Age == 0 {
                if age, ageErr := GetAge(person.UserName); ageErr == nil {
                    existingUser.Age = int(age)
                    update = true
                }
            }
            if existingUser.Gender == "" {
                if gender, genderErr := GetGender(person.UserName); genderErr == nil {
                    existingUser.Gender = gender
                    update = true
                }
            }
            if existingUser.Nationality == "" {
                if nationality, natErr := GetNationality(person.UserName); natErr == nil {
                    existingUser.Nationality = nationality
                    update = true
                }
            }

            // Если есть что обновлять, выполняем запрос UPDATE
            if update {
                // Обновляем данные пользователя в базе
                _, err = db.Exec("UPDATE people SET age = $1, gender = $2, nationality = $3 WHERE id = $4",
                    existingUser.Age, existingUser.Gender, existingUser.Nationality, existingUser.ID)
                if err != nil {
                    logger.Error("Ошибка при обновлении данных пользователя: ", err.Error())
                    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                    return
                }
                c.JSON(http.StatusOK, gin.H{"message": "Данные пользователя обновлены"})
            } else {
                // Пользователь существует, но данные не обновлены
                c.JSON(http.StatusOK, gin.H{"message": "Пользователь уже существует"})
            }
            return
            } else if err != sql.ErrNoRows {
                // Ошибка при запросе к БД
                logger.Error("Ошибка при запросе к БД: ", err.Error())
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }
        
        // Если пользователь не найден, добавляем его в базу данных
        var newID uint
        err = db.QueryRow("INSERT INTO people (name, surname, patronymic, age, gender, nationality) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
            person.UserName, person.Surname, person.Patronymic, person.Age, person.Gender, person.Nationality).Scan(&newID)
        if err != nil {
            logger.Error("Ошибка при добавлении пользователя в БД: ", err.Error())
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        // Обновляем person, включая новый ID
        person.ID = newID

        logger.Debug("Завершение обработки запроса на добавление пользователя", "AddPerson")

        c.JSON(http.StatusOK, person)
    }
}
