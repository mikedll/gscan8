package main

import (
	"fmt"
	"testing"
)

func TestSearch(*testing.T) {
	
	openDbForProject(false)
	defer closeDbForProject()

	results, err := searchGistFiles(1, "insert")
	
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}
