package main

import (
	"edgeServer/configuration_service"
	"edgeServer/execution_service"
	"edgeServer/monitoring_service"
	"edgeServer/storage_service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
	"strings"
)

func main() {
	configuration_service.LoadConfiguration()
	storage_service.InitializeStorageHandler(true)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/execute/:fname", executeFunction)
	e.GET("/stop", stopFunction)
	log.Println("new!")
	go monitoring_service.MonitorContext()

	e.Logger.Fatal(e.Start(":" + configuration_service.GetMyServerPort()))
}

type ResponseExecution struct {
	Result string `json:"result"`
}

func executeFunction(c echo.Context) error {
	log.Printf("Execute %s with parameter %s", c.Param("fname"), c.QueryParam("param"))
	ans := execution_service.ExecuteAndDetachFunctionWasmer(c.Param("fname"), c.QueryParam("param"))
	log.Printf("Finished %s sending results", c.Param("fname"))
	if strings.TrimSpace(ans) == ""{
		ans = configuration_service.GetCloudServerAddr()
	}
	response := ResponseExecution{Result: ans}
	return c.JSON(http.StatusOK, response)
}

func stopFunction(c echo.Context) error {
	fName := c.QueryParam("name")
	pid := c.QueryParam("pid")
	log.Printf("Stopping function %s with pid %s", fName, pid)
	execution_service.StopFunction(fName, pid)
	return c.String(http.StatusOK, "")
}
