package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/NicolasMartino/data-analysis/src/models"
	"github.com/NicolasMartino/data-analysis/src/utils"
)

var inputPath string
var outputPath string

func main() {
	//Create Dirs if  none exist
	inputPath, outputPath = utils.CreateDirs()

	println("Simple data analysis program")
	println("The goal is to get the status of an url with http get request from a cmd line parameter or a csv file\n")

	//clean output directory
	deletedFileNames := utils.DeleteAllFilesFromDirectory(outputPath)

	if len(deletedFileNames) > 0 {
		fmt.Printf("Cleaned directory: %v\n", outputPath)
		fmt.Printf("Deleted %v files with names %+v", len(deletedFileNames), deletedFileNames)
	}

	//CMDS
	getCmd := *flag.NewFlagSet("get", flag.ExitOnError)
	getUrl := getCmd.String("url", "", "Find status for the given url")

	fromCSVCmd := *flag.NewFlagSet("from-csv", flag.ExitOnError)

	fromCSVCmd.Usage = func() {
		fmt.Printf("From CSV command usage\nInput files should be placed in the input folder: %s\nOutput files will be placed in the output folder: %s\n", inputPath, outputPath)
		fromCSVCmd.PrintDefaults()
		fmt.Printf("example for one file: from-csv -filename=test_file -csv-separator=|\n")
		fmt.Printf("example for all files input: from-csv -csv-separator=|\n")
	}

	filename := fromCSVCmd.String("filename", "", "Csv input filename")
	inputUrlColumn := fromCSVCmd.Int("url-column", 0, "Input column where to find the urls")
	csvSeparator := fromCSVCmd.String("csv-separator", ",", "Csv separator used")

	fmt.Println(os.Args)
	if len(os.Args) < 2 {
		printHelp(getCmd, fromCSVCmd)
		os.Exit(1)
	}

	fmt.Printf("Operation : %v\n", os.Args[0:])

	switch os.Args[1] {
	case getCmd.Name():
		if err := getCmd.Parse(os.Args[2:]); err == nil {
			result := handleGet(*getUrl)
			fmt.Printf("Reponded with status : %d\n\n", result.Status)
		}
	case fromCSVCmd.Name():
		if err := fromCSVCmd.Parse(os.Args[2:]); err == nil {
			handleFromCSV(*filename, *inputUrlColumn, *csvSeparator, handleGet)
		}
	default:
		printHelp(getCmd, fromCSVCmd)
	}

	fmt.Printf("Program has finished with no unrecoverable errors")
}

func handleGet(givenUrl string) (model models.UrlInfo) {
	fmt.Printf("Called get command with url: %v\n", givenUrl)

	_, err := url.ParseRequestURI(givenUrl)
	if err != nil {
		utils.Check(err)
		return
	}

	resp, err := http.Get(givenUrl)
	utils.Check(err)

	body, err := io.ReadAll(resp.Body)
	utils.Check(err)

	model = models.UrlInfo{
		Status:     resp.StatusCode,
		RequestUrl: resp.Request.URL.String(),
		Body:       string(body),
	}
	return
}

//TODO should split and readfromfile method should take a os.file as input
func handleFromCSV(filename string, inputUrlColumn int, csvSeparator string, fetcher models.UrlInfoFetcher) {
	//handle file writing
	writingChannel := models.Channel{
		Values: make(chan models.UrlInfo),
		Err:    make(chan error),
		Done:   make(chan bool),
	}

	go writeCsvFromChannel(&writingChannel)

	fmt.Printf("Called from-csv command with params [%v, %v, %v]\n", filename, inputUrlColumn, csvSeparator)

	extension := "csv"
	filepath := filepath.Join(inputPath, fmt.Sprintf("%v.%v", filename, extension))

	f, err := os.Open(filepath)
	utils.CheckFatal(err)
	// remember to close the file at the end of the program
	defer f.Close()

	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	csvReader.Comma = ([]rune(csvSeparator))[0]

	for {
		record, err := csvReader.Read()
		cleanedRecord := utils.CleanBom(record)
		if err == io.EOF {
			break
		}
		if err != nil {
			utils.CheckFatal(err)
		}
		// do something with read line
		if inputUrlColumn > len(cleanedRecord)-1 {
			//utils.Check(errors.New("could not find url column"))
			utils.CheckFatal(errors.New("could not find url column"))
		}
		result := fetcher(cleanedRecord[inputUrlColumn])
		writingChannel.Values <- result
	}
}

//Create file write from channel to disk asynchronously
func writeCsvFromChannel(writeChannel *models.Channel) (err error) {
	path := filepath.Join(outputPath, "results.csv")
	file, err := os.Create(path)
	utils.Check(err)

	defer file.Close()

	writer := csv.NewWriter(file)
	//TODO reuse csv separator
	separatorAsRune := []rune(";")
	writer.Comma = separatorAsRune[0]

	defer writer.Flush()

	//Write headers
	record := []string{"URL", "Status"}
	err = writer.Write(record)
	utils.Check(err)
	writer.Flush()

	for {
		select {
		case <-writeChannel.Done:
			postWrite, err := file.Stat()
			utils.Check(err)
			fmt.Printf("wrote %d bytes to file %s\n", postWrite.Size(), file.Name())
			return err
		case lineToWrite := <-writeChannel.Values:
			if lineToWrite.RequestUrl != "" {
				record := []string{lineToWrite.RequestUrl, fmt.Sprint(lineToWrite.Status)}
				err := writer.Write(record)
				utils.Check(err)
				writer.Flush()
			}
		}
	}
}

func printHelp(getCmd flag.FlagSet, fromCSVCmd flag.FlagSet) {
	fmt.Printf("Trying to help...\n")
	fmt.Printf("This program can fetch one url or more if they are present in a csv file in the folder input\n")
	fmt.Printf("This program can print the results to console or place them in a csv file in the folder output\n\n")
	fmt.Printf("List of commands (cmd)\n")
	fmt.Printf("Cmd get: fetches url status from value as a cmd line parameter\n")
	getCmd.PrintDefaults()
	fmt.Printf("example: get -url=https://google.com\n")
	fmt.Printf("Cmd from-csv: fetches url status from csv values\n\n")
	fromCSVCmd.Usage()
}
