package config

import (
	"github.com/There-is-Go-alternative/GoMicroServices/account/internal/utils"
)

type Config struct {
	Host            string `json:"host,omitempty"`
	Port            string `json:"port,omitempty"`
	AccountEndpoint string `json:"account_endpoint,omitempty"`
}

func ParseConfigFromPath(path string) (conf *Config, err error) {
	err = utils.DecodeJSONFromFile(path, &conf)
	return
}
