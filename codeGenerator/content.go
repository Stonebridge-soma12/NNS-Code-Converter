package codeGenerator

import (
	"encoding/json"
	"fmt"
)

const (
	ErrUnsupportedCategoryType = "unsupported category type"
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
		return fmt.Errorf("JSON Error : %s with field %s", err.Error(), "content input")
	}
	err = json.Unmarshal(data["output"], &c.Output)
	if err != nil {
		return fmt.Errorf("JSON Error : %s with field %s", err.Error(), "content output")
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

		switch l.Category {
		case "Layer":
			err = l.Param.Keras.BindKeras(l.Type, layer["param"])
		case "Math":
			err = l.Param.Math.BindMath(l.Type, layer["param"])
		default:
			return fmt.Errorf(ErrUnsupportedCategoryType)
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