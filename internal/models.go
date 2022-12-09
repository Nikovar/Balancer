package internal

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Servers []struct {
		Url string `yaml:"url"`
	} `yaml:"servers"`
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
