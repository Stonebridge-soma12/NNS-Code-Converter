package CodeGenerator

import (
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
	dense := &Dense {
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