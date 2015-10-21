package main

/*
#include <freerdp/graphics.h>
#include "webrdp.h"

extern BOOL webRDPend_paint(rdpContext* context);
extern BOOL webRDPdesktop_resize(rdpContext* context);

extern BOOL webRDPdstblt(rdpContext* context, DSTBLT_ORDER* dstblt);
extern BOOL webRDPpatblt(rdpContext* context, PATBLT_ORDER* patblt);
extern BOOL webRDPscrblt(rdpContext* context, SCRBLT_ORDER* scrblt);
extern BOOL webRDPmemblt(rdpContext* context, MEMBLT_ORDER* memblt);
extern BOOL webRDPopaquerect(rdpContext* context, OPAQUE_RECT_ORDER* opaque_rect);
extern BOOL webRDPpalette_update(rdpContext* context, PALETTE_UPDATE* palette);

static void web_gdi_register_update_callbacks(rdpUpdate* update) {
	rdpPrimaryUpdate* primary = update->primary;

	update->Palette = webRDPpalette_update;
	update->BeginPaint = NULL;
	update->EndPaint = webRDPend_paint;
	update->DesktopResize = webRDPdesktop_resize;

	primary->DstBlt = webRDPdstblt;
	primary->PatBlt = webRDPpatblt;
	primary->ScrBlt = webRDPscrblt;
	primary->OpaqueRect = webRDPopaquerect;
	primary->DrawNineGrid = NULL;
	primary->MultiDstBlt = NULL;
	primary->MultiPatBlt = NULL;
	primary->MultiScrBlt = NULL;
	primary->MultiOpaqueRect = NULL;
	primary->MultiDrawNineGrid = NULL;
	primary->LineTo = NULL;
	primary->Polyline = NULL;
	primary->MemBlt = webRDPmemblt;
	primary->Mem3Blt = NULL;
	primary->SaveBitmap = NULL;
	primary->GlyphIndex = NULL;
	primary->FastIndex = NULL;
	primary->FastGlyph = NULL;
	primary->PolygonSC = NULL;
	primary->PolygonCB = NULL;
	primary->EllipseSC = NULL;
	primary->EllipseCB = NULL;
}

static webContext* convert2webContext(rdpContext* context) {
	webContext* xfc = (webContext*) context;
	return xfc;
}

static void *getWSChan(rdpContext* context) {
	webContext* xfc = (webContext*) context;
	return xfc->wsChan;
}
*/
import "C"
import (
	"fmt"
	"log"
)

//export webRDPend_paint
func webRDPend_paint(context *C.rdpContext) C.BOOL {
	log.Println("webRDPend_paint")
	return C.TRUE
}

//export webRDPdesktop_resize
func webRDPdesktop_resize(context *C.rdpContext) C.BOOL {
	log.Println("webRDPdesktop_resize")
	return C.TRUE
}

//export webRDPdstblt
func webRDPdstblt(context *C.rdpContext, dstblt *C.DSTBLT_ORDER) C.BOOL {
	log.Println("webRDPdstblt")
	t := (*chan string)(C.getWSChan(context))
	log.Printf("t:%p", t)
	s := fmt.Sprintf("01%04x%04x%04x%04x", dstblt.nLeftRect, dstblt.nTopRect, dstblt.nWidth, dstblt.nHeight)
	*t <- s
	return C.TRUE
}

//export webRDPpatblt
func webRDPpatblt(context *C.rdpContext, patblt *C.PATBLT_ORDER) C.BOOL {
	log.Println("webRDPpatblt")
	return C.TRUE
}

//export webRDPscrblt
func webRDPscrblt(context *C.rdpContext, scrblt *C.SCRBLT_ORDER) C.BOOL {
	log.Println("webRDPscrblt")
	return C.TRUE
}

//export webRDPmemblt
func webRDPmemblt(context *C.rdpContext, memblt *C.MEMBLT_ORDER) C.BOOL {
	log.Println("webRDPmemblt")
	return C.TRUE
}

//export webRDPopaquerect
func webRDPopaquerect(context *C.rdpContext, opaque_rect *C.OPAQUE_RECT_ORDER) C.BOOL {
	log.Println("webRDPopaquerect")
	color := C.freerdp_color_convert_var(opaque_rect.color, 32, 32, C.convert2webContext(context).clrconv)
	log.Printf("webRDPopaquerect:%x==>%x", opaque_rect.color, color)
	t := (*chan string)(C.getWSChan(context))
	log.Printf("t:%p", t)
	s := fmt.Sprintf("02#%06x%04x%04x%04x%04x", opaque_rect.color, opaque_rect.nLeftRect, opaque_rect.nTopRect, opaque_rect.nWidth, opaque_rect.nHeight)
	*t <- s
	return C.TRUE
}

//export webRDPpalette_update
func webRDPpalette_update(context *C.rdpContext, palette *C.PALETTE_UPDATE) C.BOOL {
	log.Println("webRDPpalette_update")
	return C.TRUE
}

func webGdiRegisterUpdateCallbacks(update *C.rdpUpdate) {
	C.web_gdi_register_update_callbacks(update)
}
