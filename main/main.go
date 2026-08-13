package main

import (
	"dps_overlay/overlay"
)

func main() {
	//clientData, err := data.GiveData("https://127.0.0.1:2999")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(clientData)
	overlay.RunOverlay()
}
