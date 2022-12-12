package internal

import (
	"errors"
	"log"
	"net/http"
	"time"
)

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

func GetMin(addresses []string, queue map[string]map[*http.Request]bool, activeAddress map[string]bool) (string, error) {
	min := -1
	var result string
	for _, addr := range addresses {
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

//////

func GetMinRef(addresses []string, serverStats map[string]ServerProps) (string, error) {

	return "", nil
}
