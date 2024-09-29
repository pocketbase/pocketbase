package core_test

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/list"
)

func TestExpandRecords(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		testName                  string
		collectionIdOrName        string
		recordIds                 []string
		expands                   []string
		fetchFunc                 core.ExpandFetchFunc
		expectNonemptyExpandProps int
		expectExpandFailures      int
	}{
		{
			"empty records",
			"",
			[]string{},
			[]string{"self_rel_one", "self_rel_many.self_rel_one"},
			func(c *core.Collection, ids []string) ([]*core.Record, error) {
				return app.FindRecordsByIds(c.Id, ids, nil)
			},
			0,
			0,
		},
		{
			"empty expand",
			"demo4",
			[]string{"i9naidtvr6qsgb4", "qzaqccwrmva4o1n"},
			[]string{},
			func(c *core.Collection, ids []string) ([]*core.Record, error) {
				return app.FindRecordsByIds(c.Id, ids, nil)
			},
			0,
			0,
		},
		{
			"fetchFunc with error",
			"demo4",
			[]string{"i9naidtvr6qsgb4", "qzaqccwrmva4o1n"},
			[]string{"self_rel_one", "self_rel_many.self_rel_one"},
			func(c *core.Collection, ids []string) ([]*core.Record, error) {
				return nil, errors.New("test error")
			},
			0,
			2,
		},
		{
			"missing relation field",
			"demo4",
			[]string{"i9naidtvr6qsgb4", "qzaqccwrmva4o1n"},
			[]string{"missing"},
			func(c *core.Collection, ids []string) ([]*core.Record, error) {
				return app.FindRecordsByIds(c.Id, ids, nil)
			},
			0,
			1,
		},
		{
			"existing, but non-relation type field",
			"demo4",
			[]string{"i9naidtvr6qsgb4", "qzaqccwrmva4o1n"},
			[]string{"title"},
			func(c *core.Collection, ids []string) ([]*core.Record, error) {
				return app.FindRecordsByIds(c.Id, ids, nil)
			},
			0,
			1,
		},
		{
			"invalid/missing second level expand",
			"demo4",
			[]string{"i9naidtvr6qsgb4", "qzaqccwrmva4o1n"},
			[]string{"rel_one_no_cascade.title"},
			func(c *core.Collection, ids []string) ([]*core.Record, error) {
				return app.FindRecordsByIds(c.Id, ids, nil)
			},
			0,
			1,
		},
		{
			"with nil fetchfunc",
			"users",
			[]string{
				"bgs820n361vj1qd",
				"4q1xlclmfloku33",
				"oap640cot4yru2s", // no rels
			},
			[]string{"rel"},
			nil,
			2,
			0,
		},
		{
			"expand normalizations",
			"demo4",
			[]string{"i9naidtvr6qsgb4", "qzaqccwrmva4o1n"},
			[]string{
				"self_rel_one", "self_rel_many.self_rel_many.rel_one_no_cascade",
				"self_rel_many.self_rel_one.self_rel_many.self_rel_one.rel_one_no_cascade",
				"self_rel_many", "self_rel_many.",
				"  self_rel_many  ", "",
			},
			func(c *core.Collection, ids []string) ([]*core.Record, error) {
				return app.FindRecordsByIds(c.Id, ids, nil)
			},
			9,
			0,
		},
		{
			"single expand",
			"users",
			[]string{
				"bgs820n361vj1qd",
				"4q1xlclmfloku33",
				"oap640cot4yru2s", // no rels
			},
			[]string{"rel"},
			func(c *core.Collection, ids []string) ([]*core.Record, error) {
				return app.FindRecordsByIds(c.Id, ids, nil)
			},
			2,
			0,
		},
		{
			"with nil fetchfunc",
			"users",
			[]string{
				"bgs820n361vj1qd",
				"4q1xlclmfloku33",
				"oap640cot4yru2s", // no rels
			},
			[]string{"rel"},
			nil,
			2,
			0,
		},
		{
			"maxExpandDepth reached",
			"demo4",
			[]string{"qzaqccwrmva4o1n"},
			[]string{"self_rel_many.self_rel_many.self_rel_many.self_rel_many.self_rel_many.self_rel_many.self_rel_many.self_rel_many.self_rel_many.self_rel_many.self_rel_many.self_rel_many"},
			func(c *core.Collection, ids []string) ([]*core.Record, error) {
				return app.FindRecordsByIds(c.Id, ids, nil)
			},
			6,
			0,
		},
		{
			"simple back single relation field expand (deprecated syntax)",
			"demo3",
			[]string{"lcl9d87w22ml6jy"},
			[]string{"demo4(rel_one_no_cascade_required)"},
			func(c *core.Collection, ids []string) ([]*core.Record, error) {
				return app.FindRecordsByIds(c.Id, ids, nil)
			},
			1,
			0,
		},
		{
			"simple back expand via single relation field",
			"demo3",
			[]string{"lcl9d87w22ml6jy"},
			[]string{"demo4_via_rel_one_no_cascade_required"},
			func(c *core.Collection, ids []string) ([]*core.Record, error) {
				return app.FindRecordsByIds(c.Id, ids, nil)
			},
			1,
			0,
		},
		{
			"nested back expand via single relation field",
			"demo3",
			[]string{"lcl9d87w22ml6jy"},
			[]string{
				"demo4_via_rel_one_no_cascade_required.self_rel_many.self_rel_many.self_rel_one",
			},
			func(c *core.Collection, ids []string) ([]*core.Record, error) {
				return app.FindRecordsByIds(c.Id, ids, nil)
			},
			5,
			0,
		},
		{
			"nested back expand via multiple relation field",
			"demo3",
			[]string{"lcl9d87w22ml6jy"},
			[]string{
				"demo4_via_rel_many_no_cascade_required.self_rel_many.rel_many_no_cascade_required.demo4_via_rel_many_no_cascade_required",
			},
			func(c *core.Collection, ids []string) ([]*core.Record, error) {
				return app.FindRecordsByIds(c.Id, ids, nil)
			},
			7,
			0,
		},
		{
			"expand multiple relations sharing a common path",
			"demo4",
			[]string{"qzaqccwrmva4o1n"},
			[]string{
				"rel_one_no_cascade",
				"rel_many_no_cascade",
				"self_rel_many.self_rel_one.rel_many_cascade",
				"self_rel_many.self_rel_one.rel_many_no_cascade_required",
			},
			func(c *core.Collection, ids []string) ([]*core.Record, error) {
				return app.FindRecordsByIds(c.Id, ids, nil)
			},
			5,
			0,
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			ids := list.ToUniqueStringSlice(s.recordIds)
			records, _ := app.FindRecordsByIds(s.collectionIdOrName, ids)
			failed := app.ExpandRecords(records, s.expands, s.fetchFunc)

			if len(failed) != s.expectExpandFailures {
				t.Errorf("Expected %d failures, got %d\n%v", s.expectExpandFailures, len(failed), failed)
			}

			encoded, _ := json.Marshal(records)
			encodedStr := string(encoded)
			totalExpandProps := strings.Count(encodedStr, `"`+core.FieldNameExpand+`":`)
			totalEmptyExpands := strings.Count(encodedStr, `"`+core.FieldNameExpand+`":{}`)
			totalNonemptyExpands := totalExpandProps - totalEmptyExpands

			if s.expectNonemptyExpandProps != totalNonemptyExpands {
				t.Errorf("Expected %d expand props, got %d\n%v", s.expectNonemptyExpandProps, totalNonemptyExpands, encodedStr)
			}
		})
	}
}

func TestExpandRecord(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		testName                  string
		collectionIdOrName        string
		recordId                  string
		expands                   []string
		fetchFunc                 core.ExpandFetchFunc
		expectNonemptyExpandProps int
		expectExpandFailures      int
	}{
		{
			"empty expand",
			"demo4",
			"i9naidtvr6qsgb4",
			[]string{},
			func(c *core.Collection, ids []string) ([]*core.Record, error) {
				return app.FindRecordsByIds(c.Id, ids, nil)
			},
			0,
			0,
		},
		{
			"fetchFunc with error",
			"demo4",
			"i9naidtvr6qsgb4",
			[]string{"self_rel_one", "self_rel_many.self_rel_one"},
			func(c *core.Collection, ids []string) ([]*core.Record, error) {
				return nil, errors.New("test error")
			},
			0,
			2,
		},
		{
			"missing relation field",
			"demo4",
			"i9naidtvr6qsgb4",
			[]string{"missing"},
			func(c *core.Collection, ids []string) ([]*core.Record, error) {
				return app.FindRecordsByIds(c.Id, ids, nil)
			},
			0,
			1,
		},
		{
			"existing, but non-relation type field",
			"demo4",
			"i9naidtvr6qsgb4",
			[]string{"title"},
			func(c *core.Collection, ids []string) ([]*core.Record, error) {
				return app.FindRecordsByIds(c.Id, ids, nil)
			},
			0,
			1,
		},
		{
			"invalid/missing second level expand",
			"demo4",
			"qzaqccwrmva4o1n",
			[]string{"rel_one_no_cascade.title"},
			func(c *core.Collection, ids []string) ([]*core.Record, error) {
				return app.FindRecordsByIds(c.Id, ids, nil)
			},
			0,
			1,
		},
		{
			"expand normalizations",
			"demo4",
			"qzaqccwrmva4o1n",
			[]string{
				"self_rel_one", "self_rel_many.self_rel_many.rel_one_no_cascade",
				"self_rel_many.self_rel_one.self_rel_many.self_rel_one.rel_one_no_cascade",
				"self_rel_many", "self_rel_many.",
				"  self_rel_many  ", "",
			},
			func(c *core.Collection, ids []string) ([]*core.Record, error) {
				return app.FindRecordsByIds(c.Id, ids, nil)
			},
			8,
			0,
		},
		{
			"no rels to expand",
			"users",
			"oap640cot4yru2s",
			[]string{"rel"},
			func(c *core.Collection, ids []string) ([]*core.Record, error) {
				return app.FindRecordsByIds(c.Id, ids, nil)
			},
			0,
			0,
		},
		{
			"maxExpandDepth reached",
			"demo4",
			"qzaqccwrmva4o1n",
			[]string{"self_rel_many.self_rel_many.self_rel_many.self_rel_many.self_rel_many.self_rel_many.self_rel_many.self_rel_many.self_rel_many.self_rel_many.self_rel_many.self_rel_many"},
			func(c *core.Collection, ids []string) ([]*core.Record, error) {
				return app.FindRecordsByIds(c.Id, ids, nil)
			},
			6,
			0,
		},
		{
			"simple indirect expand via single relation field (deprecated syntax)",
			"demo3",
			"lcl9d87w22ml6jy",
			[]string{"demo4(rel_one_no_cascade_required)"},
			func(c *core.Collection, ids []string) ([]*core.Record, error) {
				return app.FindRecordsByIds(c.Id, ids, nil)
			},
			1,
			0,
		},
		{
			"simple indirect expand via single relation field",
			"demo3",
			"lcl9d87w22ml6jy",
			[]string{"demo4_via_rel_one_no_cascade_required"},
			func(c *core.Collection, ids []string) ([]*core.Record, error) {
				return app.FindRecordsByIds(c.Id, ids, nil)
			},
			1,
			0,
		},
		{
			"nested indirect expand via single relation field",
			"demo3",
			"lcl9d87w22ml6jy",
			[]string{
				"demo4(rel_one_no_cascade_required).self_rel_many.self_rel_many.self_rel_one",
			},
			func(c *core.Collection, ids []string) ([]*core.Record, error) {
				return app.FindRecordsByIds(c.Id, ids, nil)
			},
			5,
			0,
		},
		{
			"nested indirect expand via single relation field",
			"demo3",
			"lcl9d87w22ml6jy",
			[]string{
				"demo4_via_rel_many_no_cascade_required.self_rel_many.rel_many_no_cascade_required.demo4_via_rel_many_no_cascade_required",
			},
			func(c *core.Collection, ids []string) ([]*core.Record, error) {
				return app.FindRecordsByIds(c.Id, ids, nil)
			},
			7,
			0,
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			record, _ := app.FindRecordById(s.collectionIdOrName, s.recordId)
			failed := app.ExpandRecord(record, s.expands, s.fetchFunc)

			if len(failed) != s.expectExpandFailures {
				t.Errorf("Expected %d failures, got %d\n%v", s.expectExpandFailures, len(failed), failed)
			}

			encoded, _ := json.Marshal(record)
			encodedStr := string(encoded)
			totalExpandProps := strings.Count(encodedStr, `"`+core.FieldNameExpand+`":`)
			totalEmptyExpands := strings.Count(encodedStr, `"`+core.FieldNameExpand+`":{}`)
			totalNonemptyExpands := totalExpandProps - totalEmptyExpands

			if s.expectNonemptyExpandProps != totalNonemptyExpands {
				t.Errorf("Expected %d expand props, got %d\n%v", s.expectNonemptyExpandProps, totalNonemptyExpands, encodedStr)
			}
		})
	}
}

func TestBackRelationExpandSingeVsArrayResult(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	record, err := app.FindRecordById("demo3", "7nwo8tuiatetxdm")
	if err != nil {
		t.Fatal(err)
	}

	// non-unique indirect expand
	{
		errs := app.ExpandRecord(record, []string{"demo4_via_rel_one_cascade"}, func(c *core.Collection, ids []string) ([]*core.Record, error) {
			return app.FindRecordsByIds(c.Id, ids, nil)
		})
		if len(errs) > 0 {
			t.Fatal(errs)
		}

		result, ok := record.Expand()["demo4_via_rel_one_cascade"].([]*core.Record)
		if !ok {
			t.Fatalf("Expected the expanded result to be a slice, got %v", result)
		}
	}

	// unique indirect expand
	{
		// mock a unique constraint for the rel_one_cascade field
		// ---
		demo4, err := app.FindCollectionByNameOrId("demo4")
		if err != nil {
			t.Fatal(err)
		}

		demo4.Indexes = append(demo4.Indexes, "create unique index idx_unique_expand on demo4 (rel_one_cascade)")

		if err := app.Save(demo4); err != nil {
			t.Fatalf("Failed to mock unique constraint: %v", err)
		}
		// ---

		errs := app.ExpandRecord(record, []string{"demo4_via_rel_one_cascade"}, func(c *core.Collection, ids []string) ([]*core.Record, error) {
			return app.FindRecordsByIds(c.Id, ids, nil)
		})
		if len(errs) > 0 {
			t.Fatal(errs)
		}

		result, ok := record.Expand()["demo4_via_rel_one_cascade"].(*core.Record)
		if !ok {
			t.Fatalf("Expected the expanded result to be a single model, got %v", result)
		}
	}
}
