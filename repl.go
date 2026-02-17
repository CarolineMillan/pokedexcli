package main

import (
	"strings"
)

func cleanInput(text string) []string {
	//	var slice []string
	lower := strings.ToLower(text)
	slice := strings.Fields(lower)
	return slice
}
