package database_cache

import "Lighthouse/internal/database/spec/interfaces"

func Lru(db *interfaces.DatabaseConnectorStrategy) *interfaces.DatabaseConnectorStrategy {
	return db
}
