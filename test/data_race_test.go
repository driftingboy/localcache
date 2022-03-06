package test

import (
	"fmt"
	"testing"
)

func Test_Race(t *testing.T) {
	c := make(chan bool)
	m := make(map[string]string)
	go func() {
		m["1"] = "a" // First conflicting access.
		c <- true
	}()
	<-c          // 满足 happend before
	m["2"] = "b" // Second conflicting access.
	// <-c // data race panic
	for k, v := range m {
		fmt.Println(k, v)
	}
}
