package main

import (
	"codeconverter/CodeGenerator"
	"codeconverter/Config"
	"codeconverter/MessageQ"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"strconv"
)

func MakeModel(c echo.Context) error {
	var project CodeGenerator.Project
	err := project.BindProject(c.Request())
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

func Fit(c echo.Context) error {
	var project CodeGenerator.Project

	err := project.BindProject(c.Request())
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	err = project.SaveModel()
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	// Zip saved model
	targetBase := fmt.Sprintf("./%d/", project.UserId)
	files, err := CodeGenerator.GetFileLists(targetBase + "Model")
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	err = CodeGenerator.Zip(targetBase + "Model.zip", files)
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

	conn, err := MessageQ.CreateConnection(cfg.Account, cfg.Pw, cfg.Host, cfg.VHost)
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

func GetSavedModel(c echo.Context) error {
	userId, _ := strconv.ParseInt(c.Request().Header.Get("id"), 10, 64)
	target := fmt.Sprintf("./%d/", userId)
	return c.File(target + "Model.zip")
}