package lastore

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	liburl "net/url"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/gettext"
	"pkg.deepin.io/lib/xdg/basedir"
)

const (
	aptSourcesFile       = "/etc/apt/sources.list"
	aptSourcesOriginFile = aptSourcesFile + ".origin"
)

var disableSourceCheckFile = filepath.Join(basedir.GetUserConfigDir(), "deepin",
	"lastore-session-helper", "disable-source-check")

func (l *Lastore) checkSource() {
	const sourceCheckedFile = "/tmp/lastore-session-helper-source-checked"

	if _, err := os.Stat(sourceCheckedFile); os.IsNotExist(err) {
		ok := doCheckSource()
		logger.Info("checkSource:", ok)
		if !ok {
			l.notifySourceModified(l.createSourceModifiedActions())
		}

		err = touchFile(sourceCheckedFile)
		if err != nil {
			logger.Warning("failed to touch source-checked file:", err)
		}
	}
}

func (l *Lastore) createSourceModifiedActions() []NotifyAction {
	return []NotifyAction{
		{
			Id:   "restore",
			Name: gettext.Tr("Restore"),
			Callback: func() {
				logger.Info("restore source")
				err := l.core.RestoreSystemSource(0)
				if err != nil {
					logger.Warningf("failed to restore source:", err)
				}
			},
		},
		{
			Id:   "Cancel",
			Name: gettext.Tr("Cancel"),
			Callback: func() {
				logger.Info("cancel restore source")
			},
		},
	}
}

const (
	propNameSourceCheckEnabled = "SourceCheckEnabled"
)

func (l *Lastore) emitPropChangedSourceCheckEnabled(value bool) {
	l.service.EmitPropertyChanged(l, propNameSourceCheckEnabled, value)
}

func (l *Lastore) SetSourceCheckEnabled(val bool) *dbus.Error {
	err := l.setSourceCheckEnabled(val)
	return dbusutil.ToError(err)
}

func (l *Lastore) setSourceCheckEnabled(val bool) error {
	if l.SourceCheckEnabled == val {
		return nil
	}

	if val {
		// enable
		err := os.Remove(disableSourceCheckFile)
		if err != nil {
			return err
		}
	} else {
		// disable
		err := touchFile(disableSourceCheckFile)
		if err != nil {
			return err
		}
	}

	l.SourceCheckEnabled = val
	l.emitPropChangedSourceCheckEnabled(val)
	return nil
}

// return is source ok?
func doCheckSource() bool {
	originSources, err := loadAptSources(aptSourcesOriginFile)
	if err != nil {
		logger.Warning("failed to load origin apt sources:", err)
		return true
	}

	currentSources, err := loadAptSources(aptSourcesFile)
	if err != nil {
		logger.Warning("failed to load current apt sources:", err)
		return false
	}

	logger.Debug("origin sources:", sourcesToString(originSources))
	logger.Debug("current sources:", sourcesToString(currentSources))
	if !aptSourcesEqual(originSources, currentSources) {
		return false
	}

	return true
}

func aptSourcesEqual(a, b []*sourceLineParsed) bool {
	if len(a) != len(b) {
		return false
	}

	for i := 0; i < len(a); i++ {
		if !a[i].equal(b[i]) {
			return false
		}
	}
	return true
}

var debReg = regexp.MustCompile(`deb\s+(\[(.*)])?(.+)`)

type sourceLineParsed struct {
	options    map[string]string
	url        string
	suite      string
	components []string
}

func (s *sourceLineParsed) String() string {
	var options string
	if len(s.options) > 0 {
		var optStrSlice []string
		for key, value := range s.options {
			optStrSlice = append(optStrSlice, key+"="+value)
		}
		sort.Strings(optStrSlice)
		options = "[" + strings.Join(optStrSlice, " ") + "] "
	}
	return fmt.Sprintf("deb %s%s %s %s", options, s.url, s.suite,
		strings.Join(s.components, " "))
}

func sourcesToString(sources []*sourceLineParsed) string {
	var stringSlice []string
	for _, source := range sources {
		stringSlice = append(stringSlice, source.String())
	}
	return strings.Join(stringSlice, ", ")
}

func (s *sourceLineParsed) equal(other *sourceLineParsed) bool {
	return reflect.DeepEqual(s, other)
}

var errInvalidSourceLine = errors.New("invalid source line")

func parseSourceLine(src []byte) (*sourceLineParsed, error) {
	matchResult := debReg.FindSubmatch(src)
	if matchResult == nil {
		return nil, errInvalidSourceLine
	}
	var optionMap map[string]string
	options := matchResult[2]
	optionsFields := bytes.Fields(options)

	if len(optionsFields) > 0 {
		optionMap = make(map[string]string)
	}

	for _, option := range optionsFields {
		optionParts := bytes.SplitN(option, []byte{'='}, 2)
		if len(optionParts) != 2 {
			return nil, errInvalidSourceLine
		}
		optionMap[string(optionParts[0])] = string(optionParts[1])
	}

	other := matchResult[3]
	otherFields := bytes.Fields(other)
	if len(otherFields) < 3 {
		return nil, errInvalidSourceLine
	}
	url := string(otherFields[0])
	u, err := liburl.Parse(url)
	if err != nil {
		return nil, err
	}
	switch u.Scheme {
	case "file", "cdrom", "http", "https", "ftp", "copy", "rsh", "ssh":
		// pass
	default:
		return nil, fmt.Errorf("invalid url scheme %q", u.Scheme)
	}

	suite := string(otherFields[1])
	var components []string
	for _, component := range otherFields[2:] {
		components = append(components, string(component))
	}
	sort.Strings(components)

	return &sourceLineParsed{
		options:    optionMap,
		url:        url,
		suite:      suite,
		components: components,
	}, nil
}

func loadAptSources(filename string) ([]*sourceLineParsed, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer f.Close()
	var result []*sourceLineParsed
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := bytes.TrimSpace(scanner.Bytes())
		if bytes.HasPrefix(line, []byte{'#'}) ||
			bytes.HasPrefix(line, []byte("deb-src")) ||
			len(line) == 0 {
			// ignore comment, deb-src and empty line
			continue
		}
		parsed, err := parseSourceLine(line)
		if err != nil {
			return nil, err
		}
		result = append(result, parsed)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return result, nil
}
