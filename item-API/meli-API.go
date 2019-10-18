package main

import (
	"encoding/json"
	"net/http"
	"time"
)

const CUSTOM_TIMEOUT = time.Duration(15 * time.Second)

const MELI_SITE = "https://api.mercadolibre.com/"
const MELI_ITEMS = "items/"

func getItem(itemId string) responseAPI {
	return getJson(MELI_SITE + MELI_ITEMS + itemId)
}

func getJson(url string) responseAPI {
	httpClient := &http.Client{Timeout: CUSTOM_TIMEOUT}

	resp, err := httpClient.Get(url)

	if err != nil {
		//do something
	}
	defer resp.Body.Close()

	var result JSON_default

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		//do something
	}

	return responseAPI{Json: &result, StatusCode: resp.StatusCode}
}
