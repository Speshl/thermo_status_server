package main

import (
	"ThermoServer/config"
	"ThermoServer/controllers"
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	ctx := context.Background()
	config := config.NewConfig()

	conn, err := pgxpool.Connect(ctx, config.ConnString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	thermoStatusController := controllers.NewThermoStatusController(conn)
	e.POST("/thermostatus", thermoStatusController.PostThermoStatus)
	e.POST("/thermostatus/", thermoStatusController.PostThermoStatus)

	e.GET("/thermostatus/:sourceName", thermoStatusController.GetThermoStatus)

	e.Logger.Fatal(e.Start(":" + config.HTTPPort))
}
