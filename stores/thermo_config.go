package stores

import (
	"ThermoServer/models"
	"context"
	"fmt"
	"log"
)

type ThermoConfigStoreIface interface {
	InsertThermoConfig(models.ThermoConfig) error
	SelectThermoConfig(string) (*models.ThermoConfig, error)
}

type ThermoConfigStore struct {
	conn ThermoConnIFace
}

func NewThermoConfigStore(conn ThermoConnIFace) *ThermoConfigStore {
	return &ThermoConfigStore{
		conn: conn,
	}
}

func (s *ThermoConfigStore) UpsertThermoConfig(ctx context.Context, config models.ThermoConfig) error {
	log.Printf("Inserting Thermo Config: %+v", config)
	queryString := `
	INSERT INTO config (
		source_name, enabled, target_diff_temp, target_temp, target_bypass_offset, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6)
	ON CONFLICT (source_name)
	DO UPDATE SET
		enabled = $7,
		target_diff_temp = $8,
		target_temp = $9,
		target_bypass_offset = $10,
		updated_at = $11;
	`

	_, err := s.conn.Exec(ctx, queryString,
		config.ThermoName.SourceName,
		config.Enabled,
		config.TargetDiffTemp,
		config.TargetTemp,
		config.TargetBypassOffset,
		config.UpdatedAt,
		config.Enabled,
		config.TargetDiffTemp,
		config.TargetTemp,
		config.TargetBypassOffset,
		config.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *ThermoConfigStore) SelectThermoConfig(ctx context.Context, sourceName string) (*models.ThermoConfig, error) {
	log.Printf("Selecting Thermo Config For Source: %s", sourceName)
	model := models.ThermoConfig{}

	queryString := fmt.Sprintf(`
	SELECT 
	source_name, enabled, target_diff_temp, target_temp, target_bypass_offset, updated_at
	FROM config
	WHERE source_name = '%s'
	ORDER BY updated_at DESC
	LIMIT 1;
	`, sourceName)

	err := s.conn.QueryRow(ctx, queryString).Scan(
		&model.SourceName,
		&model.Enabled,
		&model.TargetDiffTemp,
		&model.TargetTemp,
		&model.TargetBypassOffset,
		&model.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if model.UpdatedAt.IsZero() {
		log.Println("Thermo Status Not Found")
		return nil, nil
	}

	log.Printf("Selected Thermo Config: %+v", model)
	return &model, nil
}
