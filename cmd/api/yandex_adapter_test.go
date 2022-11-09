package main

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestData struct {
	latitude     float64
	longitude    float64
	restarauntId string
}

var testData = TestData{
	latitude:     55.696233,
	longitude:    37.570431,
	restarauntId: "48274",
}

func TestDoGetRequest(t *testing.T) {
	// 1. Json response
	url := "https://jsonplaceholder.typicode.com/todos/1"
	res, err := doGetRequest(url)
	require.Nil(t, err, fmt.Sprintf("json request failed: %v", err))
	assert.Equal(t, "{\"completed\":false,\"id\":1,\"title\":\"delectus aut autem\",\"userId\":1}", res.String(), fmt.Sprintf("Response: %s", res.String()))
}

func TestGetRestaraunts(t *testing.T) {
	// 1. Request 1 restaraunt
	num := 1
	restaraunts, _, err := getRestaraunts(testData.latitude, testData.longitude, num, false, nil)
	log.Println("qwe", len(restaraunts))
	require.Nil(t, err, fmt.Sprintf("failed to get restaraunts: %v", err))
	assert.Equal(t, num, len(restaraunts))

	// 2. Request 3 restaraunts
	num = 3
	restaraunts, _, err = getRestaraunts(testData.latitude, testData.longitude, num, false, nil)
	require.Nil(t, err, fmt.Sprintf("failed to get restaraunts: %v", err))

	assert.Equal(t, num, len(restaraunts))
}

func TestGetRestarauntMenu(t *testing.T) {
	// 1. Get restaraunt menu as json
	menu, err := getRestarauntMenu(testData.restarauntId)
	require.Nil(t, err, fmt.Sprintf("failed to get restaraunt menu: %v", err))
	assert.NotEqualf(t, "{}", menu.String(), "empty menu received")
}

func TestGetRandomFood(t *testing.T) {
	// 1. Get 1 food card
	foodResponse, err := GetRandomFood(1, testData.latitude, testData.longitude, false, nil)
	require.Nil(t, err, fmt.Sprintf("failed to get food cards: %v", err))
	require.True(t, foodResponse.Succeed)
	foodCards := foodResponse.FoodCards
	require.Equal(t, 1, len(foodCards), "1 food card expected")

	// 2. Get 3 food cards
	foodResponse, err = GetRandomFood(3, testData.latitude, testData.longitude, false, nil)
	require.Nil(t, err, fmt.Sprintf("failed to get food cards: %v", err))
	require.True(t, foodResponse.Succeed)
	foodCards = foodResponse.FoodCards
	require.Equal(t, 3, len(foodCards), "3 food card expected")
}
