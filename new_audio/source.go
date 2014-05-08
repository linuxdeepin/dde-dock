package main

type Source struct {
	index int32

	Name        string
	Description string
	Mute        bool   `access:"readwrite"`
	Volume      uint32 `access:"readwrite"`

	//ActivePort *SourcePortInfo
	Ports      []PortInfo
	ActivePort uint32

	Outputs []*SourceOutput
}

type SourceOutput struct {
}

func (*Source) SelectPort(port uint32) {
}
