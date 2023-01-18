package controllers

import (
	"ThermoServer/models"
	"ThermoServer/stores"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
)

type ThermoStatusControllerIFace interface {
	GetThermoStatus(echo.Context) error
	PostThermoStatus(echo.Context) error
}

type ThermoStatusController struct {
	store *stores.ThermoStatusStore
}

func NewThermoStatusController(conn *pgxpool.Pool) *ThermoStatusController {
	return &ThermoStatusController{
		store: stores.NewThermoStatusStore(conn),
	}
}

// e.GET("/users/:id", getUser)
func (tsc *ThermoStatusController) GetThermoStatus(c echo.Context) error {
	var err error
	defer func() {
		if err != nil {
			log.Println(err)
		}
	}()

	log.Println("GetThermoStatus Request Recieved")
	/*eventTime, err := time.Parse("2006-01-02T15:04:05", c.Param("eventTime"))
	if err != nil {
		return err
	}*/
	sourceName := c.Param("sourceName")

	status, err := tsc.store.SelectThermoStatus(c.Request().Context(), sourceName)
	if err != nil {
		return err
	}

	if status == nil {
		return c.JSON(http.StatusNotFound, nil)
	}

	return c.JSON(http.StatusOK, *status)
}

func (tsc *ThermoStatusController) PostThermoStatus(c echo.Context) error {
	log.Println("PostThermoStatus Request Recieved")
	var err error
	defer func() {
		if err != nil {
			log.Println(err)
		}
	}()

	status := models.ThermoStatus{}
	err = c.Bind(&status)
	if err != nil {
		return err
	}

	if status.EventTime.IsZero() {
		status.EventTime = time.Now()
	}

	err = tsc.store.InsertThermoStatus(c.Request().Context(), status)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, status)
}
