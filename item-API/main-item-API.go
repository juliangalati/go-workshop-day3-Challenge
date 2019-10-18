package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type JSON_generic map[string]interface{}

type responseAPI struct {
	Json       *JSON_generic
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
		response := getItemFull(itemId)
		c.JSON(response.StatusCode, response.Json)
	})

	r.Run()
}
