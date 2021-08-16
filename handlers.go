package main

import (
	"codeconverter/CodeGenerator"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

func MakeModel (c echo.Context) error {
	project, err := CodeGenerator.BindProject(c.Request())
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	err = CodeGenerator.GenerateModel(project.Config, project.Content, false)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// return fIle.
	return c.Attachment("./model.py", "model.py")
}

// For testing Keras remote monitor
func TrainMonitor (c echo.Context) error {
	type Epoch struct {
		Acc float64	`json:"accuracy"`
		Epoch int	`json:"epoch"`
		Loss float64	`json:"loss"`
		LearningRate	float64	`json:"lr"`
		ValAcc	float64	`json:"val_accuracy"`
		ValLoss	float64	`json:"val_loss"`
	}

	e := new(Epoch)
	err := c.Bind(e)
	if err != nil {
		return err
	}

	fmt.Println(e)

	return c.NoContent(http.StatusOK)
}
