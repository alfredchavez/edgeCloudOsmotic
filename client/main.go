package main

import (
	"encoding/json"
	"github.com/rs/xid"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
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

func generateRandomId() string {
	return xid.New().String()
}

func performRequestsUntilGetTheResult(url string, parameter int) {
	defer wg.Done()
	results := new(Result)
	fName := "untitled-"+ generateRandomId()
	param := strconv.Itoa(parameter)
	myUrl := url
	start := time.Now()
	for {
		err := getJson(myUrl + "/execute/" + fName + "?param=" + param, &results)
		if err != nil {
			log.Println("Problem getting function "+err.Error())
			elapsed := time.Since(start)
			sumTo(elapsed.Seconds())
			return
		}
		if strings.HasPrefix(results.Result, "http") {
			myUrl = results.Result
		} else {
			elapsed := time.Since(start)
			sumTo(elapsed.Seconds())
			return
		}
	}
}

var totalSum float64
var lock sync.Mutex
var wg sync.WaitGroup
var messages chan string


func sumTo(val float64) {
	lock.Lock()
	defer lock.Unlock()
	totalSum += val
}

func main() {
	totalSum = 0.0
	mainServer := "http://127.0.0.1:9000"
	messages = make(chan string)
	wg.Add(10)
	for i:=0; i< 10; i++{
		firstUrl := new(FromMain)
		getJson(mainServer + "/query_execute", &firstUrl)
		go performRequestsUntilGetTheResult(firstUrl.Url, 20000000)
		//log.Printf("Value %.4f ans %s", dur.Seconds(), ans)
	}
	wg.Wait()
	log.Printf("total time %.4f", totalSum)
}
