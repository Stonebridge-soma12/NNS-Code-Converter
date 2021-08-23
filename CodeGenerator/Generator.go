package CodeGenerator

import (
	"encoding/json"
	"net/http"
)

type Project struct {
	Config  Config  `json:"config"`
	Content Content `json:"content"`
}

const (
	importTf      = "import tensorflow as tf\n\n"
	importTfa     = "import tensorflow_addons as tfa\n\n"
	tf            = "tf"
	tfa           = "tfa"
	keras         = ".keras"
	layers        = ".layers"
	createModel   = "model = tf.keras.Model(inputs=%s, outputs=%s)\n\n"
)

func (p *Project) SaveModel() error {
	var codes []string
	codes = append(codes, importTf)
	codes = append(codes, importTfa)
	codes = append(codes, "import model\n\n")

	// Python comment.
	codes = append(codes, "model.model.save('Model')")

	// Generate train python file
	err := MakeTextFile(codes, "train.py")
	if err != nil {
		return err
	}

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


func (p *Project) BindProject(r *http.Request) error {
	project := new(Project)
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
		return err
	}

	err = project.Config.UnmarshalConfig(config)
	if err != nil {
		return err
	}

	// Unmarshalling Content.
	err = json.Unmarshal(data["content"], &cc)
	if err != nil {
		return err
	}

	err = project.Content.BindContent(cc)
	if err != nil {
		return err
	}

	return nil
}
