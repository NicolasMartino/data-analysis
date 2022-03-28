package models

import (
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

type CsvFile struct {
	Filename       string
	InputUrlColumn int
	CsvSeparator   string
}
