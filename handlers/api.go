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
        logger.Error("Ошибка при запросе возраста:", err.Error())
        return 0, err
    }
    defer resp.Body.Close()

    var data map[string]interface{}
    err = json.NewDecoder(resp.Body).Decode(&data)
    if err != nil {
        logger.Error("Ошибка при разборе JSON ответа:", err.Error())
        return 0, err
    }

    if ageValue, ok := data["age"]; ok {
        age, ok := ageValue.(float64)
        if !ok {
            return 0, errors.New("неверный тип данных для возраста")
        }
        return uint8(age), nil
    }
    return 0, errors.New("возраст не найден")
}

func GetGender(name string) (string, error) {
    url := fmt.Sprintf("https://api.genderize.io/?name=%s", name)
    resp, err := http.Get(url)
    if err != nil {
        logger.Error("Ошибка при запросе пола:", err.Error())
        return "", err
    }
    defer resp.Body.Close()

    var data map[string]interface{}
    err = json.NewDecoder(resp.Body).Decode(&data)
    if err != nil {
        logger.Error("Ошибка при разборе JSON ответа:", err.Error())
        return "", err
    }

    if genderValue, ok := data["gender"]; ok {
        gender, ok := genderValue.(string)
        if !ok {
            return "", errors.New("неверный тип данных для пола")
        }
        return gender, nil
    }
    return "", errors.New("пол не найден")
}


// Поскольку api.nationalize.io возвращает массив с некими значениями вероятностей, то в национальность
// запишем наибольшее значение "country_id"
func GetNationality(name string) (string, error) {
    resp, err := http.Get(fmt.Sprintf("https://api.nationalize.io/?name=%s", name))
    if err != nil {
        logger.Error("Ошибка при запросе национальности:", err.Error())
        return "", err
    }
    defer resp.Body.Close()

    var data map[string]interface{}
    err = json.NewDecoder(resp.Body).Decode(&data)
    if err != nil {
        logger.Error("Ошибка при разборе JSON ответа:", err.Error())
        return "", err
    }

    countries, ok := data["country"].([]interface{})
    if !ok || len(countries) == 0 {
        return "", errors.New("данные о странах не найдены")
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
        return "", errors.New("не удалось найти наиболее вероятную страну")
    }

    return mostProbableCountryID, nil
}

