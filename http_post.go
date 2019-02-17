package main 

import (
	"fmt"
	"time"
	"net/url"
	"io/ioutil"
	"net/http"
)

func main() {
	for {
		postMessage()
		time.Sleep(time.Second)
	}
}

func postMessage() {
	v := url.Values{}
	v.Set("data", "[{\"userid\":1,\"age\":18,\"name\":\"lufei\"},{\"userid\":2,\"age\":19,\"name\":\"jack\"},{\"userid\":3,\"age\":20,\"name\":\"mingren\"}]")
	v.Set("roomid", "55555")
	
	//body := ioutil.NopCloser(strings.NewReader(v.Encode()))

   response, err := http.PostForm("http://lufei.me:9600/message", v)
   defer response.Body.Close()

	if err != nil {
		fmt.Println("http PostForm failed")
		return
	}
	
	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("http response get failed")
		return
	}
	
	fmt.Println(string(content))
}