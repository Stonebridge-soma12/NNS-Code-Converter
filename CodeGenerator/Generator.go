package CodeGenerator

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Content struct {
	Output string      `json:"output"`
	Input  string      `json:"input"`
	Layers []Component `json:"layers"`
}

type Config struct {
	Optimizer    string   `json:"optimizer"`
	LearningRate float64  `json:"learning_rate"`
	Loss         string   `json:"loss"`
	Metrics      []string `json:"metrics"`
	BatchSize    int      `json:"batch_size"`
	Epochs       int      `json:"epochs"`
	Output       string   `json:"output"`
}

type Component struct {
	Category string            `json:"category"`
	Type     string            `json:"type"`
	Name     string            `json:"name"`
	Input    *string           `json:"input"`
	Output   *string           `json:"output"`
	Config   map[string]string `json:"config"`
}

type Project struct {
	Config  Config  `json:"config"`
	Content Content `json:"content"`
}

const (
	ImportTF = "import tensorflow as tf\n\n"
	TF = "tf"
	Keras = ".keras"
	Layer = ".layers"
	Math = ".math"
	CreateModel = "model = %s.Model(inputs=%s, outputs=%s)\n\n"
)


var category = map[string]string{
	"Layer": TF + Keras + Layer,
	"Math":  TF + Math,
}

func digitCheck(target string) bool {
	re, err := regexp.Compile("\\d")
	if err != nil {
		panic(err)
	}

	return re.MatchString(target)
}

func SortLayers(source []Component) []Component {
	// Sorting layer components via BFS.
	type node struct {
		idx int
		Output *string
	}

	var result []Component	// result Content slice.
	adj := make(map[string][]node)	// adjustment matrix of each nodes.
	var inputIdx int

	// setup adjustment matrix.
	for idx, layer := range source {
		// Input layer is always first.u
		var input string
		if layer.Type == "Input" {
			inputIdx = idx

			// result = append(result, layer)
		}
		input = layer.Name

		var nodeSlice []node
		if adj[input] == nil {
			nodeSlice = append(nodeSlice, node{ idx, layer.Output })
			adj[input] = nodeSlice
		} else {
			prev, _ := adj[input]
			nodeSlice = prev
			nodeSlice = append(nodeSlice, node{ idx, layer.Output })
			adj[input] = nodeSlice
		}
	}

	// Using BFS with queue
	var q Queue
	q.Push(source[inputIdx].Name)
	for !q.Empty() {
		current := q.Pop()
		for _, next := range adj[current] {
			if next.Output != nil {
				q.Push(*next.Output)
			}
			result = append(result, source[next.idx])
		}
	}

	return result
}

// Generate Layer codes from content.json
func GenLayers(content Content) ([]string, error) {
	var codes []string

	layers := SortLayers(content.Layers)

	// code converting
	for _, d := range layers {
		layer := d.Name
		layer += " = "
		layer += category[d.Category] + "."
		layer += d.Type

		layer += "("
		i := 1
		for conf := range d.Config {
			if d.Config[conf] == "" {
				return nil, errors.New("The value of layer config is empty")
			}
			var param string

			// if data is array like.
			if strings.Contains(d.Config[conf], ",") {
				// if tuple.
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

	// create model.
	model := fmt.Sprintf(CreateModel, TF+Keras, content.Input, content.Output)
	codes = append(codes, model)

	return codes, nil
}

// generate compile codes from config.json
func GenConfig(config Config) []string {
	var codes []string

	// get optimizer
	optimizer := fmt.Sprintf("%s.optimizers.%s(lr=%f)", TF+Keras, strings.Title(config.Optimizer), config.LearningRate)

	// get metrics
	var metrics string
	for i := 1; i <= len(config.Metrics); i++ {
		metrics += fmt.Sprintf("\"%s\"", config.Metrics[i-1])
		if i < len(config.Metrics) {
			metrics += ", "
		}
	}

	// get compile
	compile := fmt.Sprintf("model.compile(optimizer=%s, loss=\"%s\", metrics=[%s])\n", optimizer, config.Loss, metrics)
	codes = append(codes, compile)

	return codes
}

func GenerateModel(config Config, content Content) error {
	var codes []string
	codes = append(codes, ImportTF)

	Layers, err := GenLayers(content)
	if err != nil {
		return err
	}
	codes = append(codes, Layers...)

	codes = append(codes, GenConfig(config)...)

	// create python file
	py, err := os.Create("model.py")
	if err != nil {
		return err
	}
	defer py.Close()

	fileSize := 0
	for _, code := range codes {
		n, err := py.Write([]byte(code))
		if err != nil {
			return err
		}
		fileSize += n
	}
	fmt.Printf("Code converting is finish with %d bytes size\n", fileSize)

	return nil
}
