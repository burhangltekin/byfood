package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/burhangltekin/byfood/models"
	"github.com/burhangltekin/byfood/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestGetBooks(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		prepare      func()
		expectStatus int
		expectError  string
	}{
		{
			name: "basic get books",
			prepare: func() {
				// Setup valid DB and insert test book
				db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
				if err != nil {
					t.Fatalf("failed to connect to in-memory db: %v", err)
				}
				if err := db.AutoMigrate(&models.Book{}); err != nil {
					t.Fatalf("failed to migrate: %v", err)
				}
				utils.DB = db
				testBook := models.Book{Title: "Test Book", Author: "Test Author", Year: 2024}
				db.Create(&testBook)
			},
			expectStatus: http.StatusOK,
			expectError:  "",
		},
		{
			name: "db error",
			prepare: func() {
				db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
				if err != nil {
					t.Fatalf("failed to open db: %v", err)
				}
				sqlDB, _ := db.DB()
				sqlDB.Close()
				utils.DB = db
			},
			expectStatus: http.StatusInternalServerError,
			expectError:  "Failed to fetch books",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare()
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest(http.MethodGet, "/api/books", nil)
			GetBooks(c)
			if w.Code != tt.expectStatus {
				t.Errorf("expected status %d, got %d", tt.expectStatus, w.Code)
			}
			if tt.expectStatus == http.StatusOK {
				var books []models.Book
				if err := json.Unmarshal(w.Body.Bytes(), &books); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				if len(books) != 1 || books[0].Title != "Test Book" {
					t.Errorf("expected 1 book with title 'Test Book', got %+v", books)
				}
			} else if tt.expectError != "" {
				var resp map[string]string
				if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
					t.Fatalf("failed to unmarshal error response: %v", err)
				}
				if resp["error"] != tt.expectError {
					t.Errorf("expected error message %q, got %+v", tt.expectError, resp)
				}
			}
		})
	}
}

func TestGetBook(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to in-memory db: %v", err)
	}
	if err := db.AutoMigrate(&models.Book{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	utils.DB = db

	testBook := models.Book{Title: "Test Book", Author: "Test Author", Year: 2024}
	db.Create(&testBook)

	tests := []struct {
		name         string
		id           string
		expectStatus int
		expectTitle  string
		expectError  string
	}{
		{
			name:         "get existing book",
			id:           "1",
			expectStatus: http.StatusOK,
			expectTitle:  testBook.Title,
			expectError:  "",
		},
		{
			name:         "get non-existing book",
			id:           "999",
			expectStatus: http.StatusNotFound,
			expectTitle:  "",
			expectError:  "Book not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = gin.Params{{Key: "id", Value: tt.id}}
			c.Request, _ = http.NewRequest(http.MethodGet, "/api/books/"+tt.id, nil)
			GetBook(c)
			if w.Code != tt.expectStatus {
				t.Errorf("expected status %d, got %d", tt.expectStatus, w.Code)
			}
			if tt.expectStatus == http.StatusOK {
				var book models.Book
				if err := json.Unmarshal(w.Body.Bytes(), &book); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				if book.Title != tt.expectTitle {
					t.Errorf("expected title %q, got %q", tt.expectTitle, book.Title)
				}
			} else {
				var resp map[string]string
				if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
					t.Fatalf("failed to unmarshal error response: %v", err)
				}
				if resp["error"] != tt.expectError {
					t.Errorf("expected error message %q, got %+v", tt.expectError, resp)
				}
			}
		})
	}
}

func TestCreateBook(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to in-memory db: %v", err)
	}
	if err := db.AutoMigrate(&models.Book{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	utils.DB = db

	tests := []struct {
		name         string
		body         string
		prepare      func()
		expectStatus int
		expectTitle  string
		expectError  string
	}{
		{
			name:         "create valid book",
			body:         `{"title":"New Book","author":"Author","year":2024}`,
			prepare:      func() {},
			expectStatus: http.StatusCreated,
			expectTitle:  "New Book",
			expectError:  "",
		},
		{
			name:         "missing required field",
			body:         `{"author":"Author","year":2024}`,
			prepare:      func() {},
			expectStatus: http.StatusBadRequest,
			expectTitle:  "",
			expectError:  "bad request",
		},
		{
			name:         "invalid year",
			body:         `{"title":"Book","author":"Author","year":2200}`,
			prepare:      func() {},
			expectStatus: http.StatusBadRequest,
			expectTitle:  "",
			expectError:  "bad request",
		},
		{
			name: "db create error",
			body: `{"title":"Book","author":"Author","year":2024}`,
			prepare: func() {
				db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
				if err != nil {
					t.Fatalf("failed to open db: %v", err)
				}
				sqlDB, _ := db.DB()
				sqlDB.Close()
				utils.DB = db
			},
			expectStatus: http.StatusInternalServerError,
			expectTitle:  "",
			expectError:  "Failed to create book",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
			if err != nil {
				t.Fatalf("failed to connect to in-memory db: %v", err)
			}
			if err := db.AutoMigrate(&models.Book{}); err != nil {
				t.Fatalf("failed to migrate: %v", err)
			}
			utils.DB = db
			if tt.prepare != nil {
				tt.prepare()
			}
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest(http.MethodPost, "/api/books", strings.NewReader(tt.body))
			c.Request.Header.Set("Content-Type", "application/json")
			CreateBook(c)
			if w.Code != tt.expectStatus {
				t.Errorf("expected status %d, got %d", tt.expectStatus, w.Code)
			}
			if tt.expectStatus == http.StatusCreated {
				var book models.Book
				if err := json.Unmarshal(w.Body.Bytes(), &book); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				if book.Title != tt.expectTitle {
					t.Errorf("expected title %q, got %q", tt.expectTitle, book.Title)
				}
			} else {
				var resp map[string]interface{}
				if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
					t.Fatalf("failed to unmarshal error response: %v", err)
				}
				if tt.expectError == "bad request" {
					if resp["error"] == nil {
						t.Errorf("expected error message, got %+v", resp)
					}
				} else if tt.expectError != "" {
					if resp["error"] != tt.expectError {
						t.Errorf("expected error message %q, got %+v", tt.expectError, resp)
					}
				}
			}
		})
	}
}

func TestUpdateBook(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testBook := models.Book{Title: "Old Title", Author: "Old Author", Year: 2000}

	tests := []struct {
		name         string
		prepare      func() (*gin.Context, *httptest.ResponseRecorder)
		requestBody  string
		expectStatus int
		expectTitle  string
		expectError  string
	}{
		{
			name: "update existing book",
			prepare: func() (*gin.Context, *httptest.ResponseRecorder) {
				db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
				err := db.AutoMigrate(&models.Book{})
				if err != nil {
					return nil, nil
				}
				utils.DB = db
				db.Create(&testBook)
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Params = gin.Params{{Key: "id", Value: "1"}}
				return c, w
			},
			requestBody:  `{"title":"New Title","author":"New Author","year":2024}`,
			expectStatus: http.StatusOK,
			expectTitle:  "New Title",
			expectError:  "",
		},
		{
			name: "update non-existing book",
			prepare: func() (*gin.Context, *httptest.ResponseRecorder) {
				db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
				err := db.AutoMigrate(&models.Book{})
				if err != nil {
					return nil, nil
				}
				utils.DB = db
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Params = gin.Params{{Key: "id", Value: "999"}}
				return c, w
			},
			requestBody:  `{"title":"New Title","author":"New Author","year":2024}`,
			expectStatus: http.StatusNotFound,
			expectTitle:  "",
			expectError:  "Book not found",
		},
		{
			name: "invalid request body",
			prepare: func() (*gin.Context, *httptest.ResponseRecorder) {
				db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
				err := db.AutoMigrate(&models.Book{})
				if err != nil {
					return nil, nil
				}
				utils.DB = db
				db.Create(&testBook)
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Params = gin.Params{{Key: "id", Value: "1"}}
				return c, w
			},
			requestBody:  `{"title":123,"author":"New Author","year":2024}`,
			expectStatus: http.StatusBadRequest,
			expectTitle:  "",
			expectError:  "bad request",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := tt.prepare()
			c.Request, _ = http.NewRequest(http.MethodPut, "/api/books/"+c.Params.ByName("id"), strings.NewReader(tt.requestBody))
			c.Request.Header.Set("Content-Type", "application/json")

			UpdateBook(c)
			assert.Equal(t, tt.expectStatus, w.Code)
			if tt.expectStatus == http.StatusOK {
				var book models.Book
				err := json.Unmarshal(w.Body.Bytes(), &book)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectTitle, book.Title)
			} else {
				var resp map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				if tt.expectError == "bad request" {
					assert.NotNil(t, resp["error"])
				} else if tt.expectError != "" {
					assert.Equal(t, tt.expectError, resp["error"])
				} else {
					assert.Nil(t, resp["error"])
				}
			}
		})
	}
}

func TestDeleteBook(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testBook := models.Book{Title: "To Delete", Author: "Author", Year: 2020}

	tests := []struct {
		name         string
		prepare      func() (*gin.Context, *httptest.ResponseRecorder)
		expectStatus int
		expectError  string
	}{
		{
			name: "delete existing book",
			prepare: func() (*gin.Context, *httptest.ResponseRecorder) {
				db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
				db.AutoMigrate(&models.Book{})
				utils.DB = db
				db.Create(&testBook)
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Params = gin.Params{{Key: "id", Value: "1"}}
				c.Request, _ = http.NewRequest(http.MethodDelete, "/api/books/1", nil)
				return c, w
			},
			expectStatus: http.StatusOK,
			expectError:  "",
		},
		{
			name: "delete non-existing book",
			prepare: func() (*gin.Context, *httptest.ResponseRecorder) {
				db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
				db.AutoMigrate(&models.Book{})
				utils.DB = db
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Params = gin.Params{{Key: "id", Value: "999"}}
				c.Request, _ = http.NewRequest(http.MethodDelete, "/api/books/999", nil)
				return c, w
			},
			expectStatus: http.StatusNotFound,
			expectError:  "Book not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := tt.prepare()
			DeleteBook(c)
			assert.Equal(t, tt.expectStatus, w.Code)
			if tt.expectStatus == http.StatusOK {
				var resp map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				if len(resp) > 0 {
					assert.Contains(t, resp, "message")
				}
			} else {
				var resp map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectError, resp["error"])
			}
		})
	}
}
