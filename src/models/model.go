package models

type Extractor func(channelToWriteTo Channel)
type UrlInfoFetcher func(url string) UrlInfo

type UrlInfo struct {
	Status     int
	RequestUrl string
	Body       string
}

type Channel struct {
	Values chan UrlInfo
	Err    chan error
	Done   chan bool
}
