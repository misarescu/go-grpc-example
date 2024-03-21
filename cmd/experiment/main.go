package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

var citiesMap = map[string]Coordinates{
	"Cluj-Napoca": {
		Latitude:  46.7667,
		Longitude: 23.6300,
	},
	"Vulcan": {
		Latitude:  45.3833,
		Longitude: 23.2667,
	},
	"Carei": {
		Latitude:  47.6833,
		Longitude: 22.4667,
	},
	"Bucharest": {
		Latitude:  44.4323,
		Longitude: 26.1063,
	},
}

type WeatherResponse struct {
	Coordinates
	UtcOffset int     `json:"utc_offset_seconds"`
	TimeZone  string  `json:"timezone"`
	Elevation float64 `json:"elevation"`
	Hourly    struct {
		Time         []string  `json:"time"`
		Temperatures []float64 `json:"temperature_2m"`
	} `json:"hourly,omitempty"`
	Current struct {
		Time        string  `json:"time"`
		Interval    int     `json:"interval"`
		Temperature float64 `json:"temperature_2m"`
		Rain        float64 `json:"rain"`
		Wind        float64 `json:"wind_speed_10m"`
	} `json:"current,omitempty"`
}

func ErrorAndExit(e error) {
	log.Printf("%s\nexiting.\n", e)
	os.Exit(1)
}

func FetchWeather(city string, res chan<- string, wg *sync.WaitGroup) (WeatherResponse, error) {
	defer wg.Done()
	apiUrl := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%.2f&longitude=%.2f&current=temperature_2m,rain,wind_speed_10m&timezone=auto&forecast_days=1", citiesMap[city].Latitude, citiesMap[city].Longitude)
	// log.Printf("fetch url: %s", apiUrl)
	data := WeatherResponse{}
	response, err := http.Get(apiUrl)
	if err != nil {
		log.Println("error fetching api")
		return WeatherResponse{}, err
	}
	defer response.Body.Close()

	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		log.Println("error decoding json")
		return WeatherResponse{}, err
	}

	res <- fmt.Sprintf("%+v", data)

	return data, nil
}

func main() {
	var wg sync.WaitGroup

	res := make(chan string)

	for city := range citiesMap {
		wg.Add(1)
		go FetchWeather(city, res, &wg)
	}
	go func() {
		wg.Wait()
		close(res)
	}()

	for r := range res {
		log.Printf("%s", r)
	}
	time.Sleep(1 * time.Second)

}
