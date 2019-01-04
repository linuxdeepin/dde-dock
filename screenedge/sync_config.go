package screenedge

import (
	"encoding/json"
)

type syncConfig struct {
	m *Manager
}

type syncData struct {
	Version     string `json:"version"` // such as "1.0.0"
	LeftBottom  string `json:"left_bottom"`
	LeftTop     string `json:"left_top"`
	RightBottom string `json:"right_bottom"`
	RightTop    string `json:"right_top"`
}

const (
	syncDataVersion = "1.0"
)

func (sc *syncConfig) Get() (interface{}, error) {
	return &syncData{
		Version:     syncDataVersion,
		LeftBottom:  sc.m.settings.GetEdgeAction(BottomLeft),
		LeftTop:     sc.m.settings.GetEdgeAction(TopLeft),
		RightBottom: sc.m.settings.GetEdgeAction(BottomRight),
		RightTop:    sc.m.settings.GetEdgeAction(TopRight),
	}, nil
}

func (sc *syncConfig) Set(data []byte) error {
	var info syncData
	err := json.Unmarshal(data, &info)
	if err != nil {
		return err
	}
	sc.m.settings.SetEdgeAction(BottomLeft, info.LeftBottom)
	sc.m.settings.SetEdgeAction(TopLeft, info.LeftTop)
	sc.m.settings.SetEdgeAction(BottomRight, info.RightBottom)
	sc.m.settings.SetEdgeAction(TopRight, info.RightTop)
	return nil
}
