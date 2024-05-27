package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func checkFileExists(path string) error {
	cleanPath := filepath.Clean(path)
	components := filepath.SplitList(cleanPath)
	for i := range components {
		partialPath := filepath.Join(components[:i+1]...)
		info, err := os.Stat(partialPath)
		if os.IsNotExist(err) {
			if i == len(components)-1 {
				return fmt.Errorf("%s does not exist", filepath.Base(partialPath))
			}
			return fmt.Errorf("directory %s does not exist", filepath.Base(partialPath))
		}
		if err != nil {
			return fmt.Errorf("error checking %s: %v", partialPath, err)
		}
		if i != len(components)-1 && !info.IsDir() {
			return fmt.Errorf("%s is not a directory", filepath.Base(partialPath))
		}
	}

	return nil
}

func main() {
	dir := flag.String("dir", "", "directory containing .sql files")
	flag.Parse()

	if *dir == "" {
		fmt.Println("Please specify a directory using --dir")
		return
	}

	err := checkFileExists("config.yaml")
	var config *Config
	if err == nil {
		config = LoadConfig()
	} else {
		config = &Config{
			Checkers: []Checker{},
		}
	}

	files, err := os.ReadDir(*dir)
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
			fileChanges := processSQLFile(config, filePath)
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

func processSQLFile(config *Config, filePath string) []*Warning {
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

	changes := analyze(config, upStr)
	return changes
}
