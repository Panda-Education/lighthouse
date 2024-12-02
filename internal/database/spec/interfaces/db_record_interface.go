package interfaces

import "Lighthouse/internal/database/models"

type DbRecordInterface interface {
	ToRecord() (models.Record, error)
}
