package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"user/internal/domain"
	"user/internal/presentation/db"
	"user/internal/presentation/logger"
	"user/internal/presentation/realization"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func SetEnv() {
	gin.SetMode(gin.TestMode)

	err := logger.NewLogger()
	if err != nil {
		panic(err)
	}

	err = godotenv.Load("../../../.env")
	if err != nil {
		logger.Logger.Fatal(fmt.Sprintf("Failed to load env fail - %v", err))
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
		logger.Logger.Fatal(fmt.Sprintf("Database creating error - %v", err))
		return
	}
	DbService = dataBase

	err = DbService.CreateSchema("file://../migrations")
	if err != nil {
		logger.Logger.Fatal(fmt.Sprintf("Schema creating error - %v", err))
		return
	}

	cacheRepo := realization.NewConnectRedis(redisHost, redisPort, redisPass)
	CacheService = cacheRepo

	userService := realization.NewUserService(dataBase, cacheRepo)
	UserService = userService
}

func TestCreateHandler(t *testing.T) {
	SetEnv()

	tests := []struct {
        name         string
        input        domain.User
        expectedCode int
        expectedBody string
    }{
        {
            name: "Valid user",
            input: domain.User{
                FirstName: "John",
                LastName:  "Doe",
                Login:     "john.doe@example.com",
                Password:  "StrongPassword123!",
            },
            expectedCode: http.StatusOK,
            expectedBody: `{"id":`,
        },
        {
            name: "Missing email",
            input: domain.User{
                FirstName: "Jane",
                LastName:  "Smith",
                Login:     "",
                Password:  "StrongPassword123!",
            },
            expectedCode: http.StatusBadRequest,
            expectedBody: `{"error":"Email is required"}`,
        },
        {
            name: "Invalid email format",
            input: domain.User{
                FirstName: "Invalid",
                LastName:  "Email",
                Login:     "invalid-email",
                Password:  "StrongPassword123!",
            },
            expectedCode: http.StatusBadRequest,
            expectedBody: `{"error":"Invalid email"}`,
        },
        {
            name: "Weak password",
            input: domain.User{
                FirstName: "Weak",
                LastName:  "Password",
                Login:     "weak.password@example.com",
                Password:  "123",
            },
            expectedCode: http.StatusBadRequest,
            expectedBody: `{"error":"Invalid password"}`,
        },
        {
            name: "Duplicate email",
            input: domain.User{
                FirstName: "Duplicate",
                LastName:  "Email",
                Login:     "john.doe@example.com",
                Password:  "StrongPassword123!",
            },
            expectedCode: http.StatusBadRequest,
            expectedBody: `{"error":"User with email john.doe@example.com already exist"}`,
        },
    }

    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            router := gin.Default()
            h := NewHandlers()
            router.POST("/create", h.Create)

            body, _ := json.Marshal(test.input)
            req, _ := http.NewRequest(http.MethodPost, "/create", bytes.NewBuffer(body))
            req.Header.Set("Content-Type", "application/json")

            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)

            if w.Code != test.expectedCode {
                t.Errorf("expected %d, got %d", test.expectedCode, w.Code)
            }

            if test.expectedBody != "" && !bytes.Contains(w.Body.Bytes(), []byte(test.expectedBody)) {
                t.Errorf("expected body to contain %s, got %s", test.expectedBody, w.Body.String())
            }
        })
    }
}


func TestGetHandler(t *testing.T) {
    SetEnv()

    tests := []struct {
        name         string
        queryParam   string
        mockUser     *domain.User
        expectedCode int
        expectedBody string
    }{
        {
            name:         "Missing ID",
            queryParam:   "",
            mockUser:     nil,
            expectedCode: http.StatusBadRequest,
            expectedBody: `{"error":"Id is required"}`,
        },
        {
            name:         "Invalid ID format",
            queryParam:   "abc",
            mockUser:     nil,
            expectedCode: http.StatusBadRequest,
            expectedBody: `{"error":"Invalid id"}`,
        },
        {
            name:         "User not found",
            queryParam:   "999",
            mockUser:     nil,
            expectedCode: http.StatusBadRequest,
            expectedBody: `{"error":"User with id 999 not exist"}`,
        },
	}

    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            router := gin.Default()
            h := NewHandlers()
            router.GET("/get", h.Get)

            req, _ := http.NewRequest(http.MethodGet, "/get?id="+test.queryParam, nil)

            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)

            if w.Code != test.expectedCode {
                t.Errorf("expected %d, got %d", test.expectedCode, w.Code)
            }

            if test.expectedBody != "" && !bytes.Contains(w.Body.Bytes(), []byte(test.expectedBody)) {
                t.Errorf("expected body to contain %s, got %s", test.expectedBody, w.Body.String())
            }
        })
    }
}


func TestPutHandler(t *testing.T) {
    SetEnv()

    tests := []struct {
        name         string
        queryParam   string
        input        domain.User
        mockError    error
        expectedCode int
        expectedBody string
    }{
        {
            name:       "Valid user update",
            queryParam: "1",
            input: domain.User{
                FirstName: "Updated",
                LastName:  "User",
                Login:     "updated.user@example.com",
                Password:  "UpdatedStrongPassword123!",
            },
            mockError:    nil,
            expectedCode: http.StatusOK,
            expectedBody: "",
        },
        {
            name:         "Missing ID",
            queryParam:   "",
            input:        domain.User{},
            mockError:    nil,
            expectedCode: http.StatusBadRequest,
            expectedBody: `{"error":"Id is required"}`,
        },
        {
            name:         "Invalid ID format",
            queryParam:   "abc",
            input:        domain.User{},
            mockError:    nil,
            expectedCode: http.StatusBadRequest,
            expectedBody: `{"error":"Invalid id"}`,
        },
    }

    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            router := gin.Default()
            h := NewHandlers()
            router.PUT("/put", h.Put)

            body, _ := json.Marshal(test.input)
            req, _ := http.NewRequest(http.MethodPut, "/put?id="+test.queryParam, bytes.NewBuffer(body))
            req.Header.Set("Content-Type", "application/json")

            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)

            if w.Code != test.expectedCode {
                t.Errorf("expected %d, got %d", test.expectedCode, w.Code)
            }

            if test.expectedBody != "" && !bytes.Contains(w.Body.Bytes(), []byte(test.expectedBody)) {
                t.Errorf("expected body to contain %s, got %s", test.expectedBody, w.Body.String())
            }
        })
    }
}