package main

import (
	"fmt"
	"net/http"

	"sleepy.systems/goweb/config"
	"sleepy.systems/goweb/views"
)

func main() {
	var config config.Config
	config.Read("./config.toml")

	http.HandleFunc("/", views.HandleIndexPage)

	fmt.Printf("Server listening at %v", config.PortAsString)
	http.ListenAndServe(config.PortAsString, nil)
}
