package models

type Extractor func(channelToWriteTo Channel)
type Transformer func(url string) GetResult

type GetResult struct {
	Status     string
	RequestUrl string
	Body       string
}

type Channel struct {
	Values chan string
	Err    chan error
	Done   chan bool
}
