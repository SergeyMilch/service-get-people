CREATE TABLE IF NOT EXISTS people (
    id BIGSERIAL PRIMARY KEY,
    name TEXT,
    surname TEXT,
    patronymic TEXT,
    age INT,
    gender TEXT,
    nationality TEXT
);

CREATE INDEX idx_people_name_surname_patronymic ON people (name, surname, patronymic);
