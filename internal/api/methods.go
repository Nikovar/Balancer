package api

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

func Balance(serverStats map[string]ServerProps, addresses []string) {
	var active []string
	var inactive []string
	var target string
	var err error

	for addr, val := range serverStats {
		if val.Status {
			active = append(active, addr)
		} else {
			inactive = append(inactive, addr)
		}

	}

	for _, addr := range inactive {
		for val := range serverStats[addr].Queue {
			target, err = GetMinRef(addresses, serverStats)
			if err != nil {
				return
			}
			serverStats[target].Queue[val] = true
			delete(serverStats[addr].Queue, val)
		}
	}
}

func (c *Config) GetConfig(configPath string) error {
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
