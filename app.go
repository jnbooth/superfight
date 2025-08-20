package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	config := ParseConfig()
	SetupFrontend(config.Dir)
	hub := SetupBackend(config.Dir)
	defer hub.Shutdown()
	go hub.Run()

	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)

	fmt.Println("Listening on", addr)

	log.Fatal(http.ListenAndServe(addr, nil))
}
