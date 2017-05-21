package main

import (
	"fmt"
	"github.com/go-ini/ini"
	"github.com/sakisds/go-srtm/api"
	"log"
	"net/http"
)

func main() {
	// Read configuration
	cfg, err := ini.LooseLoad("config.ini")
	if err != nil {
		panic(err)
	}

	var address, port string
	if !cfg.Section("").HasKey("address") {
		address = ""
	} else {
		address = cfg.Section("").Key("address").String()
	}
	if !cfg.Section("").HasKey("port") {
		port = "8080"
	} else {
		port = cfg.Section("").Key("port").String()
	}

	// Elevation API
	api.InitElevation(*cfg)
	http.HandleFunc("/elevation", api.HandleElevation)
	http.HandleFunc("/elevationPath", api.HandleElevationPath)

	// Start server
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", address, port), nil))
}
