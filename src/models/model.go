package models

import (
	"io"
	"time"
)

type Extractor func(channelToWriteTo Channel)
type UrlInfoFetcher func(url string) UrlInfo

type UrlInfo struct {
	Status     int
	RequestUrl string
	Body       string
}

type CacheUrlInfo struct {
	UrlInfo    UrlInfo
	LastUpdate time.Time
}

type Channel struct {
	Values chan UrlInfo
	Err    chan error
	Done   chan bool
}

type InputCSVFile struct {
	Filename       string
	FileReader     io.Reader
	InputUrlColumn int
	CsvSeparator   string
	FilePath       string
}

type OutputCsvFile struct {
	FileWriter   io.Writer
	FilePath     string
	CsvSeparator string
	Headers      []string
}
