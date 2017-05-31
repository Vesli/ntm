package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Config configuration structure
type Config struct {
	URL        string
	Units      string
	AppID      string
	DBName     string
	DBUser     string
	DBPassword string
	DBURL      string
	FBURL      string
	FBParams   string
	Port       int
}

// LoadConfig by reading config file and unmarshal data to struct
func LoadConfig(pathToConfig string) (*Config, error) {
	data, err := ioutil.ReadFile(pathToConfig)
	if err != nil {
		return nil, fmt.Errorf("File config error: %s", err)
	}

	conf := &Config{}
	if err = json.Unmarshal(data, conf); err != nil {
		return nil, err
	}
	return conf, nil
}
