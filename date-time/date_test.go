package main

import (
	"fmt"
	"testing"
)

func TestDateTime (t *testing.T) {
	date := NewDateAndTime()
	if date == nil {
		t.Errorf("create date time failed\n")
		return
	}

	fmt.Println(date.AutoSetTime)
	fmt.Println(date.CurrentTimeZone)
	fmt.Println(date.TimeShowFormat)

	date.SetAutoSetTime(true)
	date.SetDate("2013-11-27")
	date.SetTimeZone("Asia/Shanghai")
	date.SetTime("12:11:11")
	date.SyncNtpTime()
	date.SetAutoSetTime(false)
}
