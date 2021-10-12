package config

import (
	"github.com/There-is-Go-alternative/GoMicroServices/ads/internal/utils"
)

type Config struct {
	Host       string `json:"host,omitempty"`
	Port       string `json:"port,omitempty"`
	AdEndpoint string `json:"ads_endpoint,omitempty"`
}

func ParseConfigFromPath(path string) (conf *Config, err error) {
	err = utils.DecodeJSONFromFile(path, &conf)
	return
}
