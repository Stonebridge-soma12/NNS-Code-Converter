package CodeGenerator

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
	res, err := conv2D.ToCode()
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

	res, err := dense.ToCode()
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
    "loss": "sparse_categorical_crossentropy",
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
  "content": {
    "output": "node_96afcbc0a4ba4ed9b02b579068f166f0",
    "input": "node_1605430f35f94411aaf6b97eae005e19",
    "layers": [
      {
        "category": "Layer",
        "type": "Input",
        "name": "node_1605430f35f94411aaf6b97eae005e19",
        "input": null,
        "output": "node_2fbbd8e5b0a5456faa2d47f7026b139f",
        "param": {
          "shape": [
            28,
            28,
            1
          ]
        }
      },
      {
        "category": "Layer",
        "type": "Conv2D",
        "name": "node_2fbbd8e5b0a5456faa2d47f7026b139f",
        "input": "node_1605430f35f94411aaf6b97eae005e19",
        "output": "node_39ce8c39bacb4fb392c2372fb81a0b7e",
        "param": {
          "filters": 16,
          "kernel_size": [
            16,
            16
          ],
          "padding": "same",
          "strides": [
            1,
            1
          ]
        }
      },
      {
        "category": "Layer",
        "type": "Dropout",
        "name": "node_2c8a6d78d0204888942f16317f2a079f",
        "input": "node_39ce8c39bacb4fb392c2372fb81a0b7e",
        "output": "node_71914b8774b64700b38dc3e8e7a62caa",
        "param": {
          "rate": 0.5
        }
      },
      {
        "category": "Layer",
        "type": "Activation",
        "name": "node_39ce8c39bacb4fb392c2372fb81a0b7e",
        "input": "node_2fbbd8e5b0a5456faa2d47f7026b139f",
        "output": "node_2c8a6d78d0204888942f16317f2a079f",
        "param": {
          "activation": "relu"
        }
      },
      {
        "category": "Layer",
        "type": "Flatten",
        "name": "node_71914b8774b64700b38dc3e8e7a62caa",
        "input": "node_2c8a6d78d0204888942f16317f2a079f",
        "output": "node_020cdce94de241ac9556bb0b0022c1f2",
        "param": {}
      },
      {
        "category": "Layer",
        "type": "Dense",
        "name": "node_020cdce94de241ac9556bb0b0022c1f2",
        "input": "node_71914b8774b64700b38dc3e8e7a62caa",
        "output": "node_96afcbc0a4ba4ed9b02b579068f166f0",
        "param": {
          "units": 10
        }
      },
      {
        "category": "Layer",
        "type": "Activation",
        "name": "node_96afcbc0a4ba4ed9b02b579068f166f0",
        "input": "node_020cdce94de241ac9556bb0b0022c1f2",
        "output": null,
        "param": {
          "activation": "softmax"
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
	files, err := GetFileLists("../MNIST")
	if err != nil {
		t.Error(err)
	}

	err = Zip("model.zip", files)
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