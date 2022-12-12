package api

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func Proxy(w http.ResponseWriter, r *http.Request, serverStats map[string]ServerProps) {
	var (
		err          error
		resp         *http.Response
		url          string
		addressSlice []string = make([]string, 0)
	)
	log.Println("Received message")
	for key := range serverStats {
		addressSlice = append(addressSlice, key)
	}

	for {
		client := http.Client{Timeout: 1 * time.Second}
		url, err = GetMinRef(addressSlice, serverStats)
		if err != nil {
			log.Println(err)
			continue
		}
		req, err := http.NewRequest(r.Method, url, nil)
		if err != nil {
			log.Println(err)
			continue
		}

		req.Header = r.Header
		req.Body = r.Body
		serverStats[url].Queue[r] = true
		resp, err = client.Do(req)
		log.Println("Message sent")
		if err != nil {
			serverStats[url] = ServerProps{Url: url, Status: false, Queue: serverStats[url].Queue}
			Balance(serverStats, addressSlice)
			continue
		}
		serverStats[url] = ServerProps{Url: url, Status: true, Queue: serverStats[url].Queue}
		delete(serverStats[url].Queue, r)
		break
	}
	response := fmt.Sprintf("%v", resp)
	io.WriteString(w, response)
}
