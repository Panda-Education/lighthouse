package database_cache

import (
	"Lighthouse/internal/database/models"
	"Lighthouse/internal/mock"
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestCreateLruDb(t *testing.T) {
	t.Parallel()
	db := CreateLruDb(
		mock.CreateDb(
			mock.DbLatency(time.Millisecond*10),
		),
		4,
	)

	if db == nil {
		t.Errorf("Could not create LruDb")
	}
}

func TestLruDbCache_InsertOneRecord(t *testing.T) {
	//t.Parallel()
	db := CreateLruDb(
		mock.CreateDb(
			mock.DbLatency(time.Millisecond*10),
		),
		4,
	)

	if db == nil {
		t.Errorf("Could not create LruDb")
		return
	}

	record, _ := models.CreateRecordFromString(
		"https://www.google.com",
		"google",
		"google_uid",
	)

	ctx := context.Background()

	// Expect insert of one record
	if err := db.InsertRecord(ctx, record); err != nil {
		t.Errorf("Could not insert record: %s", err)
		return
	}

	// Expect retrieval of record
	if res, err := db.FindRecord(ctx, record.Id); err != nil || res != nil {
		if err != nil {
			t.Errorf("Could not find record: %s", err)
			return
		}
		if res != record {
			t.Errorf("Mismatching records")
			return
		}
	}

	// Expect disconnect
	if err := db.Disconnect(ctx); err != nil {
		t.Errorf("Unable to disconnect: %s", err)
	}

	// Validate DB action history
	expectedDbHistory := []mock.DbAction{
		mock.InsertRecordAttempt,
		mock.ConnectAttempt,
		mock.ConnectSuccess,
		mock.InsertRecordSuccess,
		mock.DisconnectAttempt,
		mock.DisconnectSuccess,
	}

	if !reflect.DeepEqual(expectedDbHistory, db.db.(*mock.Db).History) {
		t.Errorf("History does not match. Got: %v, Expected: %v", db.db.(*mock.Db).History, expectedDbHistory)
	}

}

func TestLruDbCache_InsertCapacityPlusOne(t *testing.T) {
	t.Parallel()

	capacity := 4
	n_records := capacity + 1

	db := CreateLruDb(
		mock.CreateDb(
			mock.DbLatency(time.Millisecond*10),
		),
		capacity,
	)

	if db == nil {
		t.Errorf("Could not create LruDb")
		return
	}

	ctx := context.Background()

	// Expect insert n+1 records

	var records []*models.Record

	for i := 0; i < n_records; i++ {
		record, _ := models.CreateRecordFromString(
			"https://www.google.com",
			fmt.Sprintf("google_%v", i),
			fmt.Sprintf("google_uid_%v", i),
		)

		records = append(
			records,
			record,
		)
	}

	for i := 0; i < n_records-1; i++ {
		if err := db.InsertRecord(ctx, records[i]); err != nil {
			t.Errorf("Could not insert record: %s", err)
			return
		}

	}

	// Expect retrieval of first record
	// All these should not trigger any Db look up
	if res, err := db.FindRecord(ctx, fmt.Sprintf("google_%v", 0)); err != nil || res != nil {
		if err != nil {
			t.Errorf("Could not find record: %s", err)
			return
		}
		if res != records[0] {
			t.Errorf("Mismatching records")
			return
		}
	}

	// Insert last record and evict one record
	if err := db.InsertRecord(ctx, records[n_records-1]); err != nil {
		t.Errorf("Could not insert record: %s", err)
		return
	}

	// Find most recently evicted record
	if res, err := db.FindRecord(ctx, fmt.Sprintf("google_%v", 1)); err != nil || res != nil {
		if err != nil {
			t.Errorf("Could not find record: %s", err)
			return
		}
		if res != records[1] {
			t.Errorf("Mismatching records")
			return
		}
	}

	// Expect disconnect
	if err := db.Disconnect(ctx); err != nil {
		t.Errorf("Unable to disconnect: %s", err)
	}

	// Validate DB action history
	expectedDbHistory := []mock.DbAction{
		mock.InsertRecordAttempt,
		mock.ConnectAttempt,
		mock.ConnectSuccess,
		mock.InsertRecordSuccess,

		mock.InsertRecordAttempt,
		mock.InsertRecordSuccess,

		mock.InsertRecordAttempt,
		mock.InsertRecordSuccess,

		mock.InsertRecordAttempt,
		mock.InsertRecordSuccess,

		mock.InsertRecordAttempt,
		mock.InsertRecordSuccess,

		mock.FindRecordAttempt,
		mock.FindRecordSuccess,

		mock.DisconnectAttempt,
		mock.DisconnectSuccess,
	}

	if !reflect.DeepEqual(expectedDbHistory, db.db.(*mock.Db).History) {
		t.Errorf("History does not match. Got: %v, Expected: %v", db.db.(*mock.Db).History, expectedDbHistory)
	}

}
