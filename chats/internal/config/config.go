package config

import (
	"github.com/There-is-Go-alternative/GoMicroServices/chats/internal/utils"
)

type Config struct {
	Host         string `json:"host,omitempty"`
	Port         string `json:"port,omitempty"`
	ChatEndpoint string `json:"chat_endpoint,omitempty"`
}

func ParseConfigFromPath(path string) (conf *Config, err error) {
	err = utils.DecodeJSONFromFile(path, &conf)
	return
}
