package grub2

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/polkit"
	"strconv"
	"strings"
)

func loadJSON(file string, v interface{}) error {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	return json.Unmarshal(content, v)
}

func saveJSON(file string, v interface{}) error {
	const dirMode = 0755
	const fileMode = 0644
	err := os.MkdirAll(filepath.Dir(file), dirMode)
	if err != nil {
		return err
	}

	content, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, content, fileMode)
}

func quoteString(str string) string {
	return strconv.Quote(str)
}

type InvalidResoultionError struct {
	Resolution string
}

func (err InvalidResoultionError) Error() string {
	return fmt.Sprintf("invalid resolution %q", err.Resolution)
}

func parseResolution(v string) (w, h uint16, err error) {
	if v == "auto" {
		err = errors.New("unknown auto")
		return
	}

	arr := strings.Split(v, "x")
	if len(arr) != 2 {
		err = InvalidResoultionError{v}
		return
	}
	// parse width
	tmpw, err := strconv.ParseUint(arr[0], 10, 16)
	if err != nil {
		err = InvalidResoultionError{v}
		return
	}

	// parse height
	tmph, err := strconv.ParseUint(arr[1], 10, 16)
	if err != nil {
		err = InvalidResoultionError{v}
		return
	}

	w = uint16(tmpw)
	h = uint16(tmph)

	if w == 0 || h == 0 {
		err = InvalidResoultionError{v}
		return
	}

	return
}

func checkResolution(v string) error {
	if v == "auto" {
		return nil
	}

	_, _, err := parseResolution(v)
	return err
}

func getStringIndexInArray(a string, list []string) int {
	for i, b := range list {
		if b == a {
			return i
		}
	}
	return -1
}

func isStringInArray(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

var noCheckAuth bool

func init() {
	if os.Getenv("NO_CHECK_AUTH") == "1" {
		noCheckAuth = true
		return
	}

	polkit.Init()
}

func Tr(str string) string {
	return str
}

func checkAuthWithPid(pid uint32) (bool, error) {
	subject := polkit.NewSubject(polkit.SubjectKindUnixProcess)
	subject.SetDetail("pid", pid)
	subject.SetDetail("start-time", uint64(0))
	const actionId = DBusDest
	details := make(map[string]string)
	details["polkit.gettext_domain"] = "dde-daemon"
	details["polkit.message"] = Tr("Authentication is required to change the grub2 configuration")
	result, err := polkit.CheckAuthorization(subject, actionId, details,
		polkit.CheckAuthorizationFlagsAllowUserInteraction, "")
	if err != nil {
		return false, err
	}

	return result.IsAuthorized, nil
}

var errAuthFailed = errors.New("authentication failed")

func checkAuth(dbusMsg dbus.DMessage) error {
	if noCheckAuth {
		logger.Warning("check auth disabled")
		return nil
	}

	pid := dbusMsg.GetSenderPID()
	isAuthorized, err := checkAuthWithPid(pid)
	if err != nil {
		return err
	}
	if !isAuthorized {
		return errAuthFailed
	}
	return nil
}

func getFileMD5sum(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := md5.New()
	_, err = io.Copy(h, f)
	if err != nil {
		return "", err
	}

	sum := fmt.Sprintf("%x", h.Sum(nil))
	return sum, nil
}
