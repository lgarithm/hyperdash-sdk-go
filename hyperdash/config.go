package hyperdash

import (
	"encoding/json"
	"os"
	"path"
)

type Config struct {
	AccessToken string `json:"access_token"`
	APIKey      string `json:"api_key"`
}

func LoadConfig() (*Config, error) {
	home := os.Getenv("HOME")
	filename := path.Join(home, `.hyperdash/hyperdash.json`)
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var config Config
	if err := json.NewDecoder(f).Decode(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
