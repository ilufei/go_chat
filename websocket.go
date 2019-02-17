package main

import (
	"fmt"
	"flag"
	"log"
	"time"
	"encoding/json"
	"net/http"
	"ichat/lib"
	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "0.0.0.0:9601", "http service address")
var upgrader = websocket.Upgrader{} // use default options
var conns []*websocket.Conn

var message map[string]interface{}
type result struct {
	Roomid string
	Data []map[string]interface{}
}

func main() {
	go postMessage()

	flag.Parse()
	http.HandleFunc("/echo", echo)
	http.ListenAndServe(*addr, nil)
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

func postMessage() {
	for {
		roomids := getRoomIds()
		for _, roomid := range roomids {
			data := []map[string]interface{}{}	//清空数据
			allMessage := getMessage(roomid)
			if len(allMessage) == 0 {
				continue
			}

			for _, messageJson := range allMessage {
				if err := json.Unmarshal([]byte(messageJson), &message); err != nil {
					continue
			    }

			    data = append(data, message)
			}

			resultObject := result{roomid, data,}
			resultJson, _ := json.Marshal(resultObject)

			for _, conn := range conns {
				conn.WriteMessage(websocket.TextMessage, []byte(resultJson))
			}
		}

		time.Sleep(time.Second)
	}
}


func getMessage(roomid string) ([]string) {
	var allMessages []string 
	var messages []string
	var messageLength int64 = 50
	
	room1Key := fmt.Sprintf("%s_1", roomid)
	room2Key := fmt.Sprintf("%s_2", roomid)
	room3Key := fmt.Sprintf("%s_3", roomid)

	messages, _ = lib.RedisClient.LRange(room1Key, 0, messageLength).Result()
	if len(messages) != 0 {
		allMessages = append(allMessages, messages...)
		messageLength = messageLength - int64(len(messages))

		if messageLength <= 0 {
			_, _ = lib.RedisClient.Del(room1Key, room2Key, room3Key).Result()
			return allMessages
		}
	}
	
	messages, _ = lib.RedisClient.LRange(room2Key, 0, messageLength).Result()
	if len(messages) != 0 {
		allMessages = append(allMessages, messages...)
		messageLength = messageLength - int64(len(messages))

		if messageLength <= 0 {
			_, _ = lib.RedisClient.Del(room1Key, room2Key, room3Key).Result()
			return allMessages
		}
	}

	messages, _ = lib.RedisClient.LRange(room3Key, 0, messageLength).Result()
	if len(messages) != 0 {
		allMessages = append(allMessages, messages...)
		messageLength = messageLength - int64(len(messages))
	}

	_, _ = lib.RedisClient.Del(room1Key, room2Key, room3Key).Result()
	return allMessages
}

func getRoomIds() ([]string) {
	roomids, err := lib.RedisClient.SMembers("ichat_roomids").Result()

	if err != nil {
		return []string{}
	}

	return roomids
}




