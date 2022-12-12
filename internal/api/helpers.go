package api

import (
	"errors"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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
	min := -1
	var result string
	for _, addr := range addresses {
		if serverStats[addr].Status && min == -1 {
			min = len(serverStats[addr].Queue)
			result = addr
		}
		if min > len(serverStats[addr].Queue) {
			min = len(serverStats[addr].Queue)
			result = addr
		}
	}
	if result == "" {
		return "", errors.New("what")
	}
	return result, nil
}

func GetConfigPath() string {
	pwd, err := os.Getwd()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	split := strings.Split(pwd, string(os.PathSeparator))
	if split[len(split)-1] == "test" {
		WorkDirPath = pwd
	} else {
		log.Println(split[len(split)-1])
		pwd = filepath.Dir(pwd)
		pwd = filepath.Dir(pwd)
		WorkDirPath = pwd
	}
	absPath := filepath.Join(WorkDirPath, "configs", "settings.yml")
	return absPath
}
