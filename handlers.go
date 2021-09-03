package main

import (
	"codeconverter/CodeGenerator"
	"codeconverter/Config"
	"codeconverter/MessageQ"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/random"
	"github.com/streadway/amqp"
	"net/http"
	"os"
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
	body := project.GetTrainBody()
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	cfg, err := Config.GetConfig()
	if err != nil {
		return err
	}

	conn, err := MessageQ.CreateConnection(cfg.Account, cfg.Pw, cfg.Host)
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// define mq
	q, err := ch.QueueDeclare("Reply", false, false, false, false, nil)
	if err != nil {
		return err
	}

	// TODO: requestId (CorrelationId)를 uuid 사용하도록 수정해야함.
	requestId := random.String(32)
	publish := amqp.Publishing{
		ContentType: "application/json",
		CorrelationId: requestId,
		Body:jsonBody,
		ReplyTo: q.Name,
	}
	err = ch.Publish("", "Request", false, false, publish)
	if err != nil {
		return err
	}

	return nil
}

func GetSavedModel(c echo.Context) error {
	userId := c.Request().Header.Get("id")
	target := fmt.Sprintf("./%s/", userId)
	return c.File(target + "Model.zip")
}