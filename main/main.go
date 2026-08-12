package main

import (
	"dps_overlay/data"
	"dps_overlay/overlay"
	"fmt"
)

/*
!!! window will not be transparent unless you set either:
1) Settings->System->Display->Graphics->Add Desktop App (browse the app)
->GPU preference->set to intel integrated graphics, not nvidia
2) Nvidia Control Panel->Manage 3D settings
->OpenGL GDI compatibility->Prefer compatible
*/
func main() {
	clientData, err := data.Give_data("https://127.0.0.1:2999")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(clientData)
	overlay.StartOverlay()
}
