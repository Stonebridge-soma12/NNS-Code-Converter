package CodeGenerator

import (
	"encoding/json"
	"fmt"
)

type Layer struct {
	Category string   `json:"category"`
	Type     string   `json:"type"`
	Name     string   `json:"name"`
	Input    []string `json:"input"`
	Output   []string `json:"output"`
	Param    Param    `json:"param"`
}

const (
	categoryLayer = "Layer"
	categoryMath  = "Math"
)

const (
	ErrUnsupportedCategory = "unsupported category"
)

const (
	variableString   = "%s = %s.%s\n"
	connectionString = "%s = %s(%s)\n"
)

func UnmarshalLayer(data map[string]json.RawMessage) (Layer, error) {
	var res Layer
	err := json.Unmarshal(data["category"], &res.Category)
	if err != nil {
		return res, err
	}
	err = json.Unmarshal(data["type"], &res.Type)
	if err != nil {
		return res, err
	}
	err = json.Unmarshal(data["name"], &res.Name)
	if err != nil {
		return res, err
	}
	err = json.Unmarshal(data["input"], &res.Input)
	if err != nil {
		return res, err
	}
	err = json.Unmarshal(data["output"], &res.Output)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (l *Layer) GetVariables() (string, error) {
	var result string
	var err error

	switch l.Category {
	case categoryLayer:
		result, err = l.GenLayerVariable()
		if err != nil {
			return "", err
		}
	case categoryMath:
		result, err = l.GenMathVariable()
		if err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf(ErrUnsupportedCategory)
	}

	return result, err
}

func (l *Layer) GenMathVariable() (string, error) {
	l.Param.Math.Input = l.Input
	param, err := l.Param.Math.GetCode(l.Type)
	if err != nil {
		return "", err
	}

	result := fmt.Sprintf(variableString, l.Name, tf+keras+math, param)

	return result, nil
}

func (l *Layer) GenLayerVariable() (string, error) {
	param, err := l.Param.Keras.GetCode(l.Type)
	if err != nil {
		return "", err
	}

	result := fmt.Sprintf(variableString, l.Name, tf+keras+layers, param)

	return result, nil
}

func (l *Layer) ConnectLayer() string {
	if len(l.Input) == 0 {
		return ""
	}

	inputs := l.Input[0]
	for i := 1; i < len(l.Input); i++ {
		inputs += fmt.Sprintf(", %s", l.Input[i])
	}

	result := fmt.Sprintf(connectionString, l.Name, l.Name, inputs)

	return result
}
