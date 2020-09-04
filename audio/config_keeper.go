package audio

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"pkg.deepin.io/lib/xdg/basedir"
)

type PortConfig struct {
	Name           string
	Enabled        bool
	Volume         float64
	IncreaseVolume bool
	Balance        float64
	ReduceNoise    bool
	Mute           bool
}

type CardConfig struct {
	Name  string
	Ports map[string]*PortConfig // Name => PortConfig
}

type ConfigKeeper struct {
	ConfigMap map[string]*CardConfig // Name => CardConfig
}

var (
	configKeeper     = NewConfigKeeper()
	configKeeperFile = filepath.Join(basedir.GetUserConfigDir(), "deepin/dde-daemon/audio-config-keeper.json")
)

func NewConfigKeeper() *ConfigKeeper {
	return &ConfigKeeper{
		ConfigMap: make(map[string]*CardConfig),
	}
}

func NewCardConfig(name string) *CardConfig {
	return &CardConfig{
		Name:  name,
		Ports: make(map[string]*PortConfig),
	}
}

func NewPortConfig(name string) *PortConfig {
	return &PortConfig{
		Name:           name,
		Enabled:        true,
		Volume:         0.5,
		IncreaseVolume: false,
		Balance:        0.0,
		ReduceNoise:    false,
		Mute:           false,
	}
}

func (ck *ConfigKeeper) Save(file string) error {
	data, err := json.MarshalIndent(ck.ConfigMap, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(file, data, 0644)
}

func (ck *ConfigKeeper) Load(file string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &ck.ConfigMap)
}

func (ck *ConfigKeeper) Print() {
	data, err := json.MarshalIndent(ck.ConfigMap, "", "  ")
	if err != nil {
		logger.Warning(err)
		return
	}
	logger.Debug(string(data))
}

func (ck *ConfigKeeper) UpdateCardConfig(cardConfig *CardConfig) {
	ck.ConfigMap[cardConfig.Name] = cardConfig
}

func (ck *ConfigKeeper) RemoveCardConfig(cardName string) {
	delete(ck.ConfigMap, cardName)
}

func (ck *ConfigKeeper) GetCardAndPortConfig(cardName string, portName string) (*CardConfig, *PortConfig) {
	card, ok := ck.ConfigMap[cardName]
	if !ok {
		card = NewCardConfig(cardName)
		port := NewPortConfig(portName)
		card.UpdatePortConfig(port)
		ck.UpdateCardConfig(card)
		return card, port
	}

	port, ok := card.Ports[portName]
	if !ok {
		port = NewPortConfig(portName)
		card.UpdatePortConfig(port)
		ck.UpdateCardConfig(card)
	}
	return card, port
}

func (ck *ConfigKeeper) SetEnabled(cardName string, portName string, enabled bool) {
	card, port := ck.GetCardAndPortConfig(cardName, portName)
	port.Enabled = enabled
	card.UpdatePortConfig(port)
	ck.UpdateCardConfig(card)
}

func (ck *ConfigKeeper) SetVolume(cardName string, portName string, volume float64) {
	card, port := ck.GetCardAndPortConfig(cardName, portName)
	port.Volume = volume
	card.UpdatePortConfig(port)
	ck.UpdateCardConfig(card)
}

func (ck *ConfigKeeper) SetIncreaseVolume(cardName string, portName string, enhance bool) {
	card, port := ck.GetCardAndPortConfig(cardName, portName)
	port.IncreaseVolume = enhance
	card.UpdatePortConfig(port)
	ck.UpdateCardConfig(card)
}

func (ck *ConfigKeeper) SetBalance(cardName string, portName string, balance float64) {
	card, port := ck.GetCardAndPortConfig(cardName, portName)
	port.Balance = balance
	card.UpdatePortConfig(port)
	ck.UpdateCardConfig(card)
}

func (ck *ConfigKeeper) SetReduceNoise(cardName string, portName string, reduce bool) {
	card, port := ck.GetCardAndPortConfig(cardName, portName)
	port.ReduceNoise = reduce
	card.UpdatePortConfig(port)
	ck.UpdateCardConfig(card)
}

func (ck *ConfigKeeper) SetMute(cardName string, portName string, mute bool) {
	card, port := ck.GetCardAndPortConfig(cardName, portName)
	port.Mute = mute
	card.UpdatePortConfig(port)
	ck.UpdateCardConfig(card)
}

func (card *CardConfig) UpdatePortConfig(portConfig *PortConfig) {
	card.Ports[portConfig.Name] = portConfig
}

func (card *CardConfig) RemovePortConfig(portName string) {
	delete(card.Ports, portName)
}
