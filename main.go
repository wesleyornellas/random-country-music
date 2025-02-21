package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type Country struct {
	Name struct {
		Common string `json:"common`
	} `json:"name"`
}

type LastFRMResponse struct {
	Tracks struct {
		Track []struct {
			Name   string `json:"name"`
			Artist struct {
				Name string `json:"name"`
			} `json:"artist"`
		} `json:"track"`
	} `json:"tracks"`
}

// Response struct for the API output
type Response struct {
	Country string `json:"country"`
	Song    string `json:"song"`
	Artist  string `json:"artist"`
}

// FetchRandomCountry retrieves a random country from the RestCountries API
func FetchRandomCountry() (string, error) {
	resp, err := http.Get("https://restcountries.com/v3.1/all")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var countries []Country
	if err := json.NewDecoder(resp.Body).Decode(&countries); err != nil {
		return "", err
	}

	// Select a random country
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(countries))

	return countries[randomIndex].Name.Common, nil
}

// FetchTopSong retrieves the top song from Last.fm API
func FetchTopSong(country string) (string, string, error) {
	fmt.Printf("Trying to get top song for %s", country)

	apiKey := "c725f1f4d058ffa793984e8d42db6b3b"
	url := fmt.Sprintf("http://ws.audioscrobbler.com/2.0/?method=geo.gettoptracks&country=%s&api_key=%s&format=json", country, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	var lastFMData LastFRMResponse
	if err := json.NewDecoder(resp.Body).Decode(&lastFMData); err != nil {
		return "", "", err
	}

	if len(lastFMData.Tracks.Track) == 0 {
		return "", "", fmt.Errorf("no songs found for %s", country)
	}

	topSong := lastFMData.Tracks.Track[0]
	return topSong.Name, topSong.Artist.Name, nil

}

// Handler function to return JSON response
func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	country, err := FetchRandomCountry()
	if err != nil {
		http.Error(w, `{"error": "Failed to fetch country"}`, http.StatusInternalServerError)
		return
	}

	song, artist, err := FetchTopSong(country)
	if err != nil {
		http.Error(w, `{"error": "Failed to fetch song"}`, http.StatusInternalServerError)
		return
	}

	response := Response{
		Country: country,
		Song:    song,
		Artist:  artist,
	}

	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/random-country-music", handler)
	fmt.Println("Server is running on port 8080...")
	http.ListenAndServe(":8080", nil)
}
