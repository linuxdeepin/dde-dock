package audio

/*
#cgo pkg-config: alsa
#include <stdio.h>
#include <stdlib.h>
#include <alsa/asoundlib.h>

// if *control_id != NULL, need free *control_id
static int find_control(int card_num, char *name, char **control_id) {
    snd_hctl_t *handle;
    char card[8];
    snprintf(card, 8, "hw:%d", card_num);
    int err;
    err = snd_hctl_open(&handle, card, 0);
    if (err < 0) {
        fprintf(stderr, "Control %s open error: %s\n", card, snd_strerror(err));
        return err;
    }

    err = snd_hctl_load(handle);
    if (err < 0) {
        fprintf(stderr, "Control %s load error: %s\n", card, snd_strerror(err));
        return err;
    }

    snd_hctl_elem_t *elem;
    snd_ctl_elem_id_t *id;
    snd_ctl_elem_id_alloca(&id);
    for (elem = snd_hctl_first_elem(handle); elem != NULL; elem = snd_hctl_elem_next(elem)) {
        snd_hctl_elem_get_id(elem, id);
        char *id_str = snd_ctl_ascii_elem_id_get(id);
        if (id_str != NULL) {
            if (strstr(id_str, name) != NULL) {
                *control_id = id_str;
                break;
            } else {
                free(id_str);
            }
        }
    }

    snd_hctl_close(handle);
    return 0;
}

// if *control_id != NULL, need free *control_id
static int find_card_control(char *name, int *card_num, char **control_id) {
    *card_num = -1;
    int err;
    while(1) {
        err = snd_card_next(card_num);
        if (err < 0) {
            fprintf(stderr, "card next err: %s\n", snd_strerror(err));
            return err;
        }
        if (*card_num == -1) {
            //  no more cards are available.
            break;
        }

        err = find_control(*card_num, name, control_id);
        if (err < 0) {
            return err;
        }

        if (*control_id != NULL) {
            break;
        }
    }
    return 0;
}
*/
import "C"
import (
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"unsafe"
)

var errNotFoundControl = errors.New("not found control")

func findCardControl(name string) (card int, controlId string, err error) {
	var cCard C.int
	var cControlId *C.char
	cName := C.CString(name)

	cErr := C.find_card_control(cName, &cCard, &cControlId)
	C.free(unsafe.Pointer(cName))

	if cErr < 0 {
		err = fmt.Errorf("find_card_control err code: %d", err)
		return
	}

	if cControlId != nil {
		card = int(cCard)
		controlId = C.GoString(cControlId)
		C.free(unsafe.Pointer(cControlId))
		return
	}

	err = errNotFoundControl
	return
}

func disableAutoMuteMode() error {
	cardNum, controlId, err := findCardControl("Auto-Mute Mode")
	if err != nil {
		if err == errNotFoundControl {
			err = nil
		}
		return err
	}
	out, err := exec.Command("amixer", "-c", strconv.Itoa(cardNum), "cset", controlId, "Disabled").Output()
	if err != nil {
		return err
	}
	logger.Debugf("command amixer -c %d cset %s Disabled\n out: %s", cardNum, controlId, out)
	return nil
}
