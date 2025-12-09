package sync

import (
	"fmt"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

// InstanceInfo represents information about a PocketBase sync instance
type InstanceInfo struct {
	InstanceID    string
	NATSAddress   string // Combined host:port
}

// GetAllInstances gets all instance records from local _pbSync collection
func GetAllInstances(app *pocketbase.PocketBase) ([]*core.Record, error) {
	collection, err := app.FindCollectionByNameOrId(PBSyncCollectionName)
	if err != nil {
		return nil, fmt.Errorf("failed to find _pbSync collection: %w", err)
	}

	records := []*core.Record{}
	err = app.RecordQuery(collection.Id).All(&records)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch instances: %w", err)
	}

	return records, nil
}

// ConnectToNATS connects to a remote NATS instance with timeout
func ConnectToNATS(instanceInfo *InstanceInfo) (*nats.Conn, error) {
	// Parse natsAddress (format: "host:port")
	parts := strings.Split(instanceInfo.NATSAddress, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid natsAddress format: %s (expected host:port)", instanceInfo.NATSAddress)
	}

	natsURL := fmt.Sprintf("nats://%s", instanceInfo.NATSAddress)

	// Set connection timeout
	opts := []nats.Option{
		nats.Timeout(5 * time.Second), // 5 second connection timeout
		nats.ReconnectWait(2 * time.Second),
		nats.MaxReconnects(3),
	}

	conn, err := nats.Connect(natsURL, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS at %s: %w", natsURL, err)
	}
	return conn, nil
}

