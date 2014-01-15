package main

// #cgo amd64 386 CFLAGS: -g -Wall
// #cgo LDFLAGS: -L. -lpulse -lc
// #cgo CFLAGS: -Wno-error
// #include "stdio.h"
// #include "dde-pulse.h"
import "C"

import (
	"dlib/dbus"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"unsafe"
)

type Audio struct {
	//unexported properties
	pa            *C.pa
	cards         map[int]*Card
	sinks         map[int]*Sink
	sources       map[int]*Source
	clients       map[int]*Client
	sinkInputs    map[int]*SinkInput
	sourceOutputs map[int]*SourceOutput

	//exported properties
	HostName string
	UserName string

	//signals
	Changed       func(string)
	DeviceAdded   func(string)
	DeviceChanged func(string)
	DeviceRemoved func(string)
}

type CardProfileInfo struct {
	Name        string
	Description string
}

type Card struct {
	Index         int32
	Name          string
	owner_module  int32
	driver        string
	NProfiles     int32
	Profiles      []CardProfileInfo
	ActiveProfile *CardProfileInfo
}

type SinkPortInfo struct {
	Name        string
	Description string
	Available   int32 //0:unknown,1:no,2:yes
}

type Sink struct {
	Index       int32
	Name        string
	Description string
	driver      string
	Mute        int32
	Card        int32
	Volume      int32

	//NVolumeSteps int32

	nports     int32
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
	Index       int32
	Name        string
	Description string
	driver      string
	Mute        int32
	//NVolumeSteps int32
	Card int32
	//N_formates int32
	Volume int32

	nports     int32
	Ports      []SourcePortInfo
	ActivePort *SourcePortInfo
}

type SinkInput struct {
	Index           int32
	Name            string
	owner_module    int32
	Client          int32
	Sink            int32
	Volume          int32
	driver          string
	Mute            int32
	Has_volume      int32
	Volume_writable int32

	PropList map[string]string
}

type SourceOutput struct {
	Index           int32
	Name            string
	owner_module    int32
	Client          int32
	Source          int32
	Volume          int32
	driver          string
	Mute            int32
	Has_volume      int32
	Volume_writable int32

	PropList map[string]string
}

type Client struct {
	Index        int32
	Name         string
	owner_module int32
	driver       string
	//pa_proplist *proplist
	Prop map[string]string
}

type Volume struct {
	Channels uint32
	Values   [320]uint32
}

//func compare(x, y interface{}) bool {

//return true
//}

func getDiffProperty(x, y interface{}) map[string]interface{} {

	typ := reflect.TypeOf(x)
	valuex := reflect.ValueOf(x)
	valuey := reflect.ValueOf(y)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		valuex = valuex.Elem()
		valuey = valuey.Elem()
	}

	attrs := make(map[string]interface{})

	//only structs are supported so return an empty result if not
	if typ.Kind() != reflect.Struct {
		fmt.Print("%v type doesn't have attributes inspected\n", typ.Kind())
		return attrs
	}

	//loop through the struct's fields and set the map
	for i := 0; i < typ.NumField(); i++ {
		fielda := typ.Field(i)
		//fieldb := typ.Field(i)
		if !fielda.Anonymous {
			fvaluea := valuex.Field(i)
			fvalueb := valuey.Field(i)
			va := fvaluea.Interface()
			vb := fvalueb.Interface()
			switch fielda.Type.Kind() {
			case reflect.Int32:
				ia := va.(int32)
				ib := vb.(int32)
				if ia != ib {
					attrs[fielda.Name] = vb
				}
			case reflect.String:
				sa := va.(string)
				sb := vb.(string)
				if sa != sb {
					attrs[fielda.Name] = vb
				}
			case reflect.Map:
			case reflect.Struct:
				if !reflect.DeepEqual(va, vb) {
					attrs[fielda.Name] = vb
				}
			default:
				break
			}
		}
	}
	return attrs
}

func getCardFromC(_card C.card_t) *Card {
	card := &Card{}
	card.Index = int32(_card.index)
	card.Name = C.GoString(&_card.name[0])
	card.driver = C.GoString(&_card.driver[0])
	card.owner_module = int32(_card.owner_module)
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
	sink.driver = C.GoString(&_sink.driver[0])
	sink.Mute = int32(_sink.mute)
	sink.Name = C.GoString(&_sink.name[0])
	sink.Volume = int32(C.pa_cvolume_avg(&_sink.volume) * 100 / C.PA_VOLUME_NORM)
	//sink.NVolumeSteps = int32(_sink.n_volume_steps)
	//sink.Cvolume.Channels = uint32(_sink.volume.channels)
	//for j := 0; j < int(sink.Cvolume.Channels); j = j + 1 {
	//sink.Cvolume.Values[j] =
	//*((*uint32)(unsafe.Pointer(&_sink.volume.values[j])))
	//}
	sink.nports = int32(_sink.n_ports)
	sink.Ports = make([]SinkPortInfo, sink.nports)
	for j := 0; j < int(sink.nports); j = j + 1 {
		sink.Ports[j].Available = int32(_sink.ports[j].available)
		sink.Ports[j].Name = C.GoString(&_sink.ports[j].name[0])
		sink.Ports[j].Description = C.GoString(&_sink.ports[j].description[0])
		ret := C.strcmp((*C.char)(&_sink.ports[j].name[0]),
			(*C.char)(&_sink.active_port.name[0]))
		if ret == 0 {
			sink.ActivePort = &sink.Ports[j]
		}
	}
	if sink.nports == 0 {
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
	//source.NVolumeSteps = int32(_source.n_volume_steps)
	//source.Cvolume.Channels = uint32(_source.volume.channels)
	//for j := uint32(0); j < source.Cvolume.Channels; j = j + 1 {
	//source.Cvolume.Values[j] =
	//*((*uint32)(unsafe.Pointer(&_source.volume.values[j])))
	//}

	source.nports = int32(_source.n_ports)
	source.Ports = make([]SourcePortInfo, source.nports)
	for j := 0; j < int(source.nports); j = j + 1 {
		source.Ports[j].Available = int32(_source.ports[j].available)
		source.Ports[j].Name = C.GoString(&_source.ports[j].name[0])
		source.Ports[j].Description = C.GoString(&_source.ports[j].description[0])
		ret := C.strcmp(&_source.ports[j].name[0],
			&_source.active_port.name[0])
		if ret == 0 {
			source.ActivePort = &source.Ports[j]
		}
	}
	if source.nports == 0 {
		source.ActivePort = &SourcePortInfo{"", "", 0}
	}

	return source
}

func getSinkInputFromC(_sink_input C.sink_input_t) *SinkInput {
	sinkInput := &SinkInput{}
	sinkInput.Index = int32(_sink_input.index)
	sinkInput.Client = int32(_sink_input.client)
	sinkInput.Sink = int32(_sink_input.sink)
	sinkInput.Mute = int32(_sink_input.mute)
	sinkInput.Has_volume = int32(_sink_input.has_volume)
	sinkInput.Volume_writable = int32(_sink_input.volume_writable)
	sinkInput.Volume = int32(100 * C.pa_cvolume_avg(&_sink_input.volume) / C.PA_VOLUME_NORM)
	sinkInput.Name = C.GoString(&_sink_input.name[0])
	sinkInput.PropList = make(map[string]string)
	var prop_state *C.void = nil
	key := C.pa_proplist_iterate(
		_sink_input.proplist,
		(*unsafe.Pointer)(unsafe.Pointer(&prop_state)))
	for key != nil {
		sinkInput.PropList[C.GoString(key)] = C.GoString(C.pa_proplist_gets(_sink_input.proplist, key))
		key = C.pa_proplist_iterate(_sink_input.proplist,
			(*unsafe.Pointer)(unsafe.Pointer(&prop_state)))
	}
	//sinkInputs[i].Cvolume.Channels = uint32(audio.pa.sink_inputs[i].volume.channels)
	//for j := uint32(0); j < sinkInputs[i].Cvolume.Channels; j = j + 1 {
	//sinkInputs[i].Cvolume.Values[j] =
	//*(*uint32)(unsafe.Pointer(&audio.pa.sink_inputs[i].volume.values[j]))
	//}
	return sinkInput
}

func getSourceOutputFromC(_source_output C.source_output_t) *SourceOutput {
	sourceOutput := &SourceOutput{}
	sourceOutput.Index = int32(_source_output.index)
	sourceOutput.Name = C.GoString(&_source_output.name[0])
	sourceOutput.owner_module = int32(_source_output.owner_module)
	sourceOutput.Client = int32(_source_output.client)
	sourceOutput.Source = int32(_source_output.source)
	sourceOutput.driver = C.GoString(&_source_output.driver[0])
	sourceOutput.Mute = int32(_source_output.mute)
	//sourceOutputs[i].Cvolume.Channels = uint32(audio.pa.source_outputs[i].volume.channels)

	//for j := uint32(0); j < sourceOutputs[i].Cvolume.Channels; j = j + 1 {
	//sourceOutputs[i].Cvolume.Values[j] =
	//*((*uint32)(unsafe.Pointer(&audio.pa.source_outputs[i].volume.values[j])))
	//}

	return sourceOutput
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

	audio.getServerInfo()
	audio.getCards()
	audio.getsinks()
	audio.getSources()
	audio.getSinkInputs()
	audio.getSourceOutputs()
	return audio, nil
}

//export updateCard
func updateCard(_index C.int,
	event C.pa_subscription_event_type_t) {
	index := int(_index)
	switch event {
	case C.PA_SUBSCRIPTION_EVENT_NEW:
		//i := int32(audio.pa.cards[0].index)
		audio.cards[index] = getCardFromC(audio.pa.cards[0])
		dbus.InstallOnSession(audio.cards[index])
		break
	case C.PA_SUBSCRIPTION_EVENT_CHANGE:
		newcard := getCardFromC(audio.pa.cards[0])
		for i, _ := range audio.cards {
			if audio.cards[i].Index == newcard.Index {
				changes := getDiffProperty(audio.cards[i], newcard)
				audio.cards[index] = newcard
				dbus.InstallOnSession(audio.cards[i])
				for key, v := range changes {
					dbus.NotifyChange(audio.cards[i], key)
					fmt.Println("updating card property: %v\n", v)
				}
				break
			}
		}
		break
	case C.PA_SUBSCRIPTION_EVENT_REMOVE:
		dbus.UnInstallObject(audio.cards[index])
		delete(audio.cards, index)
		break
	}
}

//export updateSink
func updateSink(_index C.int,
	event C.pa_subscription_event_type_t) {
	index := int(_index)
	switch event {
	case C.PA_SUBSCRIPTION_EVENT_NEW:
		audio.sinks[index] = getSinkFromC(audio.pa.sinks[0])
		dbus.InstallOnSession(audio.sinks[index])
		break
	case C.PA_SUBSCRIPTION_EVENT_CHANGE:
		newsink := getSinkFromC(audio.pa.sinks[0])
		for i, _ := range audio.sinks {
			if audio.sinks[i].Index == newsink.Index {
				changes := getDiffProperty(audio.cards[i], newsink)
				audio.sinks[i] = newsink
				dbus.InstallOnSession(audio.sinks[i])
				for key, v := range changes {
					dbus.NotifyChange(audio.sinks[i], key)
					fmt.Println("updating sink property: %v\n", v)
				}
				break
			}
		}

	case C.PA_SUBSCRIPTION_EVENT_REMOVE:
		dbus.UnInstallObject(audio.sinks[index])
		delete(audio.sinks, index)
		break
	}
}

//export updateSource
func updateSource(_index C.int,
	event C.pa_subscription_event_type_t) {
	index := int(_index)
	fmt.Print("Updating source property:")
	switch event {
	case C.PA_SUBSCRIPTION_EVENT_NEW:
		audio.sources[index] = getSourceFromC(audio.pa.sources[0])
		dbus.InstallOnSession(audio.sources[index])
		break
	case C.PA_SUBSCRIPTION_EVENT_CHANGE:
		newsource := getSourceFromC(audio.pa.sources[0])
		for i, _ := range audio.sources {
			if audio.sources[i].Index == newsource.Index {
				changes := getDiffProperty(audio.sources[i], newsource)
				audio.sources[i] = newsource
				for key, _ := range changes {
					dbus.NotifyChange(audio.sources[i], key)
					fmt.Println("updating source information: %v\n", changes[key])
				}
				dbus.InstallOnSession(audio.sources[i])
				break
			}
		}
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

//export updateSinkInput
func updateSinkInput(_index C.int,
	event C.pa_subscription_event_type_t) {
	index := int(_index)
	switch event {
	case C.PA_SUBSCRIPTION_EVENT_NEW:
		//audio.sources[index] = getSourceFromC(audio.pa.sources[0])
		//dbus.InstallOnSession(audio.sources[index])
		//break
	case C.PA_SUBSCRIPTION_EVENT_CHANGE:
		//for i, _ := range audio.sources {
		//if audio.sources[i].Index == int32(audio.pa.sources[0].index) {
		//audio.sources[i] = getSourceFromC(audio.pa.sources[0])
		//dbus.InstallOnSession(audio.sources[i])
		//break
		//}
		//}
		audio.sinkInputs[index] =
			getSinkInputFromC(audio.pa.sink_inputs[0])
		dbus.InstallOnSession(audio.sinkInputs[index])
		break
	case C.PA_SUBSCRIPTION_EVENT_REMOVE:
		if audio.sinkInputs[index] != nil {
			dbus.UnInstallObject(audio.sinkInputs[index])
		}
		delete(audio.sinkInputs, index)
		break
	}
}

//export updateSourceOutput
func updateSourceOutput(_index C.int,
	event C.pa_subscription_event_type_t) {
	index := int(_index)
	switch event {
	case C.PA_SUBSCRIPTION_EVENT_NEW:

	case C.PA_SUBSCRIPTION_EVENT_CHANGE:

		audio.sourceOutputs[index] =
			getSourceOutputFromC(audio.pa.source_outputs[0])
		dbus.InstallOnSession(audio.sourceOutputs[index])
		break
	case C.PA_SUBSCRIPTION_EVENT_REMOVE:
		if audio.sourceOutputs[index] != nil {
			dbus.UnInstallObject(audio.sourceOutputs[index])
		}
		delete(audio.sourceOutputs, index)
		break
	}
}

func (audio *Audio) GetServerInfo() *Audio {
	return audio.getServerInfo()
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

func (audio *Audio) GetCards() []*Card {
	n := len(audio.cards)
	cards := make([]*Card, n)
	j := 0
	for i, _ := range audio.cards {
		cards[j] = audio.cards[i]
		j = j + 1
	}
	return cards
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

func (audio *Audio) GetSinks() []*Sink {
	n := len(audio.sinks)
	sinks := make([]*Sink, n)
	j := 0
	for i, _ := range audio.sinks {
		sinks[j] = audio.sinks[i]
		j = j + 1
	}
	return sinks
}

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

func (audio *Audio) GetSources() []*Source {
	n := len(audio.sources)
	sources := make([]*Source, n)
	j := 0
	for i, _ := range audio.sources {
		sources[j] = audio.sources[i]
		j = j + 1
	}
	return sources
}

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

func (audio *Audio) getSinkInputs() map[int]*SinkInput {
	C.pa_get_sink_input_list(audio.pa)
	n := int(audio.pa.n_sink_inputs)
	audio.sinkInputs = make(map[int]*SinkInput)

	for i := 0; i < n; i = i + 1 {
		audio.sinkInputs[i] = getSinkInputFromC(audio.pa.sink_inputs[i])
	}
	return audio.sinkInputs
}

func (audio *Audio) GetSinkInputs() []*SinkInput {
	return nil
	//return audio.getSinkInputs()
}

func (audio *Audio) getSourceOutputs() map[int]*SourceOutput {
	C.pa_get_source_output_list(audio.pa)
	n := int(audio.pa.n_source_outputs)
	audio.sourceOutputs = make(map[int]*SourceOutput)

	for i := 0; i < n; i = i + 1 {
		audio.sourceOutputs[i] =
			getSourceOutputFromC(audio.pa.source_outputs[i])
	}
	return audio.sourceOutputs
}

func (audio *Audio) GetClients() []*Client {
	C.pa_get_client_list(audio.pa)
	n := int(audio.pa.n_clients)
	var clients = make([]*Client, n)

	for i := 0; i < n; i = i + 1 {
		clients[i] = &Client{}
		clients[i].Index = int32(audio.pa.clients[i].index)
		clients[i].owner_module = int32(audio.pa.clients[i].owner_module)
		clients[i].Name = C.GoString((*C.char)(unsafe.Pointer(&audio.pa.clients[i].name[0])))
		clients[i].driver = C.GoString((*C.char)(unsafe.Pointer(&audio.pa.clients[i].driver[0])))
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

func (card *Card) GetSinks() []*Sink {
	n := len(audio.sinks)
	var sinks []*Sink = make([]*Sink, n)
	j := 0
	for _, sink := range audio.sinks {
		if sink.Card == card.Index {
			sinks[j] = sink
			j = j + 1
		}
	}
	return sinks[0:j]
}

func (card *Card) GetSources() []*Source {
	n := len(audio.sources)
	var sources []*Source = make([]*Source, n)
	j := 0
	for _, source := range audio.sources {
		sources[j] = source
		j = j + 1
	}
	return sources[0:j]
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

func (sinkInput *SinkInput) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.Audio",
		"/com/deepin/daemon/Audio/Application" +
			strconv.FormatInt(int64(sinkInput.Index), 10),
		"com.deepin.daemon.Audio.Application",
	}
}

func (sinkInput *SinkInput) GetSinkInputVolume() int32 {
	return sinkInput.Volume
}

func (sinkInput *SinkInput) setSinkInputVolume(cvolume *C.pa_cvolume) int32 {
	return int32(C.pa_set_sink_input_volume(audio.pa,
		C.int(sinkInput.Index),
		cvolume))
}

func (sinkInput *SinkInput) SetVolume(volume int32) int32 {
	var cvolume C.pa_cvolume
	cvolume.channels = C.uint8_t(2)
	for i := 0; i < 2; i = i + 1 {
		cvolume.values[i] = C.pa_volume_t(volume * C.PA_VOLUME_NORM / 100)
	}
	return sinkInput.setSinkInputVolume(&cvolume)
}

func (sinkInput *SinkInput) setSinkInputMute(mute int32) int32 {
	return int32(C.pa_set_sink_input_mute(audio.pa,
		C.int(sinkInput.Index), C.int(mute)))
}

func (sinkInput *SinkInput) SetMute(mute int32) int32 {
	return sinkInput.setSinkInputMute(mute)
}

func (sourceOutput *SourceOutput) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.Audio",
		"/com/deepin/daemon/Audio/Application" +
			strconv.FormatInt(int64(sourceOutput.Index), 10),
		"com.deepin.daemon.Audio.Application",
	}
}

func (sourceOutput *SourceOutput) GetSourceOuputVolume() int32 {
	return sourceOutput.Volume
}

func (sourceOutput *SourceOutput) setSourceOutputVolume(cvolume *C.pa_cvolume) int32 {
	return int32(C.pa_set_source_output_volume(audio.pa,
		C.int(sourceOutput.Index),
		cvolume))
}

func (sourceOutput *SourceOutput) SetVolume(volume int32) int32 {
	var cvolume C.pa_cvolume
	cvolume.channels = C.uint8_t(2)
	for i := 0; i < 2; i = i + 1 {
		cvolume.values[i] = C.pa_volume_t(volume * C.PA_VOLUME_NORM / 100)
	}
	return sourceOutput.setSourceOutputVolume(&cvolume)
}

func (sourceOutput *SourceOutput) SetMute(mute int32) int32 {
	return 1
}

func (client *Client) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.Audio",
		"/com/deepin/daemon/Audio/Client" +
			strconv.FormatInt(int64(client.Index), 10),
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
	for i, _ := range audio.sinkInputs {
		dbus.InstallOnSession(audio.sinkInputs[i])
	}
	for i, _ := range audio.sourceOutputs {
		dbus.InstallOnSession(audio.sourceOutputs[i])
	}
	dbus.DealWithUnhandledMessage()
	fmt.Println("module started\n")
	C.pa_subscribe(audio.pa)
	//select {}
}
