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
		flags := C.PTR_FLAGS_DOWN | C.PTR_FLAGS_BUTTON1
		log.Printf("ProcessRDPEvent x:%d y:%d", info.X, info.Y)
		C.webrdp_button_press(instance, C.UINT16(flags), C.int(info.X), C.int(info.Y))
	case "btnRel":
		flags := C.PTR_FLAGS_BUTTON1
		log.Printf("ProcessRDPEvent x:%d y:%d", info.X, info.Y)
		C.webrdp_button_press(instance, C.UINT16(flags), C.int(info.X), C.int(info.Y))
	}
}
