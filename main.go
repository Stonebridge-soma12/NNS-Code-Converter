package main

import (
	// "encoding/json"
	// "io/ioutil"
	"github.com/labstack/echo/v4/middleware"

	"github.com/labstack/echo/v4"
	"net/http"
)



func main() {
	e := echo.New()

	e.Use(middleware.CORS())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.POST("/make-python", MakeModel)

	e.Logger.Fatal(e.Start(":80"))
}
