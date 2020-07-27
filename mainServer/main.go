package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"mainServer/storage_service"
	"mainServer/utils"
	"net/http"
	"strconv"
)

func main() {
	storage_service.InitializeStorageHandler(true)

	e := echo.New()
	e.Use(middleware.Recover())

	e.POST("/info/:server-name", getServersStatsInfo)
	e.GET("/query_execute", callFunction)

	e.Logger.Fatal(e.Start(":9000"))

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


func getServersStatsInfo(c echo.Context) error {
	stats := GeneralStats{}
	serverName := c.Param("server-name")
	if err := c.Bind(&stats); err != nil {
		log.Printf("Could not save the posted json object %v", err)
		return err
	}
	log.Printf("Stats got from %s are %v", serverName, stats)
	storage_service.SetValue(serverName+"_mem", stats.Memory)
	storage_service.SetValue(serverName+"_cpu", stats.Cpu)
	floatMemory,_ := strconv.ParseFloat(stats.Memory, 64)
	floatCpu, _ := strconv.ParseFloat(stats.Cpu, 64)
	if floatMemory > utils.GetThreshold() || floatCpu > utils.GetThreshold(){
		addrServer := utils.GetEdgeServers()[serverName]
		// stop and migrate most expensive function
		var fPid int
		fName := "null"
		most := -1.0
		for k, v := range stats.Functions {
			if utils.DetermineFunctionWight(v.CpuUsage, v.MemoryUsage) > most {
				fName = k
				fPid = v.Pid
				most = utils.DetermineFunctionWight(v.CpuUsage, v.MemoryUsage)
			}
		}
		if fName == "null" {
			return nil
		}
		sPid := strconv.Itoa(fPid)
		_, _ = http.Get(addrServer + "/stop?" + "name=" + fName + "&pid=" + sPid)
	}
	return nil
}

type ResponseCallFunction struct {
	Url string `json:"url"`
}

func callFunction(c echo.Context) error {
	ans := "-1"
	weight := 10000.0
	for k := range utils.GetEdgeServers(){
		cpu := storage_service.GetValue(k + "_cpu")
		mem := storage_service.GetValue(k + "_mem")
		log.Println("cpu "+ cpu + " mem " + mem)
		if mem == "" || cpu == "" {
			continue
		}
		cpuf, _ := strconv.ParseFloat(cpu,32)
		memf, _ := strconv.ParseFloat(mem,32)
		if cpuf < utils.GetThreshold()  && memf < utils.GetThreshold() {
			if utils.DetermineFunctionWight(cpuf, memf) < weight {
				ans = k
				weight = utils.DetermineFunctionWight(cpuf, memf)
			}
		}
	}
	if ans == "-1" {
		ans = utils.GetCloudServer()
	} else {
		ans = utils.GetEdgeServers()[ans]
	}
	response := ResponseCallFunction{Url: ans}
	return c.JSON(http.StatusOK, response)
}