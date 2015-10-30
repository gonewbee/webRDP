package main

/*
#include "webrdp.h"

static webContext* convert2webContextC(rdpContext* context) {
	webContext* xfc = (webContext*) context;
	return xfc;
}

static INT64 getWSChan(rdpContext* context) {
	webContext* xfc = (webContext*) context;
	return xfc->chanid;
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
	log.Println("writeByChen try to send---------------")
	wschan := chans[int64(C.getWSChan(context))]
	log.Println(wschan)
	select {
	case wschan <- info:
		log.Println("send ok")
	case <-time.After(time.Second * 3):
		log.Println("send time out")
		xfc := convert2webContext(context)
		xfc.disconnect = C.TRUE
	}
}
