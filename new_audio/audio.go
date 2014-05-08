package main

type PortInfo struct {
	Name        string
	Description string
	Available   int32
}

type Audio struct {
	//Cards
	Sinks   []*Sink
	Sources []*Source

	SinkInputs    []*SinkInput
	SourceOutputs []*SourceOutput

	DefaultSink   *Sink
	DefaultSource *Sink
}
