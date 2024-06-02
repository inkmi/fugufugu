package main

import (
	"github.com/davecgh/go-spew/spew"
	pg "github.com/pganalyze/pg_query_go/v5"
)

func main() {
	sql := "ALTER TABLE at DROP COLUMN dc"
	result, err := pg.Parse(sql)
	if err != nil {
		panic(err)
	}
	spew.Dump(result)
	spew.Dump(result.Stmts[0].Stmt.GetAlterTableStmt().Cmds[0].GetAlterTableCmd().Subtype.String())
}
