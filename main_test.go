package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func TestHomepageHandler(t *testing.T) {
	mockResponse := `{"message":"Welcome to Golang Simple Jenkins"}`
	r := SetUpRouter()
	r.GET("/", HomepageHandler)
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	responseData, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, mockResponse, string(responseData))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestNewProductHandler(t *testing.T) {
	r := SetUpRouter()
	r.POST("/products", NewProductHandler)
	productID := uuid.New().String()
	product := Product{
		ID:    productID,
		Name:  "Demo Product",
		Price: 100000,
	}
	jsonValue, _ := json.Marshal(product)
	req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGetProductHandler(t *testing.T) {
	r := SetUpRouter()
	r.GET("/products", GetProductHandler)
	req, _ := http.NewRequest("GET", "/products", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var products []Product
	json.Unmarshal(w.Body.Bytes(), &products)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, products)
}

func TestUpdateProductHandler(t *testing.T) {
	r := SetUpRouter()
	r.PUT("/products/:id", UpdateProductHandler)
	product := Product{
		ID:    `P001`,
		Name:  "Demo Product",
		Price: 200000,
	}
	jsonValue, _ := json.Marshal(product)
	reqFound, _ := http.NewRequest("PUT", "/products/"+product.ID, bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, reqFound)
	assert.Equal(t, http.StatusOK, w.Code)

	reqNotFound, _ := http.NewRequest("PUT", "/products/12", bytes.NewBuffer(jsonValue))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, reqNotFound)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
