package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

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

	var changes []*Warning
	fmt.Println("Analyzing:")
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".sql" {
			filePath := filepath.Join(*dir, file.Name())
			fmt.Printf("  %s\n", file.Name())
			fileChanges := processSQLFile(filePath)
			if len(fileChanges) > 0 {
				changes = append(changes, fileChanges...)
			}
		}

	}
	if len(changes) > 0 {
		fmt.Println("Destructive changes:")
	}
	for _, c := range changes {
		fmt.Printf(" drop %s\t%s\n", c.Name, c.Object)
	}
}

func processSQLFile(filePath string) []*Warning {
	file, err := os.Open(filePath)
	if err != nil {
		return nil
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
		return nil
	}

	upStr := strings.TrimSpace(upPart.String())

	changes := analyze(upStr)
	return changes
}
