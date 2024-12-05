package database_cache

import (
	"Lighthouse/internal/database/models"
	"Lighthouse/internal/database/spec/interfaces"
	"container/list"
	"context"
	"fmt"
)

type LruDbCache struct {
	db interfaces.DatabaseConnectorStrategy

	capacity int
	cache    map[string]*list.Element
	list     *list.List
}

func CreateLruDb(
	db interfaces.DatabaseConnectorStrategy,
	capacity int,
) *LruDbCache {
	return &LruDbCache{
		db:       db,
		capacity: capacity,
		cache:    make(map[string]*list.Element),
		list:     list.New(),
	}
}

func (l *LruDbCache) Connect(ctx context.Context) error {
	return l.db.Connect(ctx)
}

func (l *LruDbCache) Disconnect(ctx context.Context) error {
	return l.db.Disconnect(ctx)
}

func (l *LruDbCache) InsertRecord(
	ctx context.Context,
	record *models.Record,
) error {

	if insertErr := l.db.InsertRecord(ctx, record); insertErr != nil {
		return insertErr
	}

	if elem, ok := l.cache[record.Id]; ok {
		l.list.MoveToFront(elem)
		return nil
	}

	insertedElem := l.list.PushFront(record)
	l.cache[record.Id] = insertedElem

	if l.list.Len() == l.capacity+1 {
		evictedElem := l.list.Back()
		delete(l.cache, evictedElem.Value.(*models.Record).Id)
	}

	return nil

}

func (l *LruDbCache) UpdateRecord(
	ctx context.Context,
	record *models.Record,
) error {

	if updateErr := l.db.UpdateRecord(ctx, record); updateErr != nil {
		return updateErr
	}

	// Update the record in the cache
	if elem, ok := l.cache[record.Id]; ok {
		if cachedRecord, ok := elem.Value.(*models.Record); ok {
			// Update all fields of the cached record
			*cachedRecord = *record
		} else {
			// Log an error or handle unexpected type
			return fmt.Errorf("cache contains an invalid type for key %v", record.Id)
		}
	}

	return nil
}

func (l *LruDbCache) DeleteRecord(
	ctx context.Context,
	id string,
) error {
	if err := l.db.DeleteRecord(ctx, id); err != nil {
		return err
	}

	// Delete the record in the cache
	if elem, ok := l.cache[id]; ok {
		if _, ok := elem.Value.(*models.Record); ok {
			// Delete record
			delete(l.cache, id)
			l.list.Remove(elem)
		} else {
			// Log an error or handle unexpected type
			return fmt.Errorf("cache contains an invalid type for key %v", id)
		}
	}

	return nil

}

func (l *LruDbCache) FindRecord(
	ctx context.Context,
	id string,
) (*models.Record, error) {

	// Find the record in the cache
	if elem, ok := l.cache[id]; ok {
		if _, ok := elem.Value.(*models.Record); ok {
			l.list.MoveToFront(elem)
			return elem.Value.(*models.Record), nil
		} else {
			// Log an error or handle unexpected type
			return nil, fmt.Errorf("cache contains an invalid type for key %v", id)
		}
	}

	record, err := l.db.FindRecord(ctx, id)
	if err != nil {
		return nil, err
	}

	elem := l.list.PushFront(record)
	l.cache[record.Id] = elem

	return record, nil
}

func (l *LruDbCache) Migrate(ctx context.Context) error {
	return l.db.Migrate(ctx)
}
