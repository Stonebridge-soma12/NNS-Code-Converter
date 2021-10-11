package codeGenerator

import (
	"encoding/json"
	"fmt"
	"strconv"
)

const (
	ErrUnsupportedKerasLayerType = "unsupported keras layer type"
)

const (
	typeInput = "Input"
	typeConv2D = "Conv2D"
	typeDense = "Dense"
	typeAvgPool2D = "AveragePooling2D"
	typeMaxPool2D = "MaxPool2D"
	typeActivation = "Activation"
	typeDropOut = "Dropout"
	typeBatchNorm = "BatchNormalization"
	typeFlatten = "Flatten"
	typeRescaling = "Rescaling"
	typeReshape = "Reshape"
)

const (
	conv1d       = `Conv1D(filters=%d, kernel_size=%d, strides=%d, padding='%s')`
	conv2d       = `Conv2D(filters=%d, kernel_size=(%d, %d), strides=(%d, %d), padding='%s')`
	dense        = `Dense(units=%d)`
	avgPooling1d = `AveragePooling1D(pool_size=%d, strides=%d, padding='%s')`
	avgPooling2d = `AveragePooling2D(pool_size=(%d, %d), strides=(%d, %d), padding='%s')`
	maxPool2d    = `MaxPool2D(pool_size=(%d, %d), strides=(%d, %d), padding='%s')`
	activation   = `Activation(activation="%s")`
	input        = `Input(shape=(%s))`
	dropout      = `Dropout(rate=%g)`
	batchNorm    = `BatchNormalization(axis=%d, momentum=%g, epsilon=%g)`
	flatten      = `Flatten()`
	rescaling    = `Rescaling(scale=%g, offset=%g)`
	reshape      = `Reshape(target_shape=(%s))`
)

func (k *Keras) BindKeras(t string, data json.RawMessage) error {
	var err error
	err = nil

	switch t {
	case typeInput:
		err = json.Unmarshal(data, &k.Input)
		if err != nil {
			return err
		}
	case typeDense:
		err = json.Unmarshal(data, &k.Dense)
		if err != nil {
			return err
		}
	case typeConv2D:
		err = json.Unmarshal(data, &k.Conv2D)
		if err != nil {
			return err
		}
	case typeAvgPool2D:
		err = json.Unmarshal(data, &k.AveragePooling2D)
		if err != nil {
			return err
		}
	case typeMaxPool2D:
		err = json.Unmarshal(data, &k.MaxPool2D)
		if err != nil {
			return err
		}
	case typeActivation:
		err = json.Unmarshal(data, &k.Activation)
		if err != nil {
			return err
		}
	case typeDropOut:
		err = json.Unmarshal(data, &k.Dropout)
		if err != nil {
			return err
		}
	case typeBatchNorm:
		err = json.Unmarshal(data, &k.BatchNormalization)
		if err != nil {
			return err
		}
	case typeFlatten:
		err = json.Unmarshal(data, &k.Flatten)
		if err != nil {
			return err
		}
	case typeRescaling:
		err = json.Unmarshal(data, &k.Rescaling)
		if err != nil {
			return err
		}
	case typeReshape:
		err = json.Unmarshal(data, &k.Reshape)
		if err != nil {
			return err
		}
	default:
		err = fmt.Errorf(ErrUnsupportedKerasLayerType)
	}

	return err
}

// GetCode converting module to code.
func (k *Keras) GetCode(t string) (string, error) {
	switch t {
	case typeInput:
		return k.Input.GetCode()
	case typeDense:
		return k.Dense.GetCode()
	case typeConv2D:
		return k.Conv2D.GetCode()
	case typeAvgPool2D:
		return k.AveragePooling2D.GetCode()
	case typeMaxPool2D:
		return k.MaxPool2D.GetCode()
	case typeActivation:
		return k.Activation.GetCode()
	case typeDropOut:
		return k.Dropout.GetCode()
	case typeBatchNorm:
		return k.BatchNormalization.GetCode()
	case typeFlatten:
		return k.Flatten.GetCode()
	case typeRescaling:
		return k.Rescaling.GetCode()
	case typeReshape:
		return k.Reshape.GetCode()
	default:
		return "", fmt.Errorf(ErrUnsupportedKerasLayerType)
	}
}

// Conv2D Convolution 2D layer
type Conv2D struct {
	Filters    *int    `json:"filters"`
	KernelSize []int   `json:"kernel_size"`
	Padding    *string `json:"padding"`
	Strides    []int   `json:"strides"`
}

func (c *Conv2D) GetCode() (string, error) {
	err := checkNil(c)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(conv2d, *c.Filters, c.KernelSize[0], c.KernelSize[1], c.Strides[0], c.Strides[1], *c.Padding), nil
}

// Dense (Affine) layer
type Dense struct {
	Units *int `json:"units"`
}

func (d *Dense) GetCode() (string, error) {
	err := checkNil(d)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(dense, *d.Units), nil
}

// AveragePooling2D layer
type AveragePooling2D struct {
	PoolSize []int   `json:"pool_size"`
	Strides  []int   `json:"strides"`
	Padding  *string `json:"padding"`
}

func (a *AveragePooling2D) GetCode() (string, error) {
	err := checkNil(a)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(avgPooling2d, a.PoolSize[0], a.PoolSize[1], a.Strides[0], a.Strides[1], *a.Padding), nil
}

// MaxPool2D layer
type MaxPool2D struct {
	PoolSize []int   `json:"pool_size"`
	Strides  []int   `json:"strides"`
	Padding  *string `json:"padding"`
}

func (m *MaxPool2D) GetCode() (string, error) {
	err := checkNil(m)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(maxPool2d, m.PoolSize[0], m.PoolSize[1], m.Strides[0], m.Strides[1], *m.Padding), nil
}

// Activation
type Activation struct {
	Activation *string `json:"activation"`
}

func (a *Activation) GetCode() (string, error) {
	err := checkNil(a)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(activation, *a.Activation), nil
}

// Input
type Input struct {
	Shape     []int `json:"shape"`
}

func (i *Input) GetCode() (string, error) {
	err := checkNil(i)
	if err != nil {
		return "", err
	}

	var shape string
	for idx := 0; idx < len(i.Shape); idx++ {
		shape += strconv.Itoa(i.Shape[idx])
		if idx < len(i.Shape)-1 {
			shape += ", "
		}
	}

	return fmt.Sprintf(input, shape), nil
}

// Dropout
type Dropout struct {
	Rate *float64 `json:"rate"`
}

func (d *Dropout) GetCode() (string, error) {
	err := checkNil(d)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(dropout, *d.Rate), nil
}

// BatchNormalization
type BatchNormalization struct {
	Axis     *int     `json:"axis"`
	Momentum *float64 `json:"momentum"`
	Epsilon  *float64 `json:"epsilon"`
}

func (b *BatchNormalization) GetCode() (string, error) {
	err := checkNil(b)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(batchNorm, *b.Axis, *b.Momentum, *b.Epsilon), nil
}

// Flatten
type Flatten struct {
	// Flatten has no parameter
}

func (f Flatten) GetCode() (string, error) {
	return fmt.Sprintf(flatten), nil
}

// Rescaling
type Rescaling struct {
	Scale  *float64 `json:"scale"`
	Offset *float64 `json:"offset"`
}

func (r *Rescaling) GetCode() (string, error) {
	err := checkNil(r)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(rescaling, *r.Scale, *r.Offset), nil
}

// Reshape
type Reshape struct {
	TargetShape []int `json:"target_shape"`
}

func (r *Reshape) GetCode() (string, error) {
	err := checkNil(r)
	if err != nil {
		return "", err
	}

	var shape string
	for idx := 0; idx < len(r.TargetShape); idx++ {
		shape += strconv.Itoa(r.TargetShape[idx])
		if idx < len(r.TargetShape)-1 {
			shape += ", "
		}
	}

	return fmt.Sprintf(reshape, shape), nil
}

type Keras struct {
	Conv2D
	Dense
	AveragePooling2D
	MaxPool2D
	Activation
	Input
	Dropout
	BatchNormalization
	Flatten
	Rescaling
	Reshape
}