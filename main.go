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
	e.Use(middleware.Logger())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	//e.GET("/getmodel", func(c echo.Context) error {
	//	return c.File("./model.py")
	//})

	e.POST("/make-python", MakeModel)
	e.POST("/publish/epoch/end", TrainMonitor)
	e.POST("/fit", Fit)

	e.Logger.Fatal(e.Start(":8080"))
}
