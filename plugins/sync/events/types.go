package events

import (
	"encoding/json"
	"time"
)

// EventType represents the type of database operation
type EventType string

const (
	EventTypeCreate              EventType = "create"
	EventTypeUpdate              EventType = "update"
	EventTypeDelete              EventType = "delete"
	EventTypeCollectionCreate    EventType = "collection_create"
	EventTypeCollectionUpdate    EventType = "collection_update"
	EventTypeCollectionDelete    EventType = "collection_delete"
	EventTypeRealtimeSubscribe   EventType = "realtime_subscribe"
	EventTypeRealtimeUnsubscribe EventType = "realtime_unsubscribe"
	EventTypeRealtimeHeartbeat   EventType = "realtime_heartbeat"
)

// EventSource represents where the event originated
type EventSource string

const (
	EventSourceLocal       EventSource = "local"       // Event originated from this instance
	EventSourceRemote      EventSource = "remote"      // Event originated from another instance
	EventSourceRepublished EventSource = "republished" // Event is being republished by this instance
)

// Event represents a database change event
type Event struct {
	Type         EventType      `json:"type"`
	InstanceID   string         `json:"instanceID"`
	Source       EventSource    `json:"source"` // Where this event originated
	Timestamp    time.Time      `json:"timestamp"`
	Collection   string         `json:"collection"`
	RecordID     string         `json:"recordID,omitempty"`
	RecordData   map[string]any `json:"recordData,omitempty"`   // For create/update
	SchemaData   map[string]any `json:"schemaData,omitempty"`   // For collection create/update
	RealtimeData map[string]any `json:"realtimeData,omitempty"` // For realtime events (resourceID, etc.)
}

// Marshal serializes the event to JSON
func (e *Event) Marshal() ([]byte, error) {
	return json.Marshal(e)
}

// Unmarshal deserializes the event from JSON
func (e *Event) Unmarshal(data []byte) error {
	return json.Unmarshal(data, e)
}

// NewCreateEvent creates a new record creation event
func NewCreateEvent(instanceID string, collection string, recordID string, recordData map[string]any) *Event {
	return &Event{
		Type:       EventTypeCreate,
		InstanceID: instanceID,
		Source:     EventSourceLocal,
		Timestamp:  time.Now(),
		Collection: collection,
		RecordID:   recordID,
		RecordData: recordData,
	}
}

// NewUpdateEvent creates a new record update event
func NewUpdateEvent(instanceID string, collection string, recordID string, recordData map[string]any) *Event {
	return &Event{
		Type:       EventTypeUpdate,
		InstanceID: instanceID,
		Source:     EventSourceLocal,
		Timestamp:  time.Now(),
		Collection: collection,
		RecordID:   recordID,
		RecordData: recordData,
	}
}

// NewDeleteEvent creates a new record deletion event
func NewDeleteEvent(instanceID string, collection string, recordID string) *Event {
	return &Event{
		Type:       EventTypeDelete,
		InstanceID: instanceID,
		Source:     EventSourceLocal,
		Timestamp:  time.Now(),
		Collection: collection,
		RecordID:   recordID,
	}
}

// NewCollectionCreateEvent creates a new collection creation event
func NewCollectionCreateEvent(instanceID string, collectionName string, schemaData map[string]any) *Event {
	return &Event{
		Type:       EventTypeCollectionCreate,
		InstanceID: instanceID,
		Source:     EventSourceLocal,
		Timestamp:  time.Now(),
		Collection: collectionName,
		SchemaData: schemaData,
	}
}

// NewCollectionUpdateEvent creates a new collection update event
func NewCollectionUpdateEvent(instanceID string, collectionName string, schemaData map[string]any) *Event {
	return &Event{
		Type:       EventTypeCollectionUpdate,
		InstanceID: instanceID,
		Source:     EventSourceLocal,
		Timestamp:  time.Now(),
		Collection: collectionName,
		SchemaData: schemaData,
	}
}

// NewCollectionDeleteEvent creates a new collection deletion event
func NewCollectionDeleteEvent(instanceID string, collectionName string) *Event {
	return &Event{
		Type:       EventTypeCollectionDelete,
		InstanceID: instanceID,
		Source:     EventSourceLocal,
		Timestamp:  time.Now(),
		Collection: collectionName,
	}
}

// NewRealtimeSubscribeEvent creates a new realtime subscription event
func NewRealtimeSubscribeEvent(instanceID string, resourceID string, activeConnections int, leaseExpires time.Time) *Event {
	return &Event{
		Type:       EventTypeRealtimeSubscribe,
		InstanceID: instanceID,
		Timestamp:  time.Now(),
		RealtimeData: map[string]any{
			"resourceID":        resourceID,
			"activeConnections": activeConnections,
			"leaseExpires":      leaseExpires,
		},
	}
}

// NewRealtimeUnsubscribeEvent creates a new realtime unsubscription event
func NewRealtimeUnsubscribeEvent(instanceID string, resourceID string) *Event {
	return &Event{
		Type:       EventTypeRealtimeUnsubscribe,
		InstanceID: instanceID,
		Timestamp:  time.Now(),
		RealtimeData: map[string]any{
			"resourceID": resourceID,
		},
	}
}

// NewRealtimeHeartbeatEvent creates a new realtime heartbeat event
func NewRealtimeHeartbeatEvent(instanceID string, resourceID string, activeConnections int, leaseExpires time.Time) *Event {
	return &Event{
		Type:       EventTypeRealtimeHeartbeat,
		InstanceID: instanceID,
		Timestamp:  time.Now(),
		RealtimeData: map[string]any{
			"resourceID":        resourceID,
			"activeConnections": activeConnections,
			"leaseExpires":      leaseExpires,
		},
	}
}

// NewRepublishedEvent creates a republished version of an existing event
func NewRepublishedEvent(originalEvent *Event) *Event {
	return &Event{
		Type:         originalEvent.Type,
		InstanceID:   originalEvent.InstanceID, // Keep original publisher
		Source:       EventSourceRepublished,
		Timestamp:    originalEvent.Timestamp,
		Collection:   originalEvent.Collection,
		RecordID:     originalEvent.RecordID,
		RecordData:   originalEvent.RecordData,
		SchemaData:   originalEvent.SchemaData,
		RealtimeData: originalEvent.RealtimeData,
	}
}

// MarkAsRemote marks an event as coming from a remote instance
func (e *Event) MarkAsRemote() *Event {
	e.Source = EventSourceRemote
	return e
}
