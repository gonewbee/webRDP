package main

/*
#include <freerdp/graphics.h>
*/
import "C"
import (
	"log"
)

//export webRdpBitmapNew
func webRdpBitmapNew(context *C.rdpContext, bitmap *C.rdpBitmap) C.BOOL {
	log.Println("webRdpBitmapNew")
	return C.TRUE
}

//export webRdpBitmapFree
func webRdpBitmapFree(context *C.rdpContext, bitmap *C.rdpBitmap) {
	log.Println("webRdpBitmapFree")
}

//export webRdpBitmapPaint
func webRdpBitmapPaint(context *C.rdpContext, bitmap *C.rdpBitmap) C.BOOL {
	log.Println("webRdpBitmapPaint")
	return C.TRUE
}

//export webRdpBitmapDecompress
func webRdpBitmapDecompress(context *C.rdpContext, bitmap *C.rdpBitmap, data *C.BYTE,
	width C.int, height C.int, bpp C.int, length C.int,
	compressed C.BOOL, codecId C.int) C.BOOL {
	log.Println("webRdpBitmapDecompress")
	return C.TRUE
}

//export webRdpBitmapSetSurface
func webRdpBitmapSetSurface(context *C.rdpContext, bitmap *C.rdpBitmap, primary C.BOOL) C.BOOL {
	log.Println("webRdpBitmapDecompress")
	return C.TRUE
}
