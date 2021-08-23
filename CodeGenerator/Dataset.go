package CodeGenerator

type DataSet struct {
	TrainURI *string `json:"train_uri"`
	ValidURI *string `json:"validation_uri"`
	Shuffle  *bool   `json:"shuffle"`
	Label    *string `json:"label"`
}

type Normalization struct {
	Usage  *bool   `json:"usage"`
	Method *string `json:"method"`
}
