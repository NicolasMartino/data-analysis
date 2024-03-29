package models

import (
	"io"
)

type LineExtractor func(channelToWriteTo Channel)
type UrlInfoFetcher func(url string) Data
type LineWritor func(input Data) []string

type Data struct {
	Status     int
	RequestUrl string
	Body       string
}

type CacheUrlInfo struct {
	UrlInfo Data
}

type Channel struct {
	Values chan Data
	Err    chan error
	Done   chan bool
}

type InputCSVFile struct {
	Filename       string
	FileReader     io.Reader
	InputUrlColumn int
	CsvSeparator   rune
	FilePath       string
}

type OutputCsvFile struct {
	FileWriter   io.Writer
	FilePath     string
	CsvSeparator rune
	Headers      []string
	LineWritor   LineWritor
}
