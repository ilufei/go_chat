package main

import (
	"flag"
	"log"
	"time"
	"net/url"
	"net/http"
	"github.com/gorilla/websocket"
)

var addrServer = flag.String("addrServer", "0.0.0.0:9602", "http service address")
var addrClient = flag.String("addrClient", "0.0.0.0:9601", "http service address")
var upgrader = websocket.Upgrader{} // use default options
var conns []*websocket.Conn

func main() {
	go getMessage()

	flag.Parse()
	http.HandleFunc("/echo", echo)
	http.ListenAndServe(*addrServer, nil)
}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	conns = append(conns, c)

	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			for i, v := range conns {
		        if v == c {
		            conns = append(conns[:i], conns[i+1:]...)
		            break
		        }
		    }

			log.Println("client read err:", err)
			break
		}
	}
}

func getMessage() {
	c := getWebsocketClient()
	defer c.Close()

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			
			//断线重连
			c = getWebsocketClient()
			continue
		}
		log.Printf("recv: %s, broadcast to client..", message)

		for _, conn := range conns {
			conn.WriteMessage(websocket.TextMessage, []byte(message))
		}
	}
}

func getWebsocketClient() (*websocket.Conn) {
	u := url.URL{Scheme: "ws", Host: *addrClient, Path: "/echo"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Printf("connect failed error :", err)
		time.Sleep(2 * time.Second)
		return getWebsocketClient()
	}

	return c
}