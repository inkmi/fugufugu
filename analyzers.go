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
		dropAnalyzer,
		alterAnalyzer,
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

func alterAnalyzer(config *Config, stmt *pg.Node) *Warning {
	if alter := stmt.GetAlterTableStmt(); alter != nil {
		relation := alter.Relation.Relname

		for _, cmd := range alter.Cmds {
			if at := cmd.GetAlterTableCmd(); at != nil {
				switch at.Subtype {
				case pg.AlterTableType_AT_DropColumn:
					return &Warning{
						Change: "D",
						Type:   "drop",
						Object: "column",
						Name:   relation + "." + at.Name,
					}
				}
			}
		}
		return nil
	}
	return nil
}

func dropAnalyzer(config *Config, stmt *pg.Node) *Warning {
	if drop := stmt.GetDropStmt(); drop != nil {
		name := drop.GetObjects()[0].GetList().Items[0].GetString_().Sval
		switch drop.RemoveType {
		case pg.ObjectType_OBJECT_TABLE:
			return &Warning{
				Change: "D",
				Type:   "drop",
				Object: "table",
				Name:   name,
			}
		case pg.ObjectType_OBJECT_VIEW:
			return &Warning{
				Change: "IC",
				Type:   "drop",
				Object: "view",
				Name:   name,
			}
		default:
			return nil
		}
	}
	return nil
}
