package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"user/internal/domain"
	"user/internal/presentation/logger"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

const (
	NOT_UNIQUE_LOGIN = "23505"
)

type Handlers struct{}

func NewHandlers() *Handlers {
	return &Handlers{}
}

func (Handlers) Create(ctx *gin.Context) {
	user := validBody(ctx)
	if user == nil {
		return
	}

	id, err := UserService.Create(*user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return
	}

	if id == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("User with email %s already exist", user.Login)})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"id": id})
}

func (Handlers) Get(ctx *gin.Context) {
	idStr := ctx.Request.URL.Query().Get("id")
	if idStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Id is required"})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	user, err := UserService.Get(domain.Id(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return
	}

	if user == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("User with id %d not exist", id)})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (Handlers) Put(ctx *gin.Context) {
	idStr := ctx.Request.URL.Query().Get("id")
	if idStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Id is required"})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	user := validBody(ctx)
	if user == nil {
		return
	}

	user.Id = domain.Id(id)
	err = UserService.Update(*user)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == NOT_UNIQUE_LOGIN {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("User with email %s already exist", user.Login)})
			return
		}

		ctx.JSON(http.StatusInternalServerError, "Internal error")
		return
	}

	ctx.Status(http.StatusOK)
}

func validBody(ctx *gin.Context) *domain.User {
	var user domain.User
	err := json.NewDecoder(ctx.Request.Body).Decode(&user)
	defer func() {
		err := ctx.Request.Body.Close()
		if err != nil {
			logger.Logger.Error("Close body error")
		}
	}()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid body"})
		return nil
	}

	if user.Login == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
		return nil
	}

	valid := IsValidEmail(user.Login)
	if !valid {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email"})
		return nil
	}

	hashPass, err := ValidPass(user.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid password"})
		return nil
	}

	user.Password = hashPass
	return &user
}
