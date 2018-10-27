package main

import (
	"github.com/cool2645/youtube-to-netdisk/carrier"
	"github.com/cool2645/youtube-to-netdisk/http"
)

func main() {
	carrier.Use(http.WebInterface{})
	carrier.Start()
}
