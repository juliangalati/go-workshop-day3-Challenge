package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

// cual es la key y id de la info extra en el json original, con key se debe agregar la info y como se obtiene
type extraInformationIndex struct {
	keyId                  string
	id                     string
	key                    string
	getterExtraInformation func(id string) responseAPI
}

//la respuesta de la API al tratar de obtener la info extra
type extraInformationAPIResponse struct {
	index    extraInformationIndex
	response responseAPI
}

const _ID = "_id"
const KEY_CATEGORY = "category"
const KEY_CATEGORY_ID = KEY_CATEGORY + _ID
const KEY_SELLER = "seller"
const KEY_SELLER_ID = KEY_SELLER + _ID
const KEY_SITE = "site"
const KEY_SITE_ID = KEY_SITE + _ID

func getItemFull(itemId string) responseAPI {
	basicItemResponse := getItem(itemId)

	if basicItemResponse.StatusCode == http.StatusOK {
		// si se encontro el item, se completa
		response, err := completeInformation(basicItemResponse.Json)
		if err != nil {
			return generalServerError(err)
		} else {
			return response
		}
	} else {
		return basicItemResponse
	}
}

func completeInformation(jsonItem *JSON_generic) (responseAPI, error) {
	// solo se intenta completar la info extra cuyo id estaba en el item (no se considera un error que falte alguno)
	extraInformationIndexes := getExtraInformationIndexes(jsonItem)
	cantExtraInfo := len(extraInformationIndexes)
	responseChan := make(chan extraInformationAPIResponse, cantExtraInfo)

	for _, extraInfo := range extraInformationIndexes {
		go func(extraInfoIndex extraInformationIndex) {
			responseChan <- extraInformationAPIResponse{index: extraInfoIndex, response: extraInfoIndex.getterExtraInformation(extraInfoIndex.id)}
		}(extraInfo)
	}

	for i := 0; i < cantExtraInfo; i++ {
		infoResponse := <-responseChan
		if infoResponse.response.StatusCode == http.StatusOK {
			replaceOldInformationWithNewInJson(jsonItem, infoResponse.index.keyId, infoResponse.index.key, infoResponse.response.Json)
		} else {
			// si un id estaba en el item y no se puede completar, falla el request entero (se considera un error)
			return responseAPI{}, errors.New(fmt.Sprintf("Getting the %v failed with status: %v", infoResponse.index.key, infoResponse.response.StatusCode))
		}
	}
	return responseAPI{Json: jsonItem, StatusCode: http.StatusOK}, nil
}

func replaceOldInformationWithNewInJson(jsonItem *JSON_generic, oldKey string, newKey string, aValue *JSON_generic) {
	delete(*jsonItem, oldKey)
	(*jsonItem)[newKey] = *aValue
}

func getExtraInformationIndexes(jsonItem *JSON_generic) (extraInformationIndexes []extraInformationIndex) {
	addExtraInformationIndex(KEY_CATEGORY_ID, KEY_CATEGORY, getCategory, jsonItem, &extraInformationIndexes)
	addExtraInformationIndex(KEY_SELLER_ID, KEY_SELLER, getUser, jsonItem, &extraInformationIndexes)
	addExtraInformationIndex(KEY_SITE_ID, KEY_SITE, getSite, jsonItem, &extraInformationIndexes)
	return
}

func addExtraInformationIndex(keyId string, key string, getterInfo func(string) responseAPI, jsonItem *JSON_generic, extraInformationIds *[]extraInformationIndex) {
	extraInfoId := getValueAsString(jsonItem, keyId)
	if extraInfoId != "" {
		*extraInformationIds = append(*extraInformationIds, extraInformationIndex{keyId: keyId, id: extraInfoId, key: key, getterExtraInformation: getterInfo})
	}
}

func getValueAsString(json *JSON_generic, key string) string {
	value := (*json)[key]

	switch value.(type) {
	case string:
		return value.(string)
	case float64:
		return strconv.FormatInt(int64(value.(float64)), 10)
	default:
		return ""
	}
}
