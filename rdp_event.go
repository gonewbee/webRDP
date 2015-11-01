package main

/*
#include "freerdp/freerdp.h"

static void webrdp_button_press(freerdp* instance, UINT16 flags, int x, int y) {
	rdpInput* input;
	input = instance->input;
	input->MouseEvent(input, flags, x, y);
}
*/
import "C"
import (
	"log"
)

func ProcessRDPEvent(instance *C.freerdp, info wsReadInfo) {
	switch info.Type {
	case "btnPre":
		var flags C.UINT16
		switch info.Button {
		case 0:
			// left
			flags = C.PTR_FLAGS_DOWN | C.PTR_FLAGS_BUTTON1
		case 1:
			// middle
			flags = C.PTR_FLAGS_DOWN | C.PTR_FLAGS_BUTTON3
		case 2:
			// right
			flags = C.PTR_FLAGS_DOWN | C.PTR_FLAGS_BUTTON2
		}
		log.Printf("ProcessRDPEvent x:%d y:%d", info.X, info.Y)
		C.webrdp_button_press(instance, flags, C.int(info.X), C.int(info.Y))
	case "btnRel":
		var flags C.UINT16
		switch info.Button {
		case 0:
			// left
			flags = C.PTR_FLAGS_BUTTON1
		case 1:
			// middle
			flags = C.PTR_FLAGS_BUTTON3
		case 2:
			// right
			flags = C.PTR_FLAGS_BUTTON2
		}
		log.Printf("ProcessRDPEvent x:%d y:%d", info.X, info.Y)
		C.webrdp_button_press(instance, flags, C.int(info.X), C.int(info.Y))
	}
}
