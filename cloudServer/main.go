package main

import (
	"edgeServer/configuration_service"
	"edgeServer/execution_service"
	"edgeServer/storage_service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

func main() {
	configuration_service.LoadConfiguration()
	storage_service.InitializeStorageHandler(true)

	e := echo.New()
	e.Use(middleware.Recover())

	e.GET("/execute/:fname", executeFunction)
	e.GET("/stop", stopFunction)
	log.Println("new!")

	e.Logger.Fatal(e.Start(":" + configuration_service.GetMyServerPort()))
}

type ResponseExecution struct {
	Result string `json:"result"`
}

func simulateDelay(){
	time.Sleep(time.Duration(rand.Intn(2-1) + 1) * time.Second)
}

func executeFunction(c echo.Context) error {
	simulateDelay()
	log.Printf("Execute %s with parameter %s", c.Param("fname"), c.QueryParam("param"))
	ans := ""
	if configuration_service.GetRuntime() == "wasmer" {
		ans = execution_service.ExecuteAndDetachFunctionWasmer(c.Param("fname"), c.QueryParam("param"))
	} else if configuration_service.GetRuntime() == "docker"  {
		ans = execution_service.ExecuteAndDetachFunctionDocker(c.Param("fname"), c.QueryParam("param"))
	}
	log.Printf("Finished %s sending results", c.Param("fname"))
	if strings.TrimSpace(ans) == "" {
		ans = storage_service.GetValue(c.Param("fname")+"_out")
	}
	response := ResponseExecution{Result: ans}
	return c.JSON(http.StatusOK, response)
}

func stopFunction(c echo.Context) error {
	migrationUri := c.QueryParam("name")
	functions := storage_service.GetAllKeysAndValues()
	keys := make([]string, 0)
	for k, _ := range functions {
		keys = append(keys, k)
	}
	ans := ""
	for _, k := range keys {
		noPost := strings.Split(k, "_")[0]
		if _, ok := functions[noPost + "_out"]; !ok {
			ans = noPost
		}
	}
	if ans == "" {
		return c.String(http.StatusOK, "")
	}
	fName := ans
	pid := functions[fName]
	log.Printf("Stopping function %s with pid %s", fName, pid)
	storage_service.SetValue(fName+"_out", migrationUri)
	if configuration_service.GetRuntime() == "wasmer" {
		execution_service.StopFunction(fName, pid)
	} else if configuration_service.GetRuntime() == "docker"  {
		execution_service.StopFunctionDocker(fName, "")
	}
	return c.String(http.StatusOK, "")
}
