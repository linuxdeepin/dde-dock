package lastore

import (
	"testing"
)

func TestSourceLineParsed_String(t *testing.T) {
	a := &sourceLineParsed{
		options: map[string]string{
			"arch":    "amd64,armel",
			"option2": "abc",
		},
		url:        "http://deb.debian.org/debian",
		suite:      "stretch",
		components: []string{"contrib", "main", "non-free"},
	}
	if a.String() != "deb [arch=amd64,armel option2=abc] http://deb.debian.org/debian"+
		" stretch contrib main non-free" {
		t.Error("Unexpect not equal")
	}

	b := &sourceLineParsed{
		options:    nil,
		url:        "http://deb.debian.org/debian",
		suite:      "stretch",
		components: []string{"main"},
	}
	if b.String() != "deb http://deb.debian.org/debian stretch main" {
		t.Error("Unexpect not equal")
	}
}

func TestAptSourcesEqual(t *testing.T) {
	a := []*sourceLineParsed{
		{
			options: map[string]string{
				"arch":    "amd64,armel",
				"option2": "abc",
			},
			url:        "http://deb.debian.org/debian",
			suite:      "stretch",
			components: []string{"contrib", "main", "non-free"},
		},
		{
			options:    nil,
			url:        "http://deb.debian.org/debian",
			suite:      "stretch",
			components: []string{"main"},
		},
	}

	b := []*sourceLineParsed{
		{
			options: map[string]string{
				"arch":    "amd64,armel",
				"option2": "abc",
			},
			url:        "http://deb.debian.org/debian",
			suite:      "stretch",
			components: []string{"contrib", "main", "non-free"},
		},
		{
			options:    nil,
			url:        "http://deb.debian.org/debian",
			suite:      "stretch",
			components: []string{"main"},
		},
	}

	if !aptSourcesEqual(a, b) {
		t.Error("a and b should be equal")
	}
}

func TestParseSourceLine(t *testing.T) {
	line1 := []byte("deb [ arch=amd64,armel option2=abc ] http://deb.debian.org/debian" +
		" stretch main contrib non-free")
	parsed, err := parseSourceLine(line1)
	if err != nil {
		t.Error("Unexpected error:", err)
	}
	t.Logf("%#v\n", parsed)
	if !parsed.equal(&sourceLineParsed{
		options: map[string]string{
			"arch":    "amd64,armel",
			"option2": "abc",
		},
		url:        "http://deb.debian.org/debian",
		suite:      "stretch",
		components: []string{"contrib", "main", "non-free"},
	}) {
		t.Error("Unexpect not equal")
	}

	line2 := []byte("deb http://deb.debian.org/debian stretch main")
	parsed, err = parseSourceLine(line2)
	if err != nil {
		t.Error("Unexpoected error:", err)
	}
	t.Logf("%#v\n", parsed)
	if !parsed.equal(&sourceLineParsed{
		options:    nil,
		url:        "http://deb.debian.org/debian",
		suite:      "stretch",
		components: []string{"main"},
	}) {
		t.Error("Unexpect not equal")
	}

	line3 := []byte("deb")
	parsed, err = parseSourceLine(line3)
	if err == nil {
		t.Error("Expect error")
	}

	line4 := []byte("deb []")
	parsed, err = parseSourceLine(line4)
	if err == nil {
		t.Error("Expect error")
	}

	line5 := []byte("deb [] http://deb.debian.org")
	parsed, err = parseSourceLine(line5)
	if err == nil {
		t.Error("Expect error")
	}

	line6 := []byte("deb [] http://deb.debian.org stretch")
	parsed, err = parseSourceLine(line6)
	if err == nil {
		t.Error("Expect error")
	}

	line7 := []byte("deb [ http://deb.debian.org stretch")
	parsed, err = parseSourceLine(line7)
	if err == nil {
		t.Error("Expect error")
	}

	line8 := []byte("deb http-xxx://deb.debian.org stretch main")
	parsed, err = parseSourceLine(line8)
	if err == nil {
		t.Error("Expect error")
	}
}
