package CodeGenerator

import (
	"encoding/json"
	"fmt"
	"strings"
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
	DataSet               DataSet     `json:"data_set"`
}

// UnmarshalConfig
func (c *Config) UnmarshalConfig(data map[string]json.RawMessage) error {
	// Unmarshal optimizer name
	err := json.Unmarshal(data["optimizer_name"], &c.OptimizerName)
	if err != nil {
		return err
	}

	//Unmarshal optimizer config
	switch strings.Title(c.OptimizerName) {
	case "Adadelta":
		err := json.Unmarshal(data["optimizer_config"], &c.OptimizerConfig.Adadelta)
		if err != nil {
			return err
		}
	case "Adadgrad":
		err := json.Unmarshal(data["optimizer_config"], &c.OptimizerConfig.Adagrad)
		if err != nil {
			return err
		}
	case "Adam":
		err := json.Unmarshal(data["optimizer_config"], &c.OptimizerConfig.Adam)
		if err != nil {
			return err
		}
	case "Adamax":
		err := json.Unmarshal(data["optimizer_config"], &c.OptimizerConfig.Adamax)
		if err != nil {
			return err
		}
	case "Nadam":
		err := json.Unmarshal(data["optimizer_config"], &c.OptimizerConfig.Nadam)
		if err != nil {
			return err
		}
	case "RMSprop":
		err := json.Unmarshal(data["optimizer_config"], &c.OptimizerConfig.RMSprop)
		if err != nil {
			return err
		}
	case "SGD":
		err := json.Unmarshal(data["optimizer_config"], &c.OptimizerConfig.SGD)
		if err != nil {
			return err
		}
	case "AdamW":
		err := json.Unmarshal(data["optimizer_config"], &c.OptimizerConfig.AdamW)
		if err != nil {
			return err
		}
	case "SGDW":
		err := json.Unmarshal(data["optimizer_config"], &c.OptimizerConfig.SGDW)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid optimizer")
	}

	err = json.Unmarshal(data["loss"], &c.Loss)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data["metrics"], &c.Metrics)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data["batch_size"], &c.BatchSize)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data["epochs"], &c.Epochs)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data["early_stop"], &c.EarlyStopping)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data["learning_rate_reduction"], &c.LearningRateReduction)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data["data_set"], &c.DataSet)
	if err != nil {
		return err
	}


	return nil
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
