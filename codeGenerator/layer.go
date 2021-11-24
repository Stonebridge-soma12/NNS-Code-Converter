package codeGenerator

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

var commentKerasMap = map[string]string{
	"Input":      "# Input layer\n",
	"Conv2D":     "# This layer creates a convolution kernel that is convolved with the layer input to produce a tensor of outputs.\n",
	"Dense":      "# Dense\n",
	"AvgPool2D":  "# AvgPool2D\n",
	"MaxPool2D":  "# MaxPool2D\n",
	"Activation": "# Activation\n",
	"DropOut":    "# Dropout\n",
	"BatchNorm":  "# BatchNorm\n",
	"Flatten":    "# Flatten\n",
	"Rescaling":  "# Rescaling\n",
	"Reshape":    "# Reshape\n",
}

var commentMathMap = map[string]string{
	"Abs":      "# abs\n",
	"Add":      "# add\n",
	"Ceil":     "# ceil\n",
	"Floor":    "# floor\n",
	"Round":    "# round\n",
	"Sqrt":     "# sqrt\n",
	"Subtract": "# subtract\n",
	"Log":      "# log\n",
}

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

func (l *Layer) GetVariables() ([]string, error) {
	var result []string
	var err error

	switch l.Category {
	case categoryLayer:
		result, err = l.GenLayerVariable()
		if err != nil {
			return nil, err
		}
	case categoryMath:
		result, err = l.GenMathVariable()
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf(ErrUnsupportedCategory)
	}

	return result, err
}

func (l *Layer) GenMathVariable() ([]string, error) {
	l.Param.Math.Input = l.Input
	param, err := l.Param.Math.GetCode(l.Type)
	if err != nil {
		return nil, err
	}
	var result []string
	result = append(result, commentMathMap[l.Type])
	result = append(result, fmt.Sprintf(variableString, l.Name, tf+math, param))

	return result, nil
}

func (l *Layer) GenLayerVariable() ([]string, error) {
	param, err := l.Param.Keras.GetCode(l.Type)
	if err != nil {
		return nil, err
	}
	var result []string
	result = append(result, commentKerasMap[l.Type])
	result = append(result, fmt.Sprintf(variableString, l.Name, tf+keras+layers, param))

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
