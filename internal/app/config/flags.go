package config

import "flag"

var FlagRunAddr string
var FlagShortLinkBaseUrl string

func ParseFlags() {

	flag.StringVar(&FlagRunAddr, "a", ":8080", "adress and port for run application")
	flag.StringVar(&FlagShortLinkBaseUrl, "b", "http://localhost:8080", "adress and port for short link redirect")

	flag.Parse()
}