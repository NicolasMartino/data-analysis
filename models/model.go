package models

type Extractor func(channelToWriteTo Channel)
type Transformer func(url string) UrlInfo

type UrlInfo struct {
	Status     string
	RequestUrl string
	Body       string
}

type Channel struct {
	Values chan UrlInfo
	Err    chan error
	Done   chan bool
}
