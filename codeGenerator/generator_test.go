package codeGenerator

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestConv2D_ToCode(t *testing.T) {
	filters := 16
	padding := "same"
	kernel := []int{16, 16}
	strides := []int{1, 1}

	conv2D := &Conv2D{
		&filters,
		kernel,
		&padding,
		strides,
	}
	res, err := conv2D.GetCode()
	if err != nil {
		fmt.Println(res)
		t.Error(err)
	} else {
		fmt.Println(res)
	}
}

func TestDense_ToCode(t *testing.T) {
	units := 10
	dense := &Dense{
		Units: &units,
	}

	res, err := dense.GetCode()
	if err != nil {
		fmt.Println(res)
		t.Error(err)
	} else {
		fmt.Println(res)
	}
}

func TestUnmarshalParam(t *testing.T) {
	data := []byte(`
{
  "config": {
    "optimizer_name": "Adam",
    "optimizer_config": {
      "learning_rate": 0.001,
      "beta_1": 0.9,
      "beta_2": 0.999,
      "epsilon": 1e-07,
      "amsgrad": false
    },
    "loss": "binary_crossentropy",
    "metrics": [
      "accuracy"
    ],
    "batch_size": 32,
    "epochs": 10,
    "early_stop": {
      "usage": true,
      "monitor": "loss",
      "patience": 2
    },
    "learning_rate_reduction": {
      "usage": true,
      "monitor": "val_accuracy",
      "patience": 2,
      "factor": 0.25,
      "min_lr": 0.0000003
    }
  },
  "data_set": {
    "train_uri": "https://s3.ap-northeast-2.amazonaws.com/image.nns/test.csv",
    "valid_uri": "",
    "shuffle": false,
    "label": "blue_win",
    "normalization": {
      "usage": true,
      "method": "MinMax"
    }
  },
  "content": {
    "output": "Activation_2",
    "input": "Input_1",
    "layers": [
      {
        "category": "Layer",
        "type": "Input",
        "name": "Input_1",
        "input": null,
        "output": [
          "Dense_1"
        ],
        "param": {
          "shape": [
            1,
            58
          ]
        }
      },
      {
        "category": "Layer",
        "type": "Dense",
        "name": "Dense_1",
        "input": [
          "Input_1"
        ],
        "output": [
          "Activation_1"
        ],
        "param": {
          "units": 256
        }
      },
      {
        "category": "Layer",
        "type": "Activation",
        "name": "Activation_1",
        "input": [
          "Dense_1"
        ],
        "output": [
          "Dense_2"
        ],
        "param": {
          "activation": "relu"
        }
      },
      {
        "category": "Layer",
        "type": "Dense",
        "name": "Dense_2",
        "input": [
          "Activation_1"
        ],
        "output": [
          "Activation_2"
        ],
        "param": {
          "units": 1
        }
      },
      {
        "category": "Layer",
        "type": "Activation",
        "name": "Activation_2",
        "input": [
          "Dense_2"
        ],
        "output": null,
        "param": {
          "activation": "sigmoid"
        }
      }
    ]
  }
}
`)

	project := Project{}
	if err := json.Unmarshal(data, &project); err != nil {
		t.Error(err)
	}
	fmt.Print(project)
}


func TestServingDir(t *testing.T) {
	dirs, err := os.ReadDir("../MNIST")
	if err != nil {
		t.Error(err)
	}

	for _, dir := range dirs {
		fmt.Println(dir.Type())
	}
}

func TestZip(t *testing.T) {
	files, err := GetFileLists("../Model")
	if err != nil {
		t.Error(err)
	}

	err = Zip("Model.zip", files)
	if err != nil {
		t.Error(err)
	}
}

func TestGetFileLists(t *testing.T) {
	files, err := GetFileLists("../MNIST")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(files)
}