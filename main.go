package main

import (
	"github.com/labstack/echo/v4"
	// "encoding/json"
	// "io/ioutil"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

func main() {
	e := echo.New()

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/api/model", GetSavedModelHandler)

	e.POST("/api/python", MakeModelHandler)
	e.POST("/api/fit", FitHandler)

	e.Logger.Fatal(e.Start(":8081"))
}
