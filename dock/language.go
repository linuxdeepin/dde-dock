// Port from glib(gkeyfile.c and gcharset.c).
// because we cannot pass NULL to glib.KeyFile.GetLocaleString like C, so the
// locale must be passed expclitly.
package dock

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

const (
	COMPONENT_CODESET = 1 << iota
	COMPONENT_MODIFIER
	COMPONENT_TERRITORY
)

var (
	alias_table map[string]string = nil
	said_before                   = false
	splitor                       = regexp.MustCompile(`\s+|:`)
)

func readAliases(filename string) {
	if alias_table == nil {
		alias_table = make(map[string]string, 0)
	}

	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	s := bufio.NewScanner(file)
	s.Split(bufio.ScanLines)
	for s.Scan() {
		lineContent := strings.TrimSpace(s.Text())
		if len(lineContent) == 0 || lineContent == "#" {
			continue
		}

		content := splitor.Split(lineContent, -1)
		if len(content) != 2 {
			continue
		}
		alias_table[content[0]] = alias_table[content[1]]
	}
}

func unaliasLang(lang string) string {
	if alias_table == nil {
		readAliases("/usr/share/locale/locale.alias")
	}

	i := 0
	for p, ok := alias_table[lang]; ok && p == lang; i++ {
		lang = p
		if i == 30 {
			if !said_before {
				logger.Warning("Too many alias levels for a locale, may indicate a loop")
			}
			said_before = true
			return lang
		}
	}

	return lang
}

func explodeLocale(locale string) (mask uint, language string, territory string, codeset string, modifier string) {
	mask = uint(0)
	uscore_pos := strings.IndexRune(locale, '_')
	dot_pos := strings.IndexRune(locale, '.')
	at_pos := strings.IndexRune(locale, '@')

	if at_pos != -1 {
		mask |= COMPONENT_MODIFIER
		modifier = locale[at_pos:]
		logger.Debug("modifier", modifier)
	} else {
		at_pos = len(locale)
	}

	if dot_pos != -1 {
		mask |= COMPONENT_CODESET
		codeset = locale[dot_pos:at_pos]
		logger.Debug("codeset", codeset)
	} else {
		dot_pos = at_pos
	}

	if uscore_pos != -1 {
		mask |= COMPONENT_TERRITORY
		territory = locale[uscore_pos:dot_pos]
		logger.Debug("territory", territory)
	} else {
		uscore_pos = dot_pos
	}

	language = locale[:uscore_pos]
	return
}

func GetLocaleVariants(locale string) []string {
	array := make([]string, 0)
	mask, language, territory, codeset, modifier := explodeLocale(locale)

	for j := uint(0); j <= mask; j++ {
		i := mask - j
		if (i & ^mask) == 0 {
			val := language
			if (i & COMPONENT_TERRITORY) != 0 {
				val = val + territory
			}

			if (i & COMPONENT_CODESET) != 0 {
				val = val + codeset
			}

			if (i & COMPONENT_MODIFIER) != 0 {
				val = val + modifier
			}

			logger.Debug("append", val)
			array = append(array, val)
		}
	}

	return array
}

func guessCategoryValue(categoryName string) (retval string) {
	retval = os.Getenv("LANGUAGE")
	if retval != "" {
		return
	}

	retval = os.Getenv("LC_ALL")
	if retval != "" {
		return
	}

	retval = os.Getenv(categoryName)
	if retval != "" {
		return
	}

	retval = os.Getenv("LANG")
	if retval != "" {
		return
	}

	retval = "C"
	return
}

func GetLanguageNames() []string {
	val := guessCategoryValue("LC_MESSAGES")
	langs := strings.Split(val, ":")
	logger.Debug(langs)
	array := make([]string, 0)
	for _, lang := range langs {
		array = GetLocaleVariants(unaliasLang(lang))
	}

	array = append(array, "C")

	return array
}
