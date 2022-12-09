package main

import (
	"errors"
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
var toSend chan *http.Request
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
		activeAddress[val.Url] = IsAlive(val.Url)
	}
	toSend = make(chan *http.Request, 128)
	http.HandleFunc("/send", Proxy)
	log.Println("Starting server")
	for i := 0; i < 40; i++ {
		go SendProxy(toSend, nil, &mutex)
	}
	err = http.ListenAndServe("127.0.0.1:3333", nil)
	if err != nil {
		log.Println(err)
	}
	log.Println("server shutdown")
}

func Proxy(w http.ResponseWriter, r *http.Request) {
	log.Println("Received message")
	toSend <- r
	io.WriteString(w, "ok")
}

func GetMin() (string, error) {
	min := -1
	var result string
	for _, addr := range Addresses {
		if activeAddress[addr] && min == -1 {
			min = len(queue[addr])
			result = addr
		}
		if min > len(queue[addr]) {
			min = len(queue[addr])
			result = addr
		}
	}
	if result == "" {
		log.Println("No response")
		return "", errors.New("what")
	}
	return result, nil
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
			target, err = GetMin()
			if err != nil {
				return
			}
			queue[target][val] = true
			delete(queue[addr], val)
		}
	}

	sum := 0
	for _, addr := range active {
		sum += len(queue[addr])
	}

	avg := sum / len(active)
	log.Println(avg)
}

func IsAlive(url string) bool {
	client := http.Client{Timeout: 1000 * time.Millisecond}
	req, err := http.NewRequest("Get", url, nil)
	if err != nil {
		log.Println("IsAlive", err)
		return false
	}
	r, err := client.Do(req)
	if err != nil {
		log.Println("IsAlive", err)
		return false
	}
	return r.StatusCode == 200
}

func Check() {
	for {
		for val := range activeAddress {
			activeAddress[val] = IsAlive(val)
		}
	}
}

func SendProxy(in, repeat chan *http.Request, mutex *sync.Mutex) error {
	client := http.Client{Timeout: 3000 * time.Millisecond}
	var reqIn *http.Request
	for {
		url, err := GetMin()
		if err != nil {
			continue
		}
		log.Println(url)

		reqIn = <-in
		queue[url][reqIn] = true
		req, err := http.NewRequest(reqIn.Method, url, nil)

		if err != nil {
			fmt.Println(err)
			continue
		}

		req.Header = reqIn.Header
		req.Body = reqIn.Body
		mutex.Lock()
		log.Println("size of requests", len(queue[url]), url)
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Proxy Do", err)
			in <- reqIn
			continue
		}
		if resp.StatusCode != 200 {
			log.Println("Proxy unreached", resp.StatusCode)
			activeAddress[url] = false
			in <- reqIn
			Balance()
		} else {
			activeAddress[url] = true
			delete(queue[url], reqIn)
		}
		mutex.Unlock()
	}
}
