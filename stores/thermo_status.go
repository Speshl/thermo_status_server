package stores

import (
	"ThermoServer/models"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type ThermoStatusStoreIface interface {
	InsertThermoStatus(models.ThermoStatus) error
	SelectThermoStatus(time.Time) (*models.ThermoStatus, error)
}

type ThermoConnIFace interface {
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
}

type ThermoStatusStore struct {
	conn ThermoConnIFace
}

func NewThermoStatusStore(conn ThermoConnIFace) *ThermoStatusStore {
	return &ThermoStatusStore{
		conn: conn,
	}
}

func (s *ThermoStatusStore) InsertThermoStatus(ctx context.Context, status models.ThermoStatus) error {
	log.Printf("Inserting Thermo Status: %+v", status)
	queryString := `
	INSERT INTO status (
		event_time, source_name, heat_on, inside_temp, inside_humidity,
		inside_heat_index,outside_temp, diff_temp
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);
	`

	_, err := s.conn.Exec(ctx, queryString,
		status.EventTime,
		status.ThermoName.SourceName,
		status.HeatOn,
		status.InsideTemp,
		status.InsideHumidity,
		status.InsideHeatIndex,
		status.OutsideTemp,
		status.DiffTemp,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *ThermoStatusStore) SelectThermoStatus(ctx context.Context, sourceName string) (*models.ThermoStatus, error) {
	log.Printf("Selecting Thermo Status For Source: %s", sourceName)
	model := models.ThermoStatus{}

	queryString := fmt.Sprintf(`
	SELECT 
		event_time, source_name, heat_on, inside_temp, inside_humidity,
		inside_heat_index, outside_temp, diff_temp
	FROM public.status
	WHERE source_name = '%s'
	ORDER BY event_time DESC
	LIMIT 1;
	`, sourceName)

	err := s.conn.QueryRow(ctx, queryString).Scan(
		&model.EventTime,
		&model.SourceName,
		&model.HeatOn,
		&model.InsideTemp,
		&model.InsideHumidity,
		&model.InsideHeatIndex,
		&model.OutsideTemp,
		&model.DiffTemp,
	)
	if err != nil {
		return nil, err
	}

	if model.EventTime.IsZero() {
		log.Println("Thermo Status Not Found")
		return nil, nil
	}

	log.Printf("Selected Thermo Status: %+v", model)
	return &model, nil
}
