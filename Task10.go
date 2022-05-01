package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)


type book struct {
	ID     string  json:"id"
	Title  string  json:"title"
	Author string  json:"author" 
	Year  string json:"year" 
	Price float32 json:"price"
  }

var books = []books{
	{ID: "1", Title: "Harry Potter", Author: "Rowling", Year: "2002"Price: 2000},
	{ID: "2", Title: " Autobiographies", Author: "ERNEST HEMINGWAY", Price: 6000.99},
	{ID: "3", Title: " The Lord of the Rings", Author: "J.R.R. Tolkien", Price: 10000.99},
	{ID: "4", Title: " Spider", Author: "Marvel", Price: 90000.99},
	{ID: "5", Title: " Iron Man", Author: "Marvel", Price: 80000.99},
	
}

func getBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)
}

func postAlbums(c *gin.Context) {
	var newBook book

	if err := c.BindJSON(&newBook); err != nil {
		return
	}

	books = append(books, newBook)
	c.IndentedJSON(http.StatusCreated, newBook)
}

func getBookByID(c *gin.Context) {
	id := c.Param("id")

	for _, boo := range books {
		if boo.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "book not found"})
}

func deleteAlbumByID(c *gin.Context) {
	id := c.Param("id")

	for i, boo := range books {
		if boo.ID == id {
			books = append(books[:i], books[i+1:]...)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "book not found"})
}

func main() {
	router := gin.Default()
	router.GET("/books", getBooks)
	router.POST("/postBooks", postBooks)
	router.GET("/books/:id", getBookByID)
	router.DELETE("/deleteBooks/:id", deleteBookByID)

	router.Run("localhost:8080")
}
