package models

import "time"

type ThermoName struct {
	SourceName string `field:"source_name"`
}

type ThermoStatus struct {
	ThermoName
	EventTime       time.Time `field:"event_time"`
	HeatOn          bool      `filed:"heat_on"`
	InsideTemp      int       `field:"inside_temp"`
	InsideHumidity  int       `field:"inside_humidity"`
	InsideHeatIndex int       `field:"inside_heat_index"`
	OutsideTemp     int       `field:"outside_temp"`
	DiffTemp        int       `field:"diff_temp"`
}

type ThermoConfig struct {
	ThermoName
	Enabled            bool      `field:"enabled"`
	TargetDiffTemp     int       `field:"target_diff_temp"`
	TargetTemp         int       `field:"target_temp"`
	TargetBypassOffset int       `filed:"target_bypass_offset"`
	LocallyUpdated     bool      `filed:"locally_updated"`
	UpdatedAt          time.Time `field:"updated_at"`
}

type ThermoFull struct {
	ThermoName
	ThermoStatus
	ThermoConfig
}

func MakeFull(status *ThermoStatus, config *ThermoConfig) *ThermoFull {
	if status == nil && config == nil {
		return nil
	}

	returnValue := ThermoFull{}
	if status != nil {
		returnValue.SourceName = status.SourceName
		returnValue.HeatOn = status.HeatOn
		returnValue.DiffTemp = status.DiffTemp
		returnValue.EventTime = status.EventTime
		returnValue.InsideHeatIndex = status.InsideHeatIndex
		returnValue.InsideHumidity = status.InsideHeatIndex
		returnValue.InsideTemp = status.InsideTemp
		returnValue.OutsideTemp = status.OutsideTemp
	}

	if config != nil {
		returnValue.SourceName = status.SourceName
		returnValue.Enabled = config.Enabled
		returnValue.TargetDiffTemp = config.TargetDiffTemp
		returnValue.TargetTemp = config.TargetTemp
		returnValue.TargetBypassOffset = config.TargetBypassOffset
		returnValue.UpdatedAt = config.UpdatedAt
	}

	return &returnValue
}
