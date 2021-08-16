package CodeGenerator

import (
	"encoding/json"
	"fmt"
	"reflect"
)

const (
	conv2d       = `Conv2D(filters=%d, kernel_size=(%d, %d), strides=(%d, %d), padding="%s")`
	dense        = `Dense(units=%d)`
	avgPooling2d = `AveragePooling2D(pool_size=(%d, %d), strides=(%d, %d), padding="%s")`
	maxPool2d    = `MaxPool2D(pool_size=(%d, %d), strides=(%d, %d), padding="%s")`
	activation   = `Activation(activation="%s"`
	input        = `Input(shape=(%s))`
	dropout      = `Dropout(rate=%g)`
	batchNorm    = `BatchNormalization(axis=%d, momentum=%g, epsilon=%g)`
	flatten      = `Flatten()`
)

type Param struct {
	Conv2D
	Dense
	AveragePooling2D
	MaxPool2D
	Activation
	Input
	Dropout
	BatchNormalization
	Flatten
}

func UnmarshalModule(data map[string]json.RawMessage) (Module, error) {
	var res Module
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

// ToCode converting module to code.
func (p *Param) ToCode(t string) (string, error) {
	switch t {
	case "Dense":
		return p.Dense.ToCode()
	case "Conv2D":
		return p.Conv2D.ToCode()
	case "AveragePooling2D":
		return p.AveragePooling2D.ToCode()
	case "MaxPool2D":
		return p.MaxPool2D.ToCode()
	case "Activation":
		return p.Activation.ToCode()
	case "Dropout":
		return p.Dropout.ToCode()
	case "BatchNormalization":
		return p.BatchNormalization.ToCode()
	case "Flatten":
		return p.Flatten.ToCode()
	default:
		return "", fmt.Errorf("The type is not available")
	}
}

// For check there is any empty fields in Param
func checkNil(object interface{}) string {
	errorString := ""

	e := reflect.ValueOf(object).Elem()
	n := e.NumField()
	for i := 0; i < n; i++ {
		value := e.Field(i)
		tType := e.Type()

		// append error which field is nil
		if reflect.ValueOf(value.Interface()).IsNil() {
			errorString += fmt.Sprintf("field %s is nil\n", tType.Field(i).Name)
		}
	}

	return errorString
}

// Conv2D Convolution 2D layer
type Conv2D struct {
	Filters    *int    `json:"filters"`
	KernelSize []int   `json:"kernel_size"`
	Padding    *string `json:"padding"`
	Strides    []int   `json:"strides"`
}

func (c *Conv2D) ToCode() (string, error) {
	err := checkNil(c)
	if err != "" {
		return "", fmt.Errorf(err)
	}

	return fmt.Sprintf(conv2d, *c.Filters, c.KernelSize[0], c.KernelSize[1], c.Strides[0], c.Strides[1], *c.Padding), nil
}

// Dense (Affine) layer
type Dense struct {
	Units *int `json:"units"`
}

func (d *Dense) ToCode() (string, error) {
	err := checkNil(d)
	if err != "" {
		return "", fmt.Errorf(err)
	}

	return fmt.Sprintf(dense, *d.Units), nil
}

// AveragePooling2D layer
type AveragePooling2D struct {
	PoolSize []int   `json:"pool_size"`
	Strides  []int   `json:"strides"`
	Padding  *string `json:"padding"`
}

func (a *AveragePooling2D) ToCode() (string, error) {
	err := checkNil(a)
	if err != "" {
		return "", fmt.Errorf(err)
	}

	return fmt.Sprintf(avgPooling2d, a.PoolSize[0], a.PoolSize[1], a.Strides[0], a.Strides[1], *a.Padding), nil
}

// MaxPool2D layer
type MaxPool2D struct {
	PoolSize []int   `json:"pool_size"`
	Strides  []int   `json:"strides"`
	Padding  *string `json:"padding"`
}

func (m *MaxPool2D) ToCode() (string, error) {
	err := checkNil(m)
	if err != "" {
		return "", fmt.Errorf(err)
	}

	return fmt.Sprintf(maxPool2d, m.PoolSize[0], m.PoolSize[1], m.Strides[0], m.Strides[1], *m.Padding), nil
}

// Activation
type Activation struct {
	Activation *string `json:"activation"`
}

func (a *Activation) ToCode() (string, error) {
	err := checkNil(a)
	if err != "" {
		return "", fmt.Errorf(err)
	}

	return fmt.Sprintf(activation, *a.Activation), nil
}

// Input
type Input struct {
	Shape []int `json:"shape"`
}

func (i *Input) ToCode() (string, error) {
	err := checkNil(i)
	if err != "" {
		return "", fmt.Errorf(err)
	}

	var shape string
	for idx := 0; idx < len(i.Shape); idx++ {
		shape += string(rune(i.Shape[idx]))
		if idx < idx-1 {
			shape += ", "
		}
	}

	return fmt.Sprintf(input, shape), nil
}

// Dropout
type Dropout struct {
	Rate float64 `json:"rate"`
}

func (d Dropout) ToCode() (string, error) {
	err := checkNil(d)
	if err != "" {
		return "", fmt.Errorf(err)
	}

	return fmt.Sprintf(dropout, d.Rate), nil
}

// BatchNormalization
type BatchNormalization struct {
	Axis     int     `json:"axis"`
	Momentum float64 `json:"momentum"`
	Epsilon  float64 `json:"epsilon"`
}

func (b BatchNormalization) ToCode() (string, error) {
	err := checkNil(b)
	if err != "" {
		return "", fmt.Errorf(err)
	}

	return fmt.Sprintf(batchNorm, b.Axis, b.Momentum, b.Epsilon), nil
}

// Flatten
type Flatten struct {
	// Flatten has no parameter
}

func (f Flatten) ToCode() (string, error) {
	return fmt.Sprintf(flatten), nil
}