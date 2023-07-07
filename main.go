package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Book struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   int    `json:"year"`
}

var books = []Book{
	{ID: "B001", Title: "Book A", Author: "Author A", Year: 2020},
	{ID: "B002", Title: "Book B", Author: "Author B", Year: 2021},
	{ID: "B003", Title: "Book C", Author: "Author C", Year: 2022},
}

func HomepageHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to Book Management"})
}

func NewBookHandler(c *gin.Context) {
	var newBook Book
	if err := c.ShouldBindJSON(&newBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	newBook.ID = uuid.New().String()
	books = append(books, newBook)
	c.JSON(http.StatusCreated, newBook)
}

func GetBooksHandler(c *gin.Context) {
	c.JSON(http.StatusOK, books)
}

func UpdateBookHandler(c *gin.Context) {
	id := c.Param("id")
	var book Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	index := -1
	for i := 0; i < len(books); i++ {
		if books[i].ID == id {
			index = i
			break
		}
	}
	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Book not found",
		})
		return
	}
	books[index] = book
	c.JSON(http.StatusOK, book)
}

func DeleteBookHandler(c *gin.Context) {
	id := c.Param("id")
	index := -1
	for i := 0; i < len(books); i++ {
		if books[i].ID == id {
			index = i
			break
		}
	}
	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Book not found",
		})
		return
	}
	books = append(books[:index], books[index+1:]...)
	c.JSON(http.StatusOK, gin.H{
		"message": "Book has been deleted",
	})
}

func main() {
	router := gin.Default()
	router.GET("/", HomepageHandler)
	router.GET("/books", GetBooksHandler)
	router.POST("/books", NewBookHandler)
	router.PUT("/books/:id", UpdateBookHandler)
	router.DELETE("/books/:id", DeleteBookHandler)
	router.Run(":8888")
}
