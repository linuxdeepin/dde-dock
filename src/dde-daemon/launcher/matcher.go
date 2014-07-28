package launcher

import (
	"fmt"
	"regexp"
	"strings"
)

func addMatcher(template string, key string, score uint32, m map[*regexp.Regexp]uint32) error {
	r, err := regexp.Compile(fmt.Sprintf(template, key))
	if err != nil {
		return err
	}

	m[r] = score
	return nil
}

const (
	POOR          uint32 = 50000
	BELOW_AVERAGE uint32 = 60000
	AVERAGE       uint32 = 70000
	ABOVE_AVERAGE uint32 = 75000
	GOOD          uint32 = 80000
	VERY_GOOD     uint32 = 85000
	EXCELLENT     uint32 = 90000
	HIGHEST       uint32 = 100000
)

// learn from synpase
// TODO:
// 1. analyse the code of synapse much deeply.
// 2. add a weight for frequency.
func getMatchers(key string) map[*regexp.Regexp]uint32 {

	// * create a couple of regexes and try to help with matching
	// * match with these regular expressions (with descending score):
	// * 1) ^query$
	// * 2) ^query
	// * 3) \bquery
	// * 4) split to words and seach \bword1.+\bword2 (if there are 2+ words)
	// * 5) query
	// * 6) split to characters and search \bq.+\bu.+\be.+\br.+\by
	// * 7) split to characters and search \bq.*u.*e.*r.*y
	m := make(map[*regexp.Regexp]uint32, 0)
	addMatcher(`(?i)^(%s)$`, key, HIGHEST, m)
	addMatcher(`(?i)^(%s)`, key, EXCELLENT, m)
	addMatcher(`(?i)\b(%s)`, key, VERY_GOOD, m)
	words := strings.Fields(key)
	if len(words) > 1 {
		addMatcher(`(?i)\b(%s)`, strings.Join(words, `).+\b(`), GOOD, m)
	}
	addMatcher("(?i)(%s)", key, BELOW_AVERAGE, m)
	charSpliter, err := regexp.Compile(`\s*`)
	if err != nil {
		logger.Warning("get char spliter failed:", err)
		return m
	}

	chars := charSpliter.Split(key, -1)
	if len(words) == 1 && len(chars) <= 5 {
		addMatcher(`(?i)\b(%s)`, strings.Join(chars, `).+\b(`),
			ABOVE_AVERAGE, m)
	}
	addMatcher(`(?i)\b(%s)`, strings.Join(chars, ").*("), POOR, m)

	return m
}
