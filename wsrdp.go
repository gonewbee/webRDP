package main

/*
#include <freerdp/graphics.h>
#include "webrdp.h"
static void *getWSChan(rdpContext* context) {
	webContext* xfc = (webContext*) context;
	return xfc->wsChan;
}
*/
import "C"
import (
	"encoding/json"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
	"strconv"
	"time"
)

type mousePos struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func wsWorker(ws *websocket.Conn, msg chan<- string, wsClosed chan<- bool) {
	var message string
	for nil == websocket.Message.Receive(ws, &message) {
		log.Println("receive:" + message)
		msg <- message
	}
	wsClosed <- true
	log.Println("wsWorker end=========")
}

// 在Channel send阻塞的情况下，可使用超时退出Goroutine
func testWorker(msg chan<- string) {
	i := 0
	for {
		select {
		case msg <- "#" + strconv.FormatInt(int64(i*0x100000), 16):
			log.Printf("testWorker:%d", i)
			time.Sleep(5 * time.Second)
			i = (i + 1) % 16
		case <-time.After(time.Second * 15):
			// Channel发送超时，退出Goroutine
			log.Println("testWorker timeout, return-------")
			return
		}
	}
	log.Println("testWorker end============")
}

func wsHandler(ws *websocket.Conn) {
	c1 := make(chan string) // 相当于make(chan string, 1)
	// c2 := make(chan string, 5) // 相当于消息队列中的最大消息数目
	wsClosed := make(chan bool)

	context := Rdp_new()
	ct := *(*chan string)(C.getWSChan(context))
	setRdpInfo(context)
	go Rdp_start(context)

	go wsWorker(ws, c1, wsClosed)

forLoop:
	for {
		// 使用Select选择要发送的数据
		select {
		case msg1 := <-c1:
			log.Println("c1 receive:" + msg1)
			var pos mousePos
			if err := json.Unmarshal([]byte(msg1), &pos); err != nil {
				log.Panic("Unmarshal error!")
				break
			}
			// 根据坐标回复颜色值
			// msg1 = "#" + strconv.FormatInt(int64(pos.X*10), 16) + strconv.FormatInt(int64(pos.Y*10), 16)
			// websocket.Message.Send(ws, msg1)
		case msg2 := <-ct:
			log.Printf("ct receive:" + msg2)
			websocket.Message.Send(ws, msg2)
		case <-wsClosed:
			log.Printf("wsClosed")
			break forLoop
		}
	}
	log.Println("wsHandler end======")
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./")))
	http.Handle("/wsDemo", websocket.Handler(wsHandler))
	log.Fatal(http.ListenAndServe(":80", nil))
}
