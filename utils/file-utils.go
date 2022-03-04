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

// Create two directories
func CreateDirs() (inputPath string, outputPath string) {
	wd, err := os.Getwd()
	Check(err)
	inputPath = filepath.Join(wd, "input")
	outputPath = filepath.Join(wd, "output")

	err = os.MkdirAll(inputPath, os.ModePerm)
	Check(err)
	err = os.MkdirAll(outputPath, os.ModePerm)
	Check(err)

	return
}

func CleanBom(s []string) []string {
	if len(s) > 0 {
		byteOrderMarkAsString := string('\uFEFF')
		s[0] = strings.TrimPrefix(s[0], byteOrderMarkAsString)
	}
	return s
}
