package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func TestHomepageHandler(t *testing.T) {
	mockResponse := `{"message":"Welcome to Book Management"}`
	r := SetUpRouter()
	r.GET("/", HomepageHandler)
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	responseData, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, mockResponse, string(responseData))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestInitializeHandler(t *testing.T) {
	r := SetUpRouter()
	r.GET("/init", InitializeHandler)
	req, _ := http.NewRequest("GET", "/init", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Books initialized", w.Body.String())
}

func TestNewBookHandler(t *testing.T) {
	r := SetUpRouter()
	r.POST("/books", NewBookHandler)
	book := Book{
		Title:  "Demo Book",
		Author: "Demo Author",
		Year:   2023,
	}
	jsonValue, _ := json.Marshal(book)
	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGetBooksHandler(t *testing.T) {
	r := SetUpRouter()
	r.GET("/books", GetBooksHandler)
	req, _ := http.NewRequest("GET", "/books", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var books []Book
	json.Unmarshal(w.Body.Bytes(), &books)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, books)
}

func TestUpdateBookHandler(t *testing.T) {
	r := SetUpRouter()
	r.PUT("/books/:id", UpdateBookHandler)
	book := Book{
		Title:  "Updated Book",
		Author: "Updated Author",
		Year:   2023,
	}
	jsonValue, _ := json.Marshal(book)
	reqFound, _ := http.NewRequest("PUT", "/books/1", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, reqFound)
	assert.Equal(t, http.StatusOK, w.Code)

	reqNotFound, _ := http.NewRequest("PUT", "/books/12", bytes.NewBuffer(jsonValue))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, reqNotFound)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteBookHandler(t *testing.T) {
	r := SetUpRouter()
	r.DELETE("/books/:id", DeleteBookHandler)
	reqFound, _ := http.NewRequest("DELETE", "/books/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, reqFound)
	assert.Equal(t, http.StatusOK, w.Code)

	reqNotFound, _ := http.NewRequest("DELETE", "/books/12", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, reqNotFound)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
