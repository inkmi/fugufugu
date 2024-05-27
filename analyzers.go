package main

import (
	"fmt"

	pg "github.com/pganalyze/pg_query_go/v5"
)

type Warning struct {
	Change string
	Type   string
	Object string
	Name   string
}

type Analyzer func(config *Config, stmt *pg.Node) *Warning

func AllAnalyzers() []Analyzer {
	return []Analyzer{
		dropTableAnalyzer,
		nameTableAnalyzer,
	}

}

// Ideas:
// - non-concurrent index
// - in-compatible changes:
//   - table name change
//   - column name change
//   - column type change
// - enforce name schema (table name, column name, ....)
// - other downtime changes

func analyze(config *Config, sql string) []*Warning {
	result, err := pg.Parse(sql)
	if err != nil {
		panic(err)
	}

	analyzers := AllAnalyzers()
	var changes []*Warning
	for _, analyzer := range analyzers {
		res := analyzer(config, result.Stmts[0].Stmt)
		if res != nil {
			changes = append(changes, res)

		}
	}
	return changes
}

func nameTableAnalyzer(config *Config, stmt *pg.Node) *Warning {
	if stmt.GetCreateStmt() != nil {
		name := stmt.GetCreateStmt().GetRelation().Relname
		valid := ValidateName(name, config.Checkers...)
		if !valid {
			fmt.Printf("%s is not a valid name\n", name)
		}
	}
	return nil
}

func dropTableAnalyzer(config *Config, stmt *pg.Node) *Warning {
	if stmt.GetDropStmt() != nil {
		name := stmt.GetDropStmt().GetObjects()[0].GetList().Items[0].GetString_().Sval
		return &Warning{
			Change: "D",
			Type:   "drop",
			Object: "table",
			Name:   name,
		}
	}
	return nil
}
