package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"time"
	"net/http"
	_ "ichat/lib"
	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "0.0.0.0:9501", "http service address")
var upgrader = websocket.Upgrader{} // use default options
var conns []*websocket.Conn

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

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	go postMessage(interrupt)

	flag.Parse()
	http.HandleFunc("/echo", echo)
	http.ListenAndServe(*addr, nil)
}

func postMessage(interrupt chan os.Signal) {
	for {
		select {
			case <-interrupt :
				log.Println("postMessage interrupt by main")
				return
			default : 
				message := getMessage()

				for _, conn := range conns {
					conn.WriteMessage(websocket.TextMessage, []byte(message))
				}

				time.Sleep(time.Second)
		}
	}
}


func getMessage() (string) {
	/*
	redisKey := Sprintf("%s_1", roomid)
	message, err := lib.RedisClient.LPop(redisKey).Result()
	*/

	data := "this is message"
	return data
}




