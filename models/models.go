package models

type Person struct {
    ID          uint   `json:"id"`
    UserName    string `json:"name"`
    Surname     string `json:"surname"`
    Patronymic  string `json:"patronymic"`
    Age         int    `json:"age"`
    Gender      string `json:"gender"`
    Nationality string `json:"nationality"`
}