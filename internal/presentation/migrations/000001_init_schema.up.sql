-- Создание таблицы users
CREATE TABLE users (
    id              SERIAL PRIMARY KEY,      -- Идентификатор пользователя
    first_name      VARCHAR(255),            -- Имя пользователя
    last_name       VARCHAR(255),            -- Фамилия пользователя
    birthday        DATE,                    -- Дата рождения
    login           VARCHAR(255) UNIQUE,    -- Имя пользователя (уникальное)
    password        VARCHAR(255)            -- Пароль пользователя в виде хэша
);
