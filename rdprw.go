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
	"encoding/binary"
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
	if nil == wschan {
		xfc := convert2webContext(context)
		xfc.disconnect = true
		return
	}
	data := make([]byte, 17)
	data[0] = info.Type
	binary.BigEndian.PutUint16(data[1:], info.Left)
	binary.BigEndian.PutUint16(data[3:], info.Top)
	binary.BigEndian.PutUint16(data[5:], info.Width)
	binary.BigEndian.PutUint16(data[7:], info.Height)
	binary.BigEndian.PutUint32(data[9:], info.Color)
	log.Printf("%x %x %x %x", data[9], data[10], data[11], data[12])
	if info.ImgLen != 0 {
		binary.BigEndian.PutUint32(data[13:], info.ImgLen)
		data = append(data, info.Img...)
		log.Println(cap(data))
	}
	select {
	case wschan <- data:
		log.Println("send ok")
	case <-time.After(time.Second * 3):
		log.Println("send time out")
		xfc := convert2webContext(context)
		xfc.disconnect = true
	}
}
