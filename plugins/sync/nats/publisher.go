package nats

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/pocketbase/pocketbase/plugins/sync/events"

	"github.com/nats-io/nats.go"
)

// Publisher handles publishing events to NATS JetStream
type Publisher struct {
	conns       []*nats.Conn
	js          []nats.JetStreamContext
	instanceIDs []string // Track which instance ID each connection belongs to
	stream      string
	subject     string
	mu          sync.RWMutex
	connected   bool
}

// NewPublisher creates a new NATS JetStream publisher from embedded server
func NewPublisher(embeddedServer *EmbeddedServer, streamName string) (*Publisher, error) {
	if embeddedServer == nil || !embeddedServer.IsRunning() {
		return nil, fmt.Errorf("embedded NATS server is not running")
	}

	conn := embeddedServer.Connection()
	js := embeddedServer.JetStream()

	if js == nil {
		return nil, fmt.Errorf("JetStream is not enabled on embedded server")
	}

	// Ensure stream exists
	subject := "pocketbase.sync.>"
	streamInfo, err := js.StreamInfo(streamName)
	if err == nats.ErrStreamNotFound {
		// Create stream if it doesn't exist
		// Use LimitsPolicy to allow multiple consumers (one per instance)
		// Each instance needs to receive all events for replication
		// MaxAge is set to 2 days - since snapshots are daily, messages older than 2 days
		// are automatically expired (they're covered by snapshots anyway)
		_, err = js.AddStream(&nats.StreamConfig{
			Name:      streamName,
			Subjects:  []string{subject},
			Retention: nats.LimitsPolicy,  // Limits policy allows multiple consumers
			MaxAge:    2 * 24 * time.Hour, // Keep messages for 2 days (auto-expire older messages)
			Storage:   nats.FileStorage,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create stream: %w", err)
		}
		log.Printf("Created NATS JetStream stream: %s", streamName)
	} else if err != nil {
		return nil, fmt.Errorf("failed to get stream info: %w", err)
	} else {
		// Stream exists, check if it has the correct retention policy
		if streamInfo.Config.Retention != nats.LimitsPolicy {
			return nil, fmt.Errorf("stream %s exists with retention policy %v, but requires LimitsPolicy. Please delete the stream and restart (it will be recreated automatically)", streamName, streamInfo.Config.Retention)
		}
		// Update MaxAge if it's set to 0 (indefinite) to enable automatic cleanup
		// This allows existing streams to benefit from the new retention policy
		if streamInfo.Config.MaxAge == 0 {
			log.Printf("Updating stream %s MaxAge from 0 (indefinite) to 2 days for automatic cleanup", streamName)
			streamInfo.Config.MaxAge = 2 * 24 * time.Hour
			if _, err := js.UpdateStream(&streamInfo.Config); err != nil {
				log.Printf("Warning: Failed to update stream MaxAge: %v (this is non-critical)", err)
			} else {
				log.Printf("Successfully updated stream %s MaxAge to 2 days", streamName)
			}
		}
	}

	pub := &Publisher{
		conns:       []*nats.Conn{conn},
		js:          []nats.JetStreamContext{js},
		instanceIDs: []string{""}, // Local connection - empty string indicates local
		stream:      streamName,
		subject:     subject,
		connected:   true,
	}

	return pub, nil
}

// AddConnection adds a connection to another NATS instance for cluster mode
func (p *Publisher) AddConnection(conn *nats.Conn) error {
	return p.AddConnectionWithInstanceID(conn, "")
}

// AddConnectionWithInstanceID adds a connection to another NATS instance with instance ID tracking
func (p *Publisher) AddConnectionWithInstanceID(conn *nats.Conn, remoteInstanceID string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	js, err := conn.JetStream()
	if err != nil {
		return fmt.Errorf("failed to get JetStream context: %w", err)
	}

	// Ensure the stream exists on the remote NATS server
	subject := "pocketbase.sync.>"
	streamInfo, err := js.StreamInfo(p.stream)
	if err == nats.ErrStreamNotFound {
		// Create stream on remote if it doesn't exist
		// MaxAge is set to 2 days - since snapshots are daily, messages older than 2 days
		// are automatically expired (they're covered by snapshots anyway)
		_, err = js.AddStream(&nats.StreamConfig{
			Name:      p.stream,
			Subjects:  []string{subject},
			Retention: nats.LimitsPolicy,
			MaxAge:    2 * 24 * time.Hour, // Keep messages for 2 days (auto-expire older messages)
			Storage:   nats.FileStorage,
		})
		if err != nil {
			return fmt.Errorf("failed to create stream on remote: %w", err)
		}
		log.Printf("Created stream on remote NATS server: %s", p.stream)
	} else if err != nil {
		return fmt.Errorf("failed to get stream info on remote: %w", err)
	} else {
		// Stream exists, check retention policy
		if streamInfo.Config.Retention != nats.LimitsPolicy {
			log.Printf("Warning: remote stream %s has retention policy %v, expected LimitsPolicy", p.stream, streamInfo.Config.Retention)
		}
	}

	p.conns = append(p.conns, conn)
	p.js = append(p.js, js)
	p.instanceIDs = append(p.instanceIDs, remoteInstanceID)
	log.Printf("Added publisher connection (total connections: %d)", len(p.conns))
	return nil
}

// Publish publishes an event to NATS JetStream (publishes to all connected instances)
func (p *Publisher) Publish(event *events.Event) error {
	p.mu.RLock()
	connected := p.connected
	jsContexts := make([]nats.JetStreamContext, len(p.js))
	copy(jsContexts, p.js)
	p.mu.RUnlock()

	if !connected || len(jsContexts) == 0 {
		return fmt.Errorf("publisher not connected to NATS")
	}

	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Use collection-specific subject for better routing
	subject := fmt.Sprintf("pocketbase.sync.%s", event.Collection)

	// Publish to all connected instances
	var lastErr error
	publishedCount := 0
	for i, js := range jsContexts {
		ack, err := js.Publish(subject, data)
		if err != nil {
			lastErr = err
			log.Printf("Warning: failed to publish to NATS instance %d: %v", i, err)
		} else {
			publishedCount++
			if i == 0 {
				log.Printf("Published event to local NATS (ack: %s): %s/%s", ack, event.Collection, event.RecordID)
			} else {
				log.Printf("Published event to remote NATS %d (ack: %s): %s/%s", i, ack, event.Collection, event.RecordID)
			}
		}
	}

	if publishedCount < len(jsContexts) {
		log.Printf("Published to %d/%d NATS instance(s) for event %s/%s", publishedCount, len(jsContexts), event.Collection, event.RecordID)
	}

	if lastErr != nil && len(jsContexts) == 1 {
		return fmt.Errorf("failed to publish event: %w", lastErr)
	}

	return nil
}

// PublishAsync publishes an event asynchronously (non-blocking)
func (p *Publisher) PublishAsync(event *events.Event) {
	go func() {
		if err := p.Publish(event); err != nil {
			log.Printf("Error publishing event asynchronously: %v", err)
		}
	}()
}

// Close closes all NATS connections
func (p *Publisher) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.connected = false

	for _, conn := range p.conns {
		if conn != nil {
			conn.Close()
		}
	}
	p.conns = nil
	p.js = nil
	p.instanceIDs = nil
	return nil
}

// RemoveConnection removes a connection by instance ID
func (p *Publisher) RemoveConnection(instanceID string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Find the index of the connection for this instance ID
	index := -1
	for i, id := range p.instanceIDs {
		if id == instanceID {
			index = i
			break
		}
	}

	if index == -1 {
		// Connection not found - might have already been removed
		return nil
	}

	log.Printf("Removing publisher connection to instance %s (index %d)", instanceID, index)

	// Close connection at this index
	if index < len(p.conns) && p.conns[index] != nil {
		p.conns[index].Close()
		p.conns = append(p.conns[:index], p.conns[index+1:]...)
	} else if index < len(p.conns) {
		p.conns = append(p.conns[:index], p.conns[index+1:]...)
	}

	if index < len(p.js) {
		p.js = append(p.js[:index], p.js[index+1:]...)
	}

	if index < len(p.instanceIDs) {
		p.instanceIDs = append(p.instanceIDs[:index], p.instanceIDs[index+1:]...)
	}

	log.Printf("Removed publisher connection to instance %s (remaining connections: %d)", instanceID, len(p.conns))
	return nil
}

// IsConnected returns whether the publisher is connected
func (p *Publisher) IsConnected() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.connected
}
