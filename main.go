package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Book struct {
	ID     uint32 `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   int    `json:"year"`
}

var books = []Book{}

func HomepageHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to Book Management"})
}

func InitializeHandler(c *gin.Context) {
	books = []Book{
		{ID: 1, Title: "Book A", Author: "Author A", Year: 2020},
		{ID: 2, Title: "Book B", Author: "Author B", Year: 2021},
		{ID: 3, Title: "Book C", Author: "Author C", Year: 2022},
	}
	c.JSON(http.StatusOK, gin.H{"message": "Books initialized"})
}

func NewBookHandler(c *gin.Context) {
	var newBook Book
	if err := c.ShouldBindJSON(&newBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	newBook.ID = uuid.New().ID()
	books = append(books, newBook)
	c.JSON(http.StatusCreated, newBook)
}

func GetBooksHandler(c *gin.Context) {
	c.JSON(http.StatusOK, books)
}

func UpdateBookHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID format",
		})
		return
	}

	var book Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	updated := false
	for i := 0; i < len(books); i++ {
		if books[i].ID == uint32(id) {
			book.ID = uint32(id)
			books[i] = book
			updated = true
			break
		}
	}

	if updated {
		c.JSON(http.StatusOK, book)
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Book not found",
		})
	}
}

func DeleteBookHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID format",
		})
		return
	}

	deleted := false
	for i := 0; i < len(books); i++ {
		if books[i].ID == uint32(id) {
			books = append(books[:i], books[i+1:]...)
			deleted = true
			break
		}
	}

	if deleted {
		c.JSON(http.StatusOK, gin.H{
			"message": "Book has been deleted",
		})
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Book not found",
		})
	}
}

func main() {
	router := gin.Default()
	router.GET("/", HomepageHandler)
	router.GET("/init", InitializeHandler)
	router.GET("/books", GetBooksHandler)
	router.POST("/books", NewBookHandler)
	router.PUT("/books/:id", UpdateBookHandler)
	router.DELETE("/books/:id", DeleteBookHandler)
	router.Run(":8888")
}
