package controllers

import (
	"log"
	"net/http"

	"github.com/burhangltekin/byfood/models"
	"github.com/burhangltekin/byfood/utils"
	"github.com/gin-gonic/gin"
)

// GetBooks godoc
// @Summary      List all books
// @Description  Get all books from the database
// @Tags         books
// @Produce      json
// @Success      200  {array}  models.Book
// @Failure      500  {object}  map[string]string
// @Router       /books [get]
func GetBooks(c *gin.Context) {
	var books []models.Book
	if err := utils.DB.Find(&books).Error; err != nil {
		log.Printf("Error fetching books: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch books"})
		return
	}
	c.JSON(http.StatusOK, books)
}

// GetBook godoc
// @Summary      Get a book by ID
// @Description  Get details of a book by its ID
// @Tags         books
// @Produce      json
// @Param        id   path      int  true  "Book ID"
// @Success      200  {object}  models.Book
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /books/{id} [get]
func GetBook(c *gin.Context) {
	id := c.Param("id")
	var book models.Book
	if err := utils.DB.First(&book, id).Error; err != nil {
		log.Printf("Book not found (id=%s): %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}
	c.JSON(http.StatusOK, book)
}

// CreateBook godoc
// @Summary      Create a new book
// @Description  Add a new book to the database
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        book  body      models.BookInput  true  "Book to create"
// @Success      201   {object}  models.Book
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /books [post]
func CreateBook(c *gin.Context) {
	var input models.BookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Invalid input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	book := models.Book{
		Title:  input.Title,
		Author: input.Author,
		Year:   input.Year,
	}
	if err := utils.DB.Create(&book).Error; err != nil {
		log.Printf("Error creating book: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create book"})
		return
	}
	c.JSON(http.StatusCreated, book)
}

// UpdateBook godoc
// @Summary      Update a book
// @Description  Update an existing book by ID
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        id    path      int                true  "Book ID"
// @Param        book  body      models.BookInput   true  "Book data"
// @Success      200   {object}  models.Book
// @Failure      400   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /books/{id} [put]
func UpdateBook(c *gin.Context) {
	id := c.Param("id")
	var book models.Book
	if err := utils.DB.First(&book, id).Error; err != nil {
		log.Printf("Book not found for update (id=%s): %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}
	var input models.BookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Invalid input for update: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	book.Title = input.Title
	book.Author = input.Author
	book.Year = input.Year
	if err := utils.DB.Save(&book).Error; err != nil {
		log.Printf("Error updating book (id=%s): %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update book"})
		return
	}
	c.JSON(http.StatusOK, book)
}

// DeleteBook godoc
// @Summary      Delete a book
// @Description  Delete a book by ID
// @Tags         books
// @Produce      json
// @Param        id   path      int  true  "Book ID"
// @Success      200  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /books/{id} [delete]
func DeleteBook(c *gin.Context) {
	id := c.Param("id")
	result := utils.DB.Delete(&models.Book{}, id)
	if result.Error != nil {
		log.Printf("Error deleting book (id=%s): %v", id, result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete book"})
		return
	}
	if result.RowsAffected == 0 {
		log.Printf("No book found to delete (id=%s)", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Book deleted"})
}
