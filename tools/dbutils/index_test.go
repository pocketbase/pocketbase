package dbutils_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/pocketbase/pocketbase/tools/dbutils"
)

func TestParseIndex(t *testing.T) {
	scenarios := []struct {
		index    string
		expected dbutils.Index
	}{
		// invalid
		{
			`invalid`,
			dbutils.Index{},
		},
		// simple
		{
			`create index indexname on tablename (col1)`,
			dbutils.Index{
				IndexName: "indexname",
				TableName: "tablename",
				Columns: []dbutils.IndexColumn{
					{Name: "col1"},
				},
			},
		},
		// all fields
		{
			`CREATE UNIQUE INDEX IF NOT EXISTS "schemaname".[indexname] on 'tablename' (
				col0,
				` + "`" + `col1` + "`" + `,
				json_extract("col2", "$.a") asc,
				"col3" collate NOCASE,
				"col4" collate RTRIM desc
			) where test = 1`,
			dbutils.Index{
				Unique:     true,
				Optional:   true,
				SchemaName: "schemaname",
				IndexName:  "indexname",
				TableName:  "tablename",
				Columns: []dbutils.IndexColumn{
					{Name: "col0"},
					{Name: "col1"},
					{Name: `json_extract("col2", "$.a")`, Sort: "ASC"},
					{Name: `col3`, Collate: "NOCASE"},
					{Name: `col4`, Collate: "RTRIM", Sort: "DESC"},
				},
				Where: "test = 1",
			},
		},
	}

	for i, s := range scenarios {
		result := dbutils.ParseIndex(s.index)

		resultRaw, err := json.Marshal(result)
		if err != nil {
			t.Fatalf("[%d] %v", i, err)
		}

		expectedRaw, err := json.Marshal(s.expected)
		if err != nil {
			t.Fatalf("[%d] %v", i, err)
		}

		if !bytes.Equal(resultRaw, expectedRaw) {
			t.Errorf("[%d] Expected \n%s \ngot \n%s", i, expectedRaw, resultRaw)
		}
	}
}

func TestIndexIsValid(t *testing.T) {
	scenarios := []struct {
		name     string
		index    dbutils.Index
		expected bool
	}{
		{
			"empty",
			dbutils.Index{},
			false,
		},
		{
			"no index name",
			dbutils.Index{
				TableName: "table",
				Columns:   []dbutils.IndexColumn{{Name: "col"}},
			},
			false,
		},
		{
			"no table name",
			dbutils.Index{
				IndexName: "index",
				Columns:   []dbutils.IndexColumn{{Name: "col"}},
			},
			false,
		},
		{
			"no columns",
			dbutils.Index{
				IndexName: "index",
				TableName: "table",
			},
			false,
		},
		{
			"min valid",
			dbutils.Index{
				IndexName: "index",
				TableName: "table",
				Columns:   []dbutils.IndexColumn{{Name: "col"}},
			},
			true,
		},
		{
			"all fields",
			dbutils.Index{
				Optional:   true,
				Unique:     true,
				SchemaName: "schema",
				IndexName:  "index",
				TableName:  "table",
				Columns:    []dbutils.IndexColumn{{Name: "col"}},
				Where:      "test = 1 OR test = 2",
			},
			true,
		},
	}

	for _, s := range scenarios {
		result := s.index.IsValid()

		if result != s.expected {
			t.Errorf("[%s] Expected %v, got %v", s.name, s.expected, result)
		}
	}
}

func TestIndexBuild(t *testing.T) {
	scenarios := []struct {
		name     string
		index    dbutils.Index
		expected string
	}{
		{
			"empty",
			dbutils.Index{},
			"",
		},
		{
			"no index name",
			dbutils.Index{
				TableName: "table",
				Columns:   []dbutils.IndexColumn{{Name: "col"}},
			},
			"",
		},
		{
			"no table name",
			dbutils.Index{
				IndexName: "index",
				Columns:   []dbutils.IndexColumn{{Name: "col"}},
			},
			"",
		},
		{
			"no columns",
			dbutils.Index{
				IndexName: "index",
				TableName: "table",
			},
			"",
		},
		{
			"min valid",
			dbutils.Index{
				IndexName: "index",
				TableName: "table",
				Columns:   []dbutils.IndexColumn{{Name: "col"}},
			},
			"CREATE INDEX `index` ON `table` (`col`)",
		},
		{
			"all fields",
			dbutils.Index{
				Optional:   true,
				Unique:     true,
				SchemaName: "schema",
				IndexName:  "index",
				TableName:  "table",
				Columns: []dbutils.IndexColumn{
					{Name: "col1", Collate: "NOCASE", Sort: "asc"},
					{Name: "col2", Sort: "desc"},
					{Name: `json_extract("col3", "$.a")`, Collate: "NOCASE"},
				},
				Where: "test = 1 OR test = 2",
			},
			"CREATE UNIQUE INDEX IF NOT EXISTS `schema`.`index` ON `table` (\n  `col1` COLLATE NOCASE ASC,\n  `col2` DESC,\n  " + `json_extract("col3", "$.a")` + " COLLATE NOCASE\n) WHERE test = 1 OR test = 2",
		},
	}

	for _, s := range scenarios {
		result := s.index.Build()

		if result != s.expected {
			t.Errorf("[%s] Expected \n%v \ngot \n%v", s.name, s.expected, result)
		}
	}
}
