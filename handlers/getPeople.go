package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/SergeyMilch/service-get-people/models"
	"github.com/SergeyMilch/service-get-people/utils/logger"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)


func GetPeople(db *sqlx.DB) gin.HandlerFunc {
    return func(c *gin.Context) {

        logger.Debug("Начало обработки запроса на получение всех пользователей", "func GetPeople")

        // Получение и обработка параметров страницы и лимита из запроса
        page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
        if err != nil || page < 1 {
            page = 1
        }
        limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
        if err != nil || limit <= 0 {
            limit = 10
        }

        // Получение параметров фильтрации
        nameFilter := c.Query("name")
        surnameFilter := c.Query("surname")

        // Формирование запроса с учетом фильтрации
        baseQuery := "FROM people"
        var args []interface{}
        var conditions []string

        if nameFilter != "" {
            conditions = append(conditions, "name = ?")
            args = append(args, nameFilter)
        }
        if surnameFilter != "" {
            conditions = append(conditions, "surname = ?")
            args = append(args, surnameFilter)
        }

        if len(conditions) > 0 {
            baseQuery += " WHERE " + strings.Join(conditions, " AND ")
        }

        var people []models.Person
        query := "SELECT * " + baseQuery + " LIMIT ? OFFSET ?"
        args = append(args, limit, (page-1)*limit)

        err = db.Select(&people, query, args...)
        if err != nil {
            logger.Error("Ошибка при получении пользователей из базы", err.Error())
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        // Получение общего количества записей с учетом фильтров для вычисления общего числа страниц
        var total int
        countQuery := "SELECT COUNT(*) " + baseQuery
        db.Get(&total, countQuery, args[:len(args)-2]...) // Исключаем параметры пагинации

        logger.Debug("Завершение обработки запроса на получение всех пользователей", "func GetPeople")

        // Отправка данных с информацией о пагинации
        c.JSON(http.StatusOK, gin.H{
            "data":       people,
            "total":      total,
            "page":       page,
            "last_page":  (total + limit - 1) / limit,
        })
    }
}