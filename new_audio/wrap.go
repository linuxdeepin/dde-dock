package main

/*
#include "dde-pulse.h"
#cgo pkg-config: libpulse
*/
import "C"
import "fmt"

func (info *paInfo) ToSink() *Sink {
	if info.Type != C.PA_SUBSCRIPTION_EVENT_SINK {
		return nil
	}
	if info.core == nil {
		panic("ToSink shoud only be called once per paInfo")
	}
	defer C.free(info.core)
	return toSinkInfo((*C.pa_sink_info)(info.core))
}

func toSinkInfo(info *C.pa_sink_info) *Sink {
	s := &Sink{}
	s.index = uint32(info.index)
	s.Name = C.GoString(info.name)
	s.Description = C.GoString(info.description)
	return s
}

func (*paInfo) ToSinkInput() *SinkInput {
	return nil
}

func GetSinkInputInfo(index uint32) {
	c := newCookie()
	C.get_sink_input_info(ctx, C.int64_t(c.id), C.uint32_t(index))

	fmt.Println("GetSinInput...", index)
	info := c.Reply()
	fmt.Println("wait GetSinInput...", index)
	info.ToSink()
}

func print_sink_info(info *C.pa_sink_info) {
	fmt.Printf("\tindex: %d\n", info.index)
	fmt.Printf("\tname: %s\n", C.GoString(info.name))
	fmt.Printf("\tdescription: %s\n", C.GoString(info.description))
	fmt.Printf("\tmute: %d\n", info.mute)
	fmt.Printf("\tvolume: channels:%d, min:%d, max:%d\n",
		info.volume.channels,
		C.pa_cvolume_min(&info.volume),
		C.pa_cvolume_max(&info.volume))
	if info.active_port != nil {
		fmt.Printf("\tactive port: name: %s\t description: %s\n",
			C.GoString(info.active_port.name), C.GoString(info.active_port.description))
	}
}

func print_sink_input(info *C.pa_sink_input_info) {
	//printf("format_info: %s\n", pa_format_info_snprint(buf, 1000, i->format));
	fmt.Printf("------------------------------\n")
	fmt.Printf("index: %d\n", info.index)
	fmt.Printf("name: %s\n", C.GoString(info.name))
	fmt.Printf("volume: channels:%d, min:%d, max:%d\n",
		info.volume.channels,
		C.pa_cvolume_min(&info.volume),
		C.pa_cvolume_max(&info.volume))
	fmt.Printf("mute: %d\n", info.mute)

	//while ((prop_key = pa_proplist_iterate(i->proplist, &prop_state))) {
	//printf("  %s: %s\n",
	//prop_key,
	//pa_proplist_gets(i->proplist, prop_key));
	//}
	//printf("format_info: %s\n", pa_format_info_snprint(buf, 1000, i->format));
	//printf("------------------------------\n");
}
func (s *Sink) forceUpdate() {
	c := newCookie()
	C.get_sink_info(ctx, C.int64_t(c.id), C.uint32_t(s.index))
	s = c.Reply().ToSink()
	s.valid = true
}

func (a *Audio) forceUpdate() {
	c := newCookie()
	C.get_sink_info_list(ctx, C.int64_t(c.id))
	fmt.Println("__________")
	for _, _info := range c.ReplyList() {
		s := _info.ToSink()
		fmt.Println(_info, s.Name)
	}
}
