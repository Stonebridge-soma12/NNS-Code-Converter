package main

import (
	"bytes"
	"codeconverter/CodeGenerator"
	"codeconverter/Config"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"os/exec"
)

func MakeModel(c echo.Context) error {
	var project CodeGenerator.Project
	err := project.BindProject(c.Request())
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	err = project.GenerateModel()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// Attach result python file
	err = c.Attachment("./model.py", "model.py")
	if err != nil {
		return err
	}

	err = os.Remove("./model.py")
	if err != nil {
		return err
	}

	return nil
}

// For testing Keras remote monitor
func TrainMonitor(c echo.Context) error {
	type Epoch struct {
		Acc          float64 `json:"accuracy"`
		Epoch        int     `json:"epoch"`
		Loss         float64 `json:"loss"`
		LearningRate float64 `json:"lr"`
		ValAcc       float64 `json:"val_accuracy"`
		ValLoss      float64 `json:"val_loss"`
	}

	e := new(Epoch)
	err := c.Bind(e)
	if err != nil {
		return err
	}

	fmt.Println(e)

	return c.NoContent(http.StatusOK)
}


func Fit(c echo.Context) error {
	var project CodeGenerator.Project

	err := project.BindProject(c.Request())
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	err = project.GenerateModel()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	err = project.SaveModel()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	cmd := exec.Command("/usr/bin/python", "train.py")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
	}

	// Zip saved model
	targetBase := fmt.Sprintf("./%s/", project.UserId)
	files, err := CodeGenerator.GetFileLists(targetBase + "Model")
	if err != nil {
		return err
	}

	err = CodeGenerator.Zip(targetBase + "Model.zip", files)
	if err != nil {
		return err
	}

	// Request to GPU server
	byteConfig, err := json.Marshal(project.Config)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(byteConfig)

	cfg, err := Config.GetConfig()
	if err != nil {
		return err
	}
	var URL string
	URL = cfg.BaseURL + cfg.Port

	req, err := http.NewRequest("POST", URL + "/run", buf)
	if err != nil {
		return err
	}
	req.Header.Add("id", project.UserId)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode > 400 {
		return c.NoContent(res.StatusCode)
	}

	return nil
}

func GetSavedModel(c echo.Context) error {
	userId := c.Request().Header.Get("id")
	target := fmt.Sprintf("./%s/", userId)
	return c.File(target + "Model.zip")
}