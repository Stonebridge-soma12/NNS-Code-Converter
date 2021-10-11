package CodeGenerator

import (
	"encoding/json"
	"fmt"
)

const (
	ErrUnsupportedMathType = "unsupported math layer type"
	ErrInsufficientNumOfInput = "number of input layers is insufficient"
)

const (
	abs   = "abs"
	add   = "add_n(%s)"
	ceil  = "ceil"
	floor = "floor"
	round = "round"
	sqrt  = "sqrt"
)

type Math struct {
	Abs
	Ceil
	Floor
	Round
	Sqrt
	Add
	Input []string
}

func (m *Math) BindMath(t string, data json.RawMessage) error {
	var err error
	err = nil

	switch t {
	case "Abs":
		err = json.Unmarshal(data, &m.Abs)
	case "Ceil":
		err = json.Unmarshal(data, &m.Ceil)
	case "Floor":
		err = json.Unmarshal(data, &m.Floor)
	case "Round":
		err = json.Unmarshal(data, &m.Round)
	case "Sqrt":
		err = json.Unmarshal(data, &m.Sqrt)
	case "Add":
		err = json.Unmarshal(data, &m.Add)
	default:
		err = fmt.Errorf("unspported math layer type")
	}

	return err
}

func (m *Math) GetCode(t string) (string, error) {
	switch t {
	case "Abs":
		return m.Abs.GetCode()
	case "Ceil":
		return m.Ceil.GetCode()
	case "Floor":
		return m.Floor.GetCode()
	case "Round":
		return m.Round.GetCode()
	case "Sqrt":
		return m.Sqrt.GetCode()
	case "Add":
		return m.Add.GetCode(m.Input)
	default:
		return "", fmt.Errorf(ErrUnsupportedMathType)
	}
}

type Abs struct {
}

func (a *Abs) GetCode() (string, error) {
	return abs, nil
}

type Ceil struct {
}

func (c *Ceil) GetCode() (string, error) {
	return ceil, nil
}

type Floor struct {
}

func (f *Floor) GetCode() (string, error) {
	return floor, nil
}

type Round struct {
}

func (r *Round) GetCode() (string, error) {
	return round, nil
}

type Sqrt struct {
}

func (s *Sqrt) GetCode() (string, error) {
	return sqrt, nil
}

type Add struct {
}

func (a *Add) GetCode(inputs []string) (string, error) {
	n := len(inputs)
	if n < 2 {
		return "", fmt.Errorf(ErrInsufficientNumOfInput)
	}

	var params string

	for i, input := range inputs {
		params += input
		if i < n - 1 {
			params += ", "
		}
	}
	params= fmt.Sprintf(add, params)

	return params, nil
}
