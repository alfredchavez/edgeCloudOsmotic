package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/labstack/echo/v4"
	"log"
	"mainServer/redis_service"
	"mainServer/utils"
	"net/http"
	"strconv"
)

func getServersInfo(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	serverName := params["server-name"]
	cpuUsage,_ := strconv.ParseFloat(r.URL.Query().Get("cpu"), 32)
	memUsage,_ := strconv.ParseFloat(r.URL.Query().Get("mem"), 32)
	log.Println(cpuUsage, memUsage)
	_ = redis_service.SetValue("mem"+serverName, r.URL.Query().Get("mem"))
	_ = redis_service.SetValue("cpu"+serverName, r.URL.Query().Get("cpu"))
	if cpuUsage > utils.GetThreshold() || memUsage > utils.GetThreshold() {
		serverAddress := utils.GetEdgeServers()[serverName]
		log.Println("server " + serverAddress)
		http.Get(serverAddress+"/stop?url="+utils.GetCloudServer())
	}
}

type ProcessStats struct {
	Pid         int     `json:"pid"`
	CpuUsage    float64 `json:"cpu_usage"`
	MemoryUsage float64 `json:"memory_usage"`
}

type GeneralStats struct {
	Memory string `json:"memory"`
	Cpu string `json:"cpu"`
	Functions map[string]ProcessStats `json:"functions"`
}

// e.POST("/info/server-name)
func getServersStatsInfo(c echo.Context) error {
	stats := GeneralStats{}
	serverName := c.Param("server-name")
	if err := c.Bind(&stats); err != nil {
		log.Printf("Could not save the posted json object %v", err)
		return err
	}

}

type ResponseCallFunction struct {
	Url string `json:"url"`
}

func callFunction(w http.ResponseWriter, r *http.Request){
	ans := "-1"
	for k := range utils.GetEdgeServers(){
		cpu, err := redis_service.GetValue("cpu"+k)
		if err != nil {
			log.Println("k"+k)
			log.Println("cannot get values from redis")
			ans = k
			continue
		}
		mem, err := redis_service.GetValue("mem"+k)
		if err != nil {
			log.Println("cannot get values from redis")
			continue
		}
		if mem == "" || cpu == "" {
			continue
		}
		log.Println("values " + mem + " " + cpu)
		cpuf, _ := strconv.ParseFloat(cpu,32)
		memf, _ := strconv.ParseFloat(mem,32)
		if cpuf < utils.GetThreshold()  && memf < utils.GetThreshold() {
			ans = k
		}
	}
	if ans == "-1" {
		ans = utils.GetCloudServer()
	} else {
		ans = utils.GetEdgeServers()[ans]
	}
	response := ResponseCallFunction{Url: ans}
	json.NewEncoder(w).Encode(&response)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/execute", callFunction).Methods("GET")
	router.HandleFunc("/info/{server-name}", getServersInfo).Methods("GET")
	_ = http.ListenAndServe(":9000", router)
}