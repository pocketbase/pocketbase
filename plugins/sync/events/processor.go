package events

import (
	"fmt"
	"log"

	"github.com/pocketbase/pocketbase/plugins/sync/sync"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/dbx"
)

// Processor handles processing incoming events and applying changes
type Processor struct {
	app            *pocketbase.PocketBase
	instanceID     string
	conflicts      map[string]*Event // Track pending conflicts: recordID -> event
}

// NewProcessor creates a new event processor
func NewProcessor(app *pocketbase.PocketBase, instanceID string) *Processor {
	return &Processor{
		app:        app,
		instanceID: instanceID,
		conflicts:  make(map[string]*Event),
	}
}

// Process processes an incoming event
func (p *Processor) Process(event *Event) error {
	// Check for conflicts (concurrent updates to same record)
	if event.Type == EventTypeUpdate || event.Type == EventTypeDelete {
		// Check if we have a pending conflict for this record
		if conflict, exists := p.conflicts[event.RecordID]; exists {
			// Resolve conflict using last-write-wins
			if event.Timestamp.After(conflict.Timestamp) {
				// New event wins, remove old conflict
				delete(p.conflicts, event.RecordID)
			} else {
				// Old conflict wins, skip this event
				log.Printf("Skipping event due to conflict resolution (old event wins): %s/%s", event.Collection, event.RecordID)
				return nil
			}
		}

		// Check if this event is concurrent with any local operations
		// For simplicity, we'll use timestamp-based conflict detection
		// In a more sophisticated implementation, we'd track pending local operations
	}

	// Apply change to local PocketBase
	if err := p.ApplyEvent(event); err != nil {
		return fmt.Errorf("failed to apply event: %w", err)
	}

	return nil
}

// ApplyEvent applies an event to the local PocketBase instance using the SDK
// This is a public method that can be called directly to avoid recursion in handlers
func (p *Processor) ApplyEvent(event *Event) error {
	switch event.Type {
	case EventTypeCreate:
		return p.applyCreate(event)
	case EventTypeUpdate:
		return p.applyUpdate(event)
	case EventTypeDelete:
		return p.applyDelete(event)
	case EventTypeCollectionCreate:
		return p.applyCollectionCreate(event)
	case EventTypeCollectionUpdate:
		return p.applyCollectionUpdate(event)
	case EventTypeCollectionDelete:
		return p.applyCollectionDelete(event)
	default:
		return fmt.Errorf("unknown event type: %s", event.Type)
	}
}

// applyCreate creates a record in PocketBase using the SDK
func (p *Processor) applyCreate(event *Event) error {
	collection, err := p.app.FindCollectionByNameOrId(event.Collection)
	if err != nil {
		return fmt.Errorf("failed to find collection %s: %w", event.Collection, err)
	}

	// Check if record with this ID already exists (upsert behavior)
	if event.RecordID != "" {
		existing := &core.Record{}
		err = p.app.RecordQuery(collection.Id).
			AndWhere(dbx.HashExp{"id": event.RecordID}).
			Limit(1).
			One(existing)
		if err == nil && existing.Id != "" {
			// Record exists, update it instead
			return p.applyUpdate(event)
		}
	}

	// Special handling for _pbSync: also check if record exists by instanceID (unique field)
	if event.Collection == sync.PBSyncCollectionName {
		instanceID, ok := event.RecordData["instanceID"].(string)
		if ok && instanceID != "" {
			existing := &core.Record{}
			err = p.app.RecordQuery(collection.Id).
				AndWhere(dbx.HashExp{"instanceID": instanceID}).
				Limit(1).
				One(existing)
			if err == nil && existing.Id != "" {
				// Record exists with this instanceID, update it instead
				// Use the existing record's ID for the update
				event.RecordID = existing.Id
				return p.applyUpdate(event)
			}
		}
	}

	// Create record
	record := core.NewRecord(collection)
	
	// Build a map of field definitions for quick lookup
	fieldMap := make(map[string]core.Field)
	for _, field := range collection.Fields {
		fieldMap[field.GetName()] = field
	}
	
	// Set all fields from the event data
	for key, value := range event.RecordData {
		// Skip system fields that shouldn't be set directly
		if key == "created" || key == "updated" {
			continue
		}
		
		// Skip tokenKey - it's instance-specific and shouldn't be synced
		if key == "tokenKey" {
			continue
		}
		
		// Get field definition
		field, hasField := fieldMap[key]
		
		// Handle relation fields - validate that referenced record exists
		if hasField {
			if relField, ok := field.(*core.RelationField); ok {
				// Skip if value is empty
				if value == nil || value == "" {
					continue
				}
				// Check if referenced record exists
				refRecordID, ok := value.(string)
				if !ok {
					// Try to extract ID from value if it's a map/object
					continue
				}
				// Validate that the referenced record exists
				refCollection, err := p.app.FindCollectionByNameOrId(relField.CollectionId)
				if err != nil {
					log.Printf("Warning: Cannot validate relation field %s: collection %s not found, skipping", key, relField.CollectionId)
					continue
				}
				refRecord := &core.Record{}
				err = p.app.RecordQuery(refCollection.Id).
					AndWhere(dbx.HashExp{"id": refRecordID}).
					Limit(1).
					One(refRecord)
				if err != nil {
					log.Printf("Warning: Relation field %s references non-existent record %s, skipping", key, refRecordID)
					continue
				}
			}
		}
		
		// Set ID explicitly if provided
		if key == "id" {
			record.Set("id", value)
		} else {
			record.Set(key, value)
		}
	}

	if err := p.app.Save(record); err != nil {
		return fmt.Errorf("failed to save record: %w", err)
	}

	return nil
}

// applyUpdate updates a record in PocketBase using the SDK
func (p *Processor) applyUpdate(event *Event) error {
	collection, err := p.app.FindCollectionByNameOrId(event.Collection)
	if err != nil {
		return fmt.Errorf("failed to find collection %s: %w", event.Collection, err)
	}

	// Find existing record
	record := &core.Record{}
	err = p.app.RecordQuery(collection.Id).
		AndWhere(dbx.HashExp{"id": event.RecordID}).
		Limit(1).
		One(record)
	
	if err != nil {
		// Record doesn't exist, create it instead
		return p.applyCreate(event)
	}

	// Build a map of field definitions for quick lookup
	fieldMap := make(map[string]core.Field)
	for _, field := range collection.Fields {
		fieldMap[field.GetName()] = field
	}
	
	// Update all fields from the event data
	for key, value := range event.RecordData {
		// Skip system fields that shouldn't be updated directly
		if key == "id" || key == "created" {
			continue
		}
		
		// Skip tokenKey - it's instance-specific and shouldn't be synced
		if key == "tokenKey" {
			continue
		}
		
		// Get field definition
		field, hasField := fieldMap[key]
		
		// Handle relation fields - validate that referenced record exists
		if hasField {
			if relField, ok := field.(*core.RelationField); ok {
				// Skip if value is empty
				if value == nil || value == "" {
					continue
				}
				// Check if referenced record exists
				refRecordID, ok := value.(string)
				if !ok {
					// Try to extract ID from value if it's a map/object
					continue
				}
				// Validate that the referenced record exists
				refCollection, err := p.app.FindCollectionByNameOrId(relField.CollectionId)
				if err != nil {
					log.Printf("Warning: Cannot validate relation field %s: collection %s not found, skipping", key, relField.CollectionId)
					continue
				}
				refRecord := &core.Record{}
				err = p.app.RecordQuery(refCollection.Id).
					AndWhere(dbx.HashExp{"id": refRecordID}).
					Limit(1).
					One(refRecord)
				if err != nil {
					log.Printf("Warning: Relation field %s references non-existent record %s, skipping", key, refRecordID)
					continue
				}
			}
		}
		
		record.Set(key, value)
	}

	if err := p.app.Save(record); err != nil {
		return fmt.Errorf("failed to save record: %w", err)
	}

	return nil
}

// applyDelete deletes a record from PocketBase using the SDK
func (p *Processor) applyDelete(event *Event) error {
	collection, err := p.app.FindCollectionByNameOrId(event.Collection)
	if err != nil {
		return fmt.Errorf("failed to find collection %s: %w", event.Collection, err)
	}

	// Find existing record
	record := &core.Record{}
	err = p.app.RecordQuery(collection.Id).
		AndWhere(dbx.HashExp{"id": event.RecordID}).
		Limit(1).
		One(record)
	
	if err != nil {
		// Record doesn't exist, consider it already deleted
		return nil
	}

	if err := p.app.Delete(record); err != nil {
		return fmt.Errorf("failed to delete record: %w", err)
	}

	return nil
}

// applyCollectionCreate creates a collection from schema event
func (p *Processor) applyCollectionCreate(event *Event) error {
	if event.SchemaData == nil {
		return fmt.Errorf("schema data is missing for collection create event")
	}

	colData := event.SchemaData

	// Skip _pbSync collection (it's managed separately)
	collectionName, ok := colData["name"].(string)
	if !ok {
		return fmt.Errorf("missing collection name in schema data")
	}
	if collectionName == sync.PBSyncCollectionName {
		return nil
	}

	// Check if collection already exists
	existing, err := p.app.FindCollectionByNameOrId(collectionName)
	if err == nil && existing != nil {
		// Collection exists, update it instead
		return p.applyCollectionUpdate(event)
	}

	// Create new collection
	return p.createCollectionFromSchema(colData)
}

// applyCollectionUpdate updates a collection from schema event
func (p *Processor) applyCollectionUpdate(event *Event) error {
	if event.SchemaData == nil {
		return fmt.Errorf("schema data is missing for collection update event")
	}

	colData := event.SchemaData

	collectionName, ok := colData["name"].(string)
	if !ok {
		return fmt.Errorf("missing collection name in schema data")
	}

	// Skip _pbSync collection (it's managed separately)
	if collectionName == sync.PBSyncCollectionName {
		return nil
	}

	// Find existing collection
	collection, err := p.app.FindCollectionByNameOrId(collectionName)
	if err != nil {
		// Collection doesn't exist, create it
		return p.createCollectionFromSchema(colData)
	}

	// Update collection
	return p.updateCollectionFromSchema(collection, colData)
}

// applyCollectionDelete deletes a collection from schema event
func (p *Processor) applyCollectionDelete(event *Event) error {
	collectionName := event.Collection

	// Skip _pbSync collection (it's managed separately)
	if collectionName == sync.PBSyncCollectionName {
		return nil
	}

	// Find collection
	collection, err := p.app.FindCollectionByNameOrId(collectionName)
	if err != nil {
		// Collection doesn't exist, consider it already deleted
		return nil
	}

	// Delete collection
	if err := p.app.Delete(collection); err != nil {
		return fmt.Errorf("failed to delete collection: %w", err)
	}

	return nil
}

// createCollectionFromSchema creates a new collection from schema data
func (p *Processor) createCollectionFromSchema(colData map[string]any) error {
	collectionType, ok := colData["type"].(string)
	if !ok {
		return fmt.Errorf("missing collection type")
	}

	var collection *core.Collection
	switch collectionType {
	case "base":
		collection = core.NewBaseCollection("")
	case "auth":
		collection = core.NewAuthCollection("")
	default:
		return fmt.Errorf("unsupported collection type: %s", collectionType)
	}

	// Set basic properties
	if name, ok := colData["name"].(string); ok {
		collection.Name = name
	}
	if system, ok := colData["system"].(bool); ok {
		collection.System = system
	}
	if listRule, ok := colData["listRule"].(string); ok {
		collection.ListRule = &listRule
	}
	if viewRule, ok := colData["viewRule"].(string); ok {
		collection.ViewRule = &viewRule
	}
	if createRule, ok := colData["createRule"].(string); ok {
		collection.CreateRule = &createRule
	}
	if updateRule, ok := colData["updateRule"].(string); ok {
		collection.UpdateRule = &updateRule
	}
	if deleteRule, ok := colData["deleteRule"].(string); ok {
		collection.DeleteRule = &deleteRule
	}

	// Add fields
	if fields, ok := colData["fields"].([]interface{}); ok {
		for _, fieldData := range fields {
			if fieldMap, ok := fieldData.(map[string]interface{}); ok {
				field := p.createFieldFromSchema(fieldMap)
				if field != nil {
					collection.Fields = append(collection.Fields, field)
				}
			}
		}
	}

	// Add indexes
	if indexes, ok := colData["indexes"].([]interface{}); ok {
		for _, idx := range indexes {
			if idxStr, ok := idx.(string); ok {
				collection.Indexes = append(collection.Indexes, idxStr)
			}
		}
	}

	return p.app.Save(collection)
}

// updateCollectionFromSchema updates an existing collection from schema data
func (p *Processor) updateCollectionFromSchema(collection *core.Collection, colData map[string]any) error {
	// Update rules
	if listRule, ok := colData["listRule"].(string); ok {
		collection.ListRule = &listRule
	}
	if viewRule, ok := colData["viewRule"].(string); ok {
		collection.ViewRule = &viewRule
	}
	if createRule, ok := colData["createRule"].(string); ok {
		collection.CreateRule = &createRule
	}
	if updateRule, ok := colData["updateRule"].(string); ok {
		collection.UpdateRule = &updateRule
	}
	if deleteRule, ok := colData["deleteRule"].(string); ok {
		collection.DeleteRule = &deleteRule
	}

	// Update fields - this is complex, so we'll do a full replacement
	// First, remove all existing fields
	collection.Fields = []core.Field{}

	// Add new fields from schema
	if fields, ok := colData["fields"].([]interface{}); ok {
		for _, fieldData := range fields {
			if fieldMap, ok := fieldData.(map[string]interface{}); ok {
				field := p.createFieldFromSchema(fieldMap)
				if field != nil {
					collection.Fields = append(collection.Fields, field)
				}
			}
		}
	}

	// Update indexes
	collection.Indexes = []string{}
	if indexes, ok := colData["indexes"].([]interface{}); ok {
		for _, idx := range indexes {
			if idxStr, ok := idx.(string); ok {
				collection.Indexes = append(collection.Indexes, idxStr)
			}
		}
	}

	return p.app.Save(collection)
}

// createFieldFromSchema creates a field from schema data (reusing logic from snapshot)
func (p *Processor) createFieldFromSchema(fieldData map[string]any) core.Field {
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

	return field
}

