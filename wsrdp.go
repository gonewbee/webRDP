package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
	"strconv"
	"time"
)

type wsReadInfo struct {
	Type   string `json:"type"`
	Ip     string `json:"ip"`
	User   string `json:"user"`
	Passwd string `json:"passwd"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
	Button int    `json:"button"`
}

type RdpDrawInfo struct {
	Type   byte   `json:"type"`
	Color  uint32 `json:"color,omitempty"`
	Left   uint16 `json:"left"`
	Top    uint16 `json:"top"`
	Width  uint16 `json:"width"`
	Height uint16 `json:"height"`
	Img    []byte `json:"img,omitempty"`
	ImgLen uint32 `json:"imglen,omitempty"`
}

var chans = make(map[int64]chan []byte)

func wsWorker(ws *websocket.Conn, msg chan<- string, wsClosed chan<- bool) {
	var message string
	for nil == websocket.Message.Receive(ws, &message) {
		log.Println("receive:" + message)
		msg <- message
	}
	log.Println("wsWorker websocket error")
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

	wschan := make(chan []byte)
	log.Println(wschan)
	log.Printf("wschan:%v", wschan)
	id, err := strconv.ParseInt(fmt.Sprintf("%v", wschan)[2:], 16, 64)
	if err != nil {
		id = time.Now().UnixNano()
	}
	log.Printf("id:0x%x", id)

	go wsWorker(ws, c1, wsClosed)
	context := Rdp_new(id)
	rdp_ran := false

forLoop:
	for {
		// 使用Select选择要发送的数据
		select {
		case msg1 := <-c1:
			log.Println("c1 receive:" + msg1)
			var info wsReadInfo
			if err := json.Unmarshal([]byte(msg1), &info); err != nil {
				log.Panic("Unmarshal error!")
				break
			}
			switch info.Type {
			case "login":
				rdp_ran = true
				chans[id] = wschan
				defer delete(chans, id)
				setRdpInfo(context, info)
				go Rdp_start(context)
				defer Rdp_stop(context)
			case "btnPre", "btnRel", "mouseMove":
				if rdp_ran {
					ProcessRDPEvent(context.instance, info)
				}
			}
			// 根据坐标回复颜色值
			// msg1 = "#" + strconv.FormatInt(int64(info.X*10), 16) + strconv.FormatInt(int64(info.Y*10), 16)
			// websocket.Message.Send(ws, msg1)
		case info := <-wschan:
			websocket.Message.Send(ws, info)
		case <-wsClosed:
			log.Printf("wsClosed")
			break forLoop
		}
	}
	if !rdp_ran {
		Rdp_free(context)
	}
	log.Println("wsHandler end======")
}

func main() {
	port := "8080"
	http.Handle("/", http.FileServer(http.Dir("./")))
	http.Handle("/wsDemo", websocket.Handler(wsHandler))
	log.Println("listen 8080")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
