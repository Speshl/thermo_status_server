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

type ThermoControllerIFace interface {
	GetThermoStatus(echo.Context) error
	PostThermoStatus(echo.Context) error
	GetThermoConfig(echo.Context) error
	PostThermoConfig(echo.Context) error
}

type ThermoController struct {
	statusStore *stores.ThermoStatusStore
	configStore *stores.ThermoConfigStore
}

func NewThermoController(conn *pgxpool.Pool) *ThermoController {
	return &ThermoController{
		statusStore: stores.NewThermoStatusStore(conn),
		configStore: stores.NewThermoConfigStore(conn),
	}
}

func (tsc *ThermoController) PostThermoSync(c echo.Context) error {
	log.Println("PostThermoSync Request Recieved")
	var err error
	defer func() {
		if err != nil {
			log.Println(err)
		}
	}()

	status := models.ThermoFull{}
	err = c.Bind(&status)
	if err != nil {
		return err
	}

	status.ThermoConfig.SourceName = status.SourceName
	status.ThermoStatus.SourceName = status.SourceName

	if status.UpdatedAt.IsZero() {
		status.UpdatedAt = time.Now()
	}

	if status.EventTime.IsZero() {
		status.EventTime = time.Now()
	}

	log.Printf("\n\nThermoFull: %+v\n\n", status)

	if status.LocallyUpdated {
		err = tsc.configStore.UpsertThermoConfig(c.Request().Context(), status.ThermoConfig)
		if err != nil {
			return err
		}
	} else {
		thermoConfig, err := tsc.configStore.SelectThermoConfig(c.Request().Context(), status.SourceName)
		if err != nil {
			return err
		}
		if thermoConfig != nil {
			status.ThermoConfig = *thermoConfig
		}
	}

	err = tsc.statusStore.InsertThermoStatus(c.Request().Context(), status.ThermoStatus)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, status)
}

// e.GET("/users/:id", getUser)
func (tsc *ThermoController) GetThermoStatus(c echo.Context) error {
	var err error
	defer func() {
		if err != nil {
			log.Println(err)
		}
	}()

	log.Println("GetThermoStatus Request Recieved")
	sourceName := c.Param("sourceName")

	status, err := tsc.statusStore.SelectThermoStatus(c.Request().Context(), sourceName)
	if err != nil {
		return err
	}

	if status == nil {
		return c.JSON(http.StatusNotFound, nil)
	}

	return c.JSON(http.StatusOK, *status)
}

func (tsc *ThermoController) PostThermoStatus(c echo.Context) error {
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

	err = tsc.statusStore.InsertThermoStatus(c.Request().Context(), status)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, status)
}

func (tsc *ThermoController) GetThermoConfig(c echo.Context) error {
	var err error
	defer func() {
		if err != nil {
			log.Println(err)
		}
	}()

	log.Println("GetThermoConfig Request Recieved")
	sourceName := c.Param("sourceName")

	config, err := tsc.configStore.SelectThermoConfig(c.Request().Context(), sourceName)
	if err != nil {
		return err
	}

	if config == nil {
		return c.JSON(http.StatusNotFound, nil)
	}

	return c.JSON(http.StatusOK, *config)
}

func (tsc *ThermoController) PostThermoConfig(c echo.Context) error {
	log.Println("PostThermoConfig Request Recieved")
	var err error
	defer func() {
		if err != nil {
			log.Println(err)
		}
	}()

	config := models.ThermoConfig{}
	err = c.Bind(&config)
	if err != nil {
		return err
	}

	if config.UpdatedAt.IsZero() {
		config.UpdatedAt = time.Now()
	}

	err = tsc.configStore.UpsertThermoConfig(c.Request().Context(), config)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, config)
}

func (tsc *ThermoController) GetThermoFull(c echo.Context) error {
	var err error
	defer func() {
		if err != nil {
			log.Println(err)
		}
	}()

	log.Println("GetThermoFull Request Recieved")
	sourceName := c.Param("sourceName")

	config, err := tsc.configStore.SelectThermoConfig(c.Request().Context(), sourceName)
	if err != nil {
		return err
	}

	status, err := tsc.statusStore.SelectThermoStatus(c.Request().Context(), sourceName)
	if err != nil {
		return err
	}

	full := models.MakeFull(status, config)
	if full == nil {
		return c.JSON(http.StatusNotFound, nil)
	}

	return c.JSON(http.StatusOK, *full)
}
