package CodeGenerator

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGenerateModel(t *testing.T) {
	type args struct {
		config  Config
		content Content
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

		})
	}
}

func TestConv2D_ToCodest(t *testing.T) {
	filters := 16
	padding := "same"
	kernel := []int{16, 16}
	strides := []int{1, 1}

	conv2D := &Conv2D{
		&filters,
		kernel,
		strides,
		&padding,
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
	data := []byte(`{
  "config": {
    "optimizer": "adam",
    "learning_rate": 0.001,
    "loss": "sparse_categorical_crossentropy",
    "metrics": ["accuracy"],
    "batch_size": 32,
    "epochs": 10
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
          "shape": [28, 28, 1]
        }
      },
      {
        "category": "Layer",
        "type": "Conv2D",
        "name": "node_2fbbd8e5b0a5456faa2d47f7026b139f",
        "input": "node_1605430f35f94411aaf6b97eae005e19",
        "output": "node_39ce8c39bacb4fb392c2372fb81a0b7e",
        "param": {
          "filters": 32,
          "kernel_size": [16, 16],
          "padding": "same",
          "strides": [1, 1]
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
}
