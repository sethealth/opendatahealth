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
	ipfsURL := os.Getenv("IPFS_URL")

	if pinsetURL == "" {
		panic("PINSET_URL is not defined")
	}
	fmt.Println("[SENTINEL] starting pinset sentinel", pinsetURL)
	fmt.Println("[SENTINEL] waiting for", WAIT)

	driver := &Driver{
		IPFSUrl:   ipfsURL,
		PinsetURL: pinsetURL,
		Elapse:    2 * time.Minute,
		Name:      os.Getenv("OPENDATA_NODE"),
	}
	for {
		err := driver.Run()
		if err != nil {
			fmt.Fprintln(os.Stderr, "SENTINEL ERROR:", err.Error())
		}
		time.Sleep(WAIT)
		fmt.Println("[SENTINEL] reconnecting in", WAIT)

	}
}
