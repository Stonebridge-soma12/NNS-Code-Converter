package codeGenerator

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	optimizer = ".optimizers"
	adaDelta  = tf + keras + optimizer + `.Adadelta(learning_rate=%g, rho=%g, epsilon=%g)`
	adaGrad   = tf + keras + optimizer + `.Adagrad(learning_rate=%g, initial_accumulator_value=%g, epsilon=%g)`
	adam      = tf + keras + optimizer + `.Adam(learning_rate=%g, beta_1=%g, beta_2=%g, epsilon=%g, amsgrad=%s)`
	adamax    = tf + keras + optimizer + `.Adamax(learning_rate=%g, beta_1=%g, beta_2=%g, epsilon=%g)`
	nadam     = tf + keras + optimizer + `.Nadam(learning_rate=%g, beta_1=%g, beta_2=%g, epsilon=%g)`
	rmsprop   = tf + keras + optimizer + `.RMSprop(learning_rate=%g, rho=%g, momentum=%g, epsilon=%g, centered=%s)`
	sgd       = tf + keras + optimizer + `.SGD(learning_rate=%g, momentum=%g, nesterov=%s)`
	adamw     = tfa + optimizer + `.AdamW(weight_decay=%g, learning_rate=%g, beta_1=%g, beta_2=%g, epsilon=%g, amsgrad=%s)`
	sgdw      = tfa + optimizer + `.SGDW(weight_decay=%g, learning_rate=%g, momentum=%g, nesterov=%s)`
)

func BoolToStr(t bool) string {
	return strings.Title(fmt.Sprintf("%t", t))
}

type Optimizer struct {
	Adadelta
	Adagrad
	Adam
	Adamax
	Nadam
	RMSprop
	SGD
	AdamW
	SGDW
}

func (o *Optimizer) BindOptimizer(name string,data json.RawMessage) error {
	//Unmarshal optimizer config
	switch strings.Title(name) {
	case "Adadelta":
		err := json.Unmarshal(data, &o.Adadelta)
		if err != nil {
			return err
		}
	case "Adadgrad":
		err := json.Unmarshal(data, &o.Adagrad)
		if err != nil {
			return err
		}
	case "Adam":
		err := json.Unmarshal(data, &o.Adam)
		if err != nil {
			return err
		}
	case "Adamax":
		err := json.Unmarshal(data, &o.Adamax)
		if err != nil {
			return err
		}
	case "Nadam":
		err := json.Unmarshal(data, &o.Nadam)
		if err != nil {
			return err
		}
	case "RMSprop":
		err := json.Unmarshal(data, &o.RMSprop)
		if err != nil {
			return err
		}
	case "SGD":
		err := json.Unmarshal(data, &o.SGD)
		if err != nil {
			return err
		}
	case "AdamW":
		err := json.Unmarshal(data, &o.AdamW)
		if err != nil {
			return err
		}
	case "SGDW":
		err := json.Unmarshal(data, &o.SGDW)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid optimizer")
	}

	return nil
}

func (o *Optimizer) ToCode(name string) (string, error) {
	switch strings.Title(name) {
	case "Adadelta":
		return o.Adadelta.ToCode()
	case "Adadgrad":
		return o.Adagrad.ToCode()
	case "Adam":
		return o.Adam.ToCode()
	case "Adamax":
		return o.Adamax.ToCode()
	case "Nadam":
		return o.Nadam.ToCode()
	case "RMSprop":
		return o.RMSprop.ToCode()
	case "SGD":
		return o.SGD.ToCode()
	case "AdamW":
		return o.AdamW.ToCode()
	case "SGDW":
		return o.SGDW.ToCode()

	default:
		return "", fmt.Errorf("Invalid optimizer")
	}
}

// Adadelta
type Adadelta struct {
	LearningRate *float64 `json:"learning_rate"`
	Decay        *float64 `json:"weight_decay"`
	Epsilon      *float64 `json:"epsilon"`
}

func (a *Adadelta) ToCode() (string, error) {
	err := checkNil(a)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(adaDelta, *a.LearningRate, *a.Decay, *a.Epsilon), nil
}

// Adagrad
type Adagrad struct {
	LearningRate    *float64 `json:"learning_rate"`
	InitAccumulator *float64 `json:"initial_accumulator_value"`
	Epsilon         *float64 `json:"epsilon"`
}

func (a *Adagrad) ToCode() (string, error) {
	err := checkNil(a)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(adaGrad, *a.LearningRate, *a.InitAccumulator, *a.Epsilon), nil
}

// Adam
type Adam struct {
	LearningRate *float64 `json:"learning_rate"`
	Beta1        *float64 `json:"beta_1"`
	Beta2        *float64 `json:"beta_2"`
	Epsilon      *float64 `json:"epsilon"`
	AmsGrad      *bool    `json:"amsgrad"`
}

func (a *Adam) ToCode() (string, error) {
	err := checkNil(a)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(adam, *a.LearningRate, *a.Beta1, *a.Beta2, *a.Epsilon, BoolToStr(*a.AmsGrad)), nil
}

// Adamax
type Adamax struct {
	LearningRate *float64 `json:"learning_rate"`
	Beta1        *float64 `json:"beta_1"`
	Beta2        *float64 `json:"beta_2"`
	Epsilon      *float64 `json:"epsilon"`
}

func (a *Adamax) ToCode() (string, error) {
	err := checkNil(a)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(adamax, *a.LearningRate, *a.Beta1, *a.Beta2, *a.Epsilon), nil
}

// Nadam
type Nadam struct {
	LearningRate *float64 `json:"learning_rate"`
	Beta1        *float64 `json:"beta_1"`
	Beta2        *float64 `json:"beta_2"`
	Epsilon      *float64 `json:"epsilon"`
}

func (n *Nadam) ToCode() (string, error) {
	err := checkNil(n)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(nadam, *n.LearningRate, *n.Beta1, *n.Beta2, *n.Epsilon), nil
}

// RMSprop
type RMSprop struct {
	LearningRate *float64 `json:"learning_rate"`
	Decay        *float64 `json:"decay"`
	Momentum     *float64 `json:"momentum"`
	Epsilon      *float64 `json:"epsilon"`
	Centered     *bool    `json:"centered"`
}

func (r *RMSprop) ToCode() (string, error) {
	err := checkNil(r)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(rmsprop, *r.LearningRate, *r.Decay, *r.Momentum, *r.Epsilon, BoolToStr(*r.Centered)), nil
}

// SGD
type SGD struct {
	LearningRate *float64 `json:"learning_rate"`
	Momentum     *float64 `json:"momentum"`
	Nesterov     *bool    `json:"nesterov"`
}

func (s *SGD) ToCode() (string, error) {
	err := checkNil(s)
	if err != nil {
		return "", nil
	}

	return fmt.Sprintf(sgd, *s.LearningRate, *s.Momentum, BoolToStr(*s.Nesterov)), nil
}

// AdamW
type AdamW struct {
	WeightDecay  *float64 `json:"weight_decay"`
	LearningRate *float64 `json:"learning_rate"`
	Beta1        *float64 `json:"beta_1"`
	Beta2        *float64 `json:"beta_2"`
	Epsilon      *float64 `json:"epsilon"`
	Amsgrad      *bool    `json:"amsgrad"`
}

func (a *AdamW) ToCode() (string, error) {
	err := checkNil(a)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(adamw, *a.WeightDecay, *a.LearningRate, *a.Beta1, *a.Beta2, *a.Epsilon, BoolToStr(*a.Amsgrad)), nil
}

// SGDW
type SGDW struct {
	WeightDecay  *float64 `json:"weight_decay"`
	LearningRate *float64 `json:"learning_rate"`
	Momentum     *float64 `json:"momentum"`
	Nesterov     *bool    `json:"nesterov"`
}

func (s *SGDW) ToCode() (string, error) {
	err := checkNil(s)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(sgdw, *s.WeightDecay, *s.LearningRate, *s.Momentum, BoolToStr(*s.Nesterov)), nil
}
