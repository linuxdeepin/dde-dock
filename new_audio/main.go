package main

import (
	"fmt"
)

func main() {
	s := &Sink{}
	s.index = 0
	s.forceUpdate()
	fmt.Println(s.Name)
	fmt.Println(s.Description)
	fmt.Println(s.Volume)
	a := &Audio{}
	a.forceUpdate()
}
