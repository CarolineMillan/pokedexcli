package main

import (
	"bufio"
	"github.com/carolinemillan/pokedexcli/internal/pokecache"
	"os"
	"time"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	c := &config{}
	// initialise config here
	c.cache = pokecache.NewCache(5 * time.Minute)
	startREPL(c, scanner)
}
