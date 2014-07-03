package mounts

import (
	"fmt"
	"os"
)

func generateUUID() string {
	f, err := os.Open("/dev/urandom")
	if err != nil {
		fmt.Println(err)
		return ""
	}

	defer f.Close()
	b := make([]byte, 16)
	f.Read(b)
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6],
		b[6:8], b[8:10], b[10:])

	return uuid
}
