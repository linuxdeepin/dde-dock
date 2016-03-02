package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestFilterList(t *testing.T) {
	var infos = []struct {
		origin    []string
		condition []string
		ret       []string
	}{
		{
			origin:    []string{"power", "audio", "dock"},
			condition: []string{"power", "dock"},
			ret:       []string{"audio"},
		},
		{
			origin:    []string{"power", "audio", "dock"},
			condition: []string{},
			ret:       []string{"power", "audio", "dock"},
		},
		{
			origin:    []string{"power", "audio", "dock"},
			condition: []string{"power", "dock", "audio"},
			ret:       []string(nil),
		},
	}

	Convey("Test filterList", t, func() {
		for _, info := range infos {
			So(filterList(info.origin, info.condition),
				ShouldResemble, info.ret)
		}
	})
}
