package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Product struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Price int32  `json:"price"`
}

var products = []Product{
	{ID: "P001", Name: "Product A", Price: 100000},
	{ID: "P002", Name: "Product B", Price: 200000},
	{ID: "P003", Name: "Product C", Price: 300000},
}

func HomepageHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to Golang Simple Jenkins"})
}

func NewProductHandler(c *gin.Context) {
	var newProduct Product
	if err := c.ShouldBindJSON(&newProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	newProduct.ID = uuid.New().String()
	products = append(products, newProduct)
	c.JSON(http.StatusCreated, newProduct)
}

func GetProductHandler(c *gin.Context) {
	c.JSON(http.StatusOK, products)
}

func UpdateProductHandler(c *gin.Context) {
	id := c.Param("id")
	var product Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	index := -1
	for i := 0; i < len(products); i++ {
		if products[i].ID == id {
			index = 1
		}
	}
	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Product not found",
		})
		return
	}
	products[index] = product
	c.JSON(http.StatusOK, product)
}

func DeleteProductHandler(c *gin.Context) {
	id := c.Param("id")
	index := -1
	for i := 0; i < len(products); i++ {
		if products[i].ID == id {
			index = 1
		}
	}
	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Product not found",
		})
		return
	}
	products = append(products[:index], products[index+1:]...)
	c.JSON(http.StatusOK, gin.H{
		"message": "Producr has been deleted",
	})
}

func main() {
	router := gin.Default()
	router.GET("/", HomepageHandler)
	router.GET("/products", GetProductHandler)
	router.POST("/products", NewProductHandler)
	router.PUT("/products/:id", UpdateProductHandler)
	router.DELETE("/products/:id", DeleteProductHandler)
	router.Run(":8888")
}
