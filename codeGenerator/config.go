package codeGenerator

import (
	"encoding/json"
	"fmt"
)

const (
	earlyStop   = "early_stop = tf.keras.callbacks.EarlyStopping(monitor='%s', patience=%d)\n"
	lrReduction = "learning_rate_reduction = tf.keras.callbacks.ReduceLROnPlateau(monitor='%s', patience=%d, verbose=1, factor=%g, min_lr=%g)\n"
)

type Config struct {
	OptimizerName         string      `json:"optimizer_name"`
	OptimizerConfig       Optimizer   `json:"optimizer_config"`
	Loss                  string      `json:"loss"`
	Metrics               []string    `json:"metrics"`
	BatchSize             int         `json:"batch_size"`
	Epochs                int         `json:"epochs"`
	Output                string      `json:"output"`
	EarlyStopping         EarlyStop   `json:"early_stop"`
	LearningRateReduction LrReduction `json:"learning_rate_reduction"`
}

// UnmarshalConfig
func (c *Config) UnmarshalConfig(data map[string]json.RawMessage) error {
	// Unmarshal optimizer name
	err := json.Unmarshal(data["optimizer_name"], &c.OptimizerName)
	if err != nil {
		return fmt.Errorf("JSON Error : %s with field %s", err.Error(), "optimizer_name")
	}

	err = c.OptimizerConfig.BindOptimizer(c.OptimizerName, data["optimizer_config"])
	if err != nil {
		return fmt.Errorf("JSON Error : %s with field %s", err.Error(), "optimizer_config")
	}

	err = json.Unmarshal(data["loss"], &c.Loss)
	if err != nil {
		return fmt.Errorf("JSON Error : %s with field %s", err.Error(), "loss")
	}
	err = json.Unmarshal(data["metrics"], &c.Metrics)
	if err != nil {
		return fmt.Errorf("JSON Error : %s with field %s", err.Error(), "metrics")
	}
	err = json.Unmarshal(data["batch_size"], &c.BatchSize)
	if err != nil {
		return fmt.Errorf("JSON Error : %s with field %s", err.Error(), "batch_size")
	}
	err = json.Unmarshal(data["epochs"], &c.Epochs)
	if err != nil {
		return fmt.Errorf("JSON Error : %s with field %s", err.Error(), "epochs")
	}
	err = json.Unmarshal(data["early_stop"], &c.EarlyStopping)
	if err != nil {
		return fmt.Errorf("JSON Error : %s with field %s", err.Error(), "early_stop")
	}
	err = json.Unmarshal(data["learning_rate_reduction"], &c.LearningRateReduction)
	if err != nil {
		return fmt.Errorf("JSON Error : %s with field %s", err.Error(), "learning_rate_reduction")
	}

	return nil
}

// generate compile codes from config.json
func (c *Config) GenConfig() ([]string, error) {
	var codes []string

	// get optimizer
	optimizer, err := c.OptimizerConfig.ToCode(c.OptimizerName)
	if err != nil {
		return nil, err
	}

	// get metrics
	var metrics string
	for i := 1; i <= len(c.Metrics); i++ {
		metrics += fmt.Sprintf("\"%s\"", c.Metrics[i-1])
		if i < len(c.Metrics) {
			metrics += ", "
		}
	}

	// get compile
	compile := fmt.Sprintf("model.compile(optimizer=%s, loss=\"%s\", metrics=[%s])\n", optimizer, c.Loss, metrics)
	codes = append(codes, compile)

	return codes, nil
}

// EarlyStop struct.
type EarlyStop struct {
	Usage    *bool   `json:"usage"`
	Monitor  *string `json:"monitor"`
	Patience *int    `json:"patience"`
}

// EarlyStop generate code.
func (e *EarlyStop) GenCode() (string, error) {
	// if not use early stopping return empty string.
	if !*e.Usage {
		return "", nil
	}

	// if Using early stopping but there is nil field, return err
	err := checkNil(e)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(earlyStop, *e.Monitor, *e.Patience), nil
}

// Learning Rate Reduction struct
type LrReduction struct {
	Usage    *bool    `json:"usage"`
	Monitor  *string  `json:"monitor"`
	Patience *int     `json:"patience"`
	Factor   *float64 `json:"factor"`
	MinLr    *float64 `json:"min_lr"`
}

// LrReduction generate code.
func (l *LrReduction) GenCode() (string, error) {
	// if not using learning rate reduction return empty string.
	if !*l.Usage {
		return "", nil
	}

	// if Using early stopping but there is nil field, return err.
	err := checkNil(l)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(lrReduction, *l.Monitor, *l.Patience, *l.Factor, *l.MinLr), nil
}
