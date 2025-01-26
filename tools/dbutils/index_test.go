package dbutils_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
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
		// simple (multiple spaces between the table and columns list)
		{
			`create index indexname on tablename   (col1)`,
			dbutils.Index{
				IndexName: "indexname",
				TableName: "tablename",
				Columns: []dbutils.IndexColumn{
					{Name: "col1"},
				},
			},
		},
		// simple (no space between the table and the columns list)
		{
			`create index indexname on tablename(col1)`,
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
		t.Run(fmt.Sprintf("scenario_%d", i), func(t *testing.T) {
			result := dbutils.ParseIndex(s.index)

			resultRaw, err := json.Marshal(result)
			if err != nil {
				t.Fatalf("Faild to marshalize parse result: %v", err)
			}

			expectedRaw, err := json.Marshal(s.expected)
			if err != nil {
				t.Fatalf("Failed to marshalize expected index: %v", err)
			}

			if !bytes.Equal(resultRaw, expectedRaw) {
				t.Errorf("Expected \n%s \ngot \n%s", expectedRaw, resultRaw)
			}
		})
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
		t.Run(s.name, func(t *testing.T) {
			result := s.index.IsValid()
			if result != s.expected {
				t.Fatalf("Expected %v, got %v", s.expected, result)
			}
		})
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
		t.Run(s.name, func(t *testing.T) {
			result := s.index.Build()
			if result != s.expected {
				t.Fatalf("Expected \n%v \ngot \n%v", s.expected, result)
			}
		})
	}
}

func TestHasSingleColumnUniqueIndex(t *testing.T) {
	scenarios := []struct {
		name     string
		column   string
		indexes  []string
		expected bool
	}{
		{
			"empty indexes",
			"test",
			nil,
			false,
		},
		{
			"empty column",
			"",
			[]string{
				"CREATE UNIQUE INDEX `index1` ON `example` (`test`)",
			},
			false,
		},
		{
			"mismatched column",
			"test",
			[]string{
				"CREATE UNIQUE INDEX `index1` ON `example` (`test2`)",
			},
			false,
		},
		{
			"non unique index",
			"test",
			[]string{
				"CREATE INDEX `index1` ON `example` (`test`)",
			},
			false,
		},
		{
			"matching columnd and unique index",
			"test",
			[]string{
				"CREATE UNIQUE INDEX `index1` ON `example` (`test`)",
			},
			true,
		},
		{
			"multiple columns",
			"test",
			[]string{
				"CREATE UNIQUE INDEX `index1` ON `example` (`test`, `test2`)",
			},
			false,
		},
		{
			"multiple indexes",
			"test",
			[]string{
				"CREATE UNIQUE INDEX `index1` ON `example` (`test`, `test2`)",
				"CREATE UNIQUE INDEX `index2` ON `example` (`test`)",
			},
			true,
		},
		{
			"partial unique index",
			"test",
			[]string{
				"CREATE UNIQUE INDEX `index` ON `example` (`test`) where test != ''",
			},
			true,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			result := dbutils.HasSingleColumnUniqueIndex(s.column, s.indexes)
			if result != s.expected {
				t.Fatalf("Expected %v got %v", s.expected, result)
			}
		})
	}
}

func TestFindSingleColumnUniqueIndex(t *testing.T) {
	scenarios := []struct {
		name     string
		column   string
		indexes  []string
		expected bool
	}{
		{
			"empty indexes",
			"test",
			nil,
			false,
		},
		{
			"empty column",
			"",
			[]string{
				"CREATE UNIQUE INDEX `index1` ON `example` (`test`)",
			},
			false,
		},
		{
			"mismatched column",
			"test",
			[]string{
				"CREATE UNIQUE INDEX `index1` ON `example` (`test2`)",
			},
			false,
		},
		{
			"non unique index",
			"test",
			[]string{
				"CREATE INDEX `index1` ON `example` (`test`)",
			},
			false,
		},
		{
			"matching columnd and unique index",
			"test",
			[]string{
				"CREATE UNIQUE INDEX `index1` ON `example` (`test`)",
			},
			true,
		},
		{
			"multiple columns",
			"test",
			[]string{
				"CREATE UNIQUE INDEX `index1` ON `example` (`test`, `test2`)",
			},
			false,
		},
		{
			"multiple indexes",
			"test",
			[]string{
				"CREATE UNIQUE INDEX `index1` ON `example` (`test`, `test2`)",
				"CREATE UNIQUE INDEX `index2` ON `example` (`test`)",
			},
			true,
		},
		{
			"partial unique index",
			"test",
			[]string{
				"CREATE UNIQUE INDEX `index` ON `example` (`test`) where test != ''",
			},
			true,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			index, exists := dbutils.FindSingleColumnUniqueIndex(s.indexes, s.column)
			if exists != s.expected {
				t.Fatalf("Expected exists %v got %v", s.expected, exists)
			}

			if !exists && len(index.Columns) > 0 {
				t.Fatal("Expected index.Columns to be empty")
			}

			if exists && !strings.EqualFold(index.Columns[0].Name, s.column) {
				t.Fatalf("Expected to find column %q in %v", s.column, index)
			}
		})
	}
}
