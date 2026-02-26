package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/carolinemillan/pokedexcli/internal/pokecache"
	"io"
	"net/http"
	"os"
	"strings"
)

func cleanInput(text string) []string {
	//	var slice []string
	lower := strings.ToLower(text)
	slice := strings.Fields(lower)
	return slice
}
func startREPL(c *config, scanner *bufio.Scanner) {
	for {
		fmt.Print("Pokedex > ")

		// check for input
		if !scanner.Scan() {
			break
		}

		// get a handle on the input and clean it up
		clean := cleanInput(scanner.Text())
		if len(clean) == 0 {
			continue
		}

		// look up the command
		commands := getCommands()
		val, ok := commands[clean[0]]
		if !ok {
			fmt.Print("Unknown command")
			continue
		}
		arg := ""
		if len(clean) > 1 {
			arg = clean[1]
		}
		err := val.callback(c, arg)
		if err != nil {
			fmt.Print(err)
		}
		fmt.Print("\n")

	}
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, string) error
}

type config struct {
	next     *string
	previous *string
	cache    *pokecache.Cache
}

type locationsResponse struct {
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"results"`
}

type pokemonsResponse struct {
	Encounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays the names of the next 20 location areas in the Pokemon world",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the names of the previous 20 location areas in the Pokemon world",
			callback:    commandMapB,
		},
		"explore": {
			name:        "explore",
			description: "Displays the names of all pokemon in the given location",
			callback:    commandExplore,
		},
	}
}

func commandExit(_ *config, _ string) error {
	fmt.Print("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(_ *config, _ string) error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n")
	for _, val := range getCommands() {
		fmt.Printf("\n%v: %v", val.name, val.description)
	}
	return nil
}

func commandMap(c *config, _ string) error {
	/// prints the next 20 locations areas in the pokemon world

	url := "https://pokeapi.co/api/v2/location-area/"
	if c.next != nil {
		url = *c.next
	}
	// check whether url is in the cache here
	body, ok := c.cache.Get(url)
	if !ok {

		res, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("Error getting location areas: %w", err)
		}
		defer res.Body.Close()

		if res.StatusCode < 200 || res.StatusCode >= 300 {
			return fmt.Errorf("bad status: %s", res.Status)
		}

		body, err = io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		// add to cache here
		c.cache.Add(url, body)
	}

	// unmarshall the json data
	var data locationsResponse
	err := json.Unmarshal(body, &data)
	if err != nil {
		return err
	}
	// get names of the loc_areas in a slice

	for i := range len(data.Results) {
		if i == 0 {
			fmt.Print(data.Results[0].Name)
		} else {
			fmt.Printf("\n%v", data.Results[i].Name)
		}
	}

	c.next = data.Next
	c.previous = data.Previous
	return nil
}

func commandMapB(c *config, _ string) error {
	/// prints the previous 20 map locations areas in the pokemon world

	if c.previous == nil {
		fmt.Print("you're on the first page")
		return nil
	}
	// check whether url is in the cache here
	body, ok := c.cache.Get(*c.previous)
	if !ok {
		res, err := http.Get(*c.previous)
		if err != nil {
			return fmt.Errorf("Error getting location areas: %w", err)
		}
		defer res.Body.Close()

		if res.StatusCode < 200 || res.StatusCode >= 300 {
			return fmt.Errorf("bad status: %s", res.Status)
		}

		body, err = io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		c.cache.Add(*c.previous, body)
	}
	// unmarshall the json data
	var data locationsResponse
	err := json.Unmarshal(body, &data)
	if err != nil {
		return err
	}
	// get names of the loc_areas in a slice

	for i := range len(data.Results) {
		if i == 0 {
			fmt.Print(data.Results[0].Name)
		} else {
			fmt.Printf("\n%v", data.Results[i].Name)
		}
	}

	c.next = data.Next
	c.previous = data.Previous
	return nil
}

func commandExplore(c *config, loc string) error {
	fmt.Printf("Exploring %s...", loc)
	/// prints the pokemon in this location

	url := "https://pokeapi.co/api/v2/location-area/" + loc

	// check whether url is in the cache here
	body, ok := c.cache.Get(url)
	if !ok {

		res, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("Error getting pokemon in %s: %w", loc, err)
		}
		defer res.Body.Close()

		if res.StatusCode < 200 || res.StatusCode >= 300 {
			return fmt.Errorf("bad status: %s", res.Status)
		}

		body, err = io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		// add to cache here
		c.cache.Add(url, body)
	}

	// unmarshall the json data
	var data pokemonsResponse
	err := json.Unmarshal(body, &data)
	if err != nil {
		return err
	}
	// get names of the pokemon in a slice

	//fmt.Printf("data: %v", data)

	for i := range len(data.Encounters) {
		fmt.Printf("\n%v", data.Encounters[i].Pokemon.Name)
	}

	return nil

}
