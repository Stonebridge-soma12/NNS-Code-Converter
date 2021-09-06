package CodeGenerator

import (
	"encoding/json"
	"fmt"
)

const (
	abs   = "abs"
	add   = "add"
	ceil  = "ceil"
	floor = "floor"
	round = "round"
	sqrt = "sqrt"
)

type Math struct {
	Abs
	Ceil
	Floor
	Round
	Sqrt
}

func (m *Math) Unmarshall(t string, data json.RawMessage) error {
	switch t {
	case "Pow":
		fmt.Println("Pow!")
	default:
		return json.Unmarshal(data, m)
	}

	return nil
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
	default:
		return "", fmt.Errorf("invalid math node type")
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