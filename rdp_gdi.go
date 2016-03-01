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
*/
import "C"
import (
	"log"
)

//export webRDPend_paint
func webRDPend_paint(context *C.rdpContext) C.BOOL {
	// log.Println("webRDPend_paint")
	return true
}

//export webRDPdesktop_resize
func webRDPdesktop_resize(context *C.rdpContext) C.BOOL {
	log.Println("webRDPdesktop_resize")
	return true
}

//export webRDPdstblt
func webRDPdstblt(context *C.rdpContext, dstblt *C.DSTBLT_ORDER) C.BOOL {
	log.Println("webRDPdstblt")
	info := RdpDrawInfo{}
	info.Type = 1
	info.Left = uint16(dstblt.nLeftRect)
	info.Top = uint16(dstblt.nTopRect)
	info.Width = uint16(dstblt.nWidth)
	info.Height = uint16(dstblt.nHeight)
	writeByChen(context, info)
	return true
}

//export webRDPpatblt
func webRDPpatblt(context *C.rdpContext, patblt *C.PATBLT_ORDER) C.BOOL {
	log.Println("webRDPpatblt")
	return true
}

//export webRDPscrblt
func webRDPscrblt(context *C.rdpContext, scrblt *C.SCRBLT_ORDER) C.BOOL {
	log.Println("webRDPscrblt")
	return true
}

//export webRDPmemblt
func webRDPmemblt(context *C.rdpContext, memblt *C.MEMBLT_ORDER) C.BOOL {
	log.Println("webRDPmemblt")
	return true
}

//export webRDPopaquerect
func webRDPopaquerect(context *C.rdpContext, opaque_rect *C.OPAQUE_RECT_ORDER) C.BOOL {
	log.Println("webRDPopaquerect")
	color := C.freerdp_color_convert_var(opaque_rect.color, 32, 32, convert2webContext(context).clrconv)
	log.Printf("webRDPopaquerect:%x==>%x", opaque_rect.color, color)
	info := RdpDrawInfo{}
	info.Type = 2
	info.Color = uint32(color)
	info.Left = uint16(opaque_rect.nLeftRect)
	info.Top = uint16(opaque_rect.nTopRect)
	info.Width = uint16(opaque_rect.nWidth)
	info.Height = uint16(opaque_rect.nHeight)
	writeByChen(context, info)
	return true
}

//export webRDPpalette_update
func webRDPpalette_update(context *C.rdpContext, palette *C.PALETTE_UPDATE) C.BOOL {
	log.Println("webRDPpalette_update")
	return true
}

func webGdiRegisterUpdateCallbacks(update *C.rdpUpdate) {
	C.web_gdi_register_update_callbacks(update)
}
