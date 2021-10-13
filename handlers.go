package main

import (
	"codeconverter/codeGenerator"
	"codeconverter/config"
	"codeconverter/messageQ"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"strconv"
)

func MakeModelHandler(c echo.Context) error {
	var project codeGenerator.Project
	err := project.BindProjectForCode(c.Request())
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	err = project.GenerateModel()
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
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

func FitHandler(c echo.Context) error {
	var project codeGenerator.Project

	err := project.BindProject(c.Request())
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	err = project.SaveModel()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// Zip saved model
	targetBase := fmt.Sprintf("./%d/", project.UserId)
	files, err := codeGenerator.GetFileLists(targetBase + "Model")
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	err = codeGenerator.Zip(targetBase + "Model.zip", files)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// Request to GPU server
	body:= project.GetTrainBody()

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	cfg, err := Config.GetConfig()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	conn, err := messageQ.CreateConnection(cfg.Account, cfg.Pw, cfg.Host, cfg.VHost)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	defer conn.Close()

	err = conn.Publish(jsonBody)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func GetSavedModelHandler(c echo.Context) error {
	userId, _ := strconv.ParseInt(c.Request().Header.Get("id"), 10, 64)
	target := fmt.Sprintf("./%d/", userId)
	return c.File(target + "Model.zip")
}