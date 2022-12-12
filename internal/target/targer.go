package target

import (
	"io"
	"log"
	"net/http"
	"time"
)

var requests *int

func StartSender(port string) {
	requests = new(int)
	*requests = 0
	go Timer()
	http.HandleFunc("/", Index)
	log.Println("server on port", "127.0.0.1:"+port, "started")
	log.Println(http.ListenAndServe("127.0.0.1:"+port, nil))
}

func Timer() {
	for {
		select {
		case <-time.After(10 * time.Second):
			log.Printf("Requests operating :%d\n", *requests)
		}
	}
}

func Index(w http.ResponseWriter, r *http.Request) {
	*requests++
	log.Println("received")
	time.Sleep(120 * time.Millisecond)
	*requests--
	io.WriteString(w, "")
}
