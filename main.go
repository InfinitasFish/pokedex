package main

import (
	"fmt"
	"bufio"
	"os"
	"internal/pokeapi"
	"internal/pokecache"
)


type cliCommand struct {
	name string
	description string
	callback func(*pokeapi.Config, *pokecache.Cache, string) error
}

func commandExit(config *pokeapi.Config, cache *pokecache.Cache, locationName string) error {
	fmt.Printf("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}

func commandHelp(config *pokeapi.Config, cache *pokecache.Cache, locationName string) error {
	fmt.Printf("Welcome to the Pokedex!\n")
	fmt.Printf("Usage:\n\n")
	fmt.Printf("help: Displays a help message\n")
	fmt.Printf("exit: Exit the Pokedex\n")
	return nil
}

func commandMap(config *pokeapi.Config, cache *pokecache.Cache, locationName string) error {
	// get next locations
	pokeapi.GetLocationsData(true, config, cache)

	// print locations
	for _, m := range config.Results {
		fmt.Printf("%v\n", m["name"])
	}
	return nil
}

func commandMapBack(config *pokeapi.Config, cache *pokecache.Cache, locationName string) error {
	// get previous locations
	pokeapi.GetLocationsData(false, config, cache)

	// print locations
	for _, m := range config.Results {
		fmt.Printf("%v\n", m["name"])
	}
	return nil
}

func commandExplore(config *pokeapi.Config, cache *pokecache.Cache, locationName string) error {
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
	}

	// tracking locations offset
	locationsConfig := &pokeapi.Config{}

	// caching locations to not repeat requests
	locationsCache := pokecache.NewCache(10)
	
	// default location name
	// kinda ugly because other commands than "explore" don't need this
	locationName := ""

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
					fmt.Printf("Location is not provided\n")
					} else {
						locationName = user_tokens[1]
					}
				}
				c.callback(locationsConfig, locationsCache, locationName)
			}
		}
	}
}

// func main() {
// 	cache := pokecache.NewCache(10)
// 	pokeapi.GetPokemonsByLocation("16", cache)
// }
