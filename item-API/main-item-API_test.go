package main

import (
	"reflect"
	"testing"
)

func TestGetExtraInfoIndexesByParam(t *testing.T) {

	t.Run("Empty param", func(t *testing.T) {
		got := getExtraInfoIndexesByParam("")
		want := EVERY_INFORMATION_AVAILABLE

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("Unknown params", func(t *testing.T) {
		got := len(getExtraInfoIndexesByParam("buyer,publisher,owner"))
		want := 0

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("seller param", func(t *testing.T) {
		got := getExtraInfoIndexesByParam("seller")
		want := []extraInformationIndex{SELLER}

		// DeepEqual non nil functions: https://github.com/golang/go/issues/8554
		/*if !reflect.DeepEqual(got, want) {*/
		if want[0].key != got[0].key {
			t.Errorf("got %v want %v", got, want)
		}
	})
}
