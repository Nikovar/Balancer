package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
	"tmp/internal"
)

var Config internal.Config
var Addresses []string
var queue map[string]map[*http.Request]bool
var activeAddress map[string]bool
var mutex sync.Mutex

func main() {
	err := Config.GetConfig("settings.yml")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	activeAddress = make(map[string]bool, len(Config.Servers))
	queue = make(map[string]map[*http.Request]bool, len(Config.Servers))

	for _, val := range Config.Servers {
		Addresses = append(Addresses, val.Url)
		queue[val.Url] = make(map[*http.Request]bool)
		activeAddress[val.Url] = internal.IsAlive(val.Url)
	}

	http.HandleFunc("/send", Proxy)
	log.Println("Starting server")

	http.ListenAndServe("127.0.0.1:3333", nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("server shutdown")
}

func Balance() {
	var active []string
	var inactive []string
	var target string
	var err error

	for addr, ok := range activeAddress {
		if ok {
			active = append(active, addr)
		} else {
			inactive = append(inactive, addr)
		}
	}

	for _, addr := range inactive {
		for val := range queue[addr] {
			target, err = internal.GetMin(Addresses, queue, activeAddress)
			if err != nil {
				return
			}
			queue[target][val] = true
			delete(queue[addr], val)
		}
	}

}

func Check() {
	for {
		for val := range activeAddress {
			activeAddress[val] = internal.IsAlive(val)
		}
	}
}

func Proxy(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		resp *http.Response
		url  string
	)
	log.Println("Received message")
	//toSend <- r
	for {
		client := http.Client{Timeout: 1 * time.Second}
		url, err = internal.GetMin(Addresses, queue, activeAddress)
		if err != nil {
			continue
		}
		req, err := http.NewRequest(r.Method, url, nil)

		if err != nil {
			log.Println(err)
			continue
		}

		req.Header = r.Header
		req.Body = r.Body
		mutex.Lock()

		queue[url][r] = true
		resp, err = client.Do(req)
		if err != nil {
			log.Println("Proxy Do", err)
			activeAddress[url] = false
			mutex.Unlock()
			Balance()
			continue
		}
		activeAddress[url] = true
		delete(queue[url], r)
		mutex.Unlock()
		break
	}
	response := fmt.Sprintf("%v", resp)
	io.WriteString(w, response)
}
