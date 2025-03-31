package main

import (
	"fmt"
	"os"
	"strconv"
	"user/internal/presentation/db"
	"user/internal/presentation/logger"
	"user/internal/presentation/realization"
	"user/internal/presentation/server"

	"github.com/joho/godotenv"
)

func main() {
	err := logger.NewLogger()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = godotenv.Load("../../../.env")
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("Failed to load env fail - %v", err))
		return
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")

	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPass := os.Getenv("REDIS_PASSWORD")

	dataBase, err := db.CreateDB(host, port, user, password, name)
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("Database creating error - %v", err))
		return
	}
	server.DbService = dataBase

	err = server.DbService.CreateSchema("file:/../../../internal/presentation/migrations")
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("Schema creating error - %v", err))
		return
	}

	cacheRepo := realization.NewConnectRedis(redisHost, redisPort, redisPass)
	server.CacheService = cacheRepo

	serverPortStr := os.Getenv("SERVER_PORT")
	serverPort, err := strconv.Atoi(serverPortStr)
	if err != nil {
		logger.Logger.Error("Invalid port")
		return
	}

	userService := realization.NewUserService(dataBase, cacheRepo)
	server.UserService = userService

	srv := server.NewServer()
	err = srv.Start(serverPort)
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("Server working error - %v", err))
	}

	srv.Shutdown()
}
