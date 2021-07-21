package main

import (
	"codeconverter/CodeGenerator"
	"github.com/labstack/echo/v4"
)

func MakeModel (c echo.Context) error {
	project := new(CodeGenerator.Project)

	// binding JSON
	err := c.Bind(project)
	if err != nil {
		panic(err)
	}

	CodeGenerator.GenerateModel(project.Config, project.Content)

	// return fIle.
	return c.Attachment("./model.py", "model.py")
}
