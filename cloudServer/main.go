package main

import (
	"edgeServer/execution_service"
	"edgeServer/monitoring_service"
	"edgeServer/utils"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

type ResponseExecuteFunction struct {
	Result string `json:"result"`
}

func executeFunction(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	functionName := params["f-name"]
	parameter := r.URL.Query().Get("param")
	response := ResponseExecuteFunction{Result: execution_service.ExecuteFunction(functionName, parameter)}
	json.NewEncoder(w).Encode(&response)
}

func stopAnyFunction(w http.ResponseWriter, r *http.Request){
	log.Println("stopping some function!")
}

func main() {
	utils.MyName = os.Args[1]
	port := os.Args[2]
	go monitoring_service.MonitorContext()
	router := mux.NewRouter()
	router.HandleFunc("/execute/{f-name}", executeFunction).Methods("GET")
	router.HandleFunc("/stop", stopAnyFunction).Methods("GET")
	log.Println(port)
	err := http.ListenAndServe(":" + port, router)
	if err != nil {
		log.Fatal(err)
	}
}
