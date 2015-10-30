package main

/*
#include <freerdp/graphics.h>
#include "webrdp.h"
#include "freerdp/gdi/gdi.h"
static int getsizeof(BYTE* p) {
	return sizeof(p);
}

static BYTE* getBYTEpalette(rdpContext* context) {
	webContext* xfc = (webContext*) context;
	return (BYTE*)xfc->palette;
}
*/
import "C"
import (
	"log"
	"unsafe"
)

//export webRdpBitmapNew
func webRdpBitmapNew(context *C.rdpContext, bitmap *C.rdpBitmap) C.BOOL {
	// log.Println("webRdpBitmapNew")
	// if bitmap.data != nil {
	// 	log.Printf("l:%d t:%d r:%d b:%d w:%d h:%d", bitmap.left, bitmap.top, bitmap.right, bitmap.bottom, bitmap.width, bitmap.height)
	// }
	return C.TRUE
}

//export webRdpBitmapFree
func webRdpBitmapFree(context *C.rdpContext, bitmap *C.rdpBitmap) {
	// log.Println("webRdpBitmapFree")
	// C._aligned_free(unsafe.Pointer(bitmap.data))
}

//export webRdpBitmapPaint
func webRdpBitmapPaint(context *C.rdpContext, bitmap *C.rdpBitmap) C.BOOL {
	log.Println("webRdpBitmapPaint")
	log.Printf("webRdpBitmapPaint length:%d", bitmap.length)
	if bitmap.data != nil {
		bs := C.GoBytes(unsafe.Pointer(bitmap.data), C.int(bitmap.length))
		log.Printf("%x %x %x %x %x %x %x %x", bs[0], bs[1], bs[2], bs[3], bs[4], bs[5], bs[6], bs[7])
		info := RdpDrawInfo{}
		info.Type = 5
		info.Bmp = bs
		info.BmpLen = int(bitmap.length)
		writeByChen(context, info)
	}
	return C.TRUE
}

//export webRdpBitmapDecompress
func webRdpBitmapDecompress(context *C.rdpContext, bitmap *C.rdpBitmap, data *C.BYTE,
	width C.int, height C.int, bpp C.int, length C.int,
	compressed C.BOOL, codecId C.int) C.BOOL {
	log.Printf("compressed:%d bpp:%d", compressed, bpp)
	size := width * height * 4
	if bitmap.data != nil {
		C._aligned_free(unsafe.Pointer(bitmap.data))
	}
	bitmap.data = (*C.BYTE)(C._aligned_malloc(C.size_t(size), 16))
	if compressed != C.FALSE {
		if bpp < 32 {
			C.freerdp_client_codecs_prepare(context.codecs, C.FREERDP_CODEC_INTERLEAVED)
			C.interleaved_decompress(context.codecs.interleaved, data, C.UINT32(length), bpp,
				&(bitmap.data), C.PIXEL_FORMAT_XRGB32, -1, 0, 0, width, height, C.getBYTEpalette(context))
		} else {
			C.freerdp_client_codecs_prepare(context.codecs, C.FREERDP_CODEC_PLANAR)
			status := C.planar_decompress(context.codecs.planar, data, C.UINT32(length),
				&(bitmap.data), C.PIXEL_FORMAT_XRGB32, -1, 0, 0, width, height, C.TRUE)
			log.Printf("webRdpBitmapDecompress status::::::%d", status)
		}
	} else {
		C.freerdp_image_flip(data, bitmap.data, width, height, bpp)
	}
	bitmap.compressed = C.FALSE
	bitmap.length = C.UINT32(size)
	bitmap.bpp = 32
	return C.TRUE
}

//export webRdpBitmapSetSurface
func webRdpBitmapSetSurface(context *C.rdpContext, bitmap *C.rdpBitmap, primary C.BOOL) C.BOOL {
	log.Println("webRdpBitmapSetSurface")
	return C.TRUE
}
