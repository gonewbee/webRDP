package main

/*
#include "webrdp.h"

static webContext* convert2webContextC(rdpContext* context) {
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
	"log"
	"time"
)

func convert2webContext(context *C.rdpContext) *C.webContext {
	return C.convert2webContextC(context)
}

func writeByChen(context *C.rdpContext, info RdpDrawInfo) {
	t := *(*chan RdpDrawInfo)(C.getWSChan(context))
	select {
	case t <- info:
		log.Println("send ok")
	case <-time.After(time.Second * 5):
		log.Println("send time out")
		xfc := convert2webContext(context)
		xfc.disconnect = C.TRUE
	}
}
