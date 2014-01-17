package main

import (
	"fmt"
	"testing"
)

func TestDesktop(t *testing.T) {
	desk, err := NewDesktopManager()
	if err != nil {
		t.Errorf("create desktop failed\n")
		return
	}

	fmt.Println(desk)
	desk.SetBottomRightAction(2)
	desk.SetTopLeftAction(1)
}
