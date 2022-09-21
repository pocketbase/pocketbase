package daos_test

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/list"
)

func TestExpandRecords(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	col, _ := app.Dao().FindCollectionByNameOrId("demo4")

	scenarios := []struct {
		recordIds            []string
		expands              []string
		fetchFunc            daos.ExpandFetchFunc
		expectExpandProps    int
		expectExpandFailures int
	}{
		// empty records
		{
			[]string{},
			[]string{"onerel", "manyrels.onerel.manyrels"},
			func(c *models.Collection, ids []string) ([]*models.Record, error) {
				return app.Dao().FindRecordsByIds(c, ids, nil)
			},
			0,
			0,
		},
		// empty expand
		{
			[]string{"b8ba58f9-e2d7-42a0-b0e7-a11efd98236b", "df55c8ff-45ef-4c82-8aed-6e2183fe1125"},
			[]string{},
			func(c *models.Collection, ids []string) ([]*models.Record, error) {
				return app.Dao().FindRecordsByIds(c, ids, nil)
			},
			0,
			0,
		},
		// empty fetchFunc
		{
			[]string{"b8ba58f9-e2d7-42a0-b0e7-a11efd98236b", "df55c8ff-45ef-4c82-8aed-6e2183fe1125"},
			[]string{"onerel", "manyrels.onerel.manyrels"},
			nil,
			0,
			2,
		},
		// fetchFunc with error
		{
			[]string{"b8ba58f9-e2d7-42a0-b0e7-a11efd98236b", "df55c8ff-45ef-4c82-8aed-6e2183fe1125"},
			[]string{"onerel", "manyrels.onerel.manyrels"},
			func(c *models.Collection, ids []string) ([]*models.Record, error) {
				return nil, errors.New("test error")
			},
			0,
			2,
		},
		// missing relation field
		{
			[]string{"b8ba58f9-e2d7-42a0-b0e7-a11efd98236b", "df55c8ff-45ef-4c82-8aed-6e2183fe1125"},
			[]string{"invalid"},
			func(c *models.Collection, ids []string) ([]*models.Record, error) {
				return app.Dao().FindRecordsByIds(c, ids, nil)
			},
			0,
			1,
		},
		// existing, but non-relation type field
		{
			[]string{"b8ba58f9-e2d7-42a0-b0e7-a11efd98236b", "df55c8ff-45ef-4c82-8aed-6e2183fe1125"},
			[]string{"title"},
			func(c *models.Collection, ids []string) ([]*models.Record, error) {
				return app.Dao().FindRecordsByIds(c, ids, nil)
			},
			0,
			1,
		},
		// invalid/missing second level expand
		{
			[]string{"b8ba58f9-e2d7-42a0-b0e7-a11efd98236b", "df55c8ff-45ef-4c82-8aed-6e2183fe1125"},
			[]string{"manyrels.invalid"},
			func(c *models.Collection, ids []string) ([]*models.Record, error) {
				return app.Dao().FindRecordsByIds(c, ids, nil)
			},
			0,
			1,
		},
		// expand normalizations
		{
			[]string{
				"b8ba58f9-e2d7-42a0-b0e7-a11efd98236b",
				"df55c8ff-45ef-4c82-8aed-6e2183fe1125",
				"b84cd893-7119-43c9-8505-3c4e22da28a9",
				"054f9f24-0a0a-4e09-87b1-bc7ff2b336a2",
			},
			[]string{"manyrels.onerel.manyrels.onerel", "manyrels.onerel", "onerel", "onerel.", " onerel ", ""},
			func(c *models.Collection, ids []string) ([]*models.Record, error) {
				return app.Dao().FindRecordsByIds(c, ids, nil)
			},
			9,
			0,
		},
		// expand multiple relations sharing a common root path
		{
			[]string{
				"i15r5aa28ad06c8",
			},
			[]string{"manyrels.onerel.manyrels.onerel", "manyrels.onerel.onerel"},
			func(c *models.Collection, ids []string) ([]*models.Record, error) {
				return app.Dao().FindRecordsByIds(c, ids, nil)
			},
			4,
			0,
		},
		// single expand
		{
			[]string{
				"b8ba58f9-e2d7-42a0-b0e7-a11efd98236b",
				"df55c8ff-45ef-4c82-8aed-6e2183fe1125",
				"b84cd893-7119-43c9-8505-3c4e22da28a9", // no manyrels
				"054f9f24-0a0a-4e09-87b1-bc7ff2b336a2", // no manyrels
			},
			[]string{"manyrels"},
			func(c *models.Collection, ids []string) ([]*models.Record, error) {
				return app.Dao().FindRecordsByIds(c, ids, nil)
			},
			2,
			0,
		},
		// maxExpandDepth reached
		{
			[]string{"b8ba58f9-e2d7-42a0-b0e7-a11efd98236b"},
			[]string{"manyrels.onerel.manyrels.onerel.manyrels.onerel.manyrels.onerel.manyrels"},
			func(c *models.Collection, ids []string) ([]*models.Record, error) {
				return app.Dao().FindRecordsByIds(c, ids, nil)
			},
			6,
			0,
		},
	}

	for i, s := range scenarios {
		ids := list.ToUniqueStringSlice(s.recordIds)
		records, _ := app.Dao().FindRecordsByIds(col, ids, nil)
		failed := app.Dao().ExpandRecords(records, s.expands, s.fetchFunc)

		if len(failed) != s.expectExpandFailures {
			t.Errorf("(%d) Expected %d failures, got %d: \n%v", i, s.expectExpandFailures, len(failed), failed)
		}

		encoded, _ := json.Marshal(records)
		encodedStr := string(encoded)
		totalExpandProps := strings.Count(encodedStr, "@expand")

		if s.expectExpandProps != totalExpandProps {
			t.Errorf("(%d) Expected %d @expand props, got %d: \n%v", i, s.expectExpandProps, totalExpandProps, encodedStr)
		}
	}
}

func TestExpandRecord(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	col, _ := app.Dao().FindCollectionByNameOrId("demo4")

	scenarios := []struct {
		recordId             string
		expands              []string
		fetchFunc            daos.ExpandFetchFunc
		expectExpandProps    int
		expectExpandFailures int
	}{
		// empty expand
		{
			"b8ba58f9-e2d7-42a0-b0e7-a11efd98236b",
			[]string{},
			func(c *models.Collection, ids []string) ([]*models.Record, error) {
				return app.Dao().FindRecordsByIds(c, ids, nil)
			},
			0,
			0,
		},
		// empty fetchFunc
		{
			"b8ba58f9-e2d7-42a0-b0e7-a11efd98236b",
			[]string{"onerel", "manyrels.onerel.manyrels"},
			nil,
			0,
			2,
		},
		// fetchFunc with error
		{
			"b8ba58f9-e2d7-42a0-b0e7-a11efd98236b",
			[]string{"onerel", "manyrels.onerel.manyrels"},
			func(c *models.Collection, ids []string) ([]*models.Record, error) {
				return nil, errors.New("test error")
			},
			0,
			2,
		},
		// invalid missing first level expand
		{
			"b8ba58f9-e2d7-42a0-b0e7-a11efd98236b",
			[]string{"invalid"},
			func(c *models.Collection, ids []string) ([]*models.Record, error) {
				return app.Dao().FindRecordsByIds(c, ids, nil)
			},
			0,
			1,
		},
		// invalid missing second level expand
		{
			"b8ba58f9-e2d7-42a0-b0e7-a11efd98236b",
			[]string{"manyrels.invalid"},
			func(c *models.Collection, ids []string) ([]*models.Record, error) {
				return app.Dao().FindRecordsByIds(c, ids, nil)
			},
			0,
			1,
		},
		// expand normalizations
		{
			"b8ba58f9-e2d7-42a0-b0e7-a11efd98236b",
			[]string{"manyrels.onerel.manyrels", "manyrels.onerel", "onerel", " onerel "},
			func(c *models.Collection, ids []string) ([]*models.Record, error) {
				return app.Dao().FindRecordsByIds(c, ids, nil)
			},
			3,
			0,
		},
		// single expand
		{
			"b8ba58f9-e2d7-42a0-b0e7-a11efd98236b",
			[]string{"manyrels"},
			func(c *models.Collection, ids []string) ([]*models.Record, error) {
				return app.Dao().FindRecordsByIds(c, ids, nil)
			},
			1,
			0,
		},
		// maxExpandDepth reached
		{
			"b8ba58f9-e2d7-42a0-b0e7-a11efd98236b",
			[]string{"manyrels.onerel.manyrels.onerel.manyrels.onerel.manyrels.onerel.manyrels"},
			func(c *models.Collection, ids []string) ([]*models.Record, error) {
				return app.Dao().FindRecordsByIds(c, ids, nil)
			},
			6,
			0,
		},
	}

	for i, s := range scenarios {
		record, _ := app.Dao().FindFirstRecordByData(col, "id", s.recordId)
		failed := app.Dao().ExpandRecord(record, s.expands, s.fetchFunc)

		if len(failed) != s.expectExpandFailures {
			t.Errorf("(%d) Expected %d failures, got %d: \n%v", i, s.expectExpandFailures, len(failed), failed)
		}

		encoded, _ := json.Marshal(record)
		encodedStr := string(encoded)
		totalExpandProps := strings.Count(encodedStr, "@expand")

		if s.expectExpandProps != totalExpandProps {
			t.Errorf("(%d) Expected %d @expand props, got %d: \n%v", i, s.expectExpandProps, totalExpandProps, encodedStr)
		}
	}
}
