package ddcci

// #cgo pkg-config: ddcutil
// #include <ddcutil_c_api.h>
import "C"
import (
	"errors"
	"fmt"
	"sync"
	"unsafe"

	"pkg.deepin.io/lib/utils"
)

type DDCCI struct {
	listPointer *C.DDCA_Display_Info_List
	listMu      sync.Mutex

	displayMap map[string]int
}

const (
	brightnessVCP = 0x10
)

func newDDCCI() (*DDCCI, error) {
	ddc := &DDCCI{
		displayMap: make(map[string]int),
	}

	status := C.ddca_set_max_tries(C.DDCA_MULTI_PART_TRIES, 5)
	if status < C.int(0) {
		return nil, fmt.Errorf("Error setting retries: %d", status)
	}

	err := ddc.RefreshDisplays()
	if err != nil {
		return nil, err
	}

	return ddc, nil
}

func (d *DDCCI) freeList() {
	if d.listPointer != nil {
		C.ddca_free_display_info_list(d.listPointer)
		d.listPointer = nil
	}
}

func (d *DDCCI) RefreshDisplays() error {
	d.listMu.Lock()
	defer d.listMu.Unlock()

	d.freeList()

	status := C.ddca_get_display_info_list2(C.bool(true), &d.listPointer)
	if status != C.int(0) {
		err := fmt.Errorf("failed to get display info list: %d", status)
		return err
	}

	for i := 0; i < int(d.listPointer.ct); i++ {
		err := d.initDisplay(i)
		if err != nil {
			logger.Warning(err)
		}
	}

	return nil
}

func (d *DDCCI) initDisplay(idx int) error {
	item := d.getDisplayInfoByIdx(idx)

	var handle C.DDCA_Display_Handle
	status := C.ddca_open_display2(item.dref, C.bool(true), &handle)
	if status != C.int(0) {
		return errors.New("failed to open monitor")
	}

	defer C.ddca_close_display(handle)

	var val C.DDCA_Non_Table_Vcp_Value
	status = C.ddca_get_non_table_vcp_value(handle, brightnessVCP, &val)
	if status != C.int(0) {
		return errors.New("failed to check DDC/CI support")
	}

	edid := C.GoBytes(unsafe.Pointer(&item.edid_bytes), 128)
	edidChecksum := getEDIDChecksum(edid)

	d.displayMap[string(edidChecksum)] = idx
	return nil
}

func (d *DDCCI) SupportBrightness(edidChecksum string) bool {
	d.listMu.Lock()
	_, ok := d.displayMap[edidChecksum]
	d.listMu.Unlock()

	return ok
}

func (d *DDCCI) GetBrightness(edidChecksum string) (brightness int, err error) {
	d.listMu.Lock()
	defer d.listMu.Unlock()

	idx, ok := d.displayMap[edidChecksum]
	if !ok {
		err = errors.New("monitor not support DDC/CI")
		return
	}

	item := d.getDisplayInfoByIdx(idx)

	var handle C.DDCA_Display_Handle
	status := C.ddca_open_display2(item.dref, C.bool(true), &handle)
	if status != C.int(0) {
		err = errors.New("failed to open monitor")
		return
	}

	defer C.ddca_close_display(handle)

	var val C.DDCA_Non_Table_Vcp_Value
	status = C.ddca_get_non_table_vcp_value(handle, brightnessVCP, &val)
	if status != C.int(0) {
		err = errors.New("failed to get brightness")
		return
	}

	brightness = int(val.sl)
	return
}

func (d *DDCCI) SetBrightness(edidChecksum string, percent int) error {
	d.listMu.Lock()
	defer d.listMu.Unlock()

	idx, ok := d.displayMap[edidChecksum]
	if !ok {
		return errors.New("monitor not support DDC/CI")
	}

	item := d.getDisplayInfoByIdx(idx)

	var handle C.DDCA_Display_Handle
	status := C.ddca_open_display2(item.dref, C.bool(true), &handle)
	if status != C.int(0) {
		return errors.New("failed to open monitor")
	}

	defer C.ddca_close_display(handle)

	status = C.ddca_set_non_table_vcp_value(handle, brightnessVCP, 0, C.uchar(percent))
	if status != C.int(0) {
		return errors.New("failed to set brightness via DDC/CI")
	}

	return nil
}

func (d *DDCCI) getDisplayInfoByIdx(idx int) *C.DDCA_Display_Info {
	start := unsafe.Pointer(uintptr(unsafe.Pointer(d.listPointer)) + uintptr(C.sizeof_DDCA_Display_Info_List))
	size := uintptr(C.sizeof_DDCA_Display_Info)

	return (*C.DDCA_Display_Info)(unsafe.Pointer(uintptr(start) + size*uintptr(idx)))
}

func getEDIDChecksum(edid []byte) string {
	if len(edid) < 128 {
		return ""
	}

	id, _ := utils.SumStrMd5(string(edid[:128]))
	return id
}
