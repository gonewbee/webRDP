package main

/*
#include "freerdp/freerdp.h"
#include "freerdp/client.h"
#include "webrdp.h"
#include "freerdp/cache/glyph.h"
#include "freerdp/channels/channels.h"

extern BOOL webfreerdp_client_global_init();
extern void webfreerdp_client_global_uninit();
extern BOOL webfreerdp_client_new(freerdp* instance, rdpContext* context);
extern void webfreerdp_client_free(freerdp* instance, rdpContext* context);
extern int webfreerdp_client_start(rdpContext* context);
extern int webfreerdp_client_stop(rdpContext* context);

extern BOOL web_pre_connect(freerdp* instance);
extern BOOL web_post_connect(freerdp* instance);
extern BOOL web_authenticate(freerdp* instance, char** username, char** password, char** domain);
extern BOOL web_verify_certificate(freerdp* instance, char* subject, char* issuer, char* fingerprint);


extern BOOL webRdpBitmapNew(rdpContext* context, rdpBitmap* bitmap);
extern void webRdpBitmapFree(rdpContext* context, rdpBitmap* bitmap);
extern BOOL webRdpBitmapPaint(rdpContext* context, rdpBitmap* bitmap);
extern BOOL webRdpBitmapDecompress(rdpContext* context, rdpBitmap* bitmap,
		BYTE* data, int width, int height, int bpp, int length,
		BOOL compressed, int codec_id);
extern BOOL webRdpBitmapSetSurface(rdpContext* context, rdpBitmap* bitmap, BOOL primary);

static int RdpClientEntry(RDP_CLIENT_ENTRY_POINTS* pEntryPoints) {
	pEntryPoints->Version = 1;
	pEntryPoints->Size = sizeof(RDP_CLIENT_ENTRY_POINTS_V1);
	pEntryPoints->GlobalInit = webfreerdp_client_global_init;
	pEntryPoints->GlobalUninit = webfreerdp_client_global_uninit;
	pEntryPoints->ContextSize = sizeof(webContext);
	pEntryPoints->ClientNew = webfreerdp_client_new;
	pEntryPoints->ClientFree = webfreerdp_client_free;
	pEntryPoints->ClientStart = webfreerdp_client_start;
	pEntryPoints->ClientStop = webfreerdp_client_stop;
	return 0;
}

static void setFuncInClient(freerdp *instance, rdpContext* context) {
	webContext* xfc = (webContext*) instance->context;
	xfc->clrconv = freerdp_clrconv_new(CLRCONV_ALPHA|CLRCONV_INVERT);
	context->channels = freerdp_channels_new();
	instance->PreConnect = web_pre_connect;
	instance->PostConnect = web_post_connect;
	instance->Authenticate = web_authenticate;
	instance->VerifyCertificate = web_verify_certificate;
}

static void setContextChan(freerdp *instance, INT64 chanid) {
	webContext* xfc = (webContext*) instance->context;
	xfc->chanid = chanid;
}

static void web_pre_connect_set(freerdp *instance) {
	rdpSettings* settings;
	settings = instance->settings;
	settings->RemoteFxCodec = 0;
    settings->FastPathOutput = 1;
    settings->ColorDepth = 32;//16;
    settings->FrameAcknowledge = 1;
    settings->LargePointerFlag = 1;
    settings->BitmapCacheV3Enabled = 0;
    settings->BitmapCachePersistEnabled = 0;

    settings->OrderSupport[NEG_DSTBLT_INDEX] = TRUE;
    settings->OrderSupport[NEG_PATBLT_INDEX] = TRUE;
    settings->OrderSupport[NEG_SCRBLT_INDEX] = TRUE;
    settings->OrderSupport[NEG_OPAQUE_RECT_INDEX] = TRUE;
    settings->OrderSupport[NEG_DRAWNINEGRID_INDEX] = FALSE;
    settings->OrderSupport[NEG_MULTIDSTBLT_INDEX] = FALSE;
    settings->OrderSupport[NEG_MULTIPATBLT_INDEX] = FALSE;
    settings->OrderSupport[NEG_MULTISCRBLT_INDEX] = FALSE;
    settings->OrderSupport[NEG_MULTIOPAQUERECT_INDEX] = TRUE;
    settings->OrderSupport[NEG_MULTI_DRAWNINEGRID_INDEX] = FALSE;
    settings->OrderSupport[NEG_LINETO_INDEX] = TRUE;
    settings->OrderSupport[NEG_POLYLINE_INDEX] = TRUE;
    settings->OrderSupport[NEG_MEMBLT_INDEX] = FALSE;

    settings->OrderSupport[NEG_MEM3BLT_INDEX] = FALSE;

    settings->OrderSupport[NEG_MEMBLT_V2_INDEX] = FALSE;
    settings->OrderSupport[NEG_MEM3BLT_V2_INDEX] = FALSE;
    settings->OrderSupport[NEG_SAVEBITMAP_INDEX] = FALSE;
    settings->OrderSupport[NEG_GLYPH_INDEX_INDEX] = TRUE;
    settings->OrderSupport[NEG_FAST_INDEX_INDEX] = TRUE;
    settings->OrderSupport[NEG_FAST_GLYPH_INDEX] = TRUE;

    settings->OrderSupport[NEG_POLYGON_SC_INDEX] = FALSE;
    settings->OrderSupport[NEG_POLYGON_CB_INDEX] = FALSE;

    settings->OrderSupport[NEG_ELLIPSE_SC_INDEX] = FALSE;
    settings->OrderSupport[NEG_ELLIPSE_CB_INDEX] = FALSE;

	settings->GlyphSupportLevel = GLYPH_SUPPORT_NONE;

	if (!instance->context->cache)
		instance->context->cache = cache_new(instance->settings);
}

static BOOL web_register_graphics(rdpGraphics* graphics) {
	rdpBitmap* bitmap = NULL;
	rdpPointer* pointer = NULL;
	rdpGlyph* glyph = NULL;
	BOOL ret = FALSE;

	if (!(bitmap = (rdpBitmap*) calloc(1, sizeof(rdpBitmap))))
		goto out;

	if (!(pointer = (rdpPointer*) calloc(1, sizeof(rdpPointer))))
		goto out;

	if (!(glyph = (rdpGlyph*) calloc(1, sizeof(rdpGlyph))))
		goto out;

	bitmap->size = sizeof(web_rdp_bitmap);
	bitmap->New = webRdpBitmapNew;
	bitmap->Free = webRdpBitmapFree;
	bitmap->Paint = webRdpBitmapPaint;
	bitmap->Decompress = webRdpBitmapDecompress;
	bitmap->SetSurface = webRdpBitmapSetSurface;

	graphics_register_bitmap(graphics, bitmap);

	ret = TRUE;

out:
	free(bitmap);
	free(pointer);
	free(glyph);

	return ret;
}

static void web_client_func(freerdp* instance) {
	BOOL status;
	DWORD nCount;
	DWORD waitStatus;
	HANDLE handles[64];
	rdpContext* context;
	webContext* xfc;

	context = instance->context;
	status = freerdp_connect(instance);

	xfc = (webContext*) instance->context;
	while (!xfc->disconnect && !freerdp_shall_disconnect(instance)) {
		nCount = 0;
		DWORD tmp = freerdp_get_event_handles(context, &handles[nCount], 64 - nCount);
		if (tmp == 0)
		{
			fprintf(stderr, "freerdp_get_event_handles failed\n");
			break;
		}
		nCount += tmp;

		waitStatus = WaitForMultipleObjects(nCount, handles, FALSE, 100);

		if (!freerdp_check_event_handles(context))
		{
			fprintf(stderr, "Failed to check FreeRDP file descriptor\n");
			break;
		}
	}
	fprintf(stdout, "web_client_func==========end\n");
	freerdp_disconnect(instance);
}
*/
import "C"

import (
	"log"
	"unsafe"
)

func test() {
	instance := C.freerdp_new()
	log.Printf("instance address:%p", instance)
	log.Printf("instance:%v", instance)
	log.Println(instance.PreConnect)
}

//export web_pre_connect
func web_pre_connect(instance *C.freerdp) C.BOOL {
	log.Println("web_pre_connect")
	C.web_pre_connect_set(instance)
	return C.TRUE
}

//export web_post_connect
func web_post_connect(instance *C.freerdp) C.BOOL {
	log.Println("web_post_connect")
	var update *C.rdpUpdate
	update = instance.context.update
	C.web_register_graphics(instance.context.graphics)
	webGdiRegisterUpdateCallbacks(update)
	// C.pointer_cache_register_callbacks(update)
	C.glyph_cache_register_callbacks(update)
	C.brush_cache_register_callbacks(update)
	C.bitmap_cache_register_callbacks(update)
	C.offscreen_cache_register_callbacks(update)
	C.palette_cache_register_callbacks(update)
	log.Println("web_post_connect end")
	return C.TRUE
}

//export web_authenticate
func web_authenticate(instance *C.freerdp, username, password, domain **C.char) C.BOOL {
	log.Println("web_authenticate")
	return C.TRUE
}

//export web_verify_certificate
func web_verify_certificate(instance *C.freerdp, subject, issuer, fingerprint *C.char) C.BOOL {
	log.Println("web_verify_certificate")
	return C.TRUE
}

//export webfreerdp_client_global_init
func webfreerdp_client_global_init() C.BOOL {
	log.Println("webfreerdp_client_global_init")
	return C.BOOL(C.TRUE)
}

//export webfreerdp_client_global_uninit
func webfreerdp_client_global_uninit() {
	log.Println("webfreerdp_client_global_uninit")
}

//export webfreerdp_client_new
func webfreerdp_client_new(instance *C.freerdp, context *C.rdpContext) C.BOOL {
	log.Println("webfreerdp_client_new")
	C.setFuncInClient(instance, context)
	return C.TRUE
}

//export webfreerdp_client_free
func webfreerdp_client_free(instance *C.freerdp, context *C.rdpContext) {
	log.Println("webfreerdp_client_free")
}

//export webfreerdp_client_start
func webfreerdp_client_start(context *C.rdpContext) C.int {
	log.Println("webfreerdp_client_start")
	var instance *C.freerdp
	instance = context.instance
	C.web_client_func(instance)
	return 0
}

//export webfreerdp_client_stop
func webfreerdp_client_stop(context *C.rdpContext) C.int {
	log.Println("webfreerdp_client_stop")
	return 0
}

func setRdpInfo(context *C.rdpContext, info wsReadInfo) {
	settings := context.instance.settings
	log.Printf("w:%d h:%d", settings.DesktopWidth, settings.DesktopHeight)
	settings.ServerHostname = C.CString(info.Ip)
	settings.Username = C.CString(info.User)
	settings.Password = C.CString(info.Passwd)
	// defer C.free(unsafe.Pointer(settings.ServerHostname))
	// defer C.free(unsafe.Pointer(settings.Username))
	// defer C.free(unsafe.Pointer(settings.Password))

	// Standard RDP
	// settings.RdpSecurity = C.TRUE
	// settings.TlsSecurity = C.FALSE
	// settings.NlaSecurity = C.FALSE
	// settings.ExtSecurity = C.FALSE
	// settings.UseRdpSecurityLayer = C.TRUE
}

func Rdp_new(chanId int64) *C.rdpContext {
	var clientEntryPoints C.RDP_CLIENT_ENTRY_POINTS
	clientEntryPoints.Size = C.DWORD(unsafe.Sizeof(clientEntryPoints))
	clientEntryPoints.Version = C.RDP_CLIENT_INTERFACE_VERSION
	log.Printf("size:%d version:%d", clientEntryPoints.Size, clientEntryPoints.Version)
	C.RdpClientEntry(&clientEntryPoints)
	log.Printf("size:%d version:%d", clientEntryPoints.Size, clientEntryPoints.Version)
	context := C.freerdp_client_context_new(&clientEntryPoints)
	C.setContextChan(context.instance, C.INT64(chanId))
	return context
}

func Rdp_free(context *C.rdpContext) {
	log.Println("Rdp_free!")
	C.freerdp_client_stop(context)
	C.freerdp_client_context_free(context)
}

func Rdp_start(context *C.rdpContext) {
	log.Println(C.GoString(context.instance.settings.ServerHostname))
	C.freerdp_client_start(context)

	log.Println("Rdp_start end!")
	Rdp_free(context)
}

func Rdp_stop(context *C.rdpContext) {
	log.Println("Rdp_stop set disconnect----------")
	xfc := convert2webContext(context)
	xfc.disconnect = C.TRUE
}
