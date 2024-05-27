package main

import (
	"bufio"
	"flag"
	"fmt"
	pg "github.com/pganalyze/pg_query_go/v5"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type DropChange struct {
	Type string
	Name string
}

func main() {
	dir := flag.String("dir", "", "directory containing .sql files")
	flag.Parse()

	if *dir == "" {
		fmt.Println("Please specify a directory using --dir")
		return
	}

	files, err := ioutil.ReadDir(*dir)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".sql" {
			filePath := filepath.Join(*dir, file.Name())
			_, _, err := processSQLFile(filePath)
			if err != nil {
				fmt.Println("Error processing file:", file.Name(), err)
				continue
			}
			fmt.Printf("Analyzing: %s\n", file.Name())
		}
	}
}

func processSQLFile(filePath string) (string, string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	var upPart, downPart strings.Builder
	var currentPart *strings.Builder

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.HasPrefix(line, "-- +goose Up"):
			currentPart = &upPart
		case strings.HasPrefix(line, "-- +goose Down"):
			currentPart = &downPart
		default:
			if currentPart != nil {
				currentPart.WriteString(line)
				currentPart.WriteString("\n")
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return "", "", err
	}

	upStr := strings.TrimSpace(upPart.String())
	downStr := strings.TrimSpace(downPart.String())

	changes := analyze(upStr)
	if len(changes) > 0 {
		fmt.Println("Destructive changes:")
	}
	for _, c := range changes {
		fmt.Printf(" drop %s\t%s\n", c.Name, c.Type)
	}
	return upStr, downStr, nil
}

func analyzeStmt(stmt *pg.Node) *DropChange {
	if stmt.GetDropStmt() != nil {
		name := stmt.GetDropStmt().GetObjects()[0].GetList().Items[0].GetString_().Sval
		return &DropChange{
			Type: "table",
			Name: name,
		}
	}
	return nil
}

func analyze(sql string) []*DropChange {
	result, err := pg.Parse(sql)
	if err != nil {
		panic(err)
	}

	var changes []*DropChange
	res := analyzeStmt(result.Stmts[0].Stmt)
	if res != nil {
		changes = append(changes, res)

	}
	return changes
}
