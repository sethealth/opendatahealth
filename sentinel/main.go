package main

import (
	"fmt"
	"os"
	"time"
)

const (
	WAIT = 10 * time.Second
)

func main() {
	pinsetURL := os.Getenv("PINSET_URL")
	if pinsetURL == "" {
		panic("PINSET_URL is not defined")
	}
	fmt.Println("[SENTINEL] starting pinset follower", pinsetURL)
	fmt.Println("[SENTINEL] waiting for", WAIT)

	driver := &Driver{
		IPFSUrl:   "localhost:5001",
		PinsetURL: pinsetURL,
		Elapse:    5 * time.Minute,
	}
	for {
		time.Sleep(WAIT)
		err := driver.Run()
		if err != nil {
			fmt.Fprintln(os.Stderr, "SENTINEL ERROR:", err.Error())
		}
		fmt.Println("[SENTINEL] reconnecting in", WAIT)

	}
}
