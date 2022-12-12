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
var mutex sync.Mutex
var ServerStats map[string]internal.ServerProps

func main() {
	err := Config.GetConfig("settings.yml")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	ServerStats = make(map[string]internal.ServerProps, len(Config.Servers))

	for _, val := range Config.Servers {
		Addresses = append(Addresses, val.Url)
		ServerStats[val.Url] = internal.ServerProps{
			Url: val.Url, Status: internal.IsAlive(val.Url), Queue: make(map[*http.Request]bool, 0)}
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

	for addr, val := range ServerStats {
		if val.Status {
			active = append(active, addr)
		} else {
			inactive = append(inactive, addr)
		}

	}

	for _, addr := range inactive {
		for val := range ServerStats[addr].Queue {
			target, err = internal.GetMinRef(Addresses, ServerStats)
			if err != nil {
				return
			}
			ServerStats[target].Queue[val] = true
			delete(ServerStats[addr].Queue, val)
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
		url, err = internal.GetMinRef(Addresses, ServerStats)
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

		ServerStats[url].Queue[r] = true
		resp, err = client.Do(req)
		if err != nil {
			log.Println("Proxy Do", err)
			ServerStats[url] = internal.ServerProps{Url: url, Status: false, Queue: ServerStats[url].Queue}
			mutex.Unlock()
			Balance()
			continue
		}
		ServerStats[url] = internal.ServerProps{Url: url, Status: true, Queue: ServerStats[url].Queue}
		delete(ServerStats[url].Queue, r)
		mutex.Unlock()
		break
	}
	response := fmt.Sprintf("%v", resp)
	io.WriteString(w, response)
}
