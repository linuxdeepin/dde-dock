package main

type Sink struct {
	index uint32
	valid bool

	Name        string
	Description string
	Mute        bool    `access:"readwrite",seter:"SetMue",getter:"GetMute"`
	Volume      uint32  `access:"readwrite"`
	Balance     float64 `access:"readwrite"`

	Ports      []PortInfo
	ActivePort int32

	channelMap []int32
	driver     string
	card       int32
	Inputs     []*SinkInput
}
type SinkInput struct {
}

func (s *Sink) hookBeforeReadProperty(name string) {
	if !s.valid {
		s.forceUpdate()
	}
}

func (sink *Sink) SelectPort(portnum int32) {
}

func (sink *Sink) SetVolume(volume uint32) {
}

func (sink *Sink) SetMute(mute bool) {
}

func (sink *Sink) SetBalance(balance float64) {
}
