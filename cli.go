package main

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Result struct {
	when  time.Time
	value bool
}

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Please pass in the host to ping")
		os.Exit(1)
	}

	host := os.Args[1]

	channel := make(chan Result)
	go pinger(channel, host)
	go reporter(channel)

	// sleep forever
	select {}
}

func pinger(reporter chan Result, host string) {

	// every second
	for now := range time.Tick(1 * time.Second) {

		// run the ping command
		err := exec.Command("ping", "-W", "1", "-c", "1", host).Run()

		// report the result
		reporter <- Result{when: now, value: err == nil}
	}
}

func reporter(source chan Result) {

	var good int
	var seconds int
	var was time.Time = time.Now()

	for {
		result := <-source
		if result.value {
			good += 1
		}
		seconds += 1

		if seconds == 1 {
			fmt.Print(result.when.Format("15:04"), " |")
		}

		if result.value {
			fmt.Print("#")
		} else {
			fmt.Print(".")
		}

		if seconds >= 60 {
			summarize(good)

			seconds = 0
			good = 0

			if was.Day() != result.when.Day() {
				fmt.Println(result.when.Format("Mon 2, 2006"))
			}
			was = result.when
		}
	}
}

func summarize(good int) {

	last_10 := math.Max(0, float64(good-50))
	first_50 := (math.Min(0, float64(good-50)) + 50) / 10

	fmt.Printf("|%s%s%s%s| %d\n", strings.Repeat("#", int(first_50)), strings.Repeat(".", 5-int(first_50)), strings.Repeat("%", int(last_10)), strings.Repeat(".", 10-int(last_10)), good)
}
