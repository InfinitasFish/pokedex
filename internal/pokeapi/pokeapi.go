package pokeapi

import (
	"fmt"
	"net/http"
	"io"
	"time"
	"encoding/json"
	"crypto/tls"
	"math/rand"
	"internal/pokecache"
)

type Config struct {
	Next *string `json:"next"`
	Previous *string `json:"previous"`
	Results []map[string]string `json:"results"`
}

// it's not guaranteed that all the remaining fields following the problematic one 
// will be unmarshaled into the target object
// so creating new struct solely for pokemons
type encounters struct {
	Encounters []map[string]interface{} `json:"pokemon_encounters"`
}

type Pokedex struct {
	CaughtPokemons map[string]Pokemon
}

// just a name for now
type Pokemon struct {
	Name string
}

type pokemonData struct {
	Experience int `json:"base_experience"`
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
	// cache.PrintEntriesTime()

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

func GetPokemonsByLocation(locationName string, cache *pokecache.Cache) ([]string, error) {
	url := "https://pokeapi.co/api/v2/location-area/" + locationName
	encountersData := encounters{}
	// default capacity 32
	pokemons := make([]string, 0, 32)

	if data, exists := cache.Get(locationName); exists {
		err := json.Unmarshal(data, &encountersData)
		if err != nil {
			return nil, err
		}

		for _, val := range encountersData.Encounters {
			pokemons = append(pokemons, val["pokemon"].(map[string]interface{})["name"].(string))
		}

		return pokemons, nil
	} else {

		client := &http.Client{
			Transport: &http.Transport{
				TLSNextProto: make(map[string]func(string, *tls.Conn) http.RoundTripper),
			},
			Timeout: 10 * time.Second,
		}

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}

		res, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		data, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		cache.Add(url, data)
		
		encountersData := encounters{}
		err = json.Unmarshal(data, &encountersData)
		if err != nil {
			return nil, err
		}
		
		for _, val := range encountersData.Encounters {
			pokemons = append(pokemons, val["pokemon"].(map[string]interface{})["name"].(string))
		}

		return pokemons, nil
	}
}

// higher base experience -> harder to catch
// suppose(!) that 500 is max base exp, anything higher -> no catch
func TryCatchPokemon(pokemonName string, cache *pokecache.Cache, pokedex *Pokedex) (bool, error) {
	var maxExp float32 = 500
	pokeData := pokemonData{}
	var url string = "https://pokeapi.co/api/v2/pokemon/" + pokemonName

	if data, exists := cache.Get(url); exists {
		err := json.Unmarshal(data, &pokeData)
		if err != nil {
			return false, err
		}
	} else {
		client := &http.Client{
			Transport: &http.Transport{
				TLSNextProto: make(map[string]func(string, *tls.Conn) http.RoundTripper),
			},
			Timeout: 10 * time.Second,
		}

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return false, err
		}

		res, err := client.Do(req)
		if err != nil {
			return false, err
		}
		defer res.Body.Close()

		data, err := io.ReadAll(res.Body)
		if err != nil {
			return false, err
		}
		cache.Add(url, data)

		err = json.Unmarshal(data, &pokeData)
		if err != nil {
			return false, err
		}
	}

	// calculating chance to catch and using random float32 in [0,1) range
	var catched bool = false
	chanceToCatch := 1 - float32(pokeData.Experience) / maxExp
	if rand.Float32() > chanceToCatch {
		catched = true
	}

	// adding pokemon to shared *Pokedex
	if catched {
		pokemon := Pokemon{Name: pokemonName}
		// don't care about duplicates
		pokedex.CaughtPokemons[pokemonName] = pokemon
	}
		
	return catched, nil
}
