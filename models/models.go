package models

type Person struct {
    ID          uint   `db:"id" json:"id"`
    UserName    string `db:"name" json:"name"`
    Surname     string `db:"surname" json:"surname"`
    Patronymic  string `db:"patronymic" json:"patronymic"`
    Age         int    `db:"age" json:"age"`
    Gender      string `db:"gender" json:"gender"`
    Nationality string `db:"nationality" json:"nationality"`
}