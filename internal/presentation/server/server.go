package server

import (
	"fmt"
	"user/internal/interfaces"
	"user/internal/presentation/db"
	"user/internal/presentation/logger"

	"github.com/gin-gonic/gin"
)

var (
	DbService    *db.DB
	CacheService interfaces.CacheRepo
	UserService  interfaces.UserRepo
)

// Server определяет сервер с сервисами
type Server struct {
	srv *gin.Engine
}

// NewServer создает новый экземпляр Server
func NewServer() *Server {
	// gin.SetMode(gin.ReleaseMode)
	srv := gin.New()

	h := NewHandlers()

	srv.POST("/users", h.Create)
	srv.GET("/users", h.Get)
	srv.PUT("/users", h.Put)

	logger.Logger.Info("Server has been created")
	return &Server{
		srv: srv,
	}
}

func (s *Server) Start(port int) error {
	err := s.srv.Run(fmt.Sprintf(":%d", port))

	return err
}

func (s *Server) Shutdown() {
	err := DbService.CloseDB()
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("PostgreSQL connection hasn't been closed: %v", err))
	}

	err = CacheService.Close()
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("Redis connection hasn't been closed: %v", err))
	}
}
