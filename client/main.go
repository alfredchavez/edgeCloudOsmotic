package main

import (
	"encoding/json"
	"net/http"
	"time"
)

var myClient = &http.Client{Timeout: 60 * time.Second}

func getJson(url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

type FromMain struct {
	Url string `json:"url"`
}

type Result struct {
	Result string `json:"result"`
}

func main() {
	mainServer := "http://127.0.0.1:9000"
	for i:=0; i< 100; i++{
		ans1 := new(FromMain)
		getJson(mainServer + "/execute", &ans1)
		println(ans1.Url)
		ans2 := new(Result)
		getJson(ans1.Url+"/execute/untitled?param=10000", &ans2)
		println(ans2.Result)
		if string(ans2.Result[0]) == "h" {
			getJson(ans2.Result+"/execute/untitled?param=10000", &ans2)
		}
		println(ans2.Result)
	}
}
