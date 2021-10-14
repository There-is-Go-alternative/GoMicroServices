package config

import (
	"github.com/There-is-Go-alternative/GoMicroServices/chats/internal/utils"
)

type Config struct {
	ChatsHost        string `json:"chats_host,omitempty"`
	ChatsPort        string `json:"chats_port,omitempty"`
	ChatsEndpoint    string `json:"chats_endpoint,omitempty"`
	MessageHost      string `json:"messages_host,omitempty"`
	MessagePort      string `json:"messages_port,omitempty"`
	MessagesEndpoint string `json:"messages_endpoint,omitempty"`
}

func ParseConfigFromPath(path string) (conf *Config, err error) {
	err = utils.DecodeJSONFromFile(path, &conf)
	return
}
