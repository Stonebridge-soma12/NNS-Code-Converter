package codeGenerator

import "encoding/json"

type DataSet struct {
	TrainURI      *string        `json:"train_uri"`
	ValidURI      *string        `json:"validation_uri"`
	Shuffle       *bool          `json:"shuffle"`
	Label         *string        `json:"label"`
	Normalization *Normalization `json:"normalization"`
	Kind          *string        `json:"kind"`
}

type Normalization struct {
	Usage  *bool   `json:"usage"`
	Method *string `json:"method"`
}

func (d *DataSet) Bind(data json.RawMessage) error {
	return json.Unmarshal(data, d)
}
