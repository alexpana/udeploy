package models

import (
	"gopkg.in/yaml.v2"
	"fmt"
)

type Config struct {
	Name         string `yaml:"name"`
	RunCommand   string `yaml:"runCommand"`
	PollInterval int    `yaml:"pollInterval"`
	Path         string
}

func ReadConfig(f string) Config {
	// TODO: implement

	var data = `
name: test
runCommand: 'python coffer.py'
pollInterval: 200
`
	var c Config
	err := yaml.Unmarshal([]byte(data), &c)

	if err != nil {
		fmt.Printf("Error occured: %v", err)
	}
	c.Path = f

	return c
}
