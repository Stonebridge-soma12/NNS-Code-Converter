package CodeGenerator

import "fmt"

const (
	earlyStop   = "early_stop = tf.keras.callbacks.EarlyStopping(monitor=%s, patience=%d)\n"
	lrReduction = "learning_rate_reduction = ReduceLROnPlateau(monitor=%s, patience=%d, verbose=1, factor=%g, min_lr=%g)\n"
)

type Config struct {
	Optimizer             string      `json:"optimizer"`
	LearningRate          float64     `json:"learning_rate"`
	Loss                  string      `json:"loss"`
	Metrics               []string    `json:"metrics"`
	BatchSize             int         `json:"batch_size"`
	Epochs                int         `json:"epochs"`
	Output                string      `json:"output"`
	EarlyStopping         EarlyStop   `json:"early_stop"`
	LearningRateReduction LrReduction `json:"learning_rate_reduction"`
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
	if err != "" {
		return "", fmt.Errorf(err)
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
	if err != "" {
		return "", fmt.Errorf(err)
	}

	return fmt.Sprintf(lrReduction, *l.Monitor, *l.Patience, *l.Factor, *l.MinLr), nil
}
