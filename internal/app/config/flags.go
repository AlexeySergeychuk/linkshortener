package config

import (
	"flag"
	"os"
)

var FlagRunAddr string
var FlagShortLinkAddr string

func ParseFlags() {

	flag.StringVar(&FlagRunAddr, "a", ":8080", "adress and port for run application")
	flag.StringVar(&FlagShortLinkAddr, "b", "http://localhost:8080", "adress and port for short link redirect")

	flag.Parse()

	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		FlagRunAddr = envRunAddr
	}

	if envShortLinkAddr := os.Getenv("BASE_URL"); envShortLinkAddr != "" {
		FlagShortLinkAddr = envShortLinkAddr
	}
}
