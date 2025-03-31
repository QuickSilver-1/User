package db

import (
	"database/sql"
	"fmt"
	"user/internal/presentation/logger"
)

// DB - структура для работы с базой данных
type DB struct {
	Db *sql.DB
}

// CreateDB создает подключение к базе данных и возвращает экземпляр DB
func CreateDB(ip, port, user, pass, nameDB string) (*DB, error) {
	logger.Logger.Debug("Database connection creating...")
	sqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", ip, port, user, pass, nameDB)
	conn, err := sql.Open("postgres", sqlInfo)

	if err != nil {
		return nil, fmt.Errorf("database connection error: %v", err)
	}

	logger.Logger.Info("Database connection has been created")
	return &DB{
		Db: conn,
	}, nil
}

// CloseDB закрывает подключение к базе данных
func (db *DB) CloseDB() error {
	logger.Logger.Debug("Closing database connection")
	return db.Db.Close()
}
