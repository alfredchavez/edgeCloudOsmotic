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
	for i:=0; i< 1; i++{
		ans1 := new(FromMain)
		getJson(mainServer + "/query_execute", &ans1)
		println("Request: ", ans1.Url)
		ans2 := new(Result)
		err := getJson(ans1.Url+"/execute/untitled?param=1000000000", &ans2)
		if err != nil {
			println(err.Error())
		}
		if ans2.Result == "h" {
		        println(ans2.Result)
			getJson(ans2.Result+"/execute/untitled?param=1000000", &ans2)
		}
		println("second time " + ans2.Result)
	}
}
