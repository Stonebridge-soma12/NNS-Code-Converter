package codeGenerator

import (
	"encoding/json"
	"fmt"
)

const (
	InputNodeIndex = 0
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
	layerIdxMap, inputCntMap := c.GetLayerNameToIdxMapAndInputCountMap()

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
	q.Push(InputNodeIndex)
	for !q.Empty() {
		current := q.Pop()
		layerConn := c.Layers[current.(int)].ConnectLayer()
		codes = append(codes, layerConn)

		for _, next := range c.Layers[current.(int)].Output {
			// 조건에 맞는 노드만 큐에 push. 즉 해당 노드의 Input 개수만큼 들어왔을 떄에만 큐에 push.
			inputCntMap[next] -= 1

			if inputCntMap[next] == 0 {
				q.Push(layerIdxMap[next])
			}
		}
	}

	// create model.
	model := fmt.Sprintf(createModel, c.Input, c.Output)
	codes = append(codes, model)

	return codes, nil
}

func (c *Content) GetLayerNameToIdxMapAndInputCountMap() (map[string]int, map[string]int) {
	indexMap := make(map[string]int)
	inputCountMap := make(map[string]int)

	idx := 1
	for _, l := range c.Layers {
		if l.Type == "Input" {
			indexMap[l.Name] = 0
		} else {
			indexMap[l.Name] = idx
			idx += 1
		}
		inputCountMap[l.Name] = len(l.Input)
	}

	return indexMap, inputCountMap
}
