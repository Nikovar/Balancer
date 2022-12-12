package sender

import (
	"log"
	"net/http"
	"sync"
	"time"
)

func Sender() {
	var wt sync.WaitGroup
	for i := 0; i < 10; i++ {
		wt.Add(1)
		go DDOS()
	}
	wt.Wait()
}

func DDOS() {
	for {
		_, err := http.Get("http://127.0.0.1:3333/send")
		if err != nil {
			log.Println(err)
		}
		log.Println("Sending message")
		time.Sleep(64 * time.Millisecond)
	}
}
