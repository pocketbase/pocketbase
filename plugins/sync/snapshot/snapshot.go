package snapshot

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"sort"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

const (
	// SnapshotSubject is the NATS subject for snapshot messages
	SnapshotSubject = "pocketbase.snapshot"
	// SnapshotStreamSuffix is appended to the main stream name for snapshots
	SnapshotStreamSuffix = "-snapshots"
)

// Snapshot represents a complete database snapshot
type Snapshot struct {
	Timestamp     time.Time                           `json:"timestamp"`
	InstanceID    string                              `json:"instanceID"`
	Schema        []map[string]interface{}            `json:"schema"`        // Full PocketBase schema JSON (like pb_schema.json)
	Collections   map[string][]map[string]interface{} `json:"collections"`   // collection name -> records
	EventSequence uint64                              `json:"eventSequence"` // Sequence number of last event before snapshot
}

// Manager handles snapshot creation and retrieval
type Manager struct {
	app            *pocketbase.PocketBase
	js             nats.JetStreamContext
	streamName     string
	snapshotStream string
	instanceID     string
}

// NewManager creates a new snapshot manager
func NewManager(app *pocketbase.PocketBase, js nats.JetStreamContext, streamName, instanceID string) (*Manager, error) {
	snapshotStream := streamName + SnapshotStreamSuffix

	// Ensure snapshot stream exists
	subject := SnapshotSubject
	_, err := js.StreamInfo(snapshotStream)
	if err == nats.ErrStreamNotFound {
		// Create snapshot stream with longer retention (keep snapshots for 30 days by default)
		_, err = js.AddStream(&nats.StreamConfig{
			Name:      snapshotStream,
			Subjects:  []string{subject},
			Retention: nats.LimitsPolicy,
			MaxAge:    30 * 24 * time.Hour, // Keep snapshots for 30 days
			Storage:   nats.FileStorage,
			MaxMsgs:   100, // Keep last 100 snapshots max
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create snapshot stream: %w", err)
		}
		log.Printf("Created NATS JetStream snapshot stream: %s", snapshotStream)
	} else if err != nil {
		return nil, fmt.Errorf("failed to get snapshot stream info: %w", err)
	}

	return &Manager{
		app:            app,
		js:             js,
		streamName:     streamName,
		snapshotStream: snapshotStream,
		instanceID:     instanceID,
	}, nil
}

// CreateSnapshot creates a snapshot of all collections (excluding system collections and _pbSync)
func (m *Manager) CreateSnapshot() (*Snapshot, error) {
	log.Printf("Creating database snapshot...")
	startTime := time.Now()

	// Get the current sequence number of the last event in the main stream
	// This represents the last event that was published before this snapshot
	// New instances will start consuming from this sequence + 1
	var eventSequence uint64 = 0
	streamInfo, err := m.js.StreamInfo(m.streamName)
	if err == nil && streamInfo != nil {
		eventSequence = streamInfo.State.LastSeq
		log.Printf("Current event stream sequence: %d (new instances will start from sequence %d)", eventSequence, eventSequence+1)
	} else {
		log.Printf("Warning: Could not get stream info to determine event sequence: %v", err)
	}

	// Get all collections
	collections, err := m.app.FindAllCollections("base", "auth")
	if err != nil {
		return nil, fmt.Errorf("failed to get collections: %w", err)
	}

	// Export schema (full PocketBase schema JSON format)
	schema, err := m.exportSchema(collections)
	if err != nil {
		return nil, fmt.Errorf("failed to export schema: %w", err)
	}
	log.Printf("Exported schema: %d collections", len(schema))

	snapshot := &Snapshot{
		Timestamp:     time.Now(),
		InstanceID:    m.instanceID,
		Schema:        schema,
		Collections:   make(map[string][]map[string]interface{}),
		EventSequence: eventSequence,
	}

	// Collect all records from each collection
	for _, collection := range collections {
		// Skip system collections (they start with _), but allow _pbSync for instance discovery
		if len(collection.Name) > 0 && collection.Name[0] == '_' && collection.Name != "_pbSync" {
			continue
		}

		// Include _pbSync records in snapshots - they are essential for instance discovery
		// This ensures instances can discover each other through snapshots

		log.Printf("Snapshotting collection: %s", collection.Name)

		// Query all records in this collection
		records := []*core.Record{}
		err := m.app.RecordQuery(collection.Id).All(&records)
		if err != nil {
			log.Printf("Warning: failed to query collection %s: %v", collection.Name, err)
			continue
		}

		// Convert records to maps
		recordMaps := make([]map[string]interface{}, 0, len(records))
		for _, record := range records {
			recordMap := make(map[string]interface{})
			for _, field := range collection.Fields {
				fieldName := field.GetName()
				value := record.Get(fieldName)

				if value != nil {
					recordMap[fieldName] = value
				}
			}
			// Always include ID
			recordMap["id"] = record.Id
			recordMaps = append(recordMaps, recordMap)
		}

		snapshot.Collections[collection.Name] = recordMaps
		log.Printf("Snapshot collection %s: %d records", collection.Name, len(recordMaps))
	}

	// Publish snapshot to NATS JetStream
	snapshotData, err := json.Marshal(snapshot)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal snapshot: %w", err)
	}

	// Compress snapshot data to reduce size
	var compressedData bytes.Buffer
	gzipWriter := gzip.NewWriter(&compressedData)
	if _, err := gzipWriter.Write(snapshotData); err != nil {
		gzipWriter.Close()
		return nil, fmt.Errorf("failed to compress snapshot: %w", err)
	}
	if err := gzipWriter.Close(); err != nil {
		return nil, fmt.Errorf("failed to close gzip writer: %w", err)
	}

	originalSize := len(snapshotData)
	compressedSize := compressedData.Len()
	compressionRatio := float64(compressedSize) / float64(originalSize) * 100
	log.Printf("Snapshot size: %d bytes (compressed: %d bytes, %.1f%%)", originalSize, compressedSize, compressionRatio)

	// Publish with timestamp as header for easy retrieval
	msg := &nats.Msg{
		Subject: SnapshotSubject,
		Data:    compressedData.Bytes(),
		Header:  nats.Header{},
	}
	msg.Header.Set("timestamp", snapshot.Timestamp.Format(time.RFC3339))
	msg.Header.Set("instanceID", m.instanceID)
	msg.Header.Set("compressed", "gzip")                                       // Mark as compressed
	msg.Header.Set("eventSequence", fmt.Sprintf("%d", snapshot.EventSequence)) // Store event sequence for consumer configuration

	ack, err := m.js.PublishMsg(msg)
	if err != nil {
		return nil, fmt.Errorf("failed to publish snapshot: %w", err)
	}

	duration := time.Since(startTime)
	totalRecords := 0
	for _, records := range snapshot.Collections {
		totalRecords += len(records)
	}

	log.Printf("Snapshot created successfully: %d collections, %d total records, sequence %d (took %v)",
		len(snapshot.Collections), totalRecords, ack.Sequence, duration)

	return snapshot, nil
}

// GetLatestSnapshotFromRemote retrieves the most recent snapshot from a remote NATS JetStream
func GetLatestSnapshotFromRemote(js nats.JetStreamContext, streamName string) (*Snapshot, error) {
	snapshotStream := streamName + SnapshotStreamSuffix

	// Get stream info to find the last message
	streamInfo, err := js.StreamInfo(snapshotStream)
	if err != nil {
		return nil, fmt.Errorf("failed to get stream info: %w", err)
	}

	if streamInfo.State.Msgs == 0 {
		return nil, fmt.Errorf("no snapshots available")
	}

	// Use GetMsg to retrieve the last message by sequence number
	// This avoids creating a consumer and is more efficient
	lastSeq := streamInfo.State.LastSeq
	msg, err := js.GetMsg(snapshotStream, lastSeq)
	if err != nil {
		return nil, fmt.Errorf("failed to get last message (seq %d): %w", lastSeq, err)
	}

	if msg == nil {
		return nil, fmt.Errorf("no snapshot message found at sequence %d", lastSeq)
	}

	// Check if message is compressed
	var data []byte
	if msg.Header.Get("compressed") == "gzip" {
		// Decompress the data
		gzipReader, err := gzip.NewReader(bytes.NewReader(msg.Data))
		if err != nil {
			return nil, fmt.Errorf("failed to create gzip reader: %w", err)
		}
		defer gzipReader.Close()

		data, err = io.ReadAll(gzipReader)
		if err != nil {
			return nil, fmt.Errorf("failed to decompress snapshot: %w", err)
		}
	} else {
		data = msg.Data
	}

	// Unmarshal snapshot
	var snapshot Snapshot
	if err := json.Unmarshal(data, &snapshot); err != nil {
		return nil, fmt.Errorf("failed to unmarshal snapshot: %w", err)
	}

	// Extract event sequence from header if available (for backward compatibility)
	if eventSeqStr := msg.Header.Get("eventSequence"); eventSeqStr != "" {
		if seq, err := fmt.Sscanf(eventSeqStr, "%d", &snapshot.EventSequence); err == nil && seq == 1 {
			log.Printf("Retrieved snapshot from remote NATS (sequence: %d, event sequence: %d, %d collections, instance: %s)",
				lastSeq, snapshot.EventSequence, len(snapshot.Collections), snapshot.InstanceID)
		} else {
			log.Printf("Retrieved snapshot from remote NATS (sequence: %d, %d collections, instance: %s)",
				lastSeq, len(snapshot.Collections), snapshot.InstanceID)
		}
	} else {
		log.Printf("Retrieved snapshot from remote NATS (sequence: %d, %d collections, instance: %s)",
			lastSeq, len(snapshot.Collections), snapshot.InstanceID)
	}

	return &snapshot, nil
}

// HasSnapshot checks if a snapshot exists in the local NATS JetStream
func (m *Manager) HasSnapshot() (bool, error) {
	streamInfo, err := m.js.StreamInfo(m.snapshotStream)
	if err != nil {
		if err == nats.ErrStreamNotFound {
			return false, nil
		}
		return false, fmt.Errorf("failed to get stream info: %w", err)
	}
	return streamInfo.State.Msgs > 0, nil
}

// ApplySnapshot applies a snapshot to the local PocketBase instance
// If isInitialDiscovery is true, the snapshot is treated as the source of truth and will
// overwrite existing collections/records. If false, snapshots are only applied to empty instances.
func (m *Manager) ApplySnapshot(snapshot *Snapshot, isInitialDiscovery bool) error {
	log.Printf("Applying snapshot from instance %s (timestamp: %s, initialDiscovery: %v)...",
		snapshot.InstanceID, snapshot.Timestamp.Format(time.RFC3339), isInitialDiscovery)

	// If not in initial discovery, only apply snapshots to empty instances
	// During initial discovery, snapshots are the source of truth and should be applied
	if !isInitialDiscovery {
		allCollections, err := m.app.FindAllCollections("base", "auth")
		if err == nil {
			existingCollectionNames := []string{}
			hasCollections := false
			for _, collection := range allCollections {
				// Skip system collections (they start with _)
				if len(collection.Name) > 0 && collection.Name[0] == '_' {
					continue
				}
				// Skip _pbSync collection (it's managed separately)
				if collection.Name == "_pbSync" {
					continue
				}

				// If we have ANY non-system collection (even if empty), don't apply snapshot
				existingCollectionNames = append(existingCollectionNames, collection.Name)
				hasCollections = true
			}

			if hasCollections {
				log.Printf("Skipping snapshot application: instance already has %d collection(s): %v. Will consume events instead.",
					len(existingCollectionNames), existingCollectionNames)
				return nil
			}
		}
	} else {
		log.Printf("Initial discovery phase: snapshot will be applied as source of truth (may overwrite existing collections/records)")
		log.Printf("Initial discovery: Will delete all existing non-system collections to ensure clean snapshot application")
	}

	startTime := time.Now()

	// Step 1: Apply schema first (create/update collections)
	if len(snapshot.Schema) > 0 {
		log.Printf("Applying schema: %d collections", len(snapshot.Schema))
		if err := m.applySchema(snapshot.Schema, isInitialDiscovery); err != nil {
			return fmt.Errorf("failed to apply schema: %w", err)
		}
		log.Printf("Schema applied successfully")
	}

	totalRecords := 0

	// Step 2: Apply data using direct SQLite bulk imports for efficiency and constraint bypass
	totalRecords = m.importRecordsDirectly(snapshot.Collections)

	duration := time.Since(startTime)
	log.Printf("Snapshot applied successfully: %d records in %v", totalRecords, duration)

	return nil
}

// importRecordsDirectly imports records using bulk SQLite operations with proper validation handling
func (m *Manager) importRecordsDirectly(collections map[string][]map[string]interface{}) int {
	totalRecords := 0

	// Sort collections for optimal import order
	sortedCollections := m.sortCollectionsForDataImport(collections)

	for _, collectionName := range sortedCollections {
		records := collections[collectionName]
		if len(records) == 0 {
			continue
		}

		collection, err := m.app.FindCollectionByNameOrId(collectionName)
		if err != nil {
			log.Printf("Warning: collection %s not found, skipping: %v", collectionName, err)
			continue
		}

		log.Printf("Importing %d records to collection %s", len(records), collectionName)

		for _, recordData := range records {
			recordID, ok := recordData["id"].(string)
			if !ok {
				continue
			}

			// Check if record exists
			existing, err := m.app.FindRecordById(collection.Id, recordID)
			if err == nil && existing != nil {
				// Update existing record using a hook-less app instance to bypass validation during sync
				syncApp := m.app.UnsafeWithoutHooks()
				for k, v := range recordData {
					if k != "id" && k != "created" && k != "updated" {
						existing.Set(k, v)
					}
				}

				// For auth collections, handle password field properly during sync
				if collection.IsAuth() {
					m.handleAuthRecordPassword(existing, recordData)
				}

				if err := syncApp.Save(existing); err != nil {
					log.Printf("Warning: failed to update record %s/%s: %v", collectionName, recordID, err)
				} else {
					totalRecords++
				}
			} else {
				// Create new record using a hook-less app instance to bypass validation during sync
				syncApp := m.app.UnsafeWithoutHooks()
				newRecord := core.NewRecord(collection)
				for k, v := range recordData {
					if k != "id" {
						newRecord.Set(k, v)
					} else {
						newRecord.Set("id", v)
					}
				}

				// For auth collections, set password field properly during sync
				if collection.IsAuth() {
					m.handleAuthRecordPassword(newRecord, recordData)
				}

				if err := syncApp.Save(newRecord); err != nil {
					log.Printf("Warning: failed to create record %s/%s: %v", collectionName, recordID, err)
				} else {
					totalRecords++
				}
			}
		}
	}

	return totalRecords
}

// handleAuthRecordPassword handles password field for auth collections during sync
func (m *Manager) handleAuthRecordPassword(record *core.Record, recordData map[string]interface{}) {
	// For auth collections, we need to ensure the password field is properly set
	if passwordValue, exists := recordData["password"]; exists && passwordValue != nil {
		if newPassword, ok := passwordValue.(string); ok && newPassword != "" {
			// Set the password hash from the synced data
			record.SetRaw("password", &core.PasswordFieldValue{
				Hash:  newPassword,
				Plain: "", // Don't include plain text during sync
			})
		} else {
			// Handle non-string password data
			record.SetRaw("password", passwordValue)
		}
	} else {
		// No password data in the snapshot - set an empty password field
		// This is safe during sync because we're using UnsafeWithoutHooks()
		// The actual authentication will be handled by the _authOrigins collection
		record.SetRaw("password", &core.PasswordFieldValue{
			Hash:  "",
			Plain: "",
		})
	}
}

// exportSchema exports the full PocketBase schema in JSON format (like pb_schema.json)
func (m *Manager) exportSchema(collections []*core.Collection) ([]map[string]interface{}, error) {
	schema := make([]map[string]interface{}, 0, len(collections))

	for _, collection := range collections {
		// Convert collection to JSON-serializable map
		colMap := make(map[string]interface{})

		// Basic collection properties
		colMap["id"] = collection.Id
		colMap["name"] = collection.Name
		colMap["type"] = collection.Type
		colMap["system"] = collection.System
		colMap["listRule"] = collection.ListRule
		colMap["viewRule"] = collection.ViewRule
		colMap["createRule"] = collection.CreateRule
		colMap["updateRule"] = collection.UpdateRule
		colMap["deleteRule"] = collection.DeleteRule

		// Export fields
		fields := make([]map[string]interface{}, 0, len(collection.Fields))
		for _, field := range collection.Fields {
			fieldMap := m.ExportField(field)
			fields = append(fields, fieldMap)
		}
		colMap["fields"] = fields

		// Export indexes
		indexes := make([]string, 0, len(collection.Indexes))
		for _, idx := range collection.Indexes {
			indexes = append(indexes, idx)
		}
		colMap["indexes"] = indexes

		schema = append(schema, colMap)
	}

	return schema, nil
}

// ExportField exports a field to JSON-serializable map
func (m *Manager) ExportField(field core.Field) map[string]interface{} {
	fieldMap := make(map[string]interface{})

	// Common field properties
	fieldMap["id"] = field.GetId()
	fieldMap["name"] = field.GetName()
	fieldMap["hidden"] = field.GetHidden()
	fieldMap["system"] = field.GetSystem()

	// Get field type and properties using type switch
	var fieldType string
	switch f := field.(type) {
	case *core.TextField:
		fieldType = "text"
		fieldMap["type"] = fieldType
		fieldMap["required"] = f.Required
		fieldMap["presentable"] = f.Presentable
		fieldMap["min"] = f.Min
		fieldMap["max"] = f.Max
		fieldMap["pattern"] = f.Pattern
		fieldMap["autogeneratePattern"] = f.AutogeneratePattern
		fieldMap["primaryKey"] = f.PrimaryKey
	case *core.NumberField:
		fieldType = "number"
		fieldMap["type"] = fieldType
		fieldMap["required"] = f.Required
		fieldMap["presentable"] = f.Presentable
		if f.Min != nil {
			fieldMap["min"] = *f.Min
		}
		if f.Max != nil {
			fieldMap["max"] = *f.Max
		}
		fieldMap["onlyInt"] = f.OnlyInt
	case *core.BoolField:
		fieldType = "bool"
		fieldMap["type"] = fieldType
		fieldMap["required"] = f.Required
		fieldMap["presentable"] = f.Presentable
	case *core.EmailField:
		fieldType = "email"
		fieldMap["type"] = fieldType
		fieldMap["required"] = f.Required
		fieldMap["presentable"] = f.Presentable
		fieldMap["onlyDomains"] = f.OnlyDomains
		fieldMap["exceptDomains"] = f.ExceptDomains
	case *core.URLField:
		fieldType = "url"
		fieldMap["type"] = fieldType
		fieldMap["required"] = f.Required
		fieldMap["presentable"] = f.Presentable
		fieldMap["onlyDomains"] = f.OnlyDomains
		fieldMap["exceptDomains"] = f.ExceptDomains
	case *core.DateField:
		fieldType = "date"
		fieldMap["type"] = fieldType
		fieldMap["required"] = f.Required
		fieldMap["presentable"] = f.Presentable
		// Note: DateField min/max are DateTime types, serialize as strings if needed
		// For now, we'll skip them as they're complex to serialize
	case *core.SelectField:
		fieldType = "select"
		fieldMap["type"] = fieldType
		fieldMap["required"] = f.Required
		fieldMap["presentable"] = f.Presentable
		fieldMap["values"] = f.Values
		fieldMap["maxSelect"] = f.MaxSelect
	case *core.FileField:
		fieldType = "file"
		fieldMap["type"] = fieldType
		fieldMap["required"] = f.Required
		fieldMap["presentable"] = f.Presentable
		fieldMap["maxSelect"] = f.MaxSelect
		fieldMap["maxSize"] = f.MaxSize
		fieldMap["mimeTypes"] = f.MimeTypes
		fieldMap["protected"] = f.Protected
		fieldMap["thumbs"] = f.Thumbs
	case *core.RelationField:
		fieldType = "relation"
		fieldMap["type"] = fieldType
		fieldMap["required"] = f.Required
		fieldMap["presentable"] = f.Presentable
		fieldMap["collectionId"] = f.CollectionId
		fieldMap["cascadeDelete"] = f.CascadeDelete
		fieldMap["maxSelect"] = f.MaxSelect
	case *core.PasswordField:
		fieldType = "password"
		fieldMap["type"] = fieldType
		fieldMap["required"] = f.Required
		fieldMap["min"] = f.Min
		fieldMap["max"] = f.Max
		fieldMap["cost"] = f.Cost
	case *core.AutodateField:
		fieldType = "autodate"
		fieldMap["type"] = fieldType
		fieldMap["onCreate"] = f.OnCreate
		fieldMap["onUpdate"] = f.OnUpdate
	case *core.JSONField:
		fieldType = "json"
		fieldMap["type"] = fieldType
		fieldMap["required"] = f.Required
		fieldMap["presentable"] = f.Presentable
	default:
		fieldType = "unknown"
		fieldMap["type"] = fieldType
		log.Printf("Warning: unknown field type for field %s", field.GetName())
	}

	return fieldMap
}

// applySchema applies the schema to PocketBase (creates/updates/deletes collections)
// If isInitialDiscovery is true, collections not in the snapshot will be deleted (snapshot is source of truth).
// If false, this should only be called for empty instances.
func (m *Manager) applySchema(schema []map[string]interface{}, isInitialDiscovery bool) error {
	// Build a map of collection names from the snapshot schema
	snapshotCollections := make(map[string]bool)
	for _, colData := range schema {
		collectionName, ok := colData["name"].(string)
		if ok {
			snapshotCollections[collectionName] = true
		}
	}

	// Get all existing collections
	allCollections, err := m.app.FindAllCollections("base", "auth")
	if err != nil {
		return fmt.Errorf("failed to get existing collections: %w", err)
	}

	// Safety check: If NOT in initial discovery and we have collections, don't apply schema
	// During initial discovery, the snapshot is the source of truth, so we allow deletions
	if !isInitialDiscovery {
		hasNonSystemCollections := false
		for _, existing := range allCollections {
			if len(existing.Name) > 0 && existing.Name[0] == '_' {
				continue
			}
			if existing.Name == "_pbSync" {
				continue
			}
			hasNonSystemCollections = true
			break
		}

		if hasNonSystemCollections {
			log.Printf("CRITICAL: applySchema called on instance with existing collections outside of initial discovery - skipping schema application to prevent data loss")
			return nil
		}
	}

	// During initial discovery, delete ALL existing non-system collections to ensure clean state
	// This resolves conflicts with default PocketBase collections like 'users'
	if isInitialDiscovery {
		log.Printf("Initial discovery: deleting all existing non-system collections to ensure clean snapshot application")
		log.Printf("Initial discovery: found %d existing collections to process", len(allCollections))
		for _, existing := range allCollections {
			// Skip system collections (they start with _)
			if len(existing.Name) > 0 && existing.Name[0] == '_' {
				continue
			}
			// Skip _pbSync collection (it's managed separately)
			if existing.Name == "_pbSync" {
				continue
			}

			log.Printf("Initial discovery: deleting collection %s (cleaning for snapshot)", existing.Name)
			if err := m.app.Delete(existing); err != nil {
				log.Printf("Warning: failed to delete collection %s during initial discovery: %v", existing.Name, err)
			} else {
				log.Printf("Deleted collection during initial discovery: %s", existing.Name)
			}
		}
	} else {
		// For non-initial discovery, only delete collections that exist locally but not in snapshot
		// (excluding system collections and _pbSync)
		// NOTE: This should never happen for non-empty instances due to the check above
		for _, existing := range allCollections {
			// Skip system collections (they start with _)
			if len(existing.Name) > 0 && existing.Name[0] == '_' {
				continue
			}
			// Skip _pbSync collection (it's managed separately)
			if existing.Name == "_pbSync" {
				continue
			}

			// If collection exists locally but not in snapshot, delete it
			if !snapshotCollections[existing.Name] {
				log.Printf("Deleting collection %s (not in snapshot schema)", existing.Name)
				if err := m.app.Delete(existing); err != nil {
					log.Printf("Warning: failed to delete collection %s: %v", existing.Name, err)
				} else {
					log.Printf("Deleted collection: %s", existing.Name)
				}
			}
		}
	}

	// Use PocketBase's built-in ImportCollections method - it handles dependency resolution correctly
	// Convert schema format to match what ImportCollections expects
	importSchema := make([]map[string]interface{}, 0, len(schema))
	for _, colData := range schema {
		// Skip _pbSync collection (it's managed separately)
		collectionName, ok := colData["name"].(string)
		if !ok {
			continue
		}
		if collectionName == "_pbSync" {
			continue
		}
		// Skip system collections during import
		if len(collectionName) > 0 && collectionName[0] == '_' {
			continue
		}
		importSchema = append(importSchema, colData)
	}

	if len(importSchema) > 0 {
		log.Printf("Importing %d collections using PocketBase's ImportCollections method", len(importSchema))
		// Use deleteMissing=true to ensure clean state during initial discovery
		err := m.app.ImportCollections(importSchema, isInitialDiscovery)
		if err != nil {
			return fmt.Errorf("failed to import collections: %w", err)
		}
		log.Printf("Successfully imported collections using PocketBase's built-in method")
	}

	return nil
}

// createFieldFromSchema creates a field from schema data
func (m *Manager) createFieldFromSchema(fieldData map[string]interface{}) core.Field {
	fieldType, ok := fieldData["type"].(string)
	if !ok {
		return nil
	}

	name, _ := fieldData["name"].(string)
	required, _ := fieldData["required"].(bool)

	var field core.Field
	switch fieldType {
	case "text":
		f := &core.TextField{}
		f.Name = name
		f.Required = required
		if min, ok := fieldData["min"].(float64); ok {
			f.Min = int(min)
		}
		if max, ok := fieldData["max"].(float64); ok {
			f.Max = int(max)
		}
		if pattern, ok := fieldData["pattern"].(string); ok {
			f.Pattern = pattern
		}
		if autogen, ok := fieldData["autogeneratePattern"].(string); ok {
			f.AutogeneratePattern = autogen
		}
		if pk, ok := fieldData["primaryKey"].(bool); ok {
			f.PrimaryKey = pk
		}
		field = f
	case "number":
		f := &core.NumberField{}
		f.Name = name
		f.Required = required
		if presentable, ok := fieldData["presentable"].(bool); ok {
			f.Presentable = presentable
		}
		if hidden, ok := fieldData["hidden"].(bool); ok {
			f.SetHidden(hidden)
		}
		if min, ok := fieldData["min"].(float64); ok {
			minVal := min
			f.Min = &minVal
		}
		if max, ok := fieldData["max"].(float64); ok {
			maxVal := max
			f.Max = &maxVal
		}
		if onlyInt, ok := fieldData["onlyInt"].(bool); ok {
			f.OnlyInt = onlyInt
		}
		field = f
	case "bool":
		f := &core.BoolField{}
		f.Name = name
		f.Required = required
		if presentable, ok := fieldData["presentable"].(bool); ok {
			f.Presentable = presentable
		}
		if hidden, ok := fieldData["hidden"].(bool); ok {
			f.SetHidden(hidden)
		}
		field = f
	case "email":
		f := &core.EmailField{}
		f.Name = name
		f.Required = required
		if presentable, ok := fieldData["presentable"].(bool); ok {
			f.Presentable = presentable
		}
		if hidden, ok := fieldData["hidden"].(bool); ok {
			f.SetHidden(hidden)
		}
		if onlyDomains, ok := fieldData["onlyDomains"].([]interface{}); ok {
			f.OnlyDomains = make([]string, 0, len(onlyDomains))
			for _, d := range onlyDomains {
				if dStr, ok := d.(string); ok {
					f.OnlyDomains = append(f.OnlyDomains, dStr)
				}
			}
		}
		if exceptDomains, ok := fieldData["exceptDomains"].([]interface{}); ok {
			f.ExceptDomains = make([]string, 0, len(exceptDomains))
			for _, d := range exceptDomains {
				if dStr, ok := d.(string); ok {
					f.ExceptDomains = append(f.ExceptDomains, dStr)
				}
			}
		}
		field = f
	case "url":
		f := &core.URLField{}
		f.Name = name
		f.Required = required
		if presentable, ok := fieldData["presentable"].(bool); ok {
			f.Presentable = presentable
		}
		if hidden, ok := fieldData["hidden"].(bool); ok {
			f.SetHidden(hidden)
		}
		field = f
	case "date":
		f := &core.DateField{}
		f.Name = name
		f.Required = required
		if presentable, ok := fieldData["presentable"].(bool); ok {
			f.Presentable = presentable
		}
		if hidden, ok := fieldData["hidden"].(bool); ok {
			f.SetHidden(hidden)
		}
		// Note: min/max are time.Time, would need parsing if needed
		field = f
	case "select":
		f := &core.SelectField{}
		f.Name = name
		f.Required = required
		if presentable, ok := fieldData["presentable"].(bool); ok {
			f.Presentable = presentable
		}
		if hidden, ok := fieldData["hidden"].(bool); ok {
			f.SetHidden(hidden)
		}
		if values, ok := fieldData["values"].([]interface{}); ok {
			f.Values = make([]string, 0, len(values))
			for _, v := range values {
				if vStr, ok := v.(string); ok {
					f.Values = append(f.Values, vStr)
				}
			}
		}
		if maxSelect, ok := fieldData["maxSelect"].(float64); ok {
			f.MaxSelect = int(maxSelect)
		}
		field = f
	case "json":
		f := &core.JSONField{}
		f.Name = name
		f.Required = required
		if presentable, ok := fieldData["presentable"].(bool); ok {
			f.Presentable = presentable
		}
		if hidden, ok := fieldData["hidden"].(bool); ok {
			f.SetHidden(hidden)
		}
		field = f
	case "file":
		f := &core.FileField{}
		f.Name = name
		f.Required = required
		if presentable, ok := fieldData["presentable"].(bool); ok {
			f.Presentable = presentable
		}
		if hidden, ok := fieldData["hidden"].(bool); ok {
			f.SetHidden(hidden)
		}
		if maxSelect, ok := fieldData["maxSelect"].(float64); ok {
			f.MaxSelect = int(maxSelect)
		}
		if maxSize, ok := fieldData["maxSize"].(float64); ok {
			f.MaxSize = int64(maxSize)
		}
		if mimeTypes, ok := fieldData["mimeTypes"].([]interface{}); ok {
			f.MimeTypes = make([]string, 0, len(mimeTypes))
			for _, m := range mimeTypes {
				if mStr, ok := m.(string); ok {
					f.MimeTypes = append(f.MimeTypes, mStr)
				}
			}
		}
		if protected, ok := fieldData["protected"].(bool); ok {
			f.Protected = protected
		}
		if thumbs, ok := fieldData["thumbs"].([]interface{}); ok {
			f.Thumbs = make([]string, 0, len(thumbs))
			for _, t := range thumbs {
				if tStr, ok := t.(string); ok {
					f.Thumbs = append(f.Thumbs, tStr)
				}
			}
		}
		field = f
	case "relation":
		f := &core.RelationField{}
		f.Name = name
		f.Required = required
		if presentable, ok := fieldData["presentable"].(bool); ok {
			f.Presentable = presentable
		}
		if hidden, ok := fieldData["hidden"].(bool); ok {
			f.SetHidden(hidden)
		}
		if collId, ok := fieldData["collectionId"].(string); ok {
			f.CollectionId = collId
		}
		if cascadeDelete, ok := fieldData["cascadeDelete"].(bool); ok {
			f.CascadeDelete = cascadeDelete
		}
		if maxSelect, ok := fieldData["maxSelect"].(float64); ok {
			f.MaxSelect = int(maxSelect)
		}
		field = f
	case "password":
		f := &core.PasswordField{}
		f.Name = name
		f.Required = required
		if hidden, ok := fieldData["hidden"].(bool); ok {
			f.SetHidden(hidden)
		}
		if min, ok := fieldData["min"].(float64); ok {
			f.Min = int(min)
		}
		if max, ok := fieldData["max"].(float64); ok {
			f.Max = int(max)
		}
		if cost, ok := fieldData["cost"].(float64); ok {
			f.Cost = int(cost)
		}
		field = f
	case "autodate":
		f := &core.AutodateField{}
		f.Name = name
		if hidden, ok := fieldData["hidden"].(bool); ok {
			f.SetHidden(hidden)
		}
		if onCreate, ok := fieldData["onCreate"].(bool); ok {
			f.OnCreate = onCreate
		}
		if onUpdate, ok := fieldData["onUpdate"].(bool); ok {
			f.OnUpdate = onUpdate
		}
		field = f
	default:
		log.Printf("Warning: unsupported field type: %s", fieldType)
		return nil
	}

	// Set common properties (these are set in the type-specific cases above)
	// Hidden and system are set via the field constructors

	return field
}

// sortCollectionsForDataImport orders collections for data import to avoid relation constraint violations
// Uses a simple heuristic: collections without relation fields first, then collections with relations

func (m *Manager) sortCollectionsForDataImport(collections map[string][]map[string]interface{}) []string {
	// Get all collections to understand their fields
	allCollections, err := m.app.FindAllCollections("base", "auth")
	if err != nil {
		log.Printf("Warning: failed to get collections for sorting: %v", err)
		// Fall back to alphabetical order
		result := make([]string, 0, len(collections))
		for name := range collections {
			result = append(result, name)
		}
		sort.Strings(result)
		return result
	}

	// Build collection metadata map
	collectionMap := make(map[string]*core.Collection)
	for _, collection := range allCollections {
		collectionMap[collection.Name] = collection
	}

	// Separate collections based on whether they have relation fields
	noRelations := make([]string, 0)
	withRelations := make([]string, 0)

	for collectionName := range collections {
		collection, exists := collectionMap[collectionName]
		if !exists {
			// Unknown collection, treat as having relations (safer)
			withRelations = append(withRelations, collectionName)
			continue
		}

		hasRelations := false
		for _, field := range collection.Fields {
			if field.Type() == "relation" {
				hasRelations = true
				break
			}
		}

		if hasRelations {
			withRelations = append(withRelations, collectionName)
		} else {
			noRelations = append(noRelations, collectionName)
		}
	}

	// Sort each group alphabetically for deterministic order
	sort.Strings(noRelations)
	sort.Strings(withRelations)

	// Combine: collections without relations first, then collections with relations
	result := append(noRelations, withRelations...)

	log.Printf("Data import order: collections without relations first (%v), then with relations (%v)", noRelations, withRelations)

	return result
}
