package pokeapi

import (
	"fmt"
	"net/http"
	"io"
	"time"
	"encoding/json"
	"crypto/tls"
	"internal/pokecache"
)

type Config struct {
	Next *string `json:"next"`
	Previous *string `json:"previous"`
	Results []map[string]string `json:"results"`
}

func GetLocationsData(isnext bool, config *Config, cache *pokecache.Cache) error {
	// use query params to get next or prev 20 locs
	// https://pokeapi.co/api/v2/location-area/?offset=20&limit=20
	var offset int = 0
	var limit int = 20
	var url string = "https://pokeapi.co/api/v2/location-area/"

	// custom HTTP/1.1 client, because HTTPS/2 hangs randomly SOMEWHY
	client := &http.Client{
		Transport: &http.Transport{
			TLSNextProto: make(map[string]func(string, *tls.Conn) http.RoundTripper),
		},
		Timeout: 10 * time.Second,
	}

	// first call
	if config.Next == nil && config.Previous == nil {
		url += fmt.Sprintf("?offset=%d&limit=%d", offset, limit)
	// next locations
	} else if isnext {
		if config.Next == nil {
			return fmt.Errorf("Error, end of the locations")
		}
		url = *config.Next
	// previous locations
	} else {
		if config.Previous == nil {
			return fmt.Errorf("Error, start of the locations")
		}
		url = *config.Previous
	}

	// debug caching
	cache.PrintEntriesTime()

	// checking for url in cache
	if val, exists := cache.Get(url); exists {
		// matching []byte into shared Config struct
		err := json.Unmarshal(val, config)
		if err != nil {
			return err
		}
		// successfull early return
		return nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	// adding new entry to cache
	cache.Add(url, data)

	// matching []byte into shared Config struct
	// hopefully config wont break on error
	err = json.Unmarshal(data, config)
	if err != nil {
		return err
	}

	return nil
}
