package gorm_pg_adapter

import (
	"Lighthouse/internal/database/models"
	"context"
	"errors"
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

	db *gorm.DB
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

func (g *GormPgAdapter) ensureDbConnection(ctx context.Context) error {

	if g.db == nil {
		if conErr := g.Connect(ctx); conErr != nil {
			return conErr
		}
	}

	return nil
}

func (g *GormPgAdapter) Connect(ctx context.Context) error {
	db, err := gorm.Open(postgres.Open(g.createDsnString()), &gorm.Config{})
	if err != nil {
		return err
	}
	g.db = db
	return nil
}

func (g *GormPgAdapter) Disconnect(ctx context.Context) error {
	if g.db == nil {
		return nil // Nothing to disconnect
	}

	// Get the underlying SQL DB
	sqlDB, err := g.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying SQL DB: %w", err)
	}

	// Close the database connection
	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close the database connection: %w", err)
	}

	g.db = nil // Set to nil to indicate the connection is closed
	return nil
}

func (g *GormPgAdapter) InsertRecord(
	ctx context.Context,
	record *models.Record,
) error {
	if err := g.ensureDbConnection(ctx); err != nil {
		return err
	}

	gormPgRecord, conversionErr := ConvertRecordToDbRecord(record)
	if conversionErr != nil {
		return conversionErr
	}

	result := g.db.Create(gormPgRecord)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (g *GormPgAdapter) UpdateRecord(
	ctx context.Context,
	record *models.Record,
) error {
	if err := g.ensureDbConnection(ctx); err != nil {
		return err
	}

	return nil
}

func (g *GormPgAdapter) DeleteRecord(
	ctx context.Context,
	record *models.Record,
) error {
	if err := g.ensureDbConnection(ctx); err != nil {
		return err
	}

	return nil
}

func (g *GormPgAdapter) FindRecord(
	ctx context.Context,
	id string,
) (*models.Record, error) {
	if err := g.ensureDbConnection(ctx); err != nil {
		return nil, err
	}

	var dbRecord GormPgRecord
	result := g.db.Where("id = ?", id).First(&dbRecord)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, result.Error
		}
	}

	r, convertErr := dbRecord.ToRecord()
	if convertErr != nil {
		return nil, convertErr
	}

	return r, nil
}

func (g *GormPgAdapter) Migrate(ctx context.Context) error {
	if err := g.ensureDbConnection(ctx); err != nil {
		return err
	}

	if err := g.db.
		WithContext(ctx).
		AutoMigrate(
			&GormPgRecord{},
		); err != nil {
		return err
	}

	return nil
}
