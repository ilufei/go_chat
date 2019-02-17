package main

import (
	"fmt"
	"encoding/json"
	"net/http"
	"ichat/lib"
)

type output struct {
	Code int
	Message string
	Data interface{} 
}

var result = &output{
	Code : 200,
	Message : "sucess",
}

func main() {
    defer func() {
        if err := recover(); err != nil {
            lib.Error("error catched : ", err)
        }
    }()

	mux := http.NewServeMux()
	mux.HandleFunc("/message", handle)
	
	server := &http.Server{
		Addr : "0.0.0.0:9600",
		Handler : mux,
		MaxHeaderBytes : 1 << 20,
	}

	server.ListenAndServe()
}

func handle(w http.ResponseWriter, r *http.Request) {
	if r.PostFormValue("data") == "" || r.PostFormValue("roomid") == "" {
		result = &output{500, "data error or roomid error", "",}
		response(w, result)
		return
	}

	dataJson := r.PostFormValue("data")
	roomid := r.PostFormValue("roomid")
	priority := r.PostFormValue("priority")

	if priority == "" {
		priority = "1"
	}

	var data []map[string]interface{}

	if err := json.Unmarshal([]byte(dataJson), &data); err != nil {
		result = &output{500, "data json decode failed", "",}
		response(w, result)
		return
    }

    if len(data) == 0 {
 		result = &output{500, "data not allow empty", "",}
		response(w, result)
		return
    }

    for _, content := range data {
		redisKey := fmt.Sprintf("%s_%s", roomid, priority)
		contentJson, _ := json.Marshal(content)
		_, err := lib.RedisClient.RPush(redisKey, contentJson).Result()

		if err != nil {
			result = &output{500, "push to redis failed", "",}
			response(w, result)
			return
		}
    }

	result = &output{200, "sucess", "",}
	response(w, result)
	return
}

func response(w http.ResponseWriter, result *output) {
	jsonStr, _ := json.Marshal(result)
	fmt.Fprintf(w, string(jsonStr))
	return
}