package main

import (
	"fmt"
	pg "github.com/pganalyze/pg_query_go/v5"
)

func main() {
	result, err := pg.Parse("DROP TABLE x")
	if err != nil {
		panic(err)
	}

	if result.Stmts[0].Stmt.GetDropStmt() != nil {
		fmt.Printf("DANGER: Data loss, DROP TABLE: %s\n", result.Stmts[0].Stmt.GetDropStmt().GetObjects()[0].GetList().Items[0].GetString_().Sval)
	}
}
