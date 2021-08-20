package CodeGenerator

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
)

type Content struct {
	Output  string   `json:"output"`
	Input   string   `json:"input"`
	Modules []Module `json:"layers"`
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

	result += m.Name
	result += " = "

	result += tf + keras + layers + "." + param
	if m.Input != nil {
		result += "(" + *m.Input + ")\n"
	} else {
		result += "\n"
	}

	return result, err
}

type Project struct {
	Config  Config  `json:"config"`
	Content Content `json:"content"`
}

const (
	importTf      = "import tensorflow as tf\n\n"
	importTfa     = "import tensorflow_addons as tfa\n\n"
	tf            = "tf"
	tfa           = "tfa"
	keras         = ".keras"
	layers        = ".layers"
	createModel   = "model = tf.keras.Model(inputs=%s, outputs=%s)\n\n"
	fitModel      = "model.model.fit(%s, %s, epochs=%d, batch_size=%d, validation_split=%g, callbacks=%s)\n"
	remoteMonitor = "remote_monitor = " + tf + keras + ".callbacks.RemoteMonitor(root='%s', path='%s', field='data', headers=None, send_as_json=True)\n"
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
func (c *Content) GenLayers() ([]string, error) {
	var codes []string

	layers := SortLayers(c.Modules)

	// code converting
	for _, d := range layers {
		layer, err := d.ToCode()
		if err != nil {
			return nil, err
		}

		codes = append(codes, layer)
	}

	// create model.
	model := fmt.Sprintf(createModel, c.Input, c.Output)
	codes = append(codes, model)

	return codes, nil
}

// generate compile codes from config.json
func (c *Config) GenConfig() ([]string, error) {
	var codes []string

	// get optimizer
	optimizer, err := c.OptimizerConfig.ToCode(c.OptimizerName)
	if err != nil {
		return nil, err
	}

	// get metrics
	var metrics string
	for i := 1; i <= len(c.Metrics); i++ {
		metrics += fmt.Sprintf("\"%s\"", c.Metrics[i-1])
		if i < len(c.Metrics) {
			metrics += ", "
		}
	}

	// get compile
	compile := fmt.Sprintf("model.compile(optimizer=%s, loss=\"%s\", metrics=[%s])\n", optimizer, c.Loss, metrics)
	codes = append(codes, compile)

	return codes, nil
}

func (c *Config) GenFit() error {
	var codes []string
	codes = append(codes, importTf)
	codes = append(codes, importTfa)
	codes = append(codes, "import model\n\n")

	// Python comment.
	codes = append(codes, "\n# Callback functions are below if use them.\n")

	es, err := c.EarlyStopping.GenCode()
	if err != nil {
		return err
	}
	codes = append(codes, es)

	lrr, err := c.LearningRateReduction.GenCode()
	if err != nil {
		return err
	}
	codes = append(codes, lrr)

	rm := fmt.Sprintf(
		remoteMonitor,
		"http://localohst:8080",
		"/publish/epoch/end",
	)

	codes = append(codes, rm)
	// add blank line
	codes = append(codes, "\n")

	// callbacks
	var callbacks string
	callbacks += "["
	callbacks += "remote_monitor"
	if *c.LearningRateReduction.Usage {
		callbacks += ", learning_rate_reduction"
	}
	if *c.EarlyStopping.Usage {
		callbacks += ", early_stop"
	}
	callbacks += "]"

	fitCode := fmt.Sprintf(fitModel, "data", "label", c.Epochs, c.BatchSize, 0.3, callbacks)
	codes = append(codes, fitCode)

	// Generate train python file
	err = MakeTextFile(codes, "train.py")
	if err != nil {
		return err
	}

	return nil
}

func GenerateModel(config Config, content Content) error {
	var codes []string
	codes = append(codes, importTf)
	codes = append(codes, importTfa)

	Layers, err := content.GenLayers()
	if err != nil {
		return err
	}
	codes = append(codes, Layers...)

	Configs, err := config.GenConfig()
	if err != nil {
		return err
	}
	codes = append(codes, Configs...)

	// create python file
	err = MakeTextFile(codes, "model.py")

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
	var config map[string]json.RawMessage
	err = json.Unmarshal(data["config"], &config)
	if err != nil {
		return nil, err
	}

	err = project.Config.UnmarshalConfig(config)
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
		case "Input":
			err = json.Unmarshal(layer["param"], &mod.Param.Input)
			if err != nil {
				return nil, err
			}
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
		case "Rescaling":
			err = json.Unmarshal(layer["param"], &mod.Param.Rescaling)
			if err != nil {
				return nil, err
			}
		case "Reshape":
			err = json.Unmarshal(layer["param"], &mod.Param.Reshape)
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
