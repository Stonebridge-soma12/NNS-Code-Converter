package main

import (
	"codeconverter/CodeGenerator"
	"github.com/labstack/echo/v4"
	"net/http"
)

func MakeModel (c echo.Context) error {
	project, err := CodeGenerator.BindProject(c.Request())
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	err = CodeGenerator.GenerateModel(project.Config, project.Content)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// return fIle.
	return c.Attachment("./model.py", "model.py")
}
