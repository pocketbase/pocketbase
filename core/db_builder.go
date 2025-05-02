package core

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/pocketbase/dbx"
)

var _ dbx.Builder = (*dualDBBuilder)(nil)

// note: expects both builder to use the same driver
type dualDBBuilder struct {
	concurrentDB    dbx.Builder
	nonconcurrentDB dbx.Builder
}

// Select implements the [dbx.Builder.Select] interface method.
func (b *dualDBBuilder) Select(cols ...string) *dbx.SelectQuery {
	return b.concurrentDB.Select(cols...)
}

// Model implements the [dbx.Builder.Model] interface method.
func (b *dualDBBuilder) Model(data interface{}) *dbx.ModelQuery {
	return b.nonconcurrentDB.Model(data)
}

// GeneratePlaceholder implements the [dbx.Builder.GeneratePlaceholder] interface method.
func (b *dualDBBuilder) GeneratePlaceholder(i int) string {
	return b.concurrentDB.GeneratePlaceholder(i)
}

// Quote implements the [dbx.Builder.Quote] interface method.
func (b *dualDBBuilder) Quote(str string) string {
	return b.concurrentDB.Quote(str)
}

// QuoteSimpleTableName implements the [dbx.Builder.QuoteSimpleTableName] interface method.
func (b *dualDBBuilder) QuoteSimpleTableName(table string) string {
	return b.concurrentDB.QuoteSimpleTableName(table)
}

// QuoteSimpleColumnName implements the [dbx.Builder.QuoteSimpleColumnName] interface method.
func (b *dualDBBuilder) QuoteSimpleColumnName(col string) string {
	return b.concurrentDB.QuoteSimpleColumnName(col)
}

// QueryBuilder implements the [dbx.Builder.QueryBuilder] interface method.
func (b *dualDBBuilder) QueryBuilder() dbx.QueryBuilder {
	return b.concurrentDB.QueryBuilder()
}

// Insert implements the [dbx.Builder.Insert] interface method.
func (b *dualDBBuilder) Insert(table string, cols dbx.Params) *dbx.Query {
	return b.nonconcurrentDB.Insert(table, cols)
}

// Upsert implements the [dbx.Builder.Upsert] interface method.
func (b *dualDBBuilder) Upsert(table string, cols dbx.Params, constraints ...string) *dbx.Query {
	return b.nonconcurrentDB.Upsert(table, cols, constraints...)
}

// Update implements the [dbx.Builder.Update] interface method.
func (b *dualDBBuilder) Update(table string, cols dbx.Params, where dbx.Expression) *dbx.Query {
	return b.nonconcurrentDB.Update(table, cols, where)
}

// Delete implements the [dbx.Builder.Delete] interface method.
func (b *dualDBBuilder) Delete(table string, where dbx.Expression) *dbx.Query {
	return b.nonconcurrentDB.Delete(table, where)
}

// CreateTable implements the [dbx.Builder.CreateTable] interface method.
func (b *dualDBBuilder) CreateTable(table string, cols map[string]string, options ...string) *dbx.Query {
	return b.nonconcurrentDB.CreateTable(table, cols, options...)
}

// RenameTable implements the [dbx.Builder.RenameTable] interface method.
func (b *dualDBBuilder) RenameTable(oldName, newName string) *dbx.Query {
	return b.nonconcurrentDB.RenameTable(oldName, newName)
}

// DropTable implements the [dbx.Builder.DropTable] interface method.
func (b *dualDBBuilder) DropTable(table string) *dbx.Query {
	return b.nonconcurrentDB.DropTable(table)
}

// TruncateTable implements the [dbx.Builder.TruncateTable] interface method.
func (b *dualDBBuilder) TruncateTable(table string) *dbx.Query {
	return b.nonconcurrentDB.TruncateTable(table)
}

// AddColumn implements the [dbx.Builder.AddColumn] interface method.
func (b *dualDBBuilder) AddColumn(table, col, typ string) *dbx.Query {
	return b.nonconcurrentDB.AddColumn(table, col, typ)
}

// DropColumn implements the [dbx.Builder.DropColumn] interface method.
func (b *dualDBBuilder) DropColumn(table, col string) *dbx.Query {
	return b.nonconcurrentDB.DropColumn(table, col)
}

// RenameColumn implements the [dbx.Builder.RenameColumn] interface method.
func (b *dualDBBuilder) RenameColumn(table, oldName, newName string) *dbx.Query {
	return b.nonconcurrentDB.RenameColumn(table, oldName, newName)
}

// AlterColumn implements the [dbx.Builder.AlterColumn] interface method.
func (b *dualDBBuilder) AlterColumn(table, col, typ string) *dbx.Query {
	return b.nonconcurrentDB.AlterColumn(table, col, typ)
}

// AddPrimaryKey implements the [dbx.Builder.AddPrimaryKey] interface method.
func (b *dualDBBuilder) AddPrimaryKey(table, name string, cols ...string) *dbx.Query {
	return b.nonconcurrentDB.AddPrimaryKey(table, name, cols...)
}

// DropPrimaryKey implements the [dbx.Builder.DropPrimaryKey] interface method.
func (b *dualDBBuilder) DropPrimaryKey(table, name string) *dbx.Query {
	return b.nonconcurrentDB.DropPrimaryKey(table, name)
}

// AddForeignKey implements the [dbx.Builder.AddForeignKey] interface method.
func (b *dualDBBuilder) AddForeignKey(table, name string, cols, refCols []string, refTable string, options ...string) *dbx.Query {
	return b.nonconcurrentDB.AddForeignKey(table, name, cols, refCols, refTable, options...)
}

// DropForeignKey implements the [dbx.Builder.DropForeignKey] interface method.
func (b *dualDBBuilder) DropForeignKey(table, name string) *dbx.Query {
	return b.nonconcurrentDB.DropForeignKey(table, name)
}

// CreateIndex implements the [dbx.Builder.CreateIndex] interface method.
func (b *dualDBBuilder) CreateIndex(table, name string, cols ...string) *dbx.Query {
	return b.nonconcurrentDB.CreateIndex(table, name, cols...)
}

// CreateUniqueIndex implements the [dbx.Builder.CreateUniqueIndex] interface method.
func (b *dualDBBuilder) CreateUniqueIndex(table, name string, cols ...string) *dbx.Query {
	return b.nonconcurrentDB.CreateUniqueIndex(table, name, cols...)
}

// DropIndex implements the [dbx.Builder.DropIndex] interface method.
func (b *dualDBBuilder) DropIndex(table, name string) *dbx.Query {
	return b.nonconcurrentDB.DropIndex(table, name)
}

// NewQuery implements the [dbx.Builder.NewQuery] interface method by
// routing the SELECT queries to the concurrent builder instance.
func (b *dualDBBuilder) NewQuery(str string) *dbx.Query {
	// note: technically INSERT/UPDATE/DELETE could also have CTE but since
	// it is rare for now this scase is ignored to avoid unnecessary complicating the checks
	trimmed := trimLeftSpaces(str)
	if hasPrefixFold(trimmed, "SELECT") || hasPrefixFold(trimmed, "WITH") {
		return b.concurrentDB.NewQuery(str)
	}

	return b.nonconcurrentDB.NewQuery(str)
}

var asciiSpace = [256]uint8{'\t': 1, '\n': 1, '\v': 1, '\f': 1, '\r': 1, ' ': 1}

// note: similar to strings.Space() but without the right trim because it is not needed in our case
func trimLeftSpaces(str string) string {
	start := 0
	for ; start < len(str); start++ {
		c := str[start]
		if c >= utf8.RuneSelf {
			return strings.TrimLeftFunc(str[start:], unicode.IsSpace)
		}
		if asciiSpace[c] == 0 {
			break
		}
	}

	return str[start:]
}

// note: the prefix is expected to be ASCII
func hasPrefixFold(str, prefix string) bool {
	if len(str) < len(prefix) {
		return false
	}

	return strings.EqualFold(str[:len(prefix)], prefix)
}
