package core

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/tools/dbutils"
	"github.com/pocketbase/pocketbase/tools/list"
)

// ExpandFetchFunc defines the function that is used to fetch the expanded relation records.
type ExpandFetchFunc func(relCollection *Collection, relIds []string) ([]*Record, error)

// ExpandRecord expands the relations of a single Record model.
//
// If optFetchFunc is not set, then a default function will be used
// that returns all relation records.
//
// Returns a map with the failed expand parameters and their errors.
func (app *BaseApp) ExpandRecord(record *Record, expands []string, optFetchFunc ExpandFetchFunc) map[string]error {
	return app.ExpandRecords([]*Record{record}, expands, optFetchFunc)
}

// ExpandRecords expands the relations of the provided Record models list.
//
// If optFetchFunc is not set, then a default function will be used
// that returns all relation records.
//
// Returns a map with the failed expand parameters and their errors.
func (app *BaseApp) ExpandRecords(records []*Record, expands []string, optFetchFunc ExpandFetchFunc) map[string]error {
	normalized := normalizeExpands(expands)

	failed := map[string]error{}

	for _, expand := range normalized {
		if err := app.expandRecords(records, expand, optFetchFunc, 1); err != nil {
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
// - if maxNestedRels(6) is reached, the function returns nil ignoring the remaining expand path
func (app *BaseApp) expandRecords(records []*Record, expandPath string, fetchFunc ExpandFetchFunc, recursionLevel int) error {
	if fetchFunc == nil {
		// load a default fetchFunc
		fetchFunc = func(relCollection *Collection, relIds []string) ([]*Record, error) {
			return app.FindRecordsByIds(relCollection.Id, relIds)
		}
	}

	if expandPath == "" || recursionLevel > maxNestedRels || len(records) == 0 {
		return nil
	}

	mainCollection := records[0].Collection()

	var relField *RelationField
	var relCollection *Collection

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
		indirectRel, _ := getCollectionByModelOrIdentifier(app, matches[1])
		if indirectRel == nil {
			return fmt.Errorf("couldn't find back-related collection %q", matches[1])
		}

		indirectRelField, _ := indirectRel.Fields.GetByName(matches[2]).(*RelationField)
		if indirectRelField == nil || indirectRelField.CollectionId != mainCollection.Id {
			return fmt.Errorf("couldn't find back-relation field %q in collection %q", matches[2], indirectRel.Name)
		}

		// add the related id(s) as a dynamic relation field value to
		// allow further expand checks at later stage in a more unified manner
		prepErr := func() error {
			q := app.DB().Select("id").
				From(indirectRel.Name).
				Limit(1000) // the limit is arbitrary chosen and may change in the future

			if indirectRelField.IsMultiple() {
				q.AndWhere(dbx.Exists(dbx.NewExp(fmt.Sprintf(
					"SELECT 1 FROM %s je WHERE je.value = {:id}",
					dbutils.JSONEach(indirectRelField.Name),
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

		// indirect/back relation
		relField = &RelationField{
			Name:         parts[0],
			MaxSelect:    2147483647,
			CollectionId: indirectRel.Id,
		}
		if _, ok := dbutils.FindSingleColumnUniqueIndex(indirectRel.Indexes, indirectRelField.GetName()); ok {
			relField.MaxSelect = 1
		}
		relCollection = indirectRel
	} else {
		// direct relation
		relField, _ = mainCollection.Fields.GetByName(parts[0]).(*RelationField)
		if relField == nil {
			return fmt.Errorf("couldn't find relation field %q in collection %q", parts[0], mainCollection.Name)
		}

		relCollection, _ = getCollectionByModelOrIdentifier(app, relField.CollectionId)
		if relCollection == nil {
			return fmt.Errorf("couldn't find related collection %q", relField.CollectionId)
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
		err := app.expandRecords(rels, parts[1], fetchFunc, recursionLevel+1)
		if err != nil {
			return err
		}
	}

	// reindex with the rel id
	indexedRels := make(map[string]*Record, len(rels))
	for _, rel := range rels {
		indexedRels[rel.Id] = rel
	}

	for _, model := range records {
		// init expand if not already
		// (this is done to ensure that the "expand" key will be returned in the response even if empty)
		if model.expand == nil {
			model.SetExpand(nil)
		}

		relIds := model.GetStringSlice(relField.Name)

		validRels := make([]*Record, 0, len(relIds))
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
		var oldExpandedRels []*Record
		switch v := expandData[relField.Name].(type) {
		case nil:
			// no old expands
		case *Record:
			oldExpandedRels = []*Record{v}
		case []*Record:
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
		if relField.IsMultiple() {
			expandData[relField.Name] = validRels
		} else {
			expandData[relField.Name] = validRels[0]
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
