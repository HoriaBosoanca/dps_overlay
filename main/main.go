package main

import (
	"dps_overlay/data"
	"dps_overlay/overlay"
	"fmt"
)

func main() {
	clientData, err := data.Give_data("https://127.0.0.1:2999")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(clientData)
	overlay.StartOverlay()
}
