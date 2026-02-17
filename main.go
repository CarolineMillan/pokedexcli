package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	c := &config{}
	startREPL(c, scanner)
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
		err := val.callback(c)
		if err != nil {
			fmt.Print(err)
		}
		fmt.Print("\n")

	}
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

type config struct {
	next     *string
	previous *string
}

type locationsResponse struct {
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"results"`
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
	}
}

func commandExit(_ *config) error {
	fmt.Print("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(_ *config) error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n")
	for _, val := range getCommands() {
		fmt.Printf("\n%v: %v", val.name, val.description)
	}
	return nil
}

func commandMap(c *config) error {

	url := "https://pokeapi.co/api/v2/location-area/"
	if c.next != nil {
		url = *c.next
	}
	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Error getting location areas: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("bad status: %s", res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	// unmarshall the json data
	var data locationsResponse
	err = json.Unmarshal(body, &data)
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

func commandMapB(c *config) error {
	// can i just update the config

	if c.previous == nil {
		fmt.Print("you're on the first page")
		return nil
	}
	res, err := http.Get(*c.previous)
	if err != nil {
		return fmt.Errorf("Error getting location areas: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("bad status: %s", res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	// unmarshall the json data
	var data locationsResponse
	err = json.Unmarshal(body, &data)
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
