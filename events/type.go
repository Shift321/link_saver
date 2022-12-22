package events

type Fetcher interface {
	Fetch(limit int) ([]Event, error)
}

type Processor {
	Process(e Event) error
	
}

type Type int

const {
	Unknown Type = iota
	Message 
}

type Event struct {
	Type Type
	Text string
}