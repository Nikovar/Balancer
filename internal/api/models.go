package api

import (
	"net/http"
)

var (
	WorkDirPath string
)

type Config struct {
	Servers []struct {
		Url string `yaml:"url"`
	} `yaml:"servers"`
}

type ServerProps struct {
	Url    string
	Status bool
	Queue  map[*http.Request]bool
}
