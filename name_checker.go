package main

import (
	"log"
	"regexp"
	"strings"
	"unicode"

	"github.com/spf13/viper"
)

type Checker interface {
	Check(name string) bool
}

type (
	SnakeCaseChecker          struct{}
	UpperSnakeCaseChecker     struct{}
	CamelCaseChecker          struct{}
	LowerCaseChecker          struct{}
	UpperCaseChecker          struct{}
	KebabCaseChecker          struct{}
	PascalCaseChecker         struct{}
	AlphanumericChecker       struct{}
	NoSQLReservedWordsChecker struct{}
	MinimumLengthChecker      struct {
		MinLength int
	}
	MaximumLengthChecker struct {
		MaxLength int
	}
	AllowedSpecialCharactersChecker struct {
		AllowedChars string
	}
	RegexChecker struct {
		Pattern string
	}
	PrefixChecker struct {
		Prefix string
	}
)

func (c SnakeCaseChecker) Check(name string) bool {
	return checkRegex(name, `^[a-z]+(_[a-z]+)*$`)
}

func (c UpperSnakeCaseChecker) Check(name string) bool {
	return checkRegex(name, `^[A-Z]+(_[A-Z]+)*$`)
}

func (c CamelCaseChecker) Check(name string) bool {
	return checkRegex(name, `^[a-z]+([A-Z][a-z]*)*$`)
}

func (c LowerCaseChecker) Check(name string) bool {
	return name == strings.ToLower(name)
}

func (c UpperCaseChecker) Check(name string) bool {
	return name == strings.ToUpper(name)
}

func (c KebabCaseChecker) Check(name string) bool {
	return checkRegex(name, `^[a-z]+(-[a-z]+)*$`)
}

func (c PascalCaseChecker) Check(name string) bool {
	return checkRegex(name, `^[A-Z][a-z]+([A-Z][a-z]*)*$`)
}

func (c AlphanumericChecker) Check(name string) bool {
	return checkRegex(name, `^[a-zA-Z0-9]+$`)
}

func (c NoSQLReservedWordsChecker) Check(name string) bool {
	reservedWords := []string{
		"SELECT", "FROM", "WHERE", "INSERT", "UPDATE", "DELETE", "CREATE", "DROP", "ALTER",
		// Add more SQL reserved words as needed
	}
	nameUpper := strings.ToUpper(name)
	for _, word := range reservedWords {
		if nameUpper == word {
			return false
		}
	}
	return true
}

func (c MinimumLengthChecker) Check(name string) bool {
	return len(name) >= c.MinLength
}

func (c MaximumLengthChecker) Check(name string) bool {
	return len(name) <= c.MaxLength
}

func (c AllowedSpecialCharactersChecker) Check(name string) bool {
	for _, ch := range name {
		if !unicode.IsLetter(ch) && !unicode.IsNumber(ch) && !strings.ContainsRune(c.AllowedChars, ch) {
			return false
		}
	}
	return true
}

func (c RegexChecker) Check(name string) bool {
	return checkRegex(name, c.Pattern)
}

func (c PrefixChecker) Check(name string) bool {
	return strings.HasPrefix(name, c.Prefix)
}

func ValidateName(name string, checkers ...Checker) bool {
	for _, checker := range checkers {
		if !checker.Check(name) {
			return false
		}
	}
	return true
}

func checkRegex(name, pattern string) bool {
	matched, _ := regexp.MatchString(pattern, name)
	return matched
}

type Config struct {
	Checkers []Checker
}

func LoadConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	return &Config{
		Checkers: LoadCheckers(),
	}
}

func LoadCheckers() []Checker {
	var checkers []Checker

	if viper.IsSet("checkers") {
		for _, c := range viper.Get("checkers").([]interface{}) {
			checkerConfig := c.(map[string]interface{})
			switch checkerConfig["type"] {
			case "snake-case":
				checkers = append(checkers, SnakeCaseChecker{})
			case "upper-snake-case":
				checkers = append(checkers, UpperSnakeCaseChecker{})
			case "camel-case":
				checkers = append(checkers, CamelCaseChecker{})
			case "lower-case":
				checkers = append(checkers, LowerCaseChecker{})
			case "upper-case":
				checkers = append(checkers, UpperCaseChecker{})
			case "kebab-case":
				checkers = append(checkers, KebabCaseChecker{})
			case "pascal-case":
				checkers = append(checkers, PascalCaseChecker{})
			case "alphanumeric":
				checkers = append(checkers, AlphanumericChecker{})
			case "no-sql-reserved-words":
				checkers = append(checkers, NoSQLReservedWordsChecker{})
			case "minimum-length":
				if val, ok := checkerConfig["minlength"].(int); ok {
					checkers = append(checkers, MinimumLengthChecker{MinLength: val})
				} else {
					log.Fatalf("minimum-length checker requires a minlength field")
				}
			case "maximum-length":
				if val, ok := checkerConfig["maxlength"].(int); ok {
					checkers = append(checkers, MaximumLengthChecker{MaxLength: val})
				} else {
					log.Fatalf("maximum-length checker requires a maxlength field")
				}
			case "allowed-special-characters":
				if val, ok := checkerConfig["allowedchars"].(string); ok {
					checkers = append(checkers, AllowedSpecialCharactersChecker{AllowedChars: val})
				} else {
					log.Fatalf("allowed-special-characters checker requires an allowedchars field")
				}
			case "regex":
				if val, ok := checkerConfig["pattern"].(string); ok {
					checkers = append(checkers, RegexChecker{Pattern: val})
				} else {
					log.Fatalf("regex checker requires a pattern field")
				}
			case "prefix":
				if val, ok := checkerConfig["prefix"].(string); ok {
					checkers = append(checkers, PrefixChecker{Prefix: val})
				} else {
					log.Fatalf("prefix checker requires a prefix field")
				}
			}
		}
	}

	return checkers
}
