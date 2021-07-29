package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	shell "github.com/ipfs/go-ipfs-api"
)

type Stats struct {
	Bw StatsBw `json:"bw"`
}

type StatsBw struct {
	TotalIn  int64
	TotalOut int64
	RateIn   float64
	RateOut  float64
}

type RequestBody struct {
	Status       string   `json:"status"`
	Stats        Stats    `json:"stats"`
	Peers        int      `json:"peers"`
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Addresses    []string `json:"Addresses"`
	AgentVersion string   `json:"version"`
}

type ErrorBody struct {
	Status string `json:"status"`
	Name   string `json:"name"`
	Error  string `json:"message"`
}

type Driver struct {
	IPFSUrl    string
	PinsetURL  string
	ErrorURL   string
	PinsetFunc func(*Driver, io.Reader) *http.Request
	Elapse     time.Duration
	Name       string

	ticket  *time.Ticker
	context context.Context
}

type Pin struct {
	Cid string `json:"cid"`
}

type Pinset struct {
	Swarms []string `json:"swarms"`
	Pins   []Pin    `json:"pins"`
}

func makeRequest(d *Driver, body io.Reader) *http.Request {
	req, _ := http.NewRequest("POST", d.PinsetURL, body)
	return req
}

func (d *Driver) loop(ipfs *shell.Shell, httpClient *http.Client) {
	defer func() {
		err := recover()
		if err != nil {
			reportError(d, fmt.Errorf("PANIC: %v", err))
		}
	}()
	fmt.Println("[SENTINEL] starting check")
	checkStart := time.Now()

	id, err := ipfs.ID()
	if err != nil {
		reportError(d, err)
		return
	}

	peers, err := getPeers(d.context, ipfs)
	if err != nil {
		reportError(d, err)
		return
	}

	savedPins, err := ipfs.Pins()
	if err != nil {
		reportError(d, err)
		return
	}

	stats, err := getStats(d.context, ipfs)
	if err != nil {
		reportError(d, err)
		return
	}

	f := makeRequest
	if d.PinsetFunc != nil {
		f = d.PinsetFunc
	}
	request := RequestBody{
		Status:       "up",
		Name:         d.Name,
		ID:           id.ID,
		Peers:        len(peers),
		AgentVersion: id.AgentVersion,
		Addresses:    id.Addresses,
		Stats:        stats,
	}
	data, _ := json.MarshalIndent(request, "", "\t")
	fmt.Println(string(data))

	req := f(d, bytes.NewBuffer(data))
	req = req.WithContext(d.context)
	res, err := httpClient.Do(req)
	if err != nil {
		reportError(d, err)
		return
	}

	var pinset Pinset
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&pinset)
	if err != nil {
		reportError(d, err)
		return
	}

	if len(pinset.Swarms) > 0 {
		toConnect := make([]string, 0, len(pinset.Swarms))
		for _, pin := range pinset.Swarms {
			if _, ok := peers[pin]; !ok {
				toConnect = append(toConnect, pin)
			}
		}
		if len(toConnect) == 0 {
			fmt.Println("[SENTINEL] nothing to connect")
		} else {
			fmt.Printf("[SENTINEL] connecting with peers %d / %d\n", len(toConnect), len(pinset.Swarms))
			connectStart := time.Now()
			err = ipfs.SwarmConnect(d.context, toConnect...)
			if err != nil {
				reportError(d, err)
			}
			fmt.Printf("[SENTINEL] c completed (%v)\n", time.Since(connectStart))
		}
	}

	toPin := make([]string, 0, len(pinset.Pins))
	for _, pin := range pinset.Pins {
		_, ok := savedPins[pin.Cid]
		if !ok {
			toPin = append(toPin, pin.Cid)
		}
	}
	if len(toPin) == 0 {
		fmt.Println("[SENTINEL] nothing to pin")
	} else {
		fmt.Printf("[SENTINEL] new pins %d / %d\n", len(toPin), len(pinset.Pins))

		var wait sync.WaitGroup
		for _, cid := range toPin {
			wait.Add(1)
			go func(newCID string) {
				fmt.Println("[SENTINEL] adding pin", newCID)
				start := time.Now()
				err = ipfs.Request("pin/add", newCID).
					Option("recursive", false).
					Exec(d.context, nil)
				if err != nil {
					reportError(d, err)
				} else {
					fmt.Printf("[SENTINEL] pin completed %s (%v)\n", newCID, time.Since(start))
				}
				wait.Done()
			}(cid)
		}
		wait.Wait()
	}
	fmt.Printf("[SENTINEL] check completed (%v)\n", time.Since(checkStart))
}

func (d *Driver) Run() error {
	ipfs := shell.NewShell(d.IPFSUrl)
	up := ipfs.IsUp()
	if !up {
		return errors.New("can not connect to IPFS node")
	}

	fmt.Println("[SENTINEL] IPFS node is up")
	httpClient := &http.Client{}

	ticker := time.NewTicker(d.Elapse)
	d.ticket = ticker
	d.context = context.Background()

	for {
		d.loop(ipfs, httpClient)
		<-ticker.C
	}
}

func getPeers(ctx context.Context, ipfs *shell.Shell) (map[string]struct{}, error) {
	peers, err := ipfs.SwarmPeers(ctx)
	if err != nil {
		return nil, err
	}
	peerMap := make(map[string]struct{}, len(peers.Peers))
	for _, p := range peers.Peers {
		peerMap[p.Addr] = struct{}{}
	}
	return peerMap, nil
}

func getStats(ctx context.Context, ipfs *shell.Shell) (Stats, error) {
	bw, err := ipfs.StatsBW(ctx)
	if err != nil {
		return Stats{}, err
	}
	return Stats{
		Bw: StatsBw{
			TotalIn:  bw.TotalIn,
			TotalOut: bw.TotalOut,
			RateIn:   bw.RateIn,
			RateOut:  bw.RateOut,
		},
	}, nil
}

func (d *Driver) Stop() {
	d.ticket.Stop()
	<-d.context.Done()
}

func reportError(d *Driver, err error) {
	fmt.Fprintln(os.Stderr, "ERROR", err.Error())
	if d.ErrorURL != "" {
		data, _ := json.Marshal(ErrorBody{
			Status: "error",
			Name:   d.Name,
			Error:  err.Error(),
		})
		http.Post(d.ErrorURL, "application/json", bytes.NewBuffer(data))
	}
}