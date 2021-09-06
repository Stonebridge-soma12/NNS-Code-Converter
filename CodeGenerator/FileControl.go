package CodeGenerator

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
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

// Get list of files in directory
func GetFileLists(target string) ([]string, error) {
	var files []string

	err := filepath.Walk(target, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			fmt.Println(path)
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}

// Reference of Zip and AddFileToZip
// https://www.python2.net/questions-62657.htm
func Zip(filename string, files []string) error {
	newZipFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	err = os.Chmod(filename, 755)

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()
	// Add files to zip
	for _, file := range files {
		if err = AddFileToZip(zipWriter, file); err != nil {
			return err
		}
	}
	return nil
}

func AddFileToZip(zipWriter *zip.Writer, filename string) error {
	fileToZip, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileToZip.Close()
	// Get the file information
	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	// Using FileInfoHeader() above only uses the basename of the file. If we want
	// to preserve the folder structure we can overwrite this with the full path.
	header.Name = filename
	// Change to deflate to gain better compression
	// see http://golang.org/pkg/archive/zip/#pkg-constants
	header.Method = zip.Deflate
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, fileToZip)
	return err
}