package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type JSON_generic map[string]interface{}

type responseAPI struct {
	Json       *JSON_generic
	StatusCode int
}

var EVERY_INFORMATION_AVAILABLE = []extraInformationIndex{CATEGORY, SELLER, SITE}
var PARAM_INFO_INDEX = map[string]extraInformationIndex{"seller": SELLER, "category": CATEGORY, "site": SITE}

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
		extraInfo := getExtraInfoIndexesByParam(c.Query("attributes"))
		response := GetItemWithExtraInfo(itemId, extraInfo)
		c.JSON(response.StatusCode, response.Json)
	})

	r.Run()
}

func getExtraInfoIndexesByParam(attributes string) (wantedInfo []extraInformationIndex) {
	infoAttributes := strings.Split(attributes, ",")
	if len(attributes) == 0 || len(infoAttributes) == 0 {
		wantedInfo = EVERY_INFORMATION_AVAILABLE
	} else {
		for _, attr := range infoAttributes {
			infoIndex, prs := PARAM_INFO_INDEX[attr]
			if prs {
				wantedInfo = append(wantedInfo, infoIndex)
			}
		}
	}
	return
}
