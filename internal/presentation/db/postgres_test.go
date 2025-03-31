package db

import (
	"testing"

	"user/internal/presentation/logger"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func init() {
	_ = logger.NewLogger()
}

// Тест для функции CreateDB
func TestCreateDB(t *testing.T) {
	db, err := CreateDB("localhost", "5432", "user", "password", "dbname")
	assert.NoError(t, err)
	assert.NotNil(t, db)
}

// Тест для функции CloseDB
func TestCloseDB_Success(t *testing.T) {
	mockDB, sqlMock, err := sqlmock.New()
	assert.NoError(t, err)
	defer func() {
		err := mockDB.Close()
		if err != nil {
			return
		}
	}()

	sqlMock.ExpectClose().WillReturnError(nil)

	db := &DB{
		Db: mockDB,
	}

	err = db.CloseDB()
	assert.NoError(t, err)

	err = sqlMock.ExpectationsWereMet()
	assert.NoError(t, err)
}
