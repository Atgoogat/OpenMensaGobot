package openmensa

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var api OpenmensaApi

func TestMain(m *testing.M) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/canteens" && r.Header.Get("Accept") == "application/json" {
			file, err := os.ReadFile("../test/data/canteens_150623.json")
			if err != nil {
				panic(err)
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(file)
			w.WriteHeader(http.StatusOK)
		} else if r.URL.Path == "/canteens/1/days" && r.Header.Get("Accept") == "application/json" {
			file, err := os.ReadFile("../test/data/canteen_days_150623.json")
			if err != nil {
				panic(err)
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(file)
			w.WriteHeader(http.StatusOK)
		} else if r.URL.Path == "/canteens/1/days/2023-09-16/meals" && r.Header.Get("Accept") == "application/json" {
			file, err := os.ReadFile("../test/data/canteen_meals_140623.json")
			if err != nil {
				panic(err)
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(file)
			w.WriteHeader(http.StatusOK)
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()
	api = NewOpenmensaApi(server.URL)

	m.Run()
}

func TestListCanteens(t *testing.T) {
	canteens, err := api.ListCanteens(1, 100)
	assert.Nil(t, err)
	assert.NotNil(t, canteens)

	magdeburg := canteens[0]
	assert.Equal(t, 1, magdeburg.Id)
}

func TestListCanteenDays(t *testing.T) {
	canteenDays, err := api.ListDaysOfCanteen(1, 1, 100)
	assert.Nil(t, err)
	assert.NotNil(t, canteenDays)

	assert.False(t, canteenDays[0].Closed)
}

func TestListCanteenDays_InvalidId(t *testing.T) {
	_, err := api.ListDaysOfCanteen(-1, 1, 100)
	assert.NotNil(t, err)

	assert.Equal(t, ErrExpectedStatusCodeOk, err)
}

func TestListMeals(t *testing.T) {
	date, _ := time.Parse("2006-01-02", "2023-09-16")
	meals, err := api.ListMealsForADay(1, date, 0, 100)
	assert.Nil(t, err)
	assert.Equal(t, 33, len(meals))
}
