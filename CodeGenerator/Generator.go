package CodeGenerator

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type Content struct {
	Output  string   `json:"output"`
	Input   string   `json:"input"`
	Modules []Module `json:"layers"`
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

type Module struct {
	Category string  `json:"category"`
	Type     string  `json:"type"`
	Name     string  `json:"name"`
	Input    *string `json:"input"`
	Output   *string `json:"output"`
	Param    Param   `json:"param"`
}

func (m *Module) ToCode() (string, error) {
	var result string
	param, err := m.Param.ToCode(m.Type)

	if m.Input != nil {
		result += m.Name
		result += " = "
	}
	result += tf + keras + "." + param
	if m.Output != nil {
		result += "(" + *m.Output + ")\n"
	}

	return result, err
}

type Project struct {
	Config  Config  `json:"config"`
	Content Content `json:"content"`
}

const (
	importTf    = "import tensorflow as tf\n\n"
	tf          = "tf"
	keras       = ".keras"
	createModel = "model = tf.keras.Model(inputs=%s, outputs=%s)\n\n"
)

func digitCheck(target string) bool {
	re, err := regexp.Compile("\\d")
	if err != nil {
		panic(err)
	}

	return re.MatchString(target)
}

func SortLayers(source []Module) []Module {
	// Sorting layer components via BFS.
	type node struct {
		idx    int
		Output *string
	}

	var result []Module            // result Content slice.
	adj := make(map[string][]node) // adjustment matrix of each nodes.
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
			nodeSlice = append(nodeSlice, node{idx, layer.Output})
			adj[input] = nodeSlice
		} else {
			prev, _ := adj[input]
			nodeSlice = prev
			nodeSlice = append(nodeSlice, node{idx, layer.Output})
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

// Generate layer codes from content.json
func GenLayers(content Content) ([]string, error) {
	var codes []string

	layers := SortLayers(content.Modules)

	// code converting
	for _, d := range layers {
		//layer := d.Name
		//layer += " = "
		layer, err := d.ToCode()
		if err != nil {
			return nil, err
		}
		//layer += params
		//if d.Input != nil {
		//	layer += fmt.Sprintf("(%s)\n", *d.Input)
		//} else {
		//	layer += "\n"
		//}

		codes = append(codes, layer)
	}

	// create model.
	model := fmt.Sprintf(createModel, content.Input, content.Output)
	codes = append(codes, model)

	return codes, nil
}

// generate compile codes from config.json
func GenConfig(config Config) []string {
	var codes []string

	// get optimizer
	optimizer := fmt.Sprintf("%s.optimizers.%s(learning_rate=%g)", tf+keras, strings.Title(config.Optimizer), config.LearningRate)

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
	codes = append(codes, importTf)

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

func BindProject(r *http.Request) (*Project, error) {
	project := new(Project)
	data := make(map[string]json.RawMessage)
	cc := make(map[string]json.RawMessage)
	var layers []map[string]json.RawMessage

	// Binding request body
	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		return nil, err
	}

	// Unmarshalling Config.
	err = json.Unmarshal(data["config"], &project.Config)
	if err != nil {
		return nil, err
	}

	// Unmarshalling Content.
	err = json.Unmarshal(data["content"], &cc)
	if err != nil {
		return nil, err
	}

	// Unmarshalling content input and output.
	err = json.Unmarshal(cc["input"], &project.Content.Input)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(cc["output"], &project.Content.Output)
	if err != nil {
		return nil, err
	}

	// Unmarshalling Modules.
	err = json.Unmarshal(cc["layers"], &layers)

	for _, layer := range layers {
		// Unmarshalling module informations except parameters.
		mod, err := UnmarshalModule(layer)
		if err != nil {
			return nil, err
		}
		switch mod.Type {
		case "Conv2D":
			err = json.Unmarshal(layer["param"], &mod.Param.Conv2D)
			if err != nil {
				return nil, err
			}
		case "Dense":
			err = json.Unmarshal(layer["param"], &mod.Param.Dense)
			if err != nil {
				return nil, err
			}
		case "AveragePooling2D":
			err = json.Unmarshal(layer["param"], &mod.Param.AveragePooling2D)
			if err != nil {
				return nil, err
			}
		case "MaxPool2D":
			err = json.Unmarshal(layer["param"], &mod.Param.MaxPool2D)
			if err != nil {
				return nil, err
			}
		case "Activation":
			err = json.Unmarshal(layer["param"], &mod.Param.Activation)
			if err != nil {
				return nil, err
			}
		case "Dropout":
			err = json.Unmarshal(layer["param"], &mod.Param.Dropout)
			if err != nil {
				return nil, err
			}
		case "BatchNormalization":
			err = json.Unmarshal(layer["param"], &mod.Param.BatchNormalization)
			if err != nil {
				return nil, err
			}
		case "Flatten":
			err = json.Unmarshal(layer["param"], &mod.Param.Flatten)
			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("inavlid node type")
		}
		project.Content.Modules = append(project.Content.Modules, mod)
	}

	return project, nil
}
