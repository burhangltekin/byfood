package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/burhangltekin/byfood/models"
	"github.com/burhangltekin/byfood/utils"
)

func mockHandler(status int, body string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(status, body)
	}
}

func SetupRoutesMock(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.GET("/books", mockHandler(200, "get books"))
		api.GET("/books/:id", mockHandler(200, "get book"))
		api.POST("/books", mockHandler(201, "created"))
		api.PUT("/books/:id", mockHandler(200, "updated"))
		api.DELETE("/books/:id", mockHandler(200, "deleted"))
	}
}

func TestSetupRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to in-memory db: %v", err)
	}
	if err := db.AutoMigrate(&models.Book{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	utils.DB = db
	testBook := models.Book{Title: "Route Book", Author: "Route Author", Year: 2024}
	db.Create(&testBook)

	tests := []struct {
		name       string
		method     string
		url        string
		body       string
		expectCode int
		checkBody  func(t *testing.T, body string)
	}{
		{
			name:       "GET /api/books",
			method:     http.MethodGet,
			url:        "/api/books",
			expectCode: 200,
			checkBody: func(t *testing.T, body string) {
				assert.Contains(t, body, "Route Book")
			},
		},
		{
			name:       "GET /api/books/:id",
			method:     http.MethodGet,
			url:        "/api/books/1",
			expectCode: 200,
			checkBody: func(t *testing.T, body string) {
				assert.Contains(t, body, "Route Book")
			},
		},
		{
			name:       "POST /api/books",
			method:     http.MethodPost,
			url:        "/api/books",
			body:       `{"title":"T","author":"A","year":2024}`,
			expectCode: 201,
			checkBody: func(t *testing.T, body string) {
				assert.Contains(t, body, "T")
			},
		},
		{
			name:       "PUT /api/books/:id",
			method:     http.MethodPut,
			url:        "/api/books/1",
			body:       `{"title":"Updated","author":"A","year":2024}`,
			expectCode: 200,
			checkBody: func(t *testing.T, body string) {
				assert.Contains(t, body, "Updated")
			},
		},
		{
			name:       "DELETE /api/books/:id",
			method:     http.MethodDelete,
			url:        "/api/books/1",
			expectCode: 200,
			checkBody: func(t *testing.T, body string) {
				// Accept empty or message
			},
		},
	}

	r := gin.New()
	SetupRoutes(r)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			var req *http.Request
			if tt.body != "" {
				req, _ = http.NewRequest(tt.method, tt.url, strings.NewReader(tt.body))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req, _ = http.NewRequest(tt.method, tt.url, nil)
			}
			r.ServeHTTP(w, req)
			assert.Equal(t, tt.expectCode, w.Code)
			if tt.checkBody != nil {
				tt.checkBody(t, w.Body.String())
			}
		})
	}
}
