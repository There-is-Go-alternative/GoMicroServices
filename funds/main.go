package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/There-is-Go-alternative/GoMicroServices/funds/internal/config"
	log "github.com/sirupsen/logrus"
)

var (
	confPath = flag.String("conf", os.Getenv("CONF_PATH"), "path to the json config file.")
)

func main() {
	flag.Parse()
	if *confPath == "" {
		log.Fatal("config path is empty")
	}

	conf, err := config.ParseConfigFromPath(*confPath)
	if err != nil {
		log.Fatalf("problem when parsing config: %v", err)
	}

	fmt.Println(conf)
}
