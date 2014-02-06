package main

import (
	"fmt"
	"time"
)

type Result struct {
	when  time.Time
	value bool
}

func main() {

	channel := make(chan Result)

	go pinger(channel)
	go reporter(channel)

	// sleep forever
	select {}
}

func pinger(reporter chan Result) {
	c := time.Tick(1 * time.Second)
	for now := range c {
		reporter <- Result{when: now, value: false}
	}
}

func reporter(source chan Result) {
	for {
		result := <-source
		fmt.Println(result.when.Format("03:04:05"))
	}
}
