package main

import (
	"fmt"
	"os"
)

type Config struct {
	Shape []int	`json:"shape"`
	Stride int	`json:"stride"`
	Padding string	`json:"padding"`
}

type Component struct {
	Type string `json:"type"`
	Name string	`json:"name"`
	Input string	`json:"input"`
	Config Config	`json:"config"`
}

func main() {
	jsonFile, err := os.Open("./test.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
}
