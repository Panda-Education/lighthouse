package interfaces

import (
	"Lighthouse/internal/database/models"
	"context"
)

type DatabaseConnectorStrategy interface {
	Connect(ctx context.Context) error
	Disconnect(ctx context.Context) error
	InsertRecord(ctx context.Context, record *models.Record) error
	UpdateRecord(ctx context.Context, record *models.Record) error
	DeleteRecord(ctx context.Context, record *models.Record) error
	FindRecord(ctx context.Context, id string) (*models.Record, error)
	Migrate(ctx context.Context) error
}
