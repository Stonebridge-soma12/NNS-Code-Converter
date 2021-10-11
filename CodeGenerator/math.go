package CodeGenerator

import (
	"encoding/json"
	"fmt"
)

const (
	ErrUnsupportedMathType    = "unsupported math layer type"
	ErrInsufficientNumOfInput = "number of input layers is insufficient"
)

const (
	typeAbs   = "Abs"
	typeCeil  = "Ceil"
	typeFloor = "Floor"
	typeRound = "Round"
	typeSqrt  = "Sqrt"
	typeAdd   = "Add"
	typeSub   = "Subtract"
	typeLog   = "Log"
)

const (
	abs      = "abs"
	add      = "add(%s)"
	ceil     = "ceil"
	floor    = "floor"
	round    = "round"
	sqrt     = "sqrt"
	subtract = "subtract"
	log      = "log"
)

type Math struct {
	Abs
	Ceil
	Floor
	Round
	Sqrt
	Add
	Sub
	Log
	Input []string
}

func (m *Math) BindMath(t string, data json.RawMessage) error {
	var err error
	err = nil

	switch t {
	case typeAbs:
		err = json.Unmarshal(data, &m.Abs)
	case typeCeil:
		err = json.Unmarshal(data, &m.Ceil)
	case typeFloor:
		err = json.Unmarshal(data, &m.Floor)
	case typeRound:
		err = json.Unmarshal(data, &m.Round)
	case typeSqrt:
		err = json.Unmarshal(data, &m.Sqrt)
	case typeAdd:
		err = json.Unmarshal(data, &m.Add)
	case typeSub:
		err = json.Unmarshal(data, &m.Sub)
	case typeLog:
		err = json.Unmarshal(data, &m.Log)
	default:
		err = fmt.Errorf(ErrUnsupportedMathType)
	}

	return err
}

func (m *Math) GetCode(t string) (string, error) {
	switch t {
	case typeAbs:
		return m.Abs.GetCode()
	case typeCeil:
		return m.Ceil.GetCode()
	case typeFloor:
		return m.Floor.GetCode()
	case typeRound:
		return m.Round.GetCode()
	case typeSqrt:
		return m.Sqrt.GetCode()
	case typeAdd:
		return m.Add.GetCode(m.Input)
	case typeSub:
		return m.Sub.GetCode(m.Input)
	case typeLog:
		return m.Log.GetCode()
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
		if i < n-1 {
			params += ", "
		}
	}
	params = fmt.Sprintf(add, params)

	return params, nil
}

type Sub struct {
}

func (s *Sub) GetCode(inputs []string) (string, error) {
	n := len(inputs)
	if n < 2 {
		return "", fmt.Errorf(ErrInsufficientNumOfInput)
	}

	var params string

	for i, input := range inputs {
		params += input
		if i < n-1 {
			params += ", "
		}
	}
	params = fmt.Sprintf(subtract, params)

	return params, nil
}

type Log struct {
}

func (l *Log) GetCode() (string, error) {
	return log, nil
}
