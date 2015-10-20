#ifndef _WEB_RDP_H
#define _WEB_RDP_H

#include "freerdp/freerdp.h"
#include "freerdp/client.h"

typedef struct {
	rdpContext context;
	void *wsChan;
} webContext;

typedef struct web_rdp_bitmap {
    rdpBitmap bitmap;
} web_rdp_bitmap;

#endif