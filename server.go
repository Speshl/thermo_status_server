package main

import (
	"ThermoServer/config"
	"ThermoServer/controllers"
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	ctx := context.Background()

	err := godotenv.Load("pi.env")
	if err != nil {
		log.Println("Skipped loading .env file")
	}

	config := config.NewConfig()

	conn, err := pgxpool.Connect(ctx, config.ConnString)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	thermoController := controllers.NewThermoController(conn)

	e.POST("/thermo/sync", thermoController.PostThermoSync)

	e.POST("/thermo/config", thermoController.PostThermoConfig)
	e.GET("/thermo/config/:sourceName", thermoController.GetThermoConfig)

	e.POST("/thermo/status", thermoController.PostThermoStatus)
	e.GET("/thermo/status/:sourceName", thermoController.GetThermoStatus)

	e.GET("/thermo/full/:sourceName", thermoController.GetThermoFull)

	e.Logger.Fatal(e.Start(":" + config.HTTPPort))
}
