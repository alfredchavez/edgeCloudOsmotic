package main

import (
	"edgeServer/configuration_service"
	"edgeServer/monitoring_service"
	"edgeServer/storage_service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

func main() {
	configuration_service.LoadConfiguration()
	storage_service.InitializeStorageHandler(true)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/execute/:f-name", executeFunction)
	e.GET("/stop/:f-name", stopFunction)

	go monitoring_service.MonitorContext()

	e.Logger.Fatal(e.Start(":" + configuration_service.GetMyServerPort()))
}

func executeFunction(c echo.Context) error {
	return c.String(http.StatusOK, "")
}

func stopFunction(c echo.Context) error {
	return c.String(http.StatusOK, "")
}
