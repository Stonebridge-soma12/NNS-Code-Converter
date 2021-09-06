package CodeGenerator

import (
	"fmt"
)

type Layer struct {
	Category string  `json:"category"`
	Type     string  `json:"type"`
	Name     string  `json:"name"`
	Input    []string `json:"input"`
	Output   []string `json:"output"`
	Param    Param   `json:"param"`
}

func (l *Layer) GetVariables() (string, error) {
	var result string
	var err error

	switch l.Category {
	case "Layer":
		result, err = l.GenLayerVariable()
		if err != nil {
			return "", err
		}
	case "Math":
		result, err = l.GenMathVariable()
		if err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("invalid node category")
	}

	return result, err
}

func (l *Layer) GenMathVariable() (string, error) {
	param, err := l.Param.Math.GetCode(l.Type)
	if err != nil {
		return "", err
	}

	result := fmt.Sprintf("%s = %s.%s\n", l.Name, tf + keras + math, param)

	return result, nil
}

func (l *Layer) GenLayerVariable() (string, error) {
	param, err := l.Param.GetCode(l.Type)
	if err != nil {
		return "", err
	}

	result := fmt.Sprintf("%s = %s.%s\n", l.Name, tf + keras + layers, param)

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

	result := fmt.Sprintf("%s = %s(%s)\n", l.Name, l.Name, inputs)

	return result
}