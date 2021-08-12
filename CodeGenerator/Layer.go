package CodeGenerator

import (
	"fmt"
	"reflect"
)

const (
	conv2d       = `Conv2D(filters=%d, kernel_size=(%d, %d), strides=(%d, %d), padding="%s")`
	dense        = `Dense(units=%d)`
	avgPooling2d = `AveragePooling2D(pool_size=(%d, %d), strides=(%d, %d), padding="%s")`
	maxPooling2d = `MaxPooling2D(pool_size=(%d, %d), strides=(%d, %d), padding="%s")`
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

// Convert Module to code
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

		fmt.Println(tType.Name(), value.Interface())

		if value.Interface() == nil {
			errorString += fmt.Sprintf("field %s is nil\n", tType.Name())
		}
	}

	return errorString
}

// Convolution 2D layer
type Conv2D struct {
	Filters    *int    `json:"filters"`
	KernelSize []int   `json:"kernel_size"`
	Strides    []int   `json:"strides"`
	Padding    *string `json:"padding"`
}

func (c *Conv2D) ToCode() (string, error) {
	err := checkNil(c)
	if err != "" {
		return "", fmt.Errorf(err)
	}

	return fmt.Sprintf(conv2d, *c.Filters, c.KernelSize[0], c.KernelSize[1], c.Strides[0], c.Strides[1], *c.Padding), nil
}

// Dense (Affain) layer
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

// Averaget pooling 2D layer
type AveragePooling2D struct {
	PoolSize [2]int `json:"pool_size"`
	Strides  [2]int `json:"strides"`
	Padding  string `json:"padding"`
}

func (a *AveragePooling2D) ToCode() (string, error) {
	err := checkNil(a)
	if err != "" {
		return "", fmt.Errorf(err)
	}

	return fmt.Sprintf(avgPooling2d, a.PoolSize[0], a.PoolSize[1], a.Strides[0], a.Strides[1], a.Padding), nil
}

// Max pooling 2D layer
type MaxPool2D struct {
	PoolSize [2]int `json:"pool_size"`
	Strides  [2]int `json:"strides"`
	Padding  [2]int `json:"padding"`
}

func (m *MaxPool2D) ToCode() (string, error) {
	err := checkNil(m)
	if err != "" {
		return "", fmt.Errorf(err)
	}

	return fmt.Sprintf(maxPooling2d, m.PoolSize[0], m.PoolSize[1], m.Strides[0], m.Strides[1], m.Padding), nil
}

// Activation
type Activation struct {
	Activation string `json:"Activation"`
}

func (a *Activation) ToCode() (string, error) {
	err := checkNil(a)
	if err != "" {
		return "", fmt.Errorf(err)
	}

	return fmt.Sprintf(activation, a.Activation), nil
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
	Epsilon  float64 `json:"Epsilon"`
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
