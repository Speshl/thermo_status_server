package models

import "time"

type ThermoStatus struct {
	EventTime      time.Time `field:"event_time"`
	SourceName     string    `field:"source_name"`
	Enabled        bool      `field:"enabled"`
	InsideTemp     int       `field:"inside_temp"`
	OutsideTemp    int       `field:"outside_temp"`
	DiffTemp       int       `field:"diff_temp"`
	TargetDiffTemp int       `field:"target_diff_temp"`
	TargetTemp     int       `field:"target_temp"`
}
