package utils

import (
	"os"
	"path/filepath"
	"strings"
)

// Delete all files in directory (folder) and return the successfully deleted files names
func DeleteAllFilesFromDirectory(path string) []string {
	d, err := os.Open(path)
	Check(err)

	defer d.Close()

	names, err := d.Readdirnames(-1)
	Check(err)

	for _, name := range names {
		os.RemoveAll(filepath.Join(path, name))
		Check(err)
	}

	return names
}

// Create two directories in the working directory from two names
func CreateDirs() (inputFileName string, outputFileName string) {
	wd, err := os.Getwd()
	Check(err)
	inputFileName = filepath.Join(wd, "input")
	outputFileName = filepath.Join(wd, "output")

	err = os.MkdirAll(inputFileName, os.ModePerm)
	Check(err)
	err = os.MkdirAll(outputFileName, os.ModePerm)
	Check(err)

	return
}

// Clean the thingy ms words adds sometimes to annoy me in csv files
func CleanBom(s []string) []string {
	if len(s) > 0 {
		byteOrderMarkAsString := string('\uFEFF')
		s[0] = strings.TrimPrefix(s[0], byteOrderMarkAsString)
	}
	return s
}
