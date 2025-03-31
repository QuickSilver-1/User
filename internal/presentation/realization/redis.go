package realization

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"
	"user/internal/domain"
	"user/internal/presentation/logger"

	"github.com/go-redis/redis/v8"
)

// RedisRepo представляет репозиторий для работы с Redis
type RedisRepo struct {
	db *redis.Client
}

// NewConnectRedis создает новое подключение к Redis
// addr - адрес Redis сервера
// pass - пароль для подключения к Redis
func NewConnectRedis(addr, port, pass string) *RedisRepo {
	r := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", addr, port),
		Password: pass,
		DB:       0,
	})

	logger.Logger.Info("Redis connection create")
	return &RedisRepo{
		db: r,
	}
}

// CreateKey создает новый ключ в Redis
// id - идентификатор ключа
// text - текст песни
func (r *RedisRepo) CreateKey(id domain.Id, user domain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	userJson, err := json.Marshal(user)
	if err != nil {
		return errors.New("object marshalling error")
	}

	res := r.db.Set(ctx, strconv.FormatUint(id, 10), userJson, 0)
	_, err = res.Result()

	if err != nil {
		return fmt.Errorf("creating Redis key error: %v", err)
	}

	logger.Logger.Debug(fmt.Sprintf("key: %d was created", id))
	return nil
}

// GetText получает значение ключа из Redis по идентификатору
// id - идентификатор ключа
func (r *RedisRepo) GetByKey(id domain.Id) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	res := r.db.Get(ctx, strconv.FormatUint(id, 10))

	err := res.Err()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("getting redis key error: %v", err)
	}

	var user domain.User
	err = json.Unmarshal([]byte(res.Val()), &user)
	if err != nil {
		return nil, errors.New("object unmurshalling error")
	}
	logger.Logger.Debug(fmt.Sprintf("key: %d was got", id))

	return &user, nil
}

// DelKey удаляет ключ из Redis
// id - идентификатор ключа
func (r *RedisRepo) DelKey(id domain.Id) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err := r.db.Del(ctx, strconv.FormatUint(id, 10)).Err()

	if err != nil {
		return fmt.Errorf("deleating redis key error: %v", err)
	}

	logger.Logger.Debug(fmt.Sprintf("key: %d was deleted", id))
	return nil
}

func (r *RedisRepo) Close() error {
	logger.Logger.Info("Redis connection was closed")
	return r.db.Close()
}
