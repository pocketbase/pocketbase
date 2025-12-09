package sync

import (
	"fmt"
	"log"
	"time"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

const (
	PBRealtimeSubscriptionsCollectionName = "_pbRealtimeSubscriptions"
)

// EnsurePBRealtimeSubscriptionsCollection creates the _pbRealtimeSubscriptions system collection if it doesn't exist
// This collection tracks active WebSocket realtime subscriptions for routing purposes
func EnsurePBRealtimeSubscriptionsCollection(app *pocketbase.PocketBase) error {
	// Retry finding collection in case database isn't fully ready
	var collection *core.Collection
	var err error
	for i := 0; i < 5; i++ {
		collection, err = app.FindCollectionByNameOrId(PBRealtimeSubscriptionsCollectionName)
		if err == nil && collection != nil {
			log.Printf("_pbRealtimeSubscriptions collection already exists")
			return nil
		}
		if i < 4 {
			time.Sleep(200 * time.Millisecond)
		}
	}

	// If we get here, collection doesn't exist, so create it
	log.Printf("Creating _pbRealtimeSubscriptions system collection...")

	// Create new base collection
	collection = core.NewBaseCollection(PBRealtimeSubscriptionsCollectionName)
	collection.System = true

	// resourceID field (text, required, indexed) - e.g., "collection:posts" or "record:posts:123"
	resourceIDField := &core.TextField{
		Name:     "resourceID",
		Required: true,
	}
	collection.Fields = append(collection.Fields, resourceIDField)
	// Add unique index for resourceID (only one active subscription per resource)
	collection.AddIndex("idx_resourceID_unique", true, "resourceID", "")

	// instanceID field (text, required) - which instance has the active subscription
	instanceIDField := &core.TextField{
		Name:     "instanceID",
		Required: true,
	}
	collection.Fields = append(collection.Fields, instanceIDField)
	// Add index for instanceID (for querying all subscriptions on an instance)
	collection.AddIndex("idx_instanceID", false, "instanceID", "")

	// activeConnections field (number) - count of active WebSocket connections
	activeConnectionsField := &core.NumberField{
		Name:     "activeConnections",
		Required: true,
	}
	collection.Fields = append(collection.Fields, activeConnectionsField)

	// lastActivity field (date) - when last connection was made
	lastActivityField := &core.DateField{
		Name:     "lastActivity",
		Required: true,
	}
	collection.Fields = append(collection.Fields, lastActivityField)

	// leaseExpires field (date) - when the lease expires (used for heartbeat/cleanup)
	leaseExpiresField := &core.DateField{
		Name:     "leaseExpires",
		Required: true,
	}
	collection.Fields = append(collection.Fields, leaseExpiresField)

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
		return fmt.Errorf("failed to save _pbRealtimeSubscriptions collection: %w", err)
	}

	log.Printf("_pbRealtimeSubscriptions collection created successfully")
	return nil
}

