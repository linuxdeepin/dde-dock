package ddcci

// #cgo pkg-config: ddcutil
// #include <ddcutil_c_api.h>
import "C"
import (
	"fmt"
	"sync"
	"unsafe"

	"pkg.deepin.io/lib/utils"
)

type ddcci struct {
	listPointer *C.DDCA_Display_Info_List
	listMu      sync.Mutex

	displayMap map[string]int
}

const (
	brightnessVCP = 0x10
)

func newDDCCI() (*ddcci, error) {
	ddc := &ddcci{
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

func (d *ddcci) freeList() {
	if d.listPointer != nil {
		C.ddca_free_display_info_list(d.listPointer)
		d.listPointer = nil
	}
}

func (d *ddcci) RefreshDisplays() error {
	d.listMu.Lock()
	defer d.listMu.Unlock()

	d.freeList()

	status := C.ddca_get_display_info_list2(C.bool(true), &d.listPointer)
	if status != C.int(0) {
		return fmt.Errorf("failed to get display info list: %d", status)
	}

	for i := 0; i < int(d.listPointer.ct); i++ {
		err := d.initDisplay(i)
		if err != nil {
			logger.Warning(err)
		}
	}

	return nil
}

func (d *ddcci) initDisplay(idx int) error {
	item := d.getDisplayInfoByIdx(idx)

	var handle C.DDCA_Display_Handle
	status := C.ddca_open_display2(item.dref, C.bool(true), &handle)
	if status != C.int(0) {
		return fmt.Errorf("failed to open monitor: %d", status)
	}

	defer C.ddca_close_display(handle)

	var val C.DDCA_Non_Table_Vcp_Value
	status = C.ddca_get_non_table_vcp_value(handle, brightnessVCP, &val)
	if status != C.int(0) {
		return fmt.Errorf("failed to check DDC/CI support: %d", status)
	}

	edid := C.GoBytes(unsafe.Pointer(&item.edid_bytes), 128)
	edidChecksum := getEDIDChecksum(edid)

	d.displayMap[string(edidChecksum)] = idx
	return nil
}

func (d *ddcci) SupportBrightness(edidChecksum string) bool {
	d.listMu.Lock()
	_, ok := d.displayMap[edidChecksum]
	d.listMu.Unlock()

	return ok
}

func (d *ddcci) GetBrightness(edidChecksum string) (brightness int, err error) {
	d.listMu.Lock()
	defer d.listMu.Unlock()

	idx, ok := d.displayMap[edidChecksum]
	if !ok {
		err = fmt.Errorf("monitor not support DDC/CI")
		return
	}

	item := d.getDisplayInfoByIdx(idx)

	var handle C.DDCA_Display_Handle
	status := C.ddca_open_display2(item.dref, C.bool(true), &handle)
	if status != C.int(0) {
		err = fmt.Errorf("failed to open monitor: %d", status)
		return
	}

	defer C.ddca_close_display(handle)

	var val C.DDCA_Non_Table_Vcp_Value
	status = C.ddca_get_non_table_vcp_value(handle, brightnessVCP, &val)
	if status != C.int(0) {
		err = fmt.Errorf("failed to get brightness: %d", status)
		return
	}

	brightness = int(val.sl)
	return
}

func (d *ddcci) SetBrightness(edidChecksum string, percent int) error {
	d.listMu.Lock()
	defer d.listMu.Unlock()

	idx, ok := d.displayMap[edidChecksum]
	if !ok {
		return fmt.Errorf("monitor not support DDC/CI")
	}

	item := d.getDisplayInfoByIdx(idx)

	var handle C.DDCA_Display_Handle
	status := C.ddca_open_display2(item.dref, C.bool(true), &handle)
	if status != C.int(0) {
		return fmt.Errorf("failed to open monitor: %d", status)
	}

	defer C.ddca_close_display(handle)

	// 开启结果验证，防止返回设置成功，但实际上没有生效的情况
	// 此方法仅对当前线程生效
	C.ddca_enable_verify(true)

	status = C.ddca_set_non_table_vcp_value(handle, brightnessVCP, 0, C.uchar(percent))
	if status != C.int(0) {
		return fmt.Errorf("failed to set brightness via DDC/CI: %d", status)
	}

	return nil
}

func (d *ddcci) getDisplayInfoByIdx(idx int) *C.DDCA_Display_Info {
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
