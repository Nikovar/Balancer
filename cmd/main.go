package main

import (
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
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

	// 	u, _ := url.Parse("http://localhost:8080")
	// rp := httputil.NewSingleHostReverseProxy(u)

	// http.HandlerFunc(rp.ServeHTTP)
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
	toSend = make(chan *http.Request, 128)
	http.HandleFunc("/send", Proxy)
	log.Println("Starting server")
	for i := 0; i < 40; i++ {
		go SendProxy(toSend, nil, &mutex)

	}
	u, _ := url.Parse("http://127.0.0.1:8080")
	rp := httputil.NewSingleHostReverseProxy(u)
	server := http.Server{Addr: "127.0.0.1:3333",
		Handler: rp}
	err = server.ListenAndServe()
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
	log.Println("Received message")
	toSend <- r
	io.WriteString(w, "ok")
}

func SendProxy(in, repeat chan *http.Request, mutex *sync.Mutex) error {
	client := http.Client{Timeout: 3000 * time.Millisecond}
	var reqIn *http.Request
	for {
		url, err := internal.GetMin(Addresses, queue, activeAddress)
		if err != nil {
			continue
		}

		reqIn = <-in
		queue[url][reqIn] = true
		req, err := http.NewRequest(reqIn.Method, url, nil)

		if err != nil {
			log.Println(err)
			continue
		}

		req.Header = reqIn.Header
		req.Body = reqIn.Body
		mutex.Lock()

		_, err = client.Do(req)
		if err != nil {
			log.Println("Proxy Do", err)
			activeAddress[url] = false
			mutex.Unlock()
			in <- reqIn
			Balance()
			continue
		}
		activeAddress[url] = true
		delete(queue[url], reqIn)
		mutex.Unlock()
	}
}
