package sync

import (
	"fmt"
	"log"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

const (
	PBSyncCollectionName = "_pbSync"
)

// EnsurePBSyncCollection creates the _pbSync system collection if it doesn't exist
func EnsurePBSyncCollection(app *pocketbase.PocketBase) error {
	// Check if collection exists by trying to find it
	collection, err := app.FindCollectionByNameOrId(PBSyncCollectionName)
	if err == nil && collection != nil {
		// Collection already exists, nothing to do
		return nil
	}

	// Collection doesn't exist, create it
	log.Printf("Creating _pbSync system collection...")

	// Create new base collection
	collection = core.NewBaseCollection(PBSyncCollectionName)
	collection.System = true

	// Add fields to the collection
	// instanceID field (text, required, unique)
	instanceIDField := &core.TextField{
		Name:     "instanceID",
		Required: true,
	}
	collection.Fields = append(collection.Fields, instanceIDField)
	// Add unique index for instanceID
	collection.AddIndex("idx_instanceID_unique", true, "instanceID", "")

	// natsAddress field (text, required) - combines host:port
	natsAddressField := &core.TextField{
		Name:     "natsAddress",
		Required: true,
	}
	collection.Fields = append(collection.Fields, natsAddressField)

	// created field (autodate, set on create)
	createdField := &core.AutodateField{
		Name:     "created",
		OnCreate: true,
		OnUpdate: false,
	}
	collection.Fields = append(collection.Fields, createdField)

	// updated field (autodate, set on create and update)
	updatedField := &core.AutodateField{
		Name:     "updated",
		OnCreate: true,
		OnUpdate: true,
	}
	collection.Fields = append(collection.Fields, updatedField)

	// Save collection
	if err := app.Save(collection); err != nil {
		return fmt.Errorf("failed to save _pbSync collection: %w", err)
	}

	log.Printf("_pbSync collection created successfully")
	return nil
}

// GetSelfRecord gets or creates the self record in _pbSync collection
func GetSelfRecord(app *pocketbase.PocketBase, instanceID, natsAddress string) (*core.Record, error) {
	// Retry finding collection in case database isn't fully ready
	var collection *core.Collection
	var err error
	for i := 0; i < 5; i++ {
		collection, err = app.FindCollectionByNameOrId(PBSyncCollectionName)
		if err == nil && collection != nil {
			break
		}
		if i < 4 {
			time.Sleep(200 * time.Millisecond)
		}
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find _pbSync collection: %w", err)
	}

	// Try to find existing record by instanceID
	record := &core.Record{}
	err = app.RecordQuery(collection.Id).
		AndWhere(dbx.HashExp{"instanceID": instanceID}).
		Limit(1).
		One(record)
	if err == nil && record.Id != "" {
		// Update fields
		record.Set("natsAddress", natsAddress)
		if err := app.Save(record); err != nil {
			return nil, fmt.Errorf("failed to update self record: %w", err)
		}
		return record, nil
	}

	// Create new record
	record = core.NewRecord(collection)
	record.Set("instanceID", instanceID)
	record.Set("natsAddress", natsAddress)

	if err := app.Save(record); err != nil {
		return nil, fmt.Errorf("failed to create self record: %w", err)
	}

	log.Printf("Created self record in _pbSync: %s", instanceID)
	return record, nil
}

