package sync

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

const (
	// LeaseDuration is how long a lease is valid before requiring a heartbeat
	LeaseDuration = 30 * time.Second
	// HeartbeatInterval is how often to send heartbeats
	HeartbeatInterval = 10 * time.Second
)

// RealtimeManager manages realtime subscription routing and leases
type RealtimeManager struct {
	app        *pocketbase.PocketBase
	instanceID string
	mu         sync.RWMutex
	// Track active subscriptions locally (resourceID -> subscription record)
	activeSubscriptions map[string]*core.Record
	// Track heartbeat timers (resourceID -> timer)
	heartbeatTimers map[string]*time.Timer
	stopChan        chan struct{}
	// Publisher functions (set by wrapper to avoid import cycle)
	publishRealtimeSubscribe   func(resourceID string, activeConnections int, leaseExpires time.Time) error
	publishRealtimeUnsubscribe func(resourceID string) error
	publishRealtimeHeartbeat   func(resourceID string, activeConnections int, leaseExpires time.Time) error
}

// NewRealtimeManager creates a new realtime manager
func NewRealtimeManager(app *pocketbase.PocketBase, instanceID string) (*RealtimeManager, error) {
	rm := &RealtimeManager{
		app:                 app,
		instanceID:          instanceID,
		activeSubscriptions: make(map[string]*core.Record),
		heartbeatTimers:     make(map[string]*time.Timer),
		stopChan:            make(chan struct{}),
	}

	// Load existing subscriptions for this instance
	if err := rm.loadExistingSubscriptions(); err != nil {
		log.Printf("Warning: Failed to load existing subscriptions: %v", err)
	}

	// Start cleanup goroutine to remove expired leases
	go rm.cleanupExpiredLeases()

	return rm, nil
}

// loadExistingSubscriptions loads all active subscriptions for this instance
func (rm *RealtimeManager) loadExistingSubscriptions() error {
	collection, err := rm.app.FindCollectionByNameOrId(PBRealtimeSubscriptionsCollectionName)
	if err != nil {
		return fmt.Errorf("failed to find _pbRealtimeSubscriptions collection: %w", err)
	}

	records := []*core.Record{}
	err = rm.app.RecordQuery(collection.Id).
		AndWhere(dbx.HashExp{"instanceID": rm.instanceID}).
		All(&records)
	if err != nil {
		return fmt.Errorf("failed to query subscriptions: %w", err)
	}

	rm.mu.Lock()
	defer rm.mu.Unlock()

	for _, record := range records {
		resourceID := record.GetString("resourceID")
		leaseExpires := record.Get("leaseExpires").(time.Time)

		// Only load subscriptions that haven't expired
		if time.Now().Before(leaseExpires) {
			rm.activeSubscriptions[resourceID] = record
			// Restart heartbeat for this subscription
			rm.startHeartbeat(resourceID, record.GetInt("activeConnections"))
		} else {
			// Clean up expired subscription
			rm.app.Delete(record)
		}
	}

	return nil
}

// HandleSubscribe is called when a realtime subscription is opened
func (rm *RealtimeManager) HandleSubscribe(resourceID string) error {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	collection, err := rm.app.FindCollectionByNameOrId(PBRealtimeSubscriptionsCollectionName)
	if err != nil {
		return fmt.Errorf("failed to find _pbRealtimeSubscriptions collection: %w", err)
	}

	// Check if subscription already exists for this resource
	var subscription *core.Record
	err = rm.app.RecordQuery(collection.Id).
		AndWhere(dbx.HashExp{"resourceID": resourceID}).
		Limit(1).
		One(subscription)

	leaseExpires := time.Now().Add(LeaseDuration)
	activeConnections := 1

	if err == nil && subscription != nil {
		// Subscription exists - check if it's for this instance
		existingInstanceID := subscription.GetString("instanceID")
		if existingInstanceID == rm.instanceID {
			// Update existing subscription
			activeConnections = subscription.GetInt("activeConnections") + 1
			subscription.Set("activeConnections", activeConnections)
			subscription.Set("leaseExpires", leaseExpires)
			subscription.Set("lastActivity", time.Now())
			if err := rm.app.Save(subscription); err != nil {
				return fmt.Errorf("failed to update subscription: %w", err)
			}
		} else {
			// Another instance has this subscription - we should route to them
			// But for now, we'll take ownership if the lease has expired
			leaseExpiresTime := subscription.Get("leaseExpires").(time.Time)
			if time.Now().After(leaseExpiresTime) {
				// Lease expired, take ownership
				subscription.Set("instanceID", rm.instanceID)
				subscription.Set("activeConnections", 1)
				subscription.Set("leaseExpires", leaseExpires)
				subscription.Set("lastActivity", time.Now())
				if err := rm.app.Save(subscription); err != nil {
					return fmt.Errorf("failed to update subscription: %w", err)
				}
			} else {
				// Another instance has active lease - return routing info
				return fmt.Errorf("resource %s is handled by instance %s", resourceID, existingInstanceID)
			}
		}
	} else {
		// Create new subscription
		subscription = core.NewRecord(collection)
		subscription.Set("resourceID", resourceID)
		subscription.Set("instanceID", rm.instanceID)
		subscription.Set("activeConnections", activeConnections)
		subscription.Set("leaseExpires", leaseExpires)
		subscription.Set("lastActivity", time.Now())
		if err := rm.app.Save(subscription); err != nil {
			return fmt.Errorf("failed to create subscription: %w", err)
		}
	}

	// Store locally and start heartbeat
	rm.activeSubscriptions[resourceID] = subscription
	rm.startHeartbeat(resourceID, activeConnections)

	// Publish subscribe event to NATS
	// Create event using events package (will be imported where this is called)
	// For now, we'll use a helper function that will be implemented in the wrapper
	// This avoids the import cycle
	if err := rm.publishRealtimeSubscribe(resourceID, activeConnections, leaseExpires); err != nil {
		log.Printf("Warning: Failed to publish realtime subscribe event: %v", err)
	}

	log.Printf("Realtime subscription opened: %s (connections: %d)", resourceID, activeConnections)
	return nil
}

// HandleUnsubscribe is called when a realtime subscription is closed
func (rm *RealtimeManager) HandleUnsubscribe(resourceID string) error {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	collection, err := rm.app.FindCollectionByNameOrId(PBRealtimeSubscriptionsCollectionName)
	if err != nil {
		return fmt.Errorf("failed to find _pbRealtimeSubscriptions collection: %w", err)
	}

	// Find subscription
	var subscription *core.Record
	err = rm.app.RecordQuery(collection.Id).
		AndWhere(dbx.HashExp{"resourceID": resourceID}).
		Limit(1).
		One(subscription)

	if err != nil {
		// Subscription not found, nothing to do
		return nil
	}

	// Check if this is our subscription
	if subscription.GetString("instanceID") != rm.instanceID {
		// Not our subscription, nothing to do
		return nil
	}

	activeConnections := subscription.GetInt("activeConnections") - 1
	if activeConnections <= 0 {
		// No more connections, delete subscription
		if err := rm.app.Delete(subscription); err != nil {
			return fmt.Errorf("failed to delete subscription: %w", err)
		}
		delete(rm.activeSubscriptions, resourceID)
		// Stop heartbeat
		if timer, ok := rm.heartbeatTimers[resourceID]; ok {
			timer.Stop()
			delete(rm.heartbeatTimers, resourceID)
		}
		// Publish unsubscribe event
		if err := rm.publishRealtimeUnsubscribe(resourceID); err != nil {
			log.Printf("Warning: Failed to publish realtime unsubscribe event: %v", err)
		}
		log.Printf("Realtime subscription closed: %s", resourceID)
	} else {
		// Update connection count
		subscription.Set("activeConnections", activeConnections)
		subscription.Set("lastActivity", time.Now())
		if err := rm.app.Save(subscription); err != nil {
			return fmt.Errorf("failed to update subscription: %w", err)
		}
		rm.activeSubscriptions[resourceID] = subscription
	}

	return nil
}

// startHeartbeat starts a heartbeat timer for a subscription
func (rm *RealtimeManager) startHeartbeat(resourceID string, activeConnections int) {
	// Stop existing timer if any
	if timer, ok := rm.heartbeatTimers[resourceID]; ok {
		timer.Stop()
	}

	// Start new timer
	timer := time.AfterFunc(HeartbeatInterval, func() {
		rm.sendHeartbeat(resourceID, activeConnections)
	})
	rm.heartbeatTimers[resourceID] = timer
}

// sendHeartbeat sends a heartbeat for a subscription
func (rm *RealtimeManager) sendHeartbeat(resourceID string, activeConnections int) {
	rm.mu.RLock()
	subscription, exists := rm.activeSubscriptions[resourceID]
	rm.mu.RUnlock()

	if !exists {
		return
	}

	// Check if subscription still belongs to us
	if subscription.GetString("instanceID") != rm.instanceID {
		return
	}

	leaseExpires := time.Now().Add(LeaseDuration)

	// Update lease
	subscription.Set("leaseExpires", leaseExpires)
	subscription.Set("lastActivity", time.Now())
	if err := rm.app.Save(subscription); err != nil {
		log.Printf("Warning: Failed to update subscription lease: %v", err)
		return
	}

	// Publish heartbeat event
	if err := rm.publishRealtimeHeartbeat(resourceID, activeConnections, leaseExpires); err != nil {
		log.Printf("Warning: Failed to publish realtime heartbeat: %v", err)
	}

	// Schedule next heartbeat
	rm.mu.Lock()
	if _, stillExists := rm.activeSubscriptions[resourceID]; stillExists {
		rm.startHeartbeat(resourceID, activeConnections)
	}
	rm.mu.Unlock()
}

// cleanupExpiredLeases periodically removes expired subscriptions
func (rm *RealtimeManager) cleanupExpiredLeases() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rm.cleanupExpired()
		case <-rm.stopChan:
			return
		}
	}
}

// cleanupExpired removes expired subscriptions
func (rm *RealtimeManager) cleanupExpired() {
	collection, err := rm.app.FindCollectionByNameOrId(PBRealtimeSubscriptionsCollectionName)
	if err != nil {
		return
	}

	records := []*core.Record{}
	err = rm.app.RecordQuery(collection.Id).
		All(&records)
	if err != nil {
		return
	}

	now := time.Now()
	for _, record := range records {
		leaseExpires := record.Get("leaseExpires").(time.Time)
		if now.After(leaseExpires) {
			// Lease expired, delete subscription
			if err := rm.app.Delete(record); err != nil {
				log.Printf("Warning: Failed to delete expired subscription: %v", err)
			} else {
				log.Printf("Cleaned up expired subscription: %s", record.GetString("resourceID"))
			}
		}
	}
}

// ProcessRealtimeEvent processes a realtime event from NATS
// eventType should be "realtime_subscribe", "realtime_unsubscribe", or "realtime_heartbeat"
func (rm *RealtimeManager) ProcessRealtimeEvent(eventType, instanceID string, realtimeData map[string]any) error {
	if instanceID == rm.instanceID {
		// Skip our own events
		return nil
	}

	resourceID, ok := realtimeData["resourceID"].(string)
	if !ok {
		return fmt.Errorf("invalid realtime event: missing resourceID")
	}

	collection, err := rm.app.FindCollectionByNameOrId(PBRealtimeSubscriptionsCollectionName)
	if err != nil {
		return fmt.Errorf("failed to find _pbRealtimeSubscriptions collection: %w", err)
	}

	switch eventType {
	case "realtime_subscribe":
		activeConnections, _ := realtimeData["activeConnections"].(int)
		var leaseExpires time.Time
		if leaseExpiresVal, ok := realtimeData["leaseExpires"]; ok {
			if t, ok := leaseExpiresVal.(time.Time); ok {
				leaseExpires = t
			} else if str, ok := leaseExpiresVal.(string); ok {
				leaseExpires, err = time.Parse(time.RFC3339, str)
				if err != nil {
					leaseExpires = time.Now().Add(LeaseDuration)
				}
			} else {
				leaseExpires = time.Now().Add(LeaseDuration)
			}
		} else {
			leaseExpires = time.Now().Add(LeaseDuration)
		}

		// Find or create subscription
		var subscription *core.Record
		err = rm.app.RecordQuery(collection.Id).
			AndWhere(dbx.HashExp{"resourceID": resourceID}).
			Limit(1).
			One(subscription)

		if err != nil {
			// Create new
			subscription = core.NewRecord(collection)
			subscription.Set("resourceID", resourceID)
		}

		subscription.Set("instanceID", instanceID)
		subscription.Set("activeConnections", activeConnections)
		subscription.Set("leaseExpires", leaseExpires)
		subscription.Set("lastActivity", time.Now())

		if err := rm.app.Save(subscription); err != nil {
			return fmt.Errorf("failed to save subscription: %w", err)
		}

	case "realtime_unsubscribe":
		// Find and delete subscription
		var subscription *core.Record
		err = rm.app.RecordQuery(collection.Id).
			AndWhere(dbx.HashExp{"resourceID": resourceID}).
			Limit(1).
			One(subscription)

		if err == nil {
			if err := rm.app.Delete(subscription); err != nil {
				return fmt.Errorf("failed to delete subscription: %w", err)
			}
		}

	case "realtime_heartbeat":
		activeConnections, _ := realtimeData["activeConnections"].(int)
		var leaseExpires time.Time
		if leaseExpiresVal, ok := realtimeData["leaseExpires"]; ok {
			if t, ok := leaseExpiresVal.(time.Time); ok {
				leaseExpires = t
			} else if str, ok := leaseExpiresVal.(string); ok {
				leaseExpires, err = time.Parse(time.RFC3339, str)
				if err != nil {
					leaseExpires = time.Now().Add(LeaseDuration)
				}
			} else {
				leaseExpires = time.Now().Add(LeaseDuration)
			}
		} else {
			leaseExpires = time.Now().Add(LeaseDuration)
		}

		// Update subscription lease
		var subscription *core.Record
		err = rm.app.RecordQuery(collection.Id).
			AndWhere(dbx.HashExp{"resourceID": resourceID}).
			Limit(1).
			One(subscription)

		if err == nil {
			subscription.Set("instanceID", instanceID)
			subscription.Set("activeConnections", activeConnections)
			subscription.Set("leaseExpires", leaseExpires)
			subscription.Set("lastActivity", time.Now())
			if err := rm.app.Save(subscription); err != nil {
				return fmt.Errorf("failed to update subscription: %w", err)
			}
		}
	}

	return nil
}

// Close stops the realtime manager
func (rm *RealtimeManager) Close() {
	close(rm.stopChan)

	rm.mu.Lock()
	defer rm.mu.Unlock()

	// Stop all heartbeats
	for _, timer := range rm.heartbeatTimers {
		timer.Stop()
	}
	rm.heartbeatTimers = make(map[string]*time.Timer)
}

// SetPublishers sets the publisher functions (called by wrapper to avoid import cycle)
func (rm *RealtimeManager) SetPublishers(
	publishSubscribe func(resourceID string, activeConnections int, leaseExpires time.Time) error,
	publishUnsubscribe func(resourceID string) error,
	publishHeartbeat func(resourceID string, activeConnections int, leaseExpires time.Time) error,
) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	rm.publishRealtimeSubscribe = publishSubscribe
	rm.publishRealtimeUnsubscribe = publishUnsubscribe
	rm.publishRealtimeHeartbeat = publishHeartbeat
}
