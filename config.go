package main

import (
	"flag"
)

type Config struct {
	Dir  string
	Host string
	Port int
}

func ParseConfig() Config {
	config := Config{}
	flag.StringVar(&config.Dir, "dir", ".", "assets directory")
	flag.StringVar(&config.Host, "h", "", "hostname to listen on")
	flag.IntVar(&config.Port, "p", 80, "port to listen on")
	flag.Parse()
	return config
}
