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
	e.GET("/model", GetSavedModel)

	e.POST("/make-python", MakeModel)
	e.POST("/publish/epoch/end", TrainMonitor)
	e.POST("/fit", Fit)

	//e.POST("/api/trainedModel", func(c echo.Context) error {
	//	form, err := c.MultipartForm()
	//	if err != nil {
	//		fmt.Println(err)
	//		return err
	//	}
	//	files := form.File["files"]
	//
	//	for _, file := range files {
	//		src, err := file.Open()
	//		if err != nil {
	//			return err
	//		}
	//		defer src.Close()
	//
	//		dst, err := os.Create(file.Filename)
	//		if err != nil {
	//			return err
	//		}
	//		defer dst.Close()
	//		if _, err = io.Copy(dst, src); err != nil {
	//			return err
	//		}
	//	}
	//
	//	return nil
	//})

	e.Logger.Fatal(e.Start(":8081"))
}
