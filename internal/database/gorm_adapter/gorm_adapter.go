package gorm_adapter

import (
	"Lighthouse/internal/database/models"
	"context"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// GormDbAdapter implements DatabaseConnectorStrategy
type GormDbAdapter struct {
	host     string
	user     string
	password string
	port     int
	dbname   string
	sslMode  string
	timezone string
}

func (g *GormDbAdapter) createDsnString() string {
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

func (g *GormDbAdapter) Connect(ctx context.Context) error {
	return nil
}

func (g *GormDbAdapter) Disconnect(ctx context.Context) error {
	return nil
}

func (g *GormDbAdapter) InsertRecord(
	ctx context.Context,
	record models.Record,
) error {
	return nil
}

func (g *GormDbAdapter) DeleteRecord(
	ctx context.Context,
	record models.Record,
) error {
	return nil
}

func (g *GormDbAdapter) FindRecord(
	ctx context.Context,
	id string,
) (models.Record, error) {
	return models.Record{}, nil
}

func (g *GormDbAdapter) Migrate(ctx context.Context) error {
	db, err := gorm.Open(postgres.Open(g.createDsnString()), &gorm.Config{})
	if err != nil {
		return err
	}

	err = db.WithContext(ctx).AutoMigrate(&GormRecord{})
	if err != nil {
		return err
	}

	return nil
}
