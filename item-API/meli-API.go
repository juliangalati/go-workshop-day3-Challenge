package main

import (
	"encoding/json"
	"net/http"
	"time"
)

const CUSTOM_TIMEOUT = time.Duration(15 * time.Second)

const MELI_BASE_SITE = "https://api.mercadolibre.com/"
const MELI_ITEMS = "items/"
const MELI_SITES = "sites/"
const MELI_USERS = "users/"
const MELI_CATEGORIES = "categories/"

func generalServerError(err error) responseAPI {
	return responseAPI{Json: JSON_generic{"message": "Server error, retry later!", "cause": err.Error()}, StatusCode: http.StatusInternalServerError}
}

func getItem(itemId string) responseAPI {
	return getJson(MELI_BASE_SITE + MELI_ITEMS + itemId)
}

func getSite(siteId string) responseAPI {
	return getJson(MELI_BASE_SITE + MELI_SITES + siteId)
}

func getUser(userId string) responseAPI {
	return getJson(MELI_BASE_SITE + MELI_USERS + userId)
}

func getCategory(categoryId string) responseAPI {
	return getJson(MELI_BASE_SITE + MELI_CATEGORIES + categoryId)
}

func getJson(url string) responseAPI {
	httpClient := &http.Client{Timeout: CUSTOM_TIMEOUT}

	resp, err := httpClient.Get(url)

	if err != nil {
		// TODO: retry a few times
		return generalServerError(err)
	}
	defer resp.Body.Close()

	var result JSON_generic

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return generalServerError(err)
	}
	return responseAPI{Json: result, StatusCode: resp.StatusCode}
}
