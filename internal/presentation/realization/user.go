package realization

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
	"user/internal/domain"
	"user/internal/interfaces"
	"user/internal/presentation/db"
	"user/internal/presentation/logger"

	"github.com/lib/pq"
)

const (
	NOT_UNIQUE_LOGIN = "23505"
)

type UserService struct {
	db    *db.DB
	cache interfaces.CacheRepo
}

func NewUserService(db *db.DB, cache interfaces.CacheRepo) *UserService {
	return &UserService{
		db:    db,
		cache: cache,
	}
}

func (s *UserService) Create(user domain.User) (*domain.Id, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var id domain.Id
	logger.Logger.Debug("Creating user...")
	err := s.db.Db.QueryRowContext(ctx, `INSERT INTO users (first_name, last_name, birthday, login, password) VALUES ($1, $2, $3, $4, $5) RETURNING id`, user.FirstName, user.LastName, user.BirthDay, user.Login, user.Password).Scan(&id)
	var pqErr *pq.Error
	if err != nil {
		if errors.As(err, &pqErr) && pqErr.Code == NOT_UNIQUE_LOGIN {
			return nil, nil
		}
		logger.Logger.Error(fmt.Sprintf("Creating user error: %v", err))
		return nil, fmt.Errorf("creating postgres user error: %v", err)
	}

	logger.Logger.Debug("The user has been created successful")
	return &id, nil
}

func (s *UserService) Get(id domain.Id) (*domain.User, error) {
	cacheUser, err := s.cache.GetByKey(id)
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("Getting Redis key error: %v", err))
		return nil, err
	}

	if cacheUser != nil {
		return cacheUser, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	logger.Logger.Debug("Getting user...")
	var user domain.User
	err = s.db.Db.QueryRowContext(ctx, `SELECT id, first_name, last_name, birthday, login FROM users WHERE id = $1`, id).Scan(&user.Id, &user.FirstName, &user.LastName, &user.BirthDay, &user.Login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		logger.Logger.Error(fmt.Sprintf("Getting user error: %v", err))
		return nil, fmt.Errorf("getting postgres user error: %v", err)
	}

	user.Password = "***"
	err = s.cache.CreateKey(id, user)
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("Creating Redis key error: %v", err))
	}

	logger.Logger.Debug("The user has been get successful")
	return &user, nil
}

func (s *UserService) Update(user domain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	logger.Logger.Debug("Updating user...")
	_, err := s.db.Db.ExecContext(ctx, `UPDATE users SET first_name = $2, last_name = $3, birthday = $4, login = $5, password = $6 WHERE id = $1`, user.Id, user.FirstName, user.LastName, user.BirthDay, user.Login, user.Password)
	var pqErr *pq.Error
	if err != nil {
		if errors.As(err, &pqErr) && pqErr.Code == NOT_UNIQUE_LOGIN {
			return err
		}
		logger.Logger.Error(fmt.Sprintf("Updating user error: %v", err))
		return fmt.Errorf("updating postgres user error: %v", err)
	}

	err = s.cache.DelKey(user.Id)
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("Deleating Redis key error: %v", err))
	}

	return nil
}
