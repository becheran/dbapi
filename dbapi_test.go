package dbapi_test

import (
	"bufio"
	"log"
	"os"
	"testing"
	"time"

	"github.com/becheran/dbapi"
	"github.com/stretchr/testify/assert"
)

const (
	evaNoBerlinAlexanderplatz = 8089001
)

func parseAPIKey() string {
	file, err := os.Open("./apikey.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	text := scanner.Text()
	if text == "" {
		log.Fatal("Expected first line to contain the API key")
	}
	return text
}

var apiKey = parseAPIKey()

func TestGetStationInfo(t *testing.T) {
	api := dbapi.API{Bearer: apiKey}

	res, err := api.StationInfo("Köln")

	assert.Nil(t, err)
	assert.Len(t, res, 1)
}

func TestGetStationInvalidBearer(t *testing.T) {
	api := dbapi.API{Bearer: "Invalid"}

	res, err := api.StationInfo("Köln")

	assert.NotNil(t, err)
	assert.Empty(t, res)
}

func TestPlan(t *testing.T) {
	api := dbapi.API{Bearer: apiKey}

	res, err := api.Plan(evaNoBerlinAlexanderplatz, time.Now())

	assert.Nil(t, err)
	assert.Empty(t, res)
}
