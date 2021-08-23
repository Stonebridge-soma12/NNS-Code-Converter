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

//// Get file list of saved model directory
//func GetModel(dir string) error {
//	var fileList []string
//
//	dirs, err := os.ReadDir(dir)
//	if err != nil {
//		return err
//	}
//
//	for _, dir := range dirs {
//		if dir.Type() == fs.ModeDir {
//
//		}
//		fmt.Println(dir.Type())
//	}
//
//	return nil
//}
