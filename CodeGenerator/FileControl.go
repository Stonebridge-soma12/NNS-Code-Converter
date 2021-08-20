package CodeGenerator

import (
	"fmt"
	"os"
)

func MakeTextFile(data []string, fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	fileSize := 0
	for _, d := range data {
		n, err := file.Write([]byte(d))
		if err != nil {
			return err
		}
		fileSize += n
	}

	fmt.Printf("Make new text file %s with size %d bytes successful.\n", fileName, fileSize)

	return nil
}
