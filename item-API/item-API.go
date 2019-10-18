package main

import (
	"net/http"
)

const _ID = "_id"
const KEY_CATEGORY = "category"
const KEY_CATEGORY_ID = KEY_CATEGORY + _ID
const KEY_SELLER = "seller"
const KEY_SELLER_ID = KEY_SELLER + _ID
const KEY_SIE = "site"
const KEY_SITE_ID = KEY_SIE + _ID

func getItemFull(itemId string) responseAPI {
	basicItemResponse := getJson(MELI_BASE_SITE + MELI_ITEMS + itemId)

	if basicItemResponse.StatusCode == http.StatusOK {
		// si se encontro item, se completa
		return completeInformation(basicItemResponse.Json)
	} else {
		return basicItemResponse
	}
}

func completeInformation(jsonItem *JSON_generic) responseAPI {
	categoryId := getValueAsString(jsonItem, KEY_CATEGORY_ID)
	if categoryId != "" {
		categoryResponse := getCategory(categoryId)
		if categoryResponse.StatusCode == http.StatusOK {
			(*jsonItem)[KEY_CATEGORY] = categoryResponse.Json
		}
	}

	return responseAPI{Json: jsonItem, StatusCode: http.StatusOK}
}

func getValueAsString(json *JSON_generic, key string) string {
	value := (*json)[key]

	switch value.(type) {
	case string:
		return value.(string)
	default:
		return ""
	}
}
