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

	"gopkg.in/yaml.v2"
)

var Config config
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
	toSend = make(chan *http.Request, 999)
	//goBalance()
	http.HandleFunc("/send", Proxy)
	log.Println("Starting server")
	//for i := 0; i < 10; i++ {
	go SendProxy(toSend, nil, &mutex)
	//}
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
	//log.Println(url, r.StatusCode)
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
	client := http.Client{Timeout: 500 * time.Millisecond}
	var prevreq *http.Request = nil
	var reqIn *http.Request
	for {
		//time.Sleep(100 * time.Millisecond)
		//mutex.Lock()
		for val := range activeAddress {
			activeAddress[val] = IsAlive(val)
		}
		//mutex.Unlock()
		url, err := GetMin()
		if err != nil {
			continue
			//return err
		}

		if prevreq == nil {
			reqIn = <-in
		} else {
			reqIn = prevreq
		}
		//log.Println(reqIn)

		//log.Println("Sending to ", url)

		req, err := http.NewRequest(reqIn.Method, url, nil)

		if err != nil {
			fmt.Println(err)
			continue
			//return err
		}

		req.Header = reqIn.Header
		req.Body = reqIn.Body
		queue[url][reqIn] = true
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Proxy Do", err)
			prevreq = reqIn
			continue
			//return err
		}
		//mutex.Lock()
		if resp.StatusCode != 200 {
			log.Println("Proxy unreached", resp.StatusCode)
			activeAddress[url] = false
			Balance()
			prevreq = reqIn
		} else {
			activeAddress[url] = true
			delete(queue[url], reqIn)
			prevreq = nil
		}
		//mutex.Unlock()
		//log.Println(activeAddress)
	}
}

func (c *config) GetConfig(configPath string) error {
	file, err := os.Open(configPath)
	if err != nil {
		return err
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	d := yaml.NewDecoder(file)

	if err := d.Decode(&c); err != nil {
		return err
	}
	return nil
}

type config struct {
	Servers []struct {
		Url string `yaml:"url"`
	} `yaml:"servers"`
}
