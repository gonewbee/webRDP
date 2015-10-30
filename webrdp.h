#ifndef _WEB_RDP_H
#define _WEB_RDP_H

#include "freerdp/freerdp.h"
#include "freerdp/client.h"

typedef struct {
	rdpContext context;
	HCLRCONV clrconv;
	INT64 chanid;
	UINT32 palette[256];
	BOOL disconnect;
} webContext;

typedef struct web_rdp_bitmap {
    rdpBitmap bitmap;
} web_rdp_bitmap;

#endif