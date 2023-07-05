package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Product struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Price int32  `json:"price"`
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
	db, err := DbConn()
	if err != nil {
		panic(err.Error())
	}
	db.Create(&newProduct)
	c.JSON(http.StatusCreated, newProduct)
}

func GetProductHandler(c *gin.Context) {
	db, err := DbConn()
	if err != nil {
		panic(err.Error())
	}
	var products []Product
	db.Find(&products)
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
	db, err := DbConn()
	if err != nil {
		panic(err.Error())
	}
	err = db.Where("id=?", id).First(&product).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Product not found",
		})
		return
	}
	db.Save(&product)
	c.JSON(http.StatusOK, product)
}

func DeleteProductHandler(c *gin.Context) {
	id := c.Param("id")
	db, err := DbConn()
	if err != nil {
		panic(err.Error())
	}
	db.Where("id=?", id).Delete(&Product{})
	c.JSON(http.StatusOK, gin.H{
		"message": "Product has been deleted",
	})
}

func InitHandler(c *gin.Context) {
	var products = []Product{
		{ID: "P001", Name: "Product A", Price: 100000},
		{ID: "P002", Name: "Product B", Price: 200000},
		{ID: "P003", Name: "Product C", Price: 300000},
	}

	db, err := DbConn()
	db.AutoMigrate(&Product{})
	if err != nil {
		panic(err.Error())
	}

	db.Create(&products)
	c.JSON(http.StatusOK, gin.H{
		"message": "Init product",
	})
}

func DbConn() (*gorm.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", host, user, password, dbName, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database")
	}
	return db, nil
}

func main() {

	router := gin.Default()
	router.GET("/", HomepageHandler)
	router.GET("/init", InitHandler)
	router.GET("/products", GetProductHandler)
	router.POST("/products", NewProductHandler)
	router.PUT("/products/:id", UpdateProductHandler)
	router.DELETE("/products/:id", DeleteProductHandler)
	err := router.Run(":8888")
	if err != nil {
		panic("Failed run server")
	}
}
