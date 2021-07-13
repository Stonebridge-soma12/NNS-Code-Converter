package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Project struct {
	Config	`json:"config"`
	Layers []Component	`json:"layers"`
}

type Content struct {
	Output string	`json:"output"`
	Layers []Component	`json:"layers"`
}

type Config struct {
	Optimizer string `json:"optimizer"`
	LearningRate int `json:"learning_rate"`
	Loss string `json:"loss"`
	Metrics []string `json:"metrics"`
	BatchSize int `json:"batch_size"`
	Epochs int `json:"epochs"`
	Output string `json:"output"`
}

type Component struct {
	Type   string                 `json:"type"`
	Name   string                 `json:"name"`
	Input  *string                 `json:"input"`
	Config map[string]string `json:"config"`
}

func digitCheck(target string) bool {
	for i := 0; i < 10; i++ {
		if strings.Contains(target, string(i + '0')) {
			return true
		}
	}
	return false
}

func main() {
	jsonFile, err := os.Open("./content.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var data Content
	var codes []string
	err = json.Unmarshal(byteValue, &data)
	if err != nil {
		fmt.Println(err)
	}

	// code converting
	for _, d := range data.Layers {
		layer := d.Name
		layer += " = tf.keras.layers."
		layer += d.Type
		layer += "("

		i := 1
		for conf := range d.Config {
			var param string

			// if data is array like.
			if strings.Contains(d.Config[conf], ",") {
				param = fmt.Sprintf("%s=(%s)", conf, d.Config[conf])
			} else {
				if digitCheck(d.Config[conf]) {
					param = fmt.Sprintf("%s=%s", conf, d.Config[conf])
				} else {
					param = fmt.Sprintf("%s=\"%s\"", conf, d.Config[conf])
				}
			}
			layer += param
			if i < len(d.Config) {
				layer += ", "
			}
			i++
		}
		layer += ")"
		if d.Input != nil {
			layer += fmt.Sprintf("(%s)\n", *d.Input)
		} else {
			layer += "\n"
		}

		codes = append(codes, layer)
	}

	// create python file
	py, err := os.Create("test.py")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer py.Close()

	fileSize := 0
	for _, code := range codes {
		n, err := py.Write([]byte(code))
		if err != nil {
			fmt.Println(err)
			return
		}
		fileSize += n
	}
	fmt.Printf("Code converting is finish with %d bytes size", fileSize)
}
