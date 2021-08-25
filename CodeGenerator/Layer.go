package CodeGenerator

import "fmt"

type Layer struct {
	Category string  `json:"category"`
	Type     string  `json:"type"`
	Name     string  `json:"name"`
	Input    *string `json:"input"`
	Output   *string `json:"output"`
	Param    Param   `json:"param"`
}

func (l *Layer) GetCode() (string, error) {
	var result string
	var err error

	switch l.Category {
	case "Layer":
		result, err = l.GetLayerCode()
		if err != nil {
			return "", err
		}
	case "Math":
		result, err = l.GetMathCode()
		if err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("invalid node category")
	}

	return result, err
}

func (l *Layer) GetLayerCode() (string, error) {
	var result string
	param, err := l.Param.GetCode(l.Type)
	if err != nil {
		return "", err
	}

	result += l.Name
	result += " = "

	result += tf + keras + layers + "." + param
	if l.Input != nil {
		result += "(" + *l.Input + ")\n"
	} else {
		result += "\n"
	}

	return result, nil
}

func (l *Layer) GetMathCode() (string, error) {
	var result string
	param, err := l.Param.Math.GetCode(l.Type)
	if err != nil {
		return "", err
	}

	result = fmt.Sprintf("%s = %s.%s(%s)\n", l.Name, tf + keras + math, param, *l.Input)

	return result, nil
}