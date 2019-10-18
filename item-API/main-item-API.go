package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type JSON_default map[string]interface{}

type responseAPI struct {
	Json       *JSON_default
	StatusCode int
}

func main() {
	r := gin.Default()

	// test ping pong
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// simple get item
	r.GET("/show/:item_id", func(c *gin.Context) {
		itemId := c.Param("item_id")
		response := getItem(itemId)
		c.JSON(response.StatusCode, response.Json)
	})

	r.Run()
}
