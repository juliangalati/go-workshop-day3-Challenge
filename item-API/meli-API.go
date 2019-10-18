package main

import (
	"encoding/json"
	"net/http"
	"time"
)

const CUSTOM_TIMEOUT = time.Duration(15 * time.Second)

const MELI_SITE = "https://api.mercadolibre.com/"
const MELI_ITEMS = "items/"

func generalServerError(err error) responseAPI {
	return responseAPI{Json: &JSON_default{"message": "Server error, retry later!", "cause": err.Error()}, StatusCode: http.StatusInternalServerError}
}

func getItem(itemId string) responseAPI {
	return getJson(MELI_SITE + MELI_ITEMS + itemId)
}

func getJson(url string) responseAPI {
	httpClient := &http.Client{Timeout: CUSTOM_TIMEOUT}

	resp, err := httpClient.Get(url)

	if err != nil {
		return generalServerError(err)
	}
	defer resp.Body.Close()

	var result JSON_default

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return generalServerError(err)
	}
	return responseAPI{Json: &result, StatusCode: resp.StatusCode}
}
