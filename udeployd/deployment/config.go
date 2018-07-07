package deployment

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	Name         string   `yaml:"name"`         // human readable name for the process
	RunCommand   string   `yaml:"runCommand"`   // command used to run the process
	RunArgs      []string `yaml:"runArgs"`      // command used to run the process
	PollInterval int      `yaml:"pollInterval"` // frequency at which to poll the process for liveness
	Path         string                         // the path of the config
}

func ReadConfig(f string) (Config, error) {
	dat, err := ioutil.ReadFile(f)
	if err != nil {
		return Config{}, err
	}

	var c Config
	err = yaml.Unmarshal(dat, &c)

	if err != nil {
		log.Printf("Error occured while reading config file %s: %v", f, err)
	}
	c.Path = f

	return c, nil
}
