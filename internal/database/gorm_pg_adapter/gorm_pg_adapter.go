package gorm_pg_adapter

import (
	"Lighthouse/internal/database/models"
	"context"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// GormPgAdapter implements DatabaseConnectorStrategy
type GormPgAdapter struct {
	host     string
	user     string
	password string
	port     int
	dbname   string
	sslMode  string
	timezone string
}

func (g *GormPgAdapter) createDsnString() string {
	return fmt.Sprintf(
		"host=%v user=%v password=%v dbname=%v port=%v sslmode=%v TimeZone=%v",
		g.host,
		g.user,
		g.password,
		g.dbname,
		g.port,
		g.sslMode,
		g.timezone,
	)
}

func (g *GormPgAdapter) Connect(ctx context.Context) error {
	return nil
}

func (g *GormPgAdapter) Disconnect(ctx context.Context) error {
	return nil
}

func (g *GormPgAdapter) InsertRecord(
	ctx context.Context,
	record *models.Record,
) error {
	return nil
}

func (g *GormPgAdapter) DeleteRecord(
	ctx context.Context,
	record *models.Record,
) error {
	return nil
}

func (g *GormPgAdapter) FindRecord(
	ctx context.Context,
	id string,
) (*models.Record, error) {
	return &models.Record{}, nil
}

func (g *GormPgAdapter) Migrate(ctx context.Context) error {
	db, err := gorm.Open(postgres.Open(g.createDsnString()), &gorm.Config{})
	if err != nil {
		return err
	}

	err = db.
		WithContext(ctx).
		AutoMigrate(
			&GormPgRecord{},
		)
	if err != nil {
		return err
	}

	return nil
}
