package mock

import (
	"Lighthouse/internal/database/models"
	"context"
	"errors"
	"fmt"
	"time"
)

/*
Db tracing history
*/

type DbAction int

const (
	ConnectAttempt DbAction = iota
	ConnectSuccess
	DisconnectAttempt
	DisconnectSuccess
	InsertRecordAttempt
	InsertRecordSuccess
	UpdateRecordAttempt
	UpdateRecordSuccess
	DeleteRecordAttempt
	DeleteRecordSuccess
	FindRecordAttempt
	FindRecordSuccess
	MigrateAttempt
	MigrateSuccess
)

/*
Db Definition
*/

type Db struct {
	records map[string]*models.Record

	// used to validating transactions and interactions
	// with middleware (like LruCache)
	History []DbAction

	// status values
	Connected bool

	// testing hyper parameters
	Latency time.Duration

	// named errors to test
	WillConnectError bool
}

/*
Db Option and creation methods
*/

type DbOption func(db *Db)

func DbWillConnectError() DbOption {
	return func(db *Db) {
		db.WillConnectError = true
	}
}

func DbLatency(latency time.Duration) DbOption {
	return func(db *Db) {
		db.Latency = latency
	}
}

func CreateDb(options ...DbOption) *Db {
	db := &Db{
		records:          make(map[string]*models.Record),
		History:          []DbAction{},
		Connected:        false,
		Latency:          time.Millisecond * 100,
		WillConnectError: false,
	}
	for _, opt := range options {
		opt(db)
	}
	return db
}

/*
Db Helper methods
*/

func (d *Db) addAction(action DbAction) {
	d.History = append(d.History, action)
}

/*
Db interface implementation
*/

func (d *Db) Connect(ctx context.Context) error {
	d.addAction(ConnectAttempt)

	// Check if the context is canceled
	if ctx.Err() != nil {
		return fmt.Errorf("context error before connecting: %w", ctx.Err())
	}

	// Simulate connection Latency
	select {
	case <-time.After(d.Latency):
		// After waiting for the Latency, check if we still have a valid context
		if ctx.Err() != nil {
			return fmt.Errorf("context error during connection: %w", ctx.Err())
		}
	case <-ctx.Done():
		// The context was canceled or timed out during our "wait"
		return fmt.Errorf("context canceled or timed out: %w", ctx.Err())
	}

	// Check if we should simulate a connection error
	if d.WillConnectError {
		return errors.New("unable to connect to database")
	}

	// Mark as Connected
	d.Connected = true
	d.addAction(ConnectSuccess)
	return nil
}

func (d *Db) Disconnect(ctx context.Context) error {
	d.addAction(DisconnectAttempt)

	// Check if the context is canceled
	if ctx.Err() != nil {
		return fmt.Errorf("context error before disconnecting: %w", ctx.Err())
	}

	// Simulate Latency
	select {
	case <-time.After(d.Latency):
		// After waiting for the Latency, check if we still have a valid context
		if ctx.Err() != nil {
			return fmt.Errorf("context error during disconnection: %w", ctx.Err())
		}
	case <-ctx.Done():
		// The context was canceled or timed out during our "wait"
		return fmt.Errorf("context canceled or timed out: %w", ctx.Err())
	}

	// Mark as Connected
	d.Connected = false
	d.addAction(DisconnectSuccess)
	return nil
}

func (d *Db) InsertRecord(ctx context.Context, record *models.Record) error {
	d.addAction(InsertRecordAttempt)

	// Check if the context is canceled
	if ctx.Err() != nil {
		return fmt.Errorf("context error before inserting: %w", ctx.Err())
	}

	if !d.Connected {
		if err := d.Connect(ctx); err != nil {
			return err
		}
	}

	// Simulate Latency
	select {
	case <-time.After(d.Latency):
		if ctx.Err() != nil {
			return fmt.Errorf("context error during Latency: %w", ctx.Err())
		}
	case <-ctx.Done():
		return fmt.Errorf("context canceled or timed out: %w", ctx.Err())
	}

	if _, exists := d.records[record.Id]; exists {
		return errors.New("record already exists")
	}

	d.records[record.Id] = record
	d.addAction(InsertRecordSuccess)
	return nil
}

func (d *Db) UpdateRecord(ctx context.Context, record *models.Record) error {
	d.addAction(UpdateRecordAttempt)

	// Check if the context is canceled
	if ctx.Err() != nil {
		return fmt.Errorf("context error before updating: %w", ctx.Err())
	}

	if !d.Connected {
		if err := d.Connect(ctx); err != nil {
			return err
		}
	}

	// Simulate Latency
	select {
	case <-time.After(d.Latency):
		if ctx.Err() != nil {
			return fmt.Errorf("context error during Latency: %w", ctx.Err())
		}
	case <-ctx.Done():
		return fmt.Errorf("context canceled or timed out: %w", ctx.Err())
	}

	_, exists := d.records[record.Id]

	if !exists {
		return fmt.Errorf("record %s does not exist", record.Id)
	}

	d.records[record.Id] = record
	d.addAction(UpdateRecordSuccess)
	return nil
}

func (d *Db) DeleteRecord(ctx context.Context, id string) error {
	d.addAction(DeleteRecordAttempt)

	// Check if the context is canceled
	if ctx.Err() != nil {
		return fmt.Errorf("context error before deleting: %w", ctx.Err())
	}

	if !d.Connected {
		if err := d.Connect(ctx); err != nil {
			return err
		}
	}

	// Simulate Latency
	select {
	case <-time.After(d.Latency):
		if ctx.Err() != nil {
			return fmt.Errorf("context error during Latency: %w", ctx.Err())
		}
	case <-ctx.Done():
		return fmt.Errorf("context canceled or timed out: %w", ctx.Err())
	}

	_, exists := d.records[id]

	if !exists {
		return fmt.Errorf("record %s does not exist", id)
	}

	delete(d.records, id)
	d.addAction(DeleteRecordSuccess)
	return nil
}

func (d *Db) FindRecord(ctx context.Context, id string) (*models.Record, error) {
	d.addAction(FindRecordAttempt)

	// Check if the context is canceled
	if ctx.Err() != nil {
		return &models.Record{}, fmt.Errorf("context error before finding: %w", ctx.Err())
	}

	if !d.Connected {
		if err := d.Connect(ctx); err != nil {
			return &models.Record{}, err
		}
	}

	// Simulate Latency
	select {
	case <-time.After(d.Latency):
		if ctx.Err() != nil {
			return &models.Record{}, fmt.Errorf("context error during Latency: %w", ctx.Err())
		}
	case <-ctx.Done():
		return &models.Record{}, fmt.Errorf("context canceled or timed out: %w", ctx.Err())
	}

	record, exists := d.records[id]

	if !exists {
		return &models.Record{}, fmt.Errorf("record %s does not exist", id)
	}

	d.addAction(FindRecordSuccess)

	return record, nil
}

func (d *Db) Migrate(ctx context.Context) error {
	d.addAction(MigrateAttempt)

	// Check if the context is canceled
	if ctx.Err() != nil {
		return fmt.Errorf("context error before migrating: %w", ctx.Err())
	}

	if !d.Connected {
		if err := d.Connect(ctx); err != nil {
			return err
		}
	}

	// Simulate Latency
	select {
	case <-time.After(d.Latency):
		if ctx.Err() != nil {
			return fmt.Errorf("context error during Latency: %w", ctx.Err())
		}
	case <-ctx.Done():
		return fmt.Errorf("context canceled or timed out: %w", ctx.Err())
	}

	// migration assumed to be ok
	d.addAction(MigrateSuccess)
	return nil
}
