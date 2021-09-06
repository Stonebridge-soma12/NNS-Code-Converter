package CodeGenerator

import (
	"encoding/json"
	"fmt"
)

type Content struct {
	Output string  `json:"output"`
	Input  string  `json:"input"`
	Layers []Layer `json:"layers"`
}

func (c *Content) BindContent(data map[string]json.RawMessage) error {
	// Unmarshalling Content.
	// Unmarshalling content input and output.
	err := json.Unmarshal(data["input"], &c.Input)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data["output"], &c.Output)
	if err != nil {
		return err
	}

	// Unmarshalling Layers.
	var layers []map[string]json.RawMessage
	err = json.Unmarshal(data["layers"], &layers)

	for _, layer := range layers {
		// Unmarshalling module informations except parameters.

		l, err := UnmarshalLayer(layer)
		if err != nil {
			return err
		}

		//if l.Category == "Math" {
		//	err = l.Param.Math.Unmarshall(l.Type, layer["param"])
		//	continue
		//}

		switch l.Type {
		case "Input":
			err = json.Unmarshal(layer["param"], &l.Param.Input)
			if err != nil {
				return err
			}
		case "Conv2D":
			err = json.Unmarshal(layer["param"], &l.Param.Conv2D)
			if err != nil {
				return err
			}
		case "Dense":
			err = json.Unmarshal(layer["param"], &l.Param.Dense)
			if err != nil {
				return err
			}
		case "AveragePooling2D":
			err = json.Unmarshal(layer["param"], &l.Param.AveragePooling2D)
			if err != nil {
				return err
			}
		case "MaxPool2D":
			err = json.Unmarshal(layer["param"], &l.Param.MaxPool2D)
			if err != nil {
				return err
			}
		case "Activation":
			err = json.Unmarshal(layer["param"], &l.Param.Activation)
			if err != nil {
				return err
			}
		case "Dropout":
			err = json.Unmarshal(layer["param"], &l.Param.Dropout)
			if err != nil {
				return err
			}
		case "BatchNormalization":
			err = json.Unmarshal(layer["param"], &l.Param.BatchNormalization)
			if err != nil {
				return err
			}
		case "Flatten":
			err = json.Unmarshal(layer["param"], &l.Param.Flatten)
			if err != nil {
				return err
			}
		case "Rescaling":
			err = json.Unmarshal(layer["param"], &l.Param.Rescaling)
			if err != nil {
				return err
			}
		case "Reshape":
			err = json.Unmarshal(layer["param"], &l.Param.Reshape)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("inavlid node type")
		}
		c.Layers = append(c.Layers, l)
	}

	return nil
}

// Generate layer codes from content.json
func (c *Content) GenLayers() ([]string, error) {
	var codes []string
	// TODO: BFS돌리며 레이어 변수 선언 및 연결
	layerIdxMap := c.GetLayerNameToIdxMap()

	// Generate layer variables
	for _, l := range c.Layers {
		layer, err := l.GetVariables()
		if err != nil {
			return nil, err
		}
		codes = append(codes, layer)
	}

	// Connect layers through BFS.
	var q Queue
	q.Push(layerIdxMap["Input_1"])
	for !q.Empty() {
		current := q.Pop()
		layerConn := c.Layers[current.(int)].ConnectLayer()
		codes = append(codes, layerConn)

		for _, next := range c.Layers[current.(int)].Output {
			q.Push(layerIdxMap[next])
		}
	}

	// create model.
	model := fmt.Sprintf(createModel, c.Input, c.Output)
	codes = append(codes, model)

	return codes, nil
}

func (c *Content) GetLayerNameToIdxMap() map[string]int {
	result := make(map[string]int)

	for i, l := range c.Layers {
		result[l.Name] = i
	}

	return result
}
//
//func SortLayers(source []Layer) []Layer {
//	// Sorting layer components via BFS.
//	type node struct {
//		idx    int
//		Output *string
//	}
//
//	var result []Layer             // result Content slice.
//	adj := make(map[string][]node) // adjustment matrix of each nodes.
//	var inputIdx int
//
//	// setup adjustment matrix.
//	for idx, layer := range source {
//		// Input layer is always first.u
//		var input string
//		if layer.Type == "Input" {
//			inputIdx = idx
//
//			// result = append(result, layer)
//		}
//		input = layer.Name
//
//		var nodeSlice []node
//		if adj[input] == nil {
//			nodeSlice = append(nodeSlice, node{idx, layer.Output})
//			adj[input] = nodeSlice
//		} else {
//			prev, _ := adj[input]
//			nodeSlice = prev
//			nodeSlice = append(nodeSlice, node{idx, layer.Output})
//			adj[input] = nodeSlice
//		}
//	}
//
//	// Using BFS with queue
//	var q Queue
//	q.Push(source[inputIdx].Name)
//	for !q.Empty() {
//		current := q.Pop()
//		for _, next := range adj[current] {
//			if next.Output != nil {
//				q.Push(*next.Output)
//			}
//			result = append(result, source[next.idx])
//		}
//	}
//
//	return result
//}