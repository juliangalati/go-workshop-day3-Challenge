package main

import (
	"net/http"
	"reflect"
	"testing"
)

func TestGetValueAsString(t *testing.T) {
	t.Run("empty JSON", func(t *testing.T) {
		aJSON := make(JSON_generic)
		got := getValueAsString(aJSON, KEY_CATEGORY_ID)
		want := ""

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("key not present", func(t *testing.T) {
		aJSON := make(JSON_generic)
		aJSON[KEY_SITE_ID] = "MLA"
		got := getValueAsString(aJSON, KEY_CATEGORY_ID)
		want := ""

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("key present, type string", func(t *testing.T) {
		aJSON := make(JSON_generic)
		aCategory := "MLA12250"
		aJSON[KEY_CATEGORY_ID] = aCategory
		got := getValueAsString(aJSON, KEY_CATEGORY_ID)
		want := aCategory

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("key present, type float64", func(t *testing.T) {
		aJSON := make(JSON_generic)
		aCategory := "MLA12250"
		aJSON[KEY_CATEGORY_ID] = aCategory
		aJSON[KEY_SELLER_ID] = float64(136879867)

		got := getValueAsString(aJSON, KEY_SELLER_ID)
		want := "136879867"

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})
}

func TestCompleteInformation(t *testing.T) {
	var categoryMock = extraInformationIndex{
		keyId: KEY_CATEGORY_ID,
		key:   KEY_CATEGORY,
		getterExtraInformation: func(id string) responseAPI {
			result := make(JSON_generic)
			result["id"] = id
			return responseAPI{Json: result, StatusCode: http.StatusOK}
		},
	}

	t.Run("complete category", func(t *testing.T) {

		aJSON := make(JSON_generic)
		aCategory := "MLA12250"
		aJSON[KEY_CATEGORY_ID] = aCategory

		wantedJSON := make(JSON_generic)
		wantedJSON[KEY_CATEGORY] = JSON_generic{"id": aCategory}

		got, _ := completeInformation(aJSON, []extraInformationIndex{categoryMock})
		want := responseAPI{
			Json:       wantedJSON,
			StatusCode: http.StatusOK,
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("complete category and site", func(t *testing.T) {
		var siteMockOk = extraInformationIndex{
			keyId: KEY_SITE_ID,
			key:   KEY_SITE,
			getterExtraInformation: func(id string) responseAPI {
				result := make(JSON_generic)
				result["id"] = id
				return responseAPI{Json: result, StatusCode: http.StatusOK}
			},
		}

		aJSON := make(JSON_generic)
		aCategory := "MLA12250"
		aSite := "MLA"
		aJSON[KEY_CATEGORY_ID] = aCategory
		aJSON[KEY_SITE_ID] = aSite

		wantedJSON := make(JSON_generic)
		wantedJSON[KEY_CATEGORY] = JSON_generic{"id": aCategory}
		wantedJSON[KEY_SITE] = JSON_generic{"id": aSite}

		got, _ := completeInformation(aJSON, []extraInformationIndex{categoryMock, siteMockOk})
		want := responseAPI{
			Json:       wantedJSON,
			StatusCode: http.StatusOK,
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("complete category and site, but site fails so everything fails", func(t *testing.T) {
		var siteMockError = extraInformationIndex{
			keyId: KEY_SITE_ID,
			key:   KEY_SITE,
			getterExtraInformation: func(id string) responseAPI {
				return responseAPI{Json: nil, StatusCode: http.StatusNotFound}
			},
		}

		aJSON := make(JSON_generic)
		aCategory := "MLA12250"
		aSite := "MLA"
		aJSON[KEY_CATEGORY_ID] = aCategory
		aJSON[KEY_SITE_ID] = aSite

		got, err := completeInformation(aJSON, []extraInformationIndex{categoryMock, siteMockError})
		want := responseAPI{}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}

		if err == nil {
			t.Errorf("error wanted")
		}
	})
}
