package main

import (
	"log"
	"net/http"
	"os"
	"runtime"
	"tmp/internal/api"
)

var Config api.Config
var Addresses []string
var ServerStats map[string]api.ServerProps

func main() {
	runtime.GOMAXPROCS(10)
	err := Config.GetConfig(api.GetConfigPath())
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	ServerStats = make(map[string]api.ServerProps, len(Config.Servers))

	for _, val := range Config.Servers {
		Addresses = append(Addresses, val.Url)
		ServerStats[val.Url] = api.ServerProps{
			Url: val.Url, Status: api.IsAlive(val.Url),
			Queue: map[*http.Request]bool{}}
	}

	http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		api.Proxy(w, r, ServerStats)
	})
	log.Println("Starting server")

	http.ListenAndServe("127.0.0.1:3333", nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("server shutdown")
}
