package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	Shape []int	`json:"shape"`
	KernelSize []int	`json:"kernel_size"`
	Stride int	`json:"stride"`
	Padding string	`json:"padding"`
}

type Component struct {
	Type string `json:"type"`
	Name string	`json:"name"`
	Input string	`json:"input"`
	Config 	map[string]interface{} `json:"config"`
}

func main() {
	jsonFile, err := os.Open("./test.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var data []Component
	err = json.Unmarshal(byteValue, &data)
	if err != nil {
		fmt.Println(err)
	}

	for _, d := range data {
		// fmt.Println(d.Config)
		for conf := range d.Config {
			// fmt.Println(conf)
			fmt.Println(conf, "=", d.Config[conf])
			//fmt.Printf("%s=%s", conf, d.Config[conf])
		}
	}

}
