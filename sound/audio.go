package main

// #cgo amd64 386 CFLAGS: -g -Wall
// #cgo LDFLAGS: -L. -lpulse -lc
// #include "stdio.h"
// #include "dde-pulse.h"
import "C"

import (
	"dlib/dbus"
	"fmt"
	"os"
	"strconv"
	"unsafe"
)

type Audio struct {
	NumDevics int32
	pa        *C.pa
	HostName  string
	UserName  string
	cards     map[int]*Card
	sinks     map[int]*Sink
	sources   map[int]*Source
	clients   map[int]*Client
	//Change func(int32)
}

type CardProfileInfo struct {
	Name        string
	Description string
}

type Card struct {
	Index         int32
	Name          string
	Owner_module  int32
	Driver        string
	NProfiles     int32
	Profiles      []CardProfileInfo
	ActiveProfile *CardProfileInfo
}

type SinkPortInfo struct {
	Name        string
	Description string
	Available   int32
}

type Sink struct {
	Index        int32
	Name         string
	Description  string
	Driver       string
	Mute         int32
	NVolumeSteps int32
	Card         int32
	Volume       int32

	NPorts     int32
	Ports      []SinkPortInfo
	ActivePort *SinkPortInfo
}

//Only capitalized first character in Capitalized structure can be exposed

type SourcePortInfo struct {
	Name        string
	Description string
	Available   int32
}

type Source struct {
	Index        int32
	Name         string
	Description  string
	Driver       string
	Mute         int32
	NVolumeSteps int32
	Card         int32
	C_ports      int32
	N_formates   int32
	Volume       int32

	NPorts     int32
	Ports      []SourcePortInfo
	ActivePort *SourcePortInfo
}

type SinkInput struct {
	Index           int32
	Name            string
	Owner_module    int32
	Client          int32
	Sink            int32
	Volume          int32
	Driver          string
	Mute            int32
	Has_volume      int32
	Volume_writable int32
}

type SourceOutput struct {
	Index           int32
	Name            string
	Owner_module    int32
	Client          int32
	Source          int32
	Volume          int32
	Driver          string
	Mute            int32
	Has_volume      int32
	Volume_writable int32
}

type Client struct {
	Index        int32
	Name         string
	Owner_module int32
	Driver       string
	//pa_proplist *proplist
	Prop map[string]string
}

type Volume struct {
	Channels uint32
	Values   [320]uint32
}

func getCardFromC(_card C.card_t) *Card {
	card := &Card{}
	card.Index = int32(_card.index)
	card.Name = C.GoString(&_card.name[0])
	card.Driver = C.GoString(&_card.driver[0])
	card.Owner_module = int32(_card.owner_module)
	card.NProfiles = int32(_card.n_profiles)

	card.Profiles = make([]CardProfileInfo, card.NProfiles)
	for j := 0; j < int(card.NProfiles); j = j + 1 {
		card.Profiles[j].Name = C.GoString(&_card.profiles[j].name[0])
		card.Profiles[j].Description = C.GoString(&_card.profiles[j].description[0])
		ret := C.strcmp((*C.char)(&_card.active_profile.name[0]),
			(*C.char)(&_card.profiles[j].name[0]))
		if ret == 0 {
			card.ActiveProfile = &card.Profiles[j]
		}
	}
	return card
}

func getSinkFromC(_sink C.sink_t) *Sink {
	sink := &Sink{}
	sink.Index = int32(_sink.index)
	sink.Card = int32(_sink.card)
	sink.Description =
		C.GoString((*C.char)(unsafe.Pointer(&_sink.description[0])))
	sink.Driver = C.GoString(&_sink.driver[0])
	sink.Mute = int32(_sink.mute)
	sink.Name = C.GoString(&_sink.name[0])
	sink.Volume = int32(C.pa_cvolume_avg(&_sink.volume) * 100 / C.PA_VOLUME_NORM)
	sink.NVolumeSteps = int32(_sink.n_volume_steps)
	//sink.Cvolume.Channels = uint32(_sink.volume.channels)
	//for j := 0; j < int(sink.Cvolume.Channels); j = j + 1 {
	//sink.Cvolume.Values[j] =
	//*((*uint32)(unsafe.Pointer(&_sink.volume.values[j])))
	//}
	sink.NPorts = int32(_sink.n_ports)
	sink.Ports = make([]SinkPortInfo, sink.NPorts)
	for j := 0; j < int(sink.NPorts); j = j + 1 {
		sink.Ports[j].Available = int32(_sink.ports[j].available)
		sink.Ports[j].Name = C.GoString(&_sink.ports[j].name[0])
		sink.Ports[j].Description = C.GoString(&_sink.ports[j].description[0])
		ret := C.strcmp((*C.char)(&_sink.ports[j].name[0]),
			(*C.char)(&_sink.active_port.name[0]))
		if ret == 0 {
			sink.ActivePort = &sink.Ports[j]
		}
	}
	if sink.NPorts == 0 {
		sink.ActivePort = &SinkPortInfo{"", "", 0}
	}
	fmt.Println("Index: " + strconv.Itoa(int((sink.Index))) + " Card:" + strconv.Itoa(int(sink.Card)))
	return sink
}

func getSourceFromC(_source C.source_t) *Source {
	source := &Source{}
	source.Index = int32(_source.index)
	source.Card = int32(_source.card)
	source.Mute = int32(_source.mute)
	source.Name = C.GoString((*C.char)(unsafe.Pointer(&_source.name[0])))
	source.Description = C.GoString(&_source.description[0])

	source.Volume = int32(100 * C.pa_cvolume_avg(&_source.volume) / C.PA_VOLUME_NORM)
	source.NVolumeSteps = int32(_source.n_volume_steps)
	//source.Cvolume.Channels = uint32(_source.volume.channels)
	//for j := uint32(0); j < source.Cvolume.Channels; j = j + 1 {
	//source.Cvolume.Values[j] =
	//*((*uint32)(unsafe.Pointer(&_source.volume.values[j])))
	//}

	source.NPorts = int32(_source.n_ports)
	source.Ports = make([]SourcePortInfo, source.NPorts)
	for j := 0; j < int(source.NPorts); j = j + 1 {
		source.Ports[j].Available = int32(_source.ports[j].available)
		source.Ports[j].Name = C.GoString(&_source.ports[j].name[0])
		source.Ports[j].Description = C.GoString(&_source.ports[j].description[0])
		ret := C.strcmp(&_source.ports[j].name[0],
			&_source.active_port.name[0])
		if ret == 0 {
			source.ActivePort = &source.Ports[j]
		}
	}
	if source.NPorts == 0 {
		source.ActivePort = &SourcePortInfo{"", "", 0}
	}

	return source
}

func NewAudio() (*Audio, error) {
	audio = &Audio{}
	audio.pa = C.pa_new()
	if audio.pa == nil {
		fmt.Fprintln(os.Stderr,
			"unable to connect to the pulseaudio server,exiting\n")
		os.Exit(-1)
	}
	audio.cards = make(map[int]*Card)
	audio.sinks = make(map[int]*Sink)
	audio.sources = make(map[int]*Source)
	audio.cards = audio.getCards()
	audio.sinks = audio.getsinks()
	audio.sources = audio.getSources()
	return audio, nil
}

func (o *Audio) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.Audio",
		"/com/deepin/daemon/Audio",
		"com.deepin.daemon.Audio",
	}
}

func (audio *Audio) getServerInfo() *Audio {
	C.pa_get_server_info(audio.pa)
	audio.HostName = C.GoString(audio.pa.server_info.host_name)
	audio.UserName = C.GoString(audio.pa.server_info.user_name)
	//fmt.Print("go: " + audio.HostName + "\n")
	fmt.Print("go: " + C.GoString((audio.pa.server_info.host_name)) + "\n")

	return audio
}

//export updateCard
func updateCard(index int,
	event C.pa_subscription_event_type_t) {
	switch event {
	case C.PA_SUBSCRIPTION_EVENT_NEW:
		//i := int32(audio.pa.cards[0].index)
		audio.cards[index] = getCardFromC(audio.pa.cards[0])
		dbus.InstallOnSession(audio.cards[index])
		break
	case C.PA_SUBSCRIPTION_EVENT_CHANGE:
		audio.cards[index] = getCardFromC(audio.pa.cards[0])
		dbus.InstallOnSession(audio.cards[index])
		//for i := 0; i < len(audio.cards); i = i + 1 {
		//if audio.cards[i].Index == int32(audio.pa.cards[0].index) {
		//audio.cards[i] = getCardFromC(audio.pa.cards[0])
		//fmt.Print("updating card property: " + audio.cards[i].ActiveProfile.Name)
		//dbus.InstallOnSession(audio.cards[i])
		//break
		//}
		//}
		break
	case C.PA_SUBSCRIPTION_EVENT_REMOVE:
		dbus.UnInstallObject(audio.cards[index])
		delete(audio.cards, index)
		break
	}
}

//export updateSink
func updateSink(index int,
	event C.pa_subscription_event_type_t) {
	switch event {
	case C.PA_SUBSCRIPTION_EVENT_NEW:
		audio.sinks[index] = getSinkFromC(audio.pa.sinks[0])
		dbus.InstallOnSession(audio.sinks[index])
		break
	case C.PA_SUBSCRIPTION_EVENT_CHANGE:
		//for i := 0; i < len(audio.sinks); i = i + 1 {
		//if audio.sinks[i].Index == int32(audio.pa.sinks[0].index) {
		//audio.sinks[i] = getSinkFromC(audio.pa.sinks[0])
		//dbus.InstallOnSession(audio.sinks[i])
		//}
		//}
		audio.sinks[index] = getSinkFromC(audio.pa.sinks[0])
		fmt.Println("Sink.mute:  \n", audio.sinks[index].Mute)
		dbus.InstallOnSession(audio.sinks[index])

	case C.PA_SUBSCRIPTION_EVENT_REMOVE:
		dbus.UnInstallObject(audio.sinks[index])
		delete(audio.sinks, index)
		break
	}
}

//export updateSource
func updateSource(index int,
	event C.pa_subscription_event_type_t) {
	fmt.Print("Updating source property:")
	switch event {
	case C.PA_SUBSCRIPTION_EVENT_NEW:
		audio.sources[index] = getSourceFromC(audio.pa.sources[0])
		dbus.InstallOnSession(audio.sources[index])
		break
	case C.PA_SUBSCRIPTION_EVENT_CHANGE:
		//for i, _ := range audio.sources {
		//if audio.sources[i].Index == int32(audio.pa.sources[0].index) {
		//audio.sources[i] = getSourceFromC(audio.pa.sources[0])
		//dbus.InstallOnSession(audio.sources[i])
		//break
		//}
		//}
		audio.sources[index] = getSourceFromC(audio.pa.sources[0])
		dbus.InstallOnSession(audio.sources[index])
		break
	case C.PA_SUBSCRIPTION_EVENT_REMOVE:
		if audio.sources[index] != nil {
			dbus.UnInstallObject(audio.sources[index])
		}
		delete(audio.sources, index)
		break
	}
}

func (audio *Audio) GetServerInfo() *Audio {
	return audio.getServerInfo()
}

func (audio *Audio) getCard() *Card {
	card := &Card{}
	return card
}

func (audio *Audio) getCards() map[int]*Card {
	C.pa_get_card_list(audio.pa)
	n := int(audio.pa.n_cards)
	//audio.cards = make([]*Card, n)

	for i := 0; i < n; i = i + 1 {
		index := int(audio.pa.cards[i].index)
		audio.cards[index] = getCardFromC(audio.pa.cards[i])
	}
	return audio.cards
}

//func (audio *Audio) GetCards() []*Card {
//return audio.getCards()
//}

func (audio *Audio) getsinks() map[int]*Sink {
	C.pa_get_device_list(audio.pa)
	n := int(audio.pa.n_sinks)
	//sinks := make([]*Sink, n)

	for i := 0; i < n; i = i + 1 {
		index := int(audio.pa.sinks[i].index)
		audio.sinks[index] = getSinkFromC(audio.pa.sinks[i])
	}

	return audio.sinks

}

//func (audio *Audio) GetSinks() []*Sink {
//return audio.getsinks()
//}

func (audio *Audio) getSources() map[int]*Source {
	C.pa_get_device_list(audio.pa)
	n := int(audio.pa.n_sources)
	//sources := make([]*Source, n)

	for i := 0; i < n; i = i + 1 {
		index := int(audio.pa.sources[i].index)
		audio.sources[index] = getSourceFromC(audio.pa.sources[i])
	}
	return audio.sources
}

//func (audio *Audio) GetSources() []*Source {
//return audio.getSources()
//}

func (audio *Audio) getSinkInputs() []*SinkInput {
	C.pa_get_sink_input_list(audio.pa)
	n := int(audio.pa.n_sink_inputs)
	sinkInputs := make([]*SinkInput, n)

	fmt.Fprint(os.Stderr, audio.pa.sink_inputs[0].index)

	for i := 0; i < n; i = i + 1 {
		sinkInputs[i] = &SinkInput{}
		sinkInputs[i].Index = int32(audio.pa.sink_inputs[i].index)
		sinkInputs[i].Client = int32(audio.pa.sink_inputs[i].client)
		sinkInputs[i].Sink = int32(audio.pa.sink_inputs[i].sink)
		sinkInputs[i].Mute = int32(audio.pa.sink_inputs[i].mute)
		sinkInputs[i].Has_volume = int32(audio.pa.sink_inputs[i].has_volume)
		sinkInputs[i].Volume_writable = int32(audio.pa.sink_inputs[i].volume_writable)
		//sinkInputs[i].Cvolume.Channels = uint32(audio.pa.sink_inputs[i].volume.channels)
		//for j := uint32(0); j < sinkInputs[i].Cvolume.Channels; j = j + 1 {
		//sinkInputs[i].Cvolume.Values[j] =
		//*(*uint32)(unsafe.Pointer(&audio.pa.sink_inputs[i].volume.values[j]))
		//}
	}
	return sinkInputs
}

func (audio *Audio) GetSinkInputs() []*SinkInput {
	return audio.getSinkInputs()
}

func (audio *Audio) GetSourceOutputs() []*SourceOutput {
	C.pa_get_source_output_list(audio.pa)
	n := int(audio.pa.n_source_outputs)
	var sourceOutputs = make([]*SourceOutput, n)

	fmt.Print(len(sourceOutputs))

	for i := 0; i < n; i = i + 1 {
		sourceOutputs[i] = &SourceOutput{}
		sourceOutputs[i].Index = int32(audio.pa.source_outputs[i].index)
		sourceOutputs[i].Name = C.GoString(&audio.pa.source_outputs[i].name[0])
		sourceOutputs[i].Owner_module = int32(audio.pa.source_outputs[i].owner_module)
		sourceOutputs[i].Client = int32(audio.pa.source_outputs[i].client)
		sourceOutputs[i].Source = int32(audio.pa.source_outputs[i].source)
		sourceOutputs[i].Driver = C.GoString(&audio.pa.source_outputs[i].driver[0])
		sourceOutputs[i].Mute = int32(audio.pa.source_outputs[i].mute)
		//sourceOutputs[i].Cvolume.Channels = uint32(audio.pa.source_outputs[i].volume.channels)

		//for j := uint32(0); j < sourceOutputs[i].Cvolume.Channels; j = j + 1 {
		//sourceOutputs[i].Cvolume.Values[j] =
		//*((*uint32)(unsafe.Pointer(&audio.pa.source_outputs[i].volume.values[j])))
		//}
	}
	return sourceOutputs
}

func (audio *Audio) GetClients() []*Client {
	C.pa_get_client_list(audio.pa)
	n := int(audio.pa.n_clients)
	var clients = make([]*Client, n)

	for i := 0; i < n; i = i + 1 {
		clients[i] = &Client{}
		clients[i].Index = int32(audio.pa.clients[i].index)
		clients[i].Owner_module = int32(audio.pa.clients[i].owner_module)
		clients[i].Name = C.GoString((*C.char)(unsafe.Pointer(&audio.pa.clients[i].name[0])))
		clients[i].Driver = C.GoString((*C.char)(unsafe.Pointer(&audio.pa.clients[i].driver[0])))
	}
	return clients
}

func (sink *Sink) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.Audio",
		"/com/deepin/daemon/Audio/Output" + strconv.FormatInt(int64(sink.Index), 10),
		"com.deepin.daemon.Audio.Output",
	}
}

func (card *Card) GetCardProfile() []CardProfileInfo {
	return card.Profiles
}

func (card *Card) setCardProfile(index C.int, port *C.char) int32 {
	return int32(C.pa_set_card_profile_by_index(
		audio.pa,
		index,
		port))
}

func (card *Card) SetCardProfile(port string) int32 {
	return card.setCardProfile(C.int(card.Index), (*C.char)(C.CString(port)))
}

func (sink *Sink) setSinkPort(port *C.char) int32 {
	ret := C.pa_set_sink_port_by_index(audio.pa, C.int(sink.Index), port)
	return int32(ret)
}

func (sink *Sink) SetSinkPort(portname string) int32 {
	port := C.CString(portname)
	return sink.setSinkPort(port)
}

func (sink *Sink) GetSinkVolume() int32 {
	return sink.Volume
}

func (sink *Sink) SetSinkVolume(volume int32) int32 {
	var cvolume C.pa_cvolume
	cvolume.channels = C.uint8_t(2)
	for i := 0; i < 2; i = i + 1 {
		cvolume.values[i] = C.pa_volume_t(volume * C.PA_VOLUME_NORM / 100)
	}

	return sink.setSinkVolume(&cvolume)
	//var cvolume C.pa_cvolume
	//cvolume.channels = C.uint8_t(2)
	//for i := uint32(0); i < 2; i = i + 1 {
	//cvolume.values[i] = *((*C.pa_volume_t)(unsafe.Pointer(&volume)))
	//}
	////for i := uint32(0); i < volume.Channels; i = i + 1 {
	////cvolume.values[i] = *((*C.pa_volume_t)(unsafe.Pointer(&volume.Values[i])))
	////}
	//return int32(C.pa_set_sink_volume_by_index(
	//audio.pa, C.int(sink.Index), &cvolume))
}

func (sink *Sink) setSinkVolume(volume *C.pa_cvolume) int32 {
	return int32(C.pa_set_sink_volume_by_index(
		audio.pa, C.int(sink.Index), volume))
}

func (sink *Sink) setSinkMute(mute int32) int32 {
	ret := C.pa_set_sink_mute_by_index(
		audio.pa, C.int(sink.Index), C.int(mute))
	return int32(ret)
}

func (sink *Sink) SetSinkMute(mute int32) int32 {
	return sink.setSinkMute(mute)
}

func (source *Source) setSourcePort(port *C.char) int32 {
	ret := C.pa_set_source_port_by_index(audio.pa, C.int(source.Index), port)
	return int32(ret)
}

func (source *Source) SetSourcePort(portname string) int32 {
	port := C.CString(portname)
	return source.setSourcePort(port)
}
func (source *Source) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.Audio",
		"/com/deepin/daemon/Audio/Input" + strconv.FormatInt(int64(source.Index), 10),
		"com.deepin.daemon.Audio.Input",
	}
}

func (source *Source) GetSourceVolume() int32 {
	return source.Volume
}

func (source *Source) setSourceVolume(volume *C.pa_cvolume) int32 {
	return int32(C.pa_set_source_volume_by_index(
		audio.pa, C.int(source.Index), volume))
}

func (source *Source) SetSourceVolume(volume int32) int32 {
	var cvolume C.pa_cvolume
	cvolume.channels = C.uint8_t(2)
	for i := 0; i < int(cvolume.channels); i = i + 1 {
		cvolume.values[i] = *(*C.pa_volume_t)(unsafe.Pointer(&volume))
	}
	return source.setSourceVolume(&cvolume)
}

func (source *Source) SetSourceMute(mute int32) int32 {
	ret := C.pa_set_source_mute_by_index(
		audio.pa, C.int(source.Index), C.int(mute))
	return int32(ret)
}

func (sink_input *SinkInput) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.Audio",
		"/com/deepin/daemon/Audio/Application" + strconv.FormatInt(int64(sink_input.Index), 10),
		"com.deepin.daemon.Audio.Application",
	}
}

func (sink_input *SinkInput) GetSink_input_volume() Volume {
	return Volume{}
}

func (sink_input *SinkInput) SetSink_input_volume(volume Volume) int32 {
	return 1
}

func (sink_input *SinkInput) SetSink_input_mute(mute int32) int32 {
	return 1
}

func (source_output *SourceOutput) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.Audio",
		"/com/deepin/daemon/Audio/Application" + strconv.FormatInt(int64(source_output.Index), 10),
		"com.deepin.daemon.Audio.Application",
	}
}

func (source_output *SourceOutput) GetSource_ouput_volume() int32 {
	return source_output.Volume
}

func (source_output *SourceOutput) SetSource_output_volume(volume Volume) Volume {
	return Volume{}
}

func (source_output *SourceOutput) SetSource_output_mute(mute int32) int32 {
	return 1
}

func (client *Client) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.Audio",
		"/com/deepin/daemon/Audio/Client" + strconv.FormatInt(int64(client.Index), 10),
		"com.deepin.daemon.Audio.Client",
	}
}

func (card *Card) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.Audio",
		"/com/deepin/daemon/Audio/Card" + strconv.FormatInt(int64(card.Index), 10),
		"com.deepin.daemon.Audio.Card",
	}
}

var audio *Audio

func main() {
	var err error
	audio, err = NewAudio()
	if err != nil {
		panic(err)
		os.Exit(-1)
	}
	dbus.InstallOnSession(audio)
	for i := range audio.cards {
		dbus.InstallOnSession(audio.cards[i])
	}
	for i := range audio.sinks {
		dbus.InstallOnSession(audio.sinks[i])
	}
	for i, _ := range audio.sources {
		dbus.InstallOnSession(audio.sources[i])
	}
	fmt.Println("module started\n")
	C.pa_subscribe(audio.pa)
	//select {}
}
