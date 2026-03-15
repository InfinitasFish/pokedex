package main

import (
	"fmt"
	"bufio"
	"os"
	"internal/pokeapi"
	"internal/pokecache"
)

// variadic callback needed...
type cliCommand struct {
	name string
	description string
	callback func(*pokeapi.Config, *pokecache.Cache, *pokeapi.Pokedex, string, string) error
}

func commandExit(config *pokeapi.Config, cache *pokecache.Cache, pokedex *pokeapi.Pokedex, locationName string, pokemonName string) error {
	fmt.Printf("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}

func commandHelp(config *pokeapi.Config, cache *pokecache.Cache, pokedex *pokeapi.Pokedex, locationName string, pokemonName string) error {
	fmt.Printf("Welcome to the Pokedex!\n")
	fmt.Printf("Usage:\n\n")
	fmt.Printf("help: Displays a help message\n")
	// xd
	fmt.Printf("exit: Exit the Pokedex\n")
	return nil
}

func commandMap(config *pokeapi.Config, cache *pokecache.Cache, pokedex *pokeapi.Pokedex, locationName string, pokemonName string) error {
	// get next locations
	pokeapi.GetLocationsData(true, config, cache)

	// print locations
	for _, m := range config.Results {
		fmt.Printf("%v\n", m["name"])
	}
	return nil
}

func commandMapBack(config *pokeapi.Config, cache *pokecache.Cache, pokedex *pokeapi.Pokedex, locationName string, pokemonName string) error {
	// get previous locations
	pokeapi.GetLocationsData(false, config, cache)

	// print locations
	for _, m := range config.Results {
		fmt.Printf("%v\n", m["name"])
	}
	return nil
}

func commandExplore(config *pokeapi.Config, cache *pokecache.Cache, pokedex *pokeapi.Pokedex, locationName string, pokemonName string) error {
	fmt.Printf("Exploring %v...\n", locationName)
	fmt.Printf("Found Pokemon:\n")

	pokemons, err := pokeapi.GetPokemonsByLocation(locationName, cache)
	if err != nil {
		return err
	}

	for _, pokemon := range pokemons {
		fmt.Printf(" - %v\n", pokemon)
	}

	return nil
}

func commandCatch(config *pokeapi.Config, cache *pokecache.Cache, pokedex *pokeapi.Pokedex, locationName string, pokemonName string) error {
	fmt.Printf("Throwing a Pokeball at %v...\n", pokemonName)
	
	catched, err := pokeapi.TryCatchPokemon(pokemonName, cache, pokedex)
	if err != nil {
		return err
	}

	if catched {
		fmt.Printf("%v was caught!\n", pokemonName)
	} else {
		fmt.Printf("%v escaped!\n", pokemonName)
	}

	return nil
}

func commandInspect(config *pokeapi.Config, cache *pokecache.Cache, pokedex *pokeapi.Pokedex, locationName string, pokemonName string) error {
	err := pokeapi.PrintPokemonData(pokemonName, pokedex)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	// possible commands
	commandRegister := map[string]cliCommand {
		"exit": cliCommand{
			name: "exit", 
			description: "Exit the Pokedex",
			callback: commandExit,
		},
		"help": cliCommand{
			name: "help", 
			description: "Pokedex usage help",
			callback: commandHelp,
		},
		"map": cliCommand{
			name: "map",
			description: "List Next 20 locations",
			callback: commandMap,
		},
		"mapb": cliCommand{
			name: "mapb",
			description: "List Previous 20 locations",
			callback: commandMapBack,
		},
		"explore": cliCommand{
			name: "explore",
			description: "List pokemons from area",
			callback: commandExplore,
		},
		"catch": cliCommand{
			name: "catch",
			description: "Try to catch a Pokemon",
			callback: commandCatch,
		},
		"inspect": cliCommand{
			name: "inspect",
			description: "Print catched Pokemon data",
			callback: commandInspect,
		},
	}

	// tracking locations offset
	locationsConfig := &pokeapi.Config{}

	// caching locations to not repeat requests
	locationsCache := pokecache.NewCache(15)

	// pokedex keeps catched pokemons
	playerPokedex := &pokeapi.Pokedex{CaughtPokemons: make(map[string]pokeapi.Pokemon, 16),}
	
	// default location and pokemon
	// kinda ugly because other commands than "explore"/"catch" don't need this
	locationName := ""
	pokemonName := ""

	// listening read eval print loop
	scanner := bufio.NewScanner(os.Stdin)
	for ;; {
		fmt.Printf("Pokedex > ")
		if scanner.Scan() {
			user_tokens := cleanInput(scanner.Text())
			if len(user_tokens) == 0 {
				fmt.Printf("Empty command\n")
				continue
			}
			command := user_tokens[0]
			if c, err := commandRegister[command]; err == false {
				fmt.Printf("Unknown command\n")
			} else {
				if c.name == "explore" {
					if len(user_tokens) < 2 {
					fmt.Printf("Location name is not provided\n")
					} else {
						locationName = user_tokens[1]
					}
				} else if c.name == "catch" || c.name == "inspect" {
					if len(user_tokens) < 2 {
					fmt.Printf("Pokemon name is not provided\n")
					} else {
						pokemonName = user_tokens[1]
					}
				}
				c.callback(locationsConfig, locationsCache, playerPokedex, locationName, pokemonName)
			}
		}
	}
}

// func main() {
// 	cache := pokecache.NewCache(10)
// 	pokeapi.GetPokemonsByLocation("16", cache)
// }
