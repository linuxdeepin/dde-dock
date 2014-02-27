package main

import "encoding/json"
import "fmt"
import "os"
import "io/ioutil"

type DisplayConfiguration struct {
	DisplayMode int16
	Outputs     []OutputConfiguration
}

var _ConfigPath = os.Getenv("HOME") + "/.config/deepin_monitors.json"

type OutputConfiguration struct {
	Name string

	Width, Height uint16
	RefreshRate   float64

	X, Y     int16
	Primary  bool
	Enabled  bool
	Rotation uint16
	Reflect  uint16
}

func LoadDisplayConfiguration(dpy *Display) DisplayConfiguration {
	f, err := os.Open(_ConfigPath)
	if err != nil {
		return generateDefaultConfig(dpy)
	}
	data, err := ioutil.ReadAll(f)
	var config DisplayConfiguration
	err = json.Unmarshal(data, &config)
	if err != nil {
		fmt.Println("Failed load displayConfiguration:", err, "Data:", string(data))
		return generateDefaultConfig(dpy)
	}
	if config.hasError(dpy) {
		return generateDefaultConfig(dpy)
	}
	return config
}

func (d DisplayConfiguration) save() {
	bytes, err := json.Marshal(d)
	if err != nil {
		panic("marshal display configuration failed:" + err.Error())
	}
	f, err := os.Create(_ConfigPath)
	if err != nil {
		fmt.Println("Couldn't save display configuration:", err)
		return
	}
	defer f.Close()
	f.Write(bytes)
}

func (d DisplayConfiguration) hasError(dpy *Display) bool {
	// check whether the display mode is in range
	if int(d.DisplayMode) > len(dpy.Outputs) || d.DisplayMode < DisplayModeMirrors {
		return true
	}
	if len(d.Outputs) != len(dpy.Outputs) {
		return true
	}
	for _, cfg := range d.Outputs {
		var currentOP *Output
		for _, op := range dpy.Outputs {
			if op.Name == cfg.Name {
				currentOP = op
				break
			}
		}
		if currentOP != nil {
			// TODO: check whether outputs configuration is valid
		} else {
			return true
		}
	}
	return false
}

func generateDefaultConfig(dpy *Display) DisplayConfiguration {
	d := DisplayConfiguration{}
	d.DisplayMode = dpy.DisplayMode
	d.Outputs = make([]OutputConfiguration, len(dpy.Outputs))
	for i, op := range dpy.Outputs {
		d.Outputs[i] = OutputConfiguration{
			X:           0,
			Y:           0,
			Width:       op.Mode.Width,
			Height:      op.Mode.Height,
			Name:        op.Name,
			Primary:     dpy.PrimaryOutput == op,
			Enabled:     true,
			Rotation:    op.Rotation,
			Reflect:     op.Reflect,
			RefreshRate: op.Mode.Rate,
		}
	}
	return d
}
func GenerateCurrentConfig(dpy *Display) DisplayConfiguration {
	d := DisplayConfiguration{}
	d.DisplayMode = dpy.DisplayMode
	d.Outputs = make([]OutputConfiguration, len(dpy.Outputs))
	for i, op := range dpy.Outputs {
		rect := op.pendingAllocation()
		d.Outputs[i] = OutputConfiguration{
			X:           rect.X,
			Y:           rect.Y,
			Width:       rect.Width,
			Height:      rect.Height,
			Name:        op.Name,
			Primary:     dpy.PrimaryOutput == op,
			Enabled:     op.Opened,
			Rotation:    op.Rotation,
			Reflect:     op.Reflect,
			RefreshRate: op.Mode.Rate,
		}
	}
	return d
}
