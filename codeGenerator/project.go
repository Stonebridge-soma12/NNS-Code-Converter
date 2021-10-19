package codeGenerator

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
)

type Project struct {
	TrainId int64   `json:"train_id"`
	UserId  int64   `header:"user_id"`
	Config  Config  `json:"config"`
	DataSet DataSet `json:"data_set"`
	Content Content `json:"content"`
}

type Train struct {
	TrainId int64   `json:"train_id"`
	UserId  int64   `json:"user_id"`
	Config  Config  `json:"config"`
	DataSet DataSet `json:"data_set"`
}

const (
	importTf    = "import tensorflow as tf\n\n"
	importTfa   = "import tensorflow_addons as tfa\n\n"
	tf          = "tf"
	tfa         = "tfa"
	keras       = ".keras"
	layers      = ".layers"
	math        = ".math"
	createModel = "model = tf.keras.Model(inputs=%s, outputs=%s)\n\n"
)

const (
	ErrInvalidJsonfield = "unexpected end of JSON input"
)

func (p *Project) BindProject(r *http.Request) error {
	data := make(map[string]json.RawMessage)
	cc := make(map[string]json.RawMessage)

	// Binding request body
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		return err
	}

	p.UserId, _ = strconv.ParseInt(r.Header.Get("id"), 10, 64)

	// Unmarshalling Config.
	var config map[string]json.RawMessage
	err = json.Unmarshal(data["config"], &config)
	if err != nil {
		return fmt.Errorf("JSON Error : %s with field %s", err.Error(), "config")
	}

	err = p.Config.UnmarshalConfig(config)
	if err != nil {
		return err
	}

	err = p.DataSet.Bind(data["data_set"])
	if err != nil {
		return fmt.Errorf("JSON Error : %s with field %s", err.Error(), "data_set")
	}

	// Unmarshalling Content.
	err = json.Unmarshal(data["content"], &cc)
	if err != nil {
		return fmt.Errorf("JSON Error : %s with field %s", err.Error(), "content")
	}

	err = json.Unmarshal(data["train_id"], &p.TrainId)
	if err != nil {
		return fmt.Errorf("JSON Error : %s with field %s", err.Error(), "train_id")
	}

	err = p.Content.BindContent(cc)
	if err != nil {
		return err
	}

	return nil
}

func (p *Project) BindProjectForCode(r *http.Request) error {
	data := make(map[string]json.RawMessage)
	cc := make(map[string]json.RawMessage)

	// Binding request body
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		return err
	}

	// Unmarshalling Config.
	var config map[string]json.RawMessage
	err = json.Unmarshal(data["config"], &config)
	if err != nil {
		return fmt.Errorf("JSON Error : %s with field %s", err.Error(), "config")
	}

	err = p.Config.UnmarshalConfig(config)
	if err != nil {
		return err
	}

	// Unmarshalling Content.
	err = json.Unmarshal(data["content"], &cc)
	if err != nil {
		return fmt.Errorf("JSON Error : %s with field %s", err.Error(), "content")
	}

	err = p.Content.BindContent(cc)
	if err != nil {
		return err
	}

	return nil
}

func (p *Project) SaveModel() error {
	err := p.GenerateModel()
	if err != nil {
		return err
	}

	err = p.GenerateSaveModel()
	if err != nil {
		return err
	}

	cmd := exec.Command("python", "train.py")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return nil
	}
	fmt.Printf("Finished saving model %d", p.UserId)

	return nil
}

func (p *Project) GenerateModel() error {
	var codes []string
	codes = append(codes, importTf)
	codes = append(codes, importTfa)

	Layers, err := p.Content.GenLayers()
	if err != nil {
		return err
	}
	codes = append(codes, Layers...)

	Configs, err := p.Config.GenConfig()
	if err != nil {
		return err
	}
	codes = append(codes, Configs...)

	// create python file
	err = MakeTextFile(codes, "model.py")

	return nil
}

func (p *Project) GenerateSaveModel() error {
	var codes []string
	codes = append(codes, importTf)
	codes = append(codes, importTfa)
	codes = append(codes, "import model\n\n")

	// Python comment.
	saveCode := fmt.Sprintf("model.model.save('./%d/Model')", p.UserId)
	codes = append(codes, saveCode)

	// Generate train python file
	err := MakeTextFile(codes, "train.py")
	if err != nil {
		return err
	}

	return nil
}

func (p *Project) GetTrainBody() Train {
	trainInfo := Train{
		TrainId: p.TrainId,
		UserId:  p.UserId,
		DataSet: p.DataSet,
		Config:  p.Config,
	}

	return trainInfo
}
