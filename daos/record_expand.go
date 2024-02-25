package daos

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tools/dbutils"
	"github.com/pocketbase/pocketbase/tools/list"
	"github.com/pocketbase/pocketbase/tools/security"
	"github.com/pocketbase/pocketbase/tools/types"
)

// MaxExpandDepth specifies the max allowed nested expand depth path.
//
// @todo Consider eventually reusing resolvers.maxNestedRels
const MaxExpandDepth = 6

// ExpandFetchFunc defines the function that is used to fetch the expanded relation records.
type ExpandFetchFunc func(relCollection *models.Collection, relIds []string) ([]*models.Record, error)

// ExpandRecord expands the relations of a single Record model.
//
// If optFetchFunc is not set, then a default function will be used
// that returns all relation records.
//
// Returns a map with the failed expand parameters and their errors.
func (dao *Dao) ExpandRecord(record *models.Record, expands []string, optFetchFunc ExpandFetchFunc) map[string]error {
	return dao.ExpandRecords([]*models.Record{record}, expands, optFetchFunc)
}

// ExpandRecords expands the relations of the provided Record models list.
//
// If optFetchFunc is not set, then a default function will be used
// that returns all relation records.
//
// Returns a map with the failed expand parameters and their errors.
func (dao *Dao) ExpandRecords(records []*models.Record, expands []string, optFetchFunc ExpandFetchFunc) map[string]error {
	normalized := normalizeExpands(expands)

	failed := map[string]error{}

	for _, expand := range normalized {
		if err := dao.expandRecords(records, expand, optFetchFunc, 1); err != nil {
			failed[expand] = err
		}
	}

	return failed
}

// Deprecated
var indirectExpandRegexOld = regexp.MustCompile(`^(\w+)\((\w+)\)$`)

var indirectExpandRegex = regexp.MustCompile(`^(\w+)_via_(\w+)$`)

// notes:
// - if fetchFunc is nil, dao.FindRecordsByIds will be used
// - all records are expected to be from the same collection
// - if MaxExpandDepth is reached, the function returns nil ignoring the remaining expand path
func (dao *Dao) expandRecords(records []*models.Record, expandPath string, fetchFunc ExpandFetchFunc, recursionLevel int) error {
	if fetchFunc == nil {
		// load a default fetchFunc
		fetchFunc = func(relCollection *models.Collection, relIds []string) ([]*models.Record, error) {
			return dao.FindRecordsByIds(relCollection.Id, relIds)
		}
	}

	if expandPath == "" || recursionLevel > MaxExpandDepth || len(records) == 0 {
		return nil
	}

	mainCollection := records[0].Collection()

	var relField *schema.SchemaField
	var relFieldOptions *schema.RelationOptions
	var relCollection *models.Collection

	parts := strings.SplitN(expandPath, ".", 2)
	var matches []string

	// @todo remove the old syntax support
	if strings.Contains(parts[0], "(") {
		matches = indirectExpandRegexOld.FindStringSubmatch(parts[0])
		if len(matches) == 3 {
			log.Printf(
				"%s expand format is deprecated and will be removed in the future. Consider replacing it with %s_via_%s.\n",
				matches[0],
				matches[1],
				matches[2],
			)
		}
	} else {
		matches = indirectExpandRegex.FindStringSubmatch(parts[0])
	}

	if len(matches) == 3 {
		indirectRel, _ := dao.FindCollectionByNameOrId(matches[1])
		if indirectRel == nil {
			return fmt.Errorf("couldn't find back-related collection %q", matches[1])
		}

		indirectRelField := indirectRel.Schema.GetFieldByName(matches[2])
		if indirectRelField == nil || indirectRelField.Type != schema.FieldTypeRelation {
			return fmt.Errorf("couldn't find back-relation field %q in collection %q", matches[2], indirectRel.Name)
		}

		indirectRelField.InitOptions()
		indirectRelFieldOptions, _ := indirectRelField.Options.(*schema.RelationOptions)
		if indirectRelFieldOptions == nil || indirectRelFieldOptions.CollectionId != mainCollection.Id {
			return fmt.Errorf("invalid back-relation field path %q", parts[0])
		}

		// add the related id(s) as a dynamic relation field value to
		// allow further expand checks at later stage in a more unified manner
		prepErr := func() error {
			q := dao.DB().Select("id").
				From(indirectRel.Name).
				Limit(1000) // the limit is arbitrary chosen and may change in the future

			if indirectRelFieldOptions.IsMultiple() {
				q.AndWhere(dbx.Exists(dbx.NewExp(fmt.Sprintf(
					"SELECT 1 FROM %s je WHERE je.value = {:id}",
					dbutils.JsonEach(indirectRelField.Name),
				))))
			} else {
				q.AndWhere(dbx.NewExp("[[" + indirectRelField.Name + "]] = {:id}"))
			}

			pq := q.Build().Prepare()

			for _, record := range records {
				var relIds []string

				err := pq.Bind(dbx.Params{"id": record.Id}).Column(&relIds)
				if err != nil {
					return errors.Join(err, pq.Close())
				}

				if len(relIds) > 0 {
					record.Set(parts[0], relIds)
				}
			}

			return pq.Close()
		}()
		if prepErr != nil {
			return prepErr
		}

		relFieldOptions = &schema.RelationOptions{
			MaxSelect:    nil,
			CollectionId: indirectRel.Id,
		}
		if dbutils.HasSingleColumnUniqueIndex(indirectRelField.Name, indirectRel.Indexes) {
			relFieldOptions.MaxSelect = types.Pointer(1)
		}
		// indirect/back relation
		relField = &schema.SchemaField{
			Id:      "_" + parts[0] + security.PseudorandomString(3),
			Type:    schema.FieldTypeRelation,
			Name:    parts[0],
			Options: relFieldOptions,
		}
		relCollection = indirectRel
	} else {
		// direct relation
		relField = mainCollection.Schema.GetFieldByName(parts[0])
		if relField == nil || relField.Type != schema.FieldTypeRelation {
			return fmt.Errorf("Couldn't find relation field %q in collection %q.", parts[0], mainCollection.Name)
		}
		relField.InitOptions()
		relFieldOptions, _ = relField.Options.(*schema.RelationOptions)
		if relFieldOptions == nil {
			return fmt.Errorf("Couldn't initialize the options of relation field %q.", parts[0])
		}

		relCollection, _ = dao.FindCollectionByNameOrId(relFieldOptions.CollectionId)
		if relCollection == nil {
			return fmt.Errorf("Couldn't find related collection %q.", relFieldOptions.CollectionId)
		}
	}

	// ---------------------------------------------------------------

	// extract the id of the relations to expand
	relIds := make([]string, 0, len(records))
	for _, record := range records {
		relIds = append(relIds, record.GetStringSlice(relField.Name)...)
	}

	// fetch rels
	rels, relsErr := fetchFunc(relCollection, relIds)
	if relsErr != nil {
		return relsErr
	}

	// expand nested fields
	if len(parts) > 1 {
		err := dao.expandRecords(rels, parts[1], fetchFunc, recursionLevel+1)
		if err != nil {
			return err
		}
	}

	// reindex with the rel id
	indexedRels := make(map[string]*models.Record, len(rels))
	for _, rel := range rels {
		indexedRels[rel.GetId()] = rel
	}

	for _, model := range records {
		relIds := model.GetStringSlice(relField.Name)

		validRels := make([]*models.Record, 0, len(relIds))
		for _, id := range relIds {
			if rel, ok := indexedRels[id]; ok {
				validRels = append(validRels, rel)
			}
		}

		if len(validRels) == 0 {
			continue // no valid relations
		}

		expandData := model.Expand()

		// normalize access to the previously expanded rel records (if any)
		var oldExpandedRels []*models.Record
		switch v := expandData[relField.Name].(type) {
		case nil:
			// no old expands
		case *models.Record:
			oldExpandedRels = []*models.Record{v}
		case []*models.Record:
			oldExpandedRels = v
		}

		// merge expands
		for _, oldExpandedRel := range oldExpandedRels {
			// find a matching rel record
			for _, rel := range validRels {
				if rel.Id != oldExpandedRel.Id {
					continue
				}

				rel.MergeExpand(oldExpandedRel.Expand())
			}
		}

		// update the expanded data
		if relFieldOptions.MaxSelect != nil && *relFieldOptions.MaxSelect <= 1 {
			expandData[relField.Name] = validRels[0]
		} else {
			expandData[relField.Name] = validRels
		}

		model.SetExpand(expandData)
	}

	return nil
}

// normalizeExpands normalizes expand strings and merges self containing paths
// (eg. ["a.b.c", "a.b", "   test  ", "  ", "test"] -> ["a.b.c", "test"]).
func normalizeExpands(paths []string) []string {
	// normalize paths
	normalized := make([]string, 0, len(paths))
	for _, p := range paths {
		p = strings.ReplaceAll(p, " ", "") // replace spaces
		p = strings.Trim(p, ".")           // trim incomplete paths
		if p != "" {
			normalized = append(normalized, p)
		}
	}

	// merge containing paths
	result := make([]string, 0, len(normalized))
	for i, p1 := range normalized {
		var skip bool
		for j, p2 := range normalized {
			if i == j {
				continue
			}
			if strings.HasPrefix(p2, p1+".") {
				// skip because there is more detailed expand path
				skip = true
				break
			}
		}
		if !skip {
			result = append(result, p1)
		}
	}

	return list.ToUniqueStringSlice(result)
}

func isRelFieldUnique(collection *models.Collection, fieldName string) bool {
	for _, idx := range collection.Indexes {
		parsed := dbutils.ParseIndex(idx)
		if parsed.Unique && len(parsed.Columns) == 1 && strings.EqualFold(parsed.Columns[0].Name, fieldName) {
			return true
		}
	}

	return false
}
