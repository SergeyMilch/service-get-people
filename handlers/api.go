package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/SergeyMilch/service-get-people/utils/logger"
)



func GetAge(name string) (uint8, error) {
    url := fmt.Sprintf("https://api.agify.io/?name=%s", name)
    resp, err := http.Get(url)
    if err != nil {
        logger.Warn("Ошибка при запросе возраста:", err.Error())
        return 0, err
    }
    defer resp.Body.Close()

    var data map[string]interface{}
    err = json.NewDecoder(resp.Body).Decode(&data)
    if err != nil {
        logger.Warn("Ошибка при разборе JSON ответа:", err.Error())
        return 0, err
    }

    if ageValue, ok := data["age"]; ok {
        age, ok := ageValue.(float64)
        if !ok {
            return 0, errors.New("неверный тип данных для возраста")
        }
        return uint8(age), nil
    }
    return 0, errors.New("age not found")
}

func GetGender(name string) (string, error) {
    url := fmt.Sprintf("https://api.genderize.io/?name=%s", name)
    resp, err := http.Get(url)
    if err != nil {
        logger.Warn("Ошибка при запросе пола:", err.Error())
        return "", err
    }
    defer resp.Body.Close()

    var data map[string]interface{}
    err = json.NewDecoder(resp.Body).Decode(&data)
    if err != nil {
        logger.Warn("Ошибка при разборе JSON ответа:", err.Error())
        return "", err
    }

    if genderValue, ok := data["gender"]; ok {
        gender, ok := genderValue.(string)
        if !ok {
            return "", errors.New("неверный тип данных для пола")
        }
        return gender, nil
    }
    return "", errors.New("gender not found")
}


// Поскольку api.nationalize.io возвращает массив с некими значениями вероятностей, то в национальность
// запишем наибольшее значение "country_id"
func GetNationality(name string) (string, error) {
    resp, err := http.Get(fmt.Sprintf("https://api.nationalize.io/?name=%s", name))
    if err != nil {
        logger.Warn("Ошибка при запросе национальности:", err.Error())
        return "", err
    }
    defer resp.Body.Close()

    var data map[string]interface{}
    err = json.NewDecoder(resp.Body).Decode(&data)
    if err != nil {
        logger.Warn("Ошибка при разборе JSON ответа:", err.Error())
        return "", err
    }

    countries, ok := data["country"].([]interface{})
    if !ok || len(countries) == 0 {
        return "", errors.New("country data not found")
    }

    var mostProbableCountryID string
    var maxProbability float64 = 0

    for _, country := range countries {
        countryData, ok := country.(map[string]interface{})
        if !ok {
            continue
        }

        probability, ok := countryData["probability"].(float64)
        if ok && probability > maxProbability {
            maxProbability = probability
            mostProbableCountryID, _ = countryData["country_id"].(string)
        }
    }

    if mostProbableCountryID == "" {
        return "", errors.New("failed to find the most likely country")
    }

    return mostProbableCountryID, nil
}

