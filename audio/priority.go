package audio

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strings"

	"pkg.deepin.io/lib/pulse"
	"pkg.deepin.io/lib/xdg/basedir"
)

const (
	PortTypeBluetooth = iota
	PortTypeHeadset
	PortTypeSpeaker
	PortTypeHdmi
)

type PortToken struct {
	CardName string
	PortName string
}

type Priorities struct {
	OutputTypePriority     []int
	OutputInstancePriority []*PortToken
	InputTypePriority      []int
	InputInstancePriority  []*PortToken
}

var (
	priorities               = NewPriorities()
	globalPrioritiesFilePath = filepath.Join(basedir.GetUserConfigDir(), "deepin/dde-daemon/priorities.json")
)

func contains(cardName string, portName string, substr string) bool {
	return strings.Contains(strings.ToLower(cardName), substr) ||
		strings.Contains(strings.ToLower(portName), substr)
}

func GetPortType(cardName string, portName string) int {
	if contains(cardName, portName, "bluez") {
		return PortTypeBluetooth
	}

	if contains(cardName, portName, "usb") {
		return PortTypeHeadset
	}

	if contains(cardName, portName, "hdmi") {
		return PortTypeHdmi
	}

	if contains(cardName, portName, "speaker") {
		return PortTypeSpeaker
	}

	return PortTypeHeadset
}

func NewPriorities() *Priorities {
	return &Priorities{
		OutputTypePriority:     make([]int, 0),
		OutputInstancePriority: make([]*PortToken, 0),
		InputTypePriority:      make([]int, 0),
		InputInstancePriority:  make([]*PortToken, 0),
	}
}

func (pr *Priorities) Save(file string) error {
	data, err := json.MarshalIndent(pr, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(file, data, 0644)
}

func (pr *Priorities) Print() {
	data, err := json.MarshalIndent(pr, "", "  ")
	if err != nil {
		logger.Warning(err)
	}

	logger.Debug(string(data))
}

func (pr *Priorities) Load(file string, cards CardList) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		logger.Warning(err)
		pr.defaultInit(cards)
		return
	}
	err = json.Unmarshal(data, pr)
	if err != nil {
		logger.Warning(err)
		pr.defaultInit(cards)
		return
	}
	pr.RemoveUnavailable(cards)
	pr.AddAvailable(cards)
}

func (pr *Priorities) RemoveUnavailable(cards CardList) {
	for i := 0; i < len(pr.InputInstancePriority); {
		portToken := pr.InputInstancePriority[i]
		if !pr.checkAvailable(cards, portToken.CardName, portToken.PortName) {
			logger.Debugf("remove input port %s %s", portToken.CardName, portToken.PortName)
			pr.removeInput(i)
		} else {
			i++
		}
	}

	for i := 0; i < len(pr.OutputInstancePriority); {
		portToken := pr.OutputInstancePriority[i]
		if !pr.checkAvailable(cards, portToken.CardName, portToken.PortName) {
			logger.Debugf("remove output port %s %s", portToken.CardName, portToken.PortName)
			pr.removeOutput(i)
		} else {
			i++
		}
	}
}

func (pr *Priorities) AddAvailable(cards CardList) {
	for _, card := range cards {
		for _, port := range card.Ports {
			if port.Available == pulse.AvailableTypeNo {
				logger.Debugf("unavailable port %s %s", card.core.Name, port.Name)
				continue
			}

			_, portConfig := configKeeper.GetCardAndPortConfig(card.core.Name, port.Name)
			if !portConfig.Enabled {
				logger.Debugf("disabled port %s %s", card.core.Name, port.Name)
				continue
			}

			if port.Direction == pulse.DirectionSink && pr.findOutput(card.core.Name, port.Name) < 0 {
				logger.Debugf("add output port %s %s", card.core.Name, port.Name)
				pr.AddOutputPort(card.core.Name, port.Name)
			} else if port.Direction == pulse.DirectionSource && pr.findInput(card.core.Name, port.Name) < 0 {
				logger.Debugf("add input port %s %s", card.core.Name, port.Name)
				pr.AddInputPort(card.core.Name, port.Name)
			}
		}
	}
}

func (pr *Priorities) AddInputPort(cardName string, portName string) {
	portType := GetPortType(cardName, portName)
	token := PortToken{cardName, portName}
	for i := 0; i < len(pr.InputInstancePriority); i++ {
		p := pr.InputInstancePriority[i]
		t := GetPortType(p.CardName, p.PortName)
		if t == portType || pr.IsInputTypeAfter(portType, t) {
			pr.insertInput(i, &token)
			return
		}
	}

	pr.InputInstancePriority = append(pr.InputInstancePriority, &token)
}

func (pr *Priorities) AddOutputPort(cardName string, portName string) {
	portType := GetPortType(cardName, portName)
	token := PortToken{cardName, portName}
	for i := 0; i < len(pr.OutputInstancePriority); i++ {
		p := pr.OutputInstancePriority[i]
		t := GetPortType(p.CardName, p.PortName)
		if t == portType || pr.IsOutputTypeAfter(portType, t) {
			pr.insertOutput(i, &token)
			return
		}
	}

	pr.OutputInstancePriority = append(pr.OutputInstancePriority, &token)
}

func (pr *Priorities) RemoveCard(cardName string) {
	for i := 0; i < len(pr.InputInstancePriority); {
		p := pr.InputInstancePriority[i]
		if p.CardName == cardName {
			pr.removeInput(i)
		} else {
			i++
		}
	}

	for i := 0; i < len(pr.OutputInstancePriority); {
		p := pr.OutputInstancePriority[i]
		if p.CardName == cardName {
			pr.removeOutput(i)
		} else {
			i++
		}
	}
}

func (pr *Priorities) RemoveInputPort(cardName string, portName string) {
	index := pr.findInput(cardName, portName)
	if index >= 0 {
		pr.removeInput(index)
	}
}

func (pr *Priorities) RemoveOutputPort(cardName string, portName string) {
	index := pr.findOutput(cardName, portName)
	if index >= 0 {
		pr.removeOutput(index)
	}
}

func (pr *Priorities) SetInputPortFirst(cardName string, portName string) {
	portType := GetPortType(cardName, portName)
	token := PortToken{cardName, portName}
	index := pr.findInput(cardName, portName)
	if index < 0 {
		return
	}
	pr.removeInput(index)
	pr.insertInput(0, &token)
	pr.setInputTypeFirst(portType)

	index = 1
	for i := 1; i < len(pr.InputInstancePriority); i++ {
		p := pr.InputInstancePriority[i]
		t := GetPortType(p.CardName, p.PortName)
		if t == portType {
			pr.removeInput(i)
			pr.insertInput(index, p)
			index++
		}
	}
}

func (pr *Priorities) SetOutputPortFirst(cardName string, portName string) {
	portType := GetPortType(cardName, portName)
	token := PortToken{cardName, portName}
	index := pr.findOutput(cardName, portName)
	if index < 0 {
		return
	}
	pr.removeOutput(index)
	pr.insertOutput(0, &token)
	pr.setOutputTypeFirst(portType)

	index = 1
	for i := 1; i < len(pr.OutputInstancePriority); i++ {
		p := pr.OutputInstancePriority[i]
		t := GetPortType(p.CardName, p.PortName)
		if t == portType {
			pr.removeOutput(i)
			pr.insertOutput(index, p)
			index++
		}
	}
}

func (pr *Priorities) GetFirstInput() (string, string) {
	if len(pr.InputInstancePriority) > 0 {
		port := pr.InputInstancePriority[0]
		return port.CardName, port.PortName
	} else {
		return "", ""
	}
}

func (pr *Priorities) GetFirstOutput() (string, string) {
	if len(pr.OutputInstancePriority) > 0 {
		port := pr.OutputInstancePriority[0]
		return port.CardName, port.PortName
	} else {
		return "", ""
	}
}

func (pr *Priorities) IsInputTypeAfter(type1 int, type2 int) bool {
	for _, t := range pr.InputTypePriority {
		if t == type1 {
			return true
		}

		if t == type2 {
			return false
		}
	}

	return false
}

func (pr *Priorities) IsOutputTypeAfter(type1 int, type2 int) bool {
	for _, t := range pr.OutputTypePriority {
		if t == type1 {
			return true
		}

		if t == type2 {
			return false
		}
	}

	return false
}

func (pr *Priorities) defaultInit(cards CardList) {
	pr.OutputTypePriority = append(pr.OutputTypePriority, PortTypeBluetooth)
	pr.OutputTypePriority = append(pr.OutputTypePriority, PortTypeHeadset)
	pr.OutputTypePriority = append(pr.OutputTypePriority, PortTypeSpeaker)
	pr.OutputTypePriority = append(pr.OutputTypePriority, PortTypeHdmi)

	pr.InputTypePriority = append(pr.InputTypePriority, PortTypeBluetooth)
	pr.InputTypePriority = append(pr.InputTypePriority, PortTypeHeadset)
	pr.InputTypePriority = append(pr.InputTypePriority, PortTypeSpeaker)
	pr.InputTypePriority = append(pr.InputTypePriority, PortTypeHdmi)

	pr.AddAvailable(cards)
}

func (pr *Priorities) checkAvailable(cards CardList, cardName string, portName string) bool {
	for _, card := range cards {
		if cardName != card.core.Name {
			continue
		}
		for _, port := range card.Ports {
			if portName != port.Name {
				continue
			}

			if port.Available == pulse.AvailableTypeYes {
				_, portConfig := configKeeper.GetCardAndPortConfig(cardName, portName)
				return portConfig.Enabled
			} else if port.Available == pulse.AvailableTypeUnknow {
				logger.Warningf("port(%s %s) available is unknown", cardName, portName)
				_, portConfig := configKeeper.GetCardAndPortConfig(cardName, portName)
				return portConfig.Enabled
			} else {
				return false
			}
		}
	}

	return false
}

func (pr *Priorities) removeInput(index int) {
	pr.InputInstancePriority = append(
		pr.InputInstancePriority[:index],
		pr.InputInstancePriority[index+1:]...,
	)
}

func (pr *Priorities) removeOutput(index int) {
	pr.OutputInstancePriority = append(
		pr.OutputInstancePriority[:index],
		pr.OutputInstancePriority[index+1:]...,
	)
}

func (pr *Priorities) insertInput(index int, portToken *PortToken) {
	tail := append([]*PortToken{}, pr.InputInstancePriority[index:]...)
	pr.InputInstancePriority = append(pr.InputInstancePriority[:index], portToken)
	pr.InputInstancePriority = append(pr.InputInstancePriority, tail...)
}

func (pr *Priorities) insertOutput(index int, portToken *PortToken) {
	tail := append([]*PortToken{}, pr.OutputInstancePriority[index:]...)
	pr.OutputInstancePriority = append(pr.OutputInstancePriority[:index], portToken)
	pr.OutputInstancePriority = append(pr.OutputInstancePriority, tail...)
}

func (pr *Priorities) setInputTypeFirst(portType int) {
	temp := []int{portType}
	for _, t := range pr.InputTypePriority {
		if t != portType {
			temp = append(temp, t)
		}
	}
	pr.InputTypePriority = temp
}

func (pr *Priorities) setOutputTypeFirst(portType int) {
	temp := []int{portType}
	for _, t := range pr.OutputTypePriority {
		if t != portType {
			temp = append(temp, t)
		}
	}
	pr.OutputTypePriority = temp
}

func (pr *Priorities) findInput(cardName string, portName string) int {
	for i := 0; i < len(pr.InputInstancePriority); i++ {
		p := pr.InputInstancePriority[i]
		if p.CardName == cardName && p.PortName == portName {
			return i
		}
	}

	return -1
}

func (pr *Priorities) findOutput(cardName string, portName string) int {
	for i := 0; i < len(pr.OutputInstancePriority); i++ {
		p := pr.OutputInstancePriority[i]
		if p.CardName == cardName && p.PortName == portName {
			return i
		}
	}

	return -1
}
