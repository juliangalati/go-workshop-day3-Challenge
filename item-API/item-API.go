package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

// cual es la key de la info extra en el json original, con key se debe agregar la info y como se obtiene
type extraInformationIndex struct {
	keyId                  string
	key                    string
	getterExtraInformation func(id string) responseAPI
}

// cual es el id de la info extra y cual es su index
type extraInformationId struct {
	index extraInformationIndex
	id    string
}

//la respuesta de la API al tratar de obtener la info extra
type extraInformationAPIResponse struct {
	infoId   extraInformationId
	response responseAPI
}

const _ID = "_id"
const KEY_CATEGORY = "category"
const KEY_CATEGORY_ID = KEY_CATEGORY + _ID
const KEY_SELLER = "seller"
const KEY_SELLER_ID = KEY_SELLER + _ID
const KEY_SITE = "site"
const KEY_SITE_ID = KEY_SITE + _ID

var SELLER = extraInformationIndex{
	keyId:                  KEY_SELLER_ID,
	key:                    KEY_SELLER,
	getterExtraInformation: getUser,
}

var CATEGORY = extraInformationIndex{
	keyId:                  KEY_CATEGORY_ID,
	key:                    KEY_CATEGORY,
	getterExtraInformation: getCategory,
}

var SITE = extraInformationIndex{
	keyId:                  KEY_SITE_ID,
	key:                    KEY_SITE,
	getterExtraInformation: getSite,
}

func getItemWithExtraInfo(itemId string, extras []extraInformationIndex) responseAPI {
	basicItemResponse := getItem(itemId)

	// si se encontro el item, se completa
	if basicItemResponse.StatusCode == http.StatusOK {
		response, err := completeInformation(basicItemResponse.Json, extras)
		if err != nil {
			return generalServerError(err)
		} else {
			return response
		}
	} else {
		return basicItemResponse
	}
}

func completeInformation(jsonItem *JSON_generic, extras []extraInformationIndex) (responseAPI, error) {
	// solo se intenta completar la info extra cuyo id estaba en el item (no se considera un error que falte alguno)
	extraInformationIds := getExtraInformationIds(jsonItem, extras)
	cantExtraInfo := len(extraInformationIds)
	responseChan := make(chan extraInformationAPIResponse, cantExtraInfo)

	for _, extraInfo := range extraInformationIds {
		go func(extraInfo extraInformationId) {
			responseChan <- extraInformationAPIResponse{infoId: extraInfo, response: extraInfo.index.getterExtraInformation(extraInfo.id)}
		}(extraInfo)
	}

	for i := 0; i < cantExtraInfo; i++ {
		infoResponse := <-responseChan
		if infoResponse.response.StatusCode == http.StatusOK {
			replaceOldInformationWithNewInJson(jsonItem, infoResponse.infoId.index.keyId, infoResponse.infoId.index.key, infoResponse.response.Json)
		} else {
			// si un id estaba en el item y no se puede completar, falla el request entero (se considera un error)
			return responseAPI{}, errors.New(fmt.Sprintf("Getting the %v failed with status: %v", infoResponse.infoId.index.key, infoResponse.response.StatusCode))
		}
	}
	return responseAPI{Json: jsonItem, StatusCode: http.StatusOK}, nil
}

func replaceOldInformationWithNewInJson(jsonItem *JSON_generic, oldKey string, newKey string, aValue *JSON_generic) {
	delete(*jsonItem, oldKey)
	(*jsonItem)[newKey] = *aValue
}

func getExtraInformationIds(jsonItem *JSON_generic, extras []extraInformationIndex) (extraInformationIds []extraInformationId) {
	for _, infoIndex := range extras {
		addExtraInformationId(infoIndex, jsonItem, &extraInformationIds)
	}
	return
}

func addExtraInformationId(infoIndex extraInformationIndex, jsonItem *JSON_generic, extraInformationIds *[]extraInformationId) {
	extraInfoId := getValueAsString(jsonItem, infoIndex.keyId)
	if extraInfoId != "" {
		*extraInformationIds = append(*extraInformationIds, extraInformationId{index: infoIndex, id: extraInfoId})
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
