package testutils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/papisz/weather"
)

// ForecastFromJSON reads JSON from file and returns a *Forecast struct
func ForecastFromJSON(filename string) *weather.Forecast {
	jsonFile, err := os.Open(path.Join("../../testdata/source", filename))
	if err != nil {
		log.Fatalf("unable to open file: %v", err)
	}

	jsonBytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatalf("unable to read file: %v", err)
	}

	forecast := &weather.Forecast{}
	if err := json.Unmarshal(jsonBytes, forecast); err != nil {
		log.Fatalf("unable to unmarshal bytes: %v", err)
	}
	return forecast
}

// JSONFileToBytes reads JSON from file and returns bytes
func JSONFileToBytes(dir, filename string) []byte {
	jsonFile, err := os.Open(path.Join(dir, filename))
	if err != nil {
		log.Fatalf("unable to open file: %v", err)
	}

	jsonBytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatalf("unable to read file: %v", err)
	}
	return jsonBytes
}
