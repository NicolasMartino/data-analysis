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
	"sync"
	"time"

	"github.com/NicolasMartino/data-analysis/src/models"
	"github.com/NicolasMartino/data-analysis/src/store"
	"github.com/NicolasMartino/data-analysis/src/utils"
)

var inputPath string
var outputPath string
var urlInfoCache store.SafeMap

func main() {
	//assign cache
	urlInfoCache = *store.NewSafeMap(250 * time.Millisecond)

	println("Simple data analysis program")
	println("The goal is to get the status of an url with http get request from a cmd line parameter or a csv file\n")

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
			result, err := handleGet(*getUrl)
			if err != nil {
				utils.Check(err)
				return
			}
			fmt.Printf("Reponded with status : %d\n\n", result.Status)
		}
	case fromCSVCmd.Name():
		if err := fromCSVCmd.Parse(os.Args[2:]); err == nil {
			//Create Dirs if  none exist
			inputPath, outputPath = utils.CreateDirs()
			csvSeparatorAsRune := []rune(*csvSeparator)[0]

			//process all files in input directory
			//each on a separate goroutine rail
			if *filename == "" {
				handleFromCSVs(csvSeparatorAsRune, inputUrlColumn)
			} else {
				handleFromCSV(csvSeparatorAsRune, *filename, *inputUrlColumn)
			}
		}
	default:
		printHelp(getCmd, fromCSVCmd)
	}

	fmt.Printf("Program is done working")
}

func handleFromCSVs(csvSeparatorAsRune rune, inputUrlColumn *int) {
	filesInput := utils.FindFiles(inputPath, ".csv")
	fmt.Printf("files input %v\n", filesInput)
	if len(filesInput) != 0 {
		var wg sync.WaitGroup

		for _, file := range filesInput {
			wg.Add(1)
			go func(file string) {
				defer wg.Done()
				handleFromCSV(csvSeparatorAsRune, file, *inputUrlColumn)
			}(file)
		}
		wg.Wait()
	}
}

// Thread safe cache results
func cacheUrlInfo(lineNumber int, givenUrl string, fetcher models.UrlInfoFetcher, writingChanel models.Channel) {
	urlInfoValue, ok := urlInfoCache.Load(givenUrl)

	if !ok {
		model, err := fetcher(givenUrl)
		if err != nil {
			utils.Check(err)
			return
		}

		urlInfoCache.Store(givenUrl, models.CacheUrlInfo{
			UrlInfo: model,
		})

		writingChanel.Values <- model
		return
	}

	fmt.Printf("Found cache value for url: %v \n", givenUrl)
	writingChanel.Values <- urlInfoValue.UrlInfo
}

// Request ressource form url
func handleGet(givenUrl string) (model models.Data, err error) {
	fmt.Printf("Called get command with url: %v\n", givenUrl)

	_, err = url.ParseRequestURI(givenUrl)
	if err != nil {
		utils.Check(err)
		return
	}

	resp, err := http.Get(givenUrl)
	if err != nil {
		utils.Check(err)
		return
	}

	return httpToDataMapper(resp), nil
}

func httpToDataMapper(resp *http.Response) (model models.Data) {

	body, err := io.ReadAll(resp.Body)
	utils.Check(err)

	model = models.Data{
		Status:     resp.StatusCode,
		RequestUrl: resp.Request.URL.String(),
		Body:       string(body),
	}

	return model
}

// Handle operations on csv files
func handleFromCSV(csvSeparatorAsRune rune, filename string, inputUrlColumn int) {

	// Create input file objects
	extension := "csv"
	var inputFilePath string
	fmt.Printf("filename:%v\n", filename)
	inputFilePath = filepath.Join(inputPath, fmt.Sprintf("%v.%v", filename, extension))

	inputFile, err := os.Open(inputFilePath)
	utils.CheckFatal(err)
	// remember to close the file at the end of the program
	defer inputFile.Close()

	inputCsvFile := models.InputCSVFile{
		FileReader:     inputFile,
		Filename:       filename,
		InputUrlColumn: inputUrlColumn,
		CsvSeparator:   csvSeparatorAsRune,
		FilePath:       inputFilePath,
	}

	//Create outputfileObjects
	outputFilePath := filepath.Join(outputPath, fmt.Sprintf("%s_results.csv", filename))
	outputFile, err := os.Create(outputFilePath)
	utils.Check(err)

	defer outputFile.Close()

	outputCsvFile := models.OutputCsvFile{
		FileWriter:   outputFile,
		CsvSeparator: csvSeparatorAsRune,
		Headers:      []string{"URL", "Status"},
		FilePath:     outputFilePath,
		LineWritor: func(input models.Data) []string {
			return []string{input.RequestUrl, fmt.Sprint(input.Status)}
		},
	}

	// print file infos
	postWrite, err := os.Stat(outputFilePath)
	utils.Check(err)
	fmt.Printf("Wrote %d bytes to file %s\n", postWrite.Size(), outputFilePath)

	//clean output directory
	deletedFileNames := utils.DeleteAllFilesFromDirectory(outputPath)

	if len(deletedFileNames) > 0 {
		fmt.Printf("Cleaned directory: %v\n", outputPath)
		fmt.Printf("Deleted %v files with names %+v\n", len(deletedFileNames), deletedFileNames)
	}

	fmt.Printf("Called from-csv command with params [%v, %v, %v]\n", inputCsvFile.Filename, inputCsvFile.InputUrlColumn, inputCsvFile.CsvSeparator)

	parseCSVLineByLine(inputCsvFile, outputCsvFile)
}

func parseCSVLineByLine(csvFile models.InputCSVFile, outputCSVFile models.OutputCsvFile) {
	//handle file writing
	writingChannel := models.Channel{
		Values: make(chan models.Data),
		Err:    make(chan error),
		Done:   make(chan bool),
	}

	go writeCSVFromChannel(&writingChannel, outputCSVFile)

	// read csv values using csv.Reader
	csvReader := csv.NewReader(csvFile.FileReader)
	csvReader.Comma = csvFile.CsvSeparator
	i := 0
	for {
		i++
		record, err := csvReader.Read()
		cleanedRecord := utils.CleanBom(record)
		if err == io.EOF {
			break
		}
		if err != nil {
			utils.CheckFatal(err)
		}
		// do something with read line
		if csvFile.InputUrlColumn > len(cleanedRecord)-1 {
			utils.CheckFatal(errors.New("could not find url column"))
		}

		cacheUrlInfo(i, cleanedRecord[csvFile.InputUrlColumn], handleGet, writingChannel)
	}

	writingChannel.Done <- true
}

//Create file write from channel to disk asynchronously
func writeCSVFromChannel(writeChannel *models.Channel, outputCSVFile models.OutputCsvFile) (err error) {

	writer := csv.NewWriter(outputCSVFile.FileWriter)
	defer writer.Flush()
	writer.Comma = outputCSVFile.CsvSeparator

	//Write headers
	record := outputCSVFile.Headers
	err = writer.Write(record)
	utils.Check(err)
	writer.Flush()

	for {
		select {
		case <-writeChannel.Done:
			fmt.Printf("Done with writing to file:%s\n", outputCSVFile.FilePath)
			return err
		case lineToWrite := <-writeChannel.Values:
			if lineToWrite.RequestUrl != "" {
				record := outputCSVFile.LineWritor(lineToWrite)
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
