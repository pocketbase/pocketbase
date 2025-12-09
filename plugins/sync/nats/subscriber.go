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

// Subscriber handles subscribing to events from NATS JetStream
type Subscriber struct {
	conns        []*nats.Conn
	js           []nats.JetStreamContext
	instanceIDs  []string // Track which instance ID each connection belongs to
	stream       string
	consumerName string
	subs         []*nats.Subscription
	instanceID   string
	eventChan    chan *events.Event
	mu           sync.RWMutex
	connected    bool
	stopChan     chan struct{}
}

// NewSubscriber creates a new NATS JetStream subscriber from embedded server
func NewSubscriber(embeddedServer *EmbeddedServer, streamName string, instanceID string) (*Subscriber, error) {
	if embeddedServer == nil || !embeddedServer.IsRunning() {
		return nil, fmt.Errorf("embedded NATS server is not running")
	}

	conn := embeddedServer.Connection()
	js := embeddedServer.JetStream()

	if js == nil {
		return nil, fmt.Errorf("JetStream is not enabled on embedded server")
	}

	// Ensure stream exists (should already exist from publisher, but check anyway)
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

	// Create or get durable consumer
	consumerName := fmt.Sprintf("pocketbase-consumer-%s", instanceID)

	// Check if consumer already exists
	_, err = js.ConsumerInfo(streamName, consumerName)
	if err == nats.ErrConsumerNotFound {
		// Consumer doesn't exist, create it
		// With LimitsPolicy, we can have multiple consumers, each receiving all messages
		_, err = js.AddConsumer(streamName, &nats.ConsumerConfig{
			Durable:       consumerName,
			DeliverPolicy: nats.DeliverAllPolicy, // Deliver all messages
			AckPolicy:     nats.AckExplicitPolicy,
			AckWait:       30 * 1000000000, // 30 seconds
			MaxDeliver:    3,               // Retry up to 3 times
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create consumer: %w", err)
		}
		log.Printf("Created NATS consumer: %s", consumerName)
	} else if err != nil {
		return nil, fmt.Errorf("failed to get consumer info: %w", err)
	} else {
		log.Printf("Using existing NATS consumer: %s", consumerName)
	}

	sub := &Subscriber{
		conns:        []*nats.Conn{conn},
		js:           []nats.JetStreamContext{js},
		instanceIDs:  []string{""}, // Local connection - empty string indicates local
		stream:       streamName,
		consumerName: consumerName,
		subs:         []*nats.Subscription{},
		instanceID:   instanceID,
		eventChan:    make(chan *events.Event, 100), // Buffered channel
		connected:    true,
		stopChan:     make(chan struct{}),
	}

	return sub, nil
}

// AddConnection adds a connection to another NATS instance for cluster mode
func (s *Subscriber) AddConnection(conn *nats.Conn) error {
	return s.AddConnectionWithInstanceID(conn, "")
}

// AddConnectionWithInstanceID adds a connection to another NATS instance with instance ID tracking
// startSequence is optional - if > 0, the consumer will start from that sequence (for snapshot-based sync)
func (s *Subscriber) AddConnectionWithInstanceID(conn *nats.Conn, remoteInstanceID string, startSequence ...uint64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	js, err := conn.JetStream()
	if err != nil {
		return fmt.Errorf("failed to get JetStream context: %w", err)
	}

	s.conns = append(s.conns, conn)
	s.js = append(s.js, js)
	s.instanceIDs = append(s.instanceIDs, remoteInstanceID)

	// If subscriber is already started, subscribe to this new connection
	// We need to create a consumer on the remote stream first
	if s.connected {
		subject := "pocketbase.sync.>"

		// Ensure the remote stream exists (it should, but check anyway)
		_, err = js.StreamInfo(s.stream)
		if err == nats.ErrStreamNotFound {
			// Create stream on remote if it doesn't exist
			_, err = js.AddStream(&nats.StreamConfig{
				Name:      s.stream,
				Subjects:  []string{subject},
				Retention: nats.LimitsPolicy,
				MaxAge:    0,
				Storage:   nats.FileStorage,
			})
			if err != nil {
				return fmt.Errorf("failed to create stream on remote: %w", err)
			}
		} else if err != nil {
			return fmt.Errorf("failed to get stream info on remote: %w", err)
		}

		// Create consumer on remote stream with our instance's consumer name
		// This allows us to receive messages from the remote instance's stream

		// Get stream info for logging
		streamInfo, streamErr := js.StreamInfo(s.stream)

		// Determine start sequence if provided (for snapshot-based sync)
		var startSeq uint64 = 0
		if len(startSequence) > 0 && startSequence[0] > 0 {
			startSeq = startSequence[0]
		}

		// Check if consumer exists
		consumerInfo, err := js.ConsumerInfo(s.stream, s.consumerName)
		if err == nats.ErrConsumerNotFound {
			// Create consumer - use DeliverByStartSequencePolicy if startSeq is provided
			consumerConfig := &nats.ConsumerConfig{
				Durable:       s.consumerName,
				AckPolicy:     nats.AckExplicitPolicy,
				AckWait:       30 * 1000000000,
				MaxDeliver:    1,    // Deliver once - no redelivery attempts
				MaxAckPending: 1000, // Allow up to 1000 unacknowledged messages
				RateLimit:     0,    // No rate limit
			}

			if startSeq > 0 {
				// Start from specific sequence (after snapshot)
				consumerConfig.DeliverPolicy = nats.DeliverByStartSequencePolicy
				consumerConfig.OptStartSeq = startSeq
				log.Printf("Creating consumer %s on remote stream %s starting from sequence %d (after snapshot)", s.consumerName, s.stream, startSeq)
			} else {
				// Get all messages from the stream (no snapshot available)
				consumerConfig.DeliverPolicy = nats.DeliverAllPolicy
				log.Printf("Creating consumer %s on remote stream %s (deliver policy: DeliverAll)", s.consumerName, s.stream)
			}

			consumerInfo, err = js.AddConsumer(s.stream, consumerConfig)
			if err != nil {
				return fmt.Errorf("failed to create consumer on remote stream: %w", err)
			}
		} else if err != nil {
			return fmt.Errorf("failed to get consumer info on remote: %w", err)
		} else {
			// Consumer exists - check if we need to update it for snapshot-based sync
			if startSeq > 0 {
				// If we have a start sequence but the consumer doesn't match, delete and recreate it
				if consumerInfo.Config.DeliverPolicy != nats.DeliverByStartSequencePolicy ||
					consumerInfo.Config.OptStartSeq != startSeq {
					log.Printf("Consumer exists but doesn't match snapshot sequence %d - deleting and recreating", startSeq)
					if err := js.DeleteConsumer(s.stream, s.consumerName); err != nil {
						log.Printf("Warning: failed to delete existing consumer: %v", err)
					} else {
						consumerConfig := &nats.ConsumerConfig{
							Durable:       s.consumerName,
							DeliverPolicy: nats.DeliverByStartSequencePolicy,
							OptStartSeq:   startSeq,
							AckPolicy:     nats.AckExplicitPolicy,
							AckWait:       30 * 1000000000,
							MaxDeliver:    1,
							MaxAckPending: 1000,
							RateLimit:     0,
						}
						consumerInfo, err = js.AddConsumer(s.stream, consumerConfig)
						if err != nil {
							return fmt.Errorf("failed to recreate consumer with start sequence: %w", err)
						}
						log.Printf("Recreated consumer %s on remote stream %s starting from sequence %d", s.consumerName, s.stream, startSeq)
					}
				} else {
					log.Printf("Using existing consumer %s on remote stream %s (starting from sequence %d, pending: %d, delivered: %d)",
						s.consumerName, s.stream, startSeq, consumerInfo.NumPending, consumerInfo.Delivered.Consumer)
				}
			} else {
				log.Printf("Using existing consumer %s on remote stream %s (pending: %d, delivered: %d)",
					s.consumerName, s.stream, consumerInfo.NumPending, consumerInfo.Delivered.Consumer)
				// If consumer exists but has pending messages that aren't being delivered,
				// delete and recreate it to reset the delivery state
				if consumerInfo.NumPending > 0 && consumerInfo.Delivered.Consumer == 0 {
					log.Printf("Consumer has %d pending messages but hasn't delivered any - deleting and recreating consumer", consumerInfo.NumPending)
					if err := js.DeleteConsumer(s.stream, s.consumerName); err != nil {
						log.Printf("Warning: failed to delete existing consumer: %v", err)
					} else {
						// Recreate the consumer with DeliverAllPolicy
						consumerInfo, err = js.AddConsumer(s.stream, &nats.ConsumerConfig{
							Durable:       s.consumerName,
							DeliverPolicy: nats.DeliverAllPolicy,
							AckPolicy:     nats.AckExplicitPolicy,
							AckWait:       30 * 1000000000,
							MaxDeliver:    1,    // Deliver once
							MaxAckPending: 1000, // Allow up to 1000 unacknowledged messages
							RateLimit:     0,    // No rate limit
						})
						if err != nil {
							return fmt.Errorf("failed to recreate consumer: %w", err)
						}
						log.Printf("Recreated consumer %s on remote stream %s", s.consumerName, s.stream)
					}
				}
			}
		}

		// Subscribe to the remote stream using PullSubscribe with explicit Bind
		// This ensures the subscription is properly bound to the consumer
		// Format: PullSubscribe(subject, consumerName, Bind(streamName, consumerName))
		sub, err := js.PullSubscribe(subject, s.consumerName, nats.Bind(s.stream, s.consumerName))
		if err != nil {
			// If Bind fails, try without it as fallback
			log.Printf("Warning: PullSubscribe with Bind failed, trying without Bind: %v", err)
			sub, err = js.PullSubscribe(subject, s.consumerName)
			if err != nil {
				return fmt.Errorf("failed to subscribe to new connection: %w", err)
			}
		}
		s.subs = append(s.subs, sub)

		// Track which instance this subscription belongs to (same index as connection)
		// The instanceID was already added to instanceIDs array above

		// Immediately try to fetch messages to "activate" the consumer
		// This helps ensure the consumer starts delivering messages
		go func(sub *nats.Subscription, subIndex int) {
			time.Sleep(100 * time.Millisecond) // Small delay to ensure subscription is ready
			msgs, err := sub.Fetch(1, nats.MaxWait(2*time.Second))
			if err == nil && len(msgs) > 0 {
				log.Printf("Subscription %d activated - fetched %d message(s)", subIndex, len(msgs))
				for _, msg := range msgs {
					s.processMessage(msg)
				}
			} else if err != nil && err != nats.ErrTimeout {
				log.Printf("Warning: activation fetch for subscription %d failed: %v", subIndex, err)
			}
		}(sub, len(s.subs)-1)

		// Small delay to ensure consumer is ready, then verify status
		go func(subIndex int) {
			time.Sleep(200 * time.Millisecond)
			consumerInfo, err := js.ConsumerInfo(s.stream, s.consumerName)
			if err == nil {
				log.Printf("Subscription %d consumer status: pending=%d, delivered=%d, redelivered=%d",
					subIndex, consumerInfo.NumPending, consumerInfo.Delivered.Consumer, consumerInfo.NumRedelivered)
			}
		}(len(s.subs) - 1)

		// Check stream info to see how many messages are in the stream (reuse existing streamInfo if available)
		if streamInfo == nil || streamErr != nil {
			streamInfo, streamErr = js.StreamInfo(s.stream)
		}
		if streamErr == nil && streamInfo != nil {
			log.Printf("Subscribed to new NATS connection (total subscriptions: %d, consumer: %s, stream: %s, stream messages: %d, first: %d, last: %d)",
				len(s.subs), s.consumerName, s.stream, streamInfo.State.Msgs, streamInfo.State.FirstSeq, streamInfo.State.LastSeq)
		} else {
			log.Printf("Subscribed to new NATS connection (total subscriptions: %d, consumer: %s, stream: %s)",
				len(s.subs), s.consumerName, s.stream)
		}
	}

	return nil
}

// Start starts the subscriber and begins receiving events
func (s *Subscriber) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.connected {
		return fmt.Errorf("subscriber not connected to NATS")
	}

	subject := "pocketbase.sync.>"

	// Subscribe to all connected instances using explicit Bind
	// This ensures subscriptions are properly bound to consumers
	for i, js := range s.js {
		sub, err := js.PullSubscribe(subject, s.consumerName, nats.Bind(s.stream, s.consumerName))
		if err != nil {
			// If Bind fails, try without it as fallback
			log.Printf("Warning: PullSubscribe with Bind failed for subscription %d, trying without Bind: %v", i, err)
			sub, err = js.PullSubscribe(subject, s.consumerName)
			if err != nil {
				return fmt.Errorf("failed to subscribe: %w", err)
			}
		}
		s.subs = append(s.subs, sub)
	}

	// Start goroutine to fetch messages
	go s.fetchMessages()

	return nil
}

// fetchMessages continuously fetches messages from all subscriptions
func (s *Subscriber) fetchMessages() {
	fetchCount := 0
	errorCounts := make(map[int]int)        // Track consecutive errors per subscription
	lastErrorLog := make(map[int]time.Time) // Track when we last logged an error for each subscription

	for {
		select {
		case <-s.stopChan:
			return
		default:
			s.mu.RLock()
			subs := make([]*nats.Subscription, len(s.subs))
			copy(subs, s.subs)
			jsContexts := make([]nats.JetStreamContext, len(s.js))
			copy(jsContexts, s.js)
			instanceIDs := make([]string, len(s.instanceIDs))
			copy(instanceIDs, s.instanceIDs)
			s.mu.RUnlock()

			if len(subs) == 0 {
				time.Sleep(100 * time.Millisecond)
				continue
			}

			for i, sub := range subs {
				if sub == nil {
					continue
				}

				// Check if subscription is valid before trying to fetch
				if !sub.IsValid() {
					s.removeSubscriptionAtIndex(i, instanceIDs)
					delete(errorCounts, i)
					delete(lastErrorLog, i)
					// Break to restart the loop with updated subscriptions
					break
				}

				// Fetch messages (pull subscriptions only support Fetch, not NextMsg)
				msgs, err := sub.Fetch(10, nats.MaxWait(1*time.Second))
				if err != nil {
					if err == nats.ErrTimeout {
						// No messages available - reset error count on success
						if errorCounts[i] > 0 {
							errorCounts[i] = 0
						}
						continue
					}

					// Check for fatal errors that indicate the subscription is dead
					errStr := err.Error()
					isFatalError := errStr == "nats: invalid subscription" ||
						errStr == "nats: subscription closed" ||
						errStr == "nats: connection closed" ||
						errStr == "nats: connection is closed" ||
						errStr == "nats: no responders available for request"

					if isFatalError {
						// Increment error count
						errorCounts[i]++

						// Remove after 3 consecutive fatal errors
						if errorCounts[i] >= 3 {
							instanceID := ""
							if i < len(instanceIDs) {
								instanceID = instanceIDs[i]
							}
							if instanceID != "" {
								log.Printf("Subscription %d (instance %s) has %d consecutive fatal errors (%v), removing connection", i, instanceID, errorCounts[i], err)
							} else {
								log.Printf("Subscription %d has %d consecutive fatal errors (%v), removing it", i, errorCounts[i], err)
							}
							s.removeSubscriptionAtIndex(i, instanceIDs)
							delete(errorCounts, i)
							delete(lastErrorLog, i)
							// Break to restart the loop with updated subscriptions
							break
						}

						// Log error only every 10 seconds per subscription to reduce spam
						now := time.Now()
						lastLog, exists := lastErrorLog[i]
						if !exists || now.Sub(lastLog) >= 10*time.Second {
							instanceID := ""
							if i < len(instanceIDs) {
								instanceID = instanceIDs[i]
							}
							if instanceID != "" {
								log.Printf("Subscription %d (instance %s) error: %v (consecutive errors: %d)", i, instanceID, err, errorCounts[i])
							} else {
								log.Printf("Subscription %d error: %v (consecutive errors: %d)", i, err, errorCounts[i])
							}
							lastErrorLog[i] = now
						}

						// Add exponential backoff: sleep longer after more errors
						backoff := time.Duration(errorCounts[i]) * 100 * time.Millisecond
						if backoff > 2*time.Second {
							backoff = 2 * time.Second
						}
						time.Sleep(backoff)
						continue
					}

					// For non-fatal errors, log occasionally but don't remove subscription
					now := time.Now()
					lastLog, exists := lastErrorLog[i]
					if !exists || now.Sub(lastLog) >= 30*time.Second {
						log.Printf("Warning: Subscription %d non-fatal error: %v", i, err)
						lastErrorLog[i] = now
					}
					continue
				}

				// Success - reset error count
				if errorCounts[i] > 0 {
					errorCounts[i] = 0
					delete(lastErrorLog, i)
				}

				if len(msgs) > 0 {
					log.Printf("Received %d message(s) from NATS subscription %d", len(msgs), i)
				}
				for _, msg := range msgs {
					s.processMessage(msg)
				}
			}

			// Small sleep to prevent CPU spinning
			time.Sleep(10 * time.Millisecond)

			// Log subscription status every 1000 fetches (roughly every 10 seconds)
			fetchCount++
			if fetchCount%1000 == 0 && len(jsContexts) > 1 {
				// Check consumer status on remote streams (only log if there are pending messages)
				for i := 1; i < len(jsContexts); i++ {
					consumerInfo, err := jsContexts[i].ConsumerInfo(s.stream, s.consumerName)
					if err == nil && consumerInfo.NumPending > 0 {
						log.Printf("Subscription %d has %d pending messages (delivered: %d)", i, consumerInfo.NumPending, consumerInfo.Delivered.Consumer)
					}
				}
			}
		}
	}
}

// removeSubscriptionAtIndex removes a subscription and its associated connection at the given index
func (s *Subscriber) removeSubscriptionAtIndex(index int, instanceIDs []string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if index < len(s.subs) && s.subs[index] != nil {
		if err := s.subs[index].Unsubscribe(); err != nil {
			// Silently ignore unsubscribe errors - subscription may already be closed
		}
		s.subs = append(s.subs[:index], s.subs[index+1:]...)
	} else if index < len(s.subs) {
		s.subs = append(s.subs[:index], s.subs[index+1:]...)
	}

	if index < len(s.conns) && s.conns[index] != nil {
		s.conns[index].Close()
		s.conns = append(s.conns[:index], s.conns[index+1:]...)
	} else if index < len(s.conns) {
		s.conns = append(s.conns[:index], s.conns[index+1:]...)
	}

	if index < len(s.js) {
		s.js = append(s.js[:index], s.js[index+1:]...)
	}

	if index < len(s.instanceIDs) {
		instanceID := s.instanceIDs[index]
		if instanceID != "" {
			log.Printf("Removed connection to instance %s", instanceID)
		}
		s.instanceIDs = append(s.instanceIDs[:index], s.instanceIDs[index+1:]...)
	}
}

// processMessage processes a single NATS message
func (s *Subscriber) processMessage(msg *nats.Msg) {
	var event events.Event
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		log.Printf("Error unmarshaling event: %v", err)
		msg.Nak() // Negative acknowledgment - will retry
		return
	}

	// Skip events from our own instance (silently, to reduce log spam)
	if event.InstanceID == s.instanceID {
		msg.Ack() // Acknowledge our own events
		return
	}

	log.Printf("Received event from instance %s: %s/%s", event.InstanceID, event.Collection, event.RecordID)

	// Mark this event as coming from a remote instance
	event.MarkAsRemote()

	// Send event to channel
	select {
	case s.eventChan <- &event:
		msg.Ack() // Acknowledge after queuing
	default:
		log.Printf("Event channel full, dropping event")
		msg.Nak() // Negative acknowledgment - will retry
	}
}

// Events returns the channel that receives events
func (s *Subscriber) Events() <-chan *events.Event {
	return s.eventChan
}

// Stop stops the subscriber
func (s *Subscriber) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	close(s.stopChan)

	for _, sub := range s.subs {
		if sub != nil {
			if err := sub.Unsubscribe(); err != nil {
				log.Printf("Warning: failed to unsubscribe: %v", err)
			}
		}
	}
	s.subs = nil

	s.connected = false
	return nil
}

// Close closes all NATS connections
func (s *Subscriber) Close() error {
	if err := s.Stop(); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, conn := range s.conns {
		if conn != nil {
			conn.Close()
		}
	}
	s.conns = nil
	s.js = nil
	s.instanceIDs = nil
	return nil
}

// RemoveConnection removes a connection by instance ID
func (s *Subscriber) RemoveConnection(instanceID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Find the index of the connection for this instance ID
	index := -1
	for i, id := range s.instanceIDs {
		if id == instanceID {
			index = i
			break
		}
	}

	if index == -1 {
		// Connection not found - might have already been removed
		return nil
	}

	log.Printf("Removing connection to instance %s (subscription %d)", instanceID, index)

	// Unsubscribe and close connection at this index
	if index < len(s.subs) && s.subs[index] != nil {
		if err := s.subs[index].Unsubscribe(); err != nil {
			log.Printf("Warning: failed to unsubscribe: %v", err)
		}
		s.subs = append(s.subs[:index], s.subs[index+1:]...)
	} else if index < len(s.subs) {
		s.subs = append(s.subs[:index], s.subs[index+1:]...)
	}

	if index < len(s.conns) && s.conns[index] != nil {
		s.conns[index].Close()
		s.conns = append(s.conns[:index], s.conns[index+1:]...)
	} else if index < len(s.conns) {
		s.conns = append(s.conns[:index], s.conns[index+1:]...)
	}

	if index < len(s.js) {
		s.js = append(s.js[:index], s.js[index+1:]...)
	}

	if index < len(s.instanceIDs) {
		s.instanceIDs = append(s.instanceIDs[:index], s.instanceIDs[index+1:]...)
	}

	log.Printf("Removed connection to instance %s (remaining connections: %d)", instanceID, len(s.conns))
	return nil
}

// IsConnected returns whether the subscriber is connected
func (s *Subscriber) IsConnected() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.connected
}

// WaitForCatchUp waits for all consumers to catch up (NumPending == 0) or timeout
// Returns true if all consumers caught up, false if timeout
func (s *Subscriber) WaitForCatchUp(timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		s.mu.RLock()
		jsContexts := make([]nats.JetStreamContext, len(s.js))
		copy(jsContexts, s.js)
		instanceIDs := make([]string, len(s.instanceIDs))
		copy(instanceIDs, s.instanceIDs)
		s.mu.RUnlock()

		allCaughtUp := true
		for i, js := range jsContexts {
			if i == 0 {
				// Skip local connection (index 0)
				continue
			}

			consumerInfo, err := js.ConsumerInfo(s.stream, s.consumerName)
			if err != nil {
				log.Printf("Warning: Failed to get consumer info for subscription %d: %v", i, err)
				continue
			}

			if consumerInfo.NumPending > 0 {
				instanceID := ""
				if i < len(instanceIDs) {
					instanceID = instanceIDs[i]
				}
				if instanceID != "" {
					log.Printf("Subscription %d (instance %s) still has %d pending messages", i, instanceID, consumerInfo.NumPending)
				} else {
					log.Printf("Subscription %d still has %d pending messages", i, consumerInfo.NumPending)
				}
				allCaughtUp = false
				break
			}
		}

		if allCaughtUp {
			log.Printf("All consumers have caught up")
			return true
		}

		// Wait a bit before checking again
		time.Sleep(500 * time.Millisecond)
	}

	log.Printf("Timeout waiting for consumers to catch up")
	return false
}

// GetPendingCountForInstance returns the number of pending messages for a specific instance
// Returns -1 if instance not found or error
func (s *Subscriber) GetPendingCountForInstance(instanceID string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for i, id := range s.instanceIDs {
		if id == instanceID && i < len(s.js) {
			consumerInfo, err := s.js[i].ConsumerInfo(s.stream, s.consumerName)
			if err == nil {
				return int(consumerInfo.NumPending)
			}
		}
	}

	return -1
}
