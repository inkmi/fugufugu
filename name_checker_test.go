package main

import (
	"testing"
)

func TestSnakeCaseChecker(t *testing.T) {
	checker := SnakeCaseChecker{}
	validNames := []string{"snake_case", "another_example"}
	invalidNames := []string{"Snake_Case", "snakeCase", "snake-case"}

	for _, name := range validNames {
		if !checker.Check(name) {
			t.Errorf("expected %q to be valid for SnakeCaseChecker", name)
		}
	}

	for _, name := range invalidNames {
		if checker.Check(name) {
			t.Errorf("expected %q to be invalid for SnakeCaseChecker", name)
		}
	}
}

func TestUpperSnakeCaseChecker(t *testing.T) {
	checker := UpperSnakeCaseChecker{}
	validNames := []string{"UPPER_SNAKE_CASE", "ANOTHER_EXAMPLE"}
	invalidNames := []string{"Upper_Snake_Case", "upper_snake_case", "UPPER-CASE"}

	for _, name := range validNames {
		if !checker.Check(name) {
			t.Errorf("expected %q to be valid for UpperSnakeCaseChecker", name)
		}
	}

	for _, name := range invalidNames {
		if checker.Check(name) {
			t.Errorf("expected %q to be invalid for UpperSnakeCaseChecker", name)
		}
	}
}

func TestCamelCaseChecker(t *testing.T) {
	checker := CamelCaseChecker{}
	validNames := []string{"camelCase", "anotherExample"}
	invalidNames := []string{"CamelCase", "camel_case", "camel-case"}

	for _, name := range validNames {
		if !checker.Check(name) {
			t.Errorf("expected %q to be valid for CamelCaseChecker", name)
		}
	}

	for _, name := range invalidNames {
		if checker.Check(name) {
			t.Errorf("expected %q to be invalid for CamelCaseChecker", name)
		}
	}
}

func TestLowerCaseChecker(t *testing.T) {
	checker := LowerCaseChecker{}
	validNames := []string{"lowercase", "anotherexample"}
	invalidNames := []string{"LowerCase"}

	for _, name := range validNames {
		if !checker.Check(name) {
			t.Errorf("expected %q to be valid for LowerCaseChecker", name)
		}
	}

	for _, name := range invalidNames {
		if checker.Check(name) {
			t.Errorf("expected %q to be invalid for LowerCaseChecker", name)
		}
	}
}

func TestUpperCaseChecker(t *testing.T) {
	checker := UpperCaseChecker{}
	validNames := []string{"UPPERCASE", "ANOTHEREXAMPLE"}
	invalidNames := []string{"UpperCase", "upper-case"}

	for _, name := range validNames {
		if !checker.Check(name) {
			t.Errorf("expected %q to be valid for UpperCaseChecker", name)
		}
	}

	for _, name := range invalidNames {
		if checker.Check(name) {
			t.Errorf("expected %q to be invalid for UpperCaseChecker", name)
		}
	}
}

func TestKebabCaseChecker(t *testing.T) {
	checker := KebabCaseChecker{}
	validNames := []string{"kebab-case", "another-example"}
	invalidNames := []string{"Kebab-Case", "kebab_case", "kebabCase"}

	for _, name := range validNames {
		if !checker.Check(name) {
			t.Errorf("expected %q to be valid for KebabCaseChecker", name)
		}
	}

	for _, name := range invalidNames {
		if checker.Check(name) {
			t.Errorf("expected %q to be invalid for KebabCaseChecker", name)
		}
	}
}

func TestPascalCaseChecker(t *testing.T) {
	checker := PascalCaseChecker{}
	validNames := []string{"PascalCase", "AnotherExample"}
	invalidNames := []string{"pascalCase", "Pascal_Case", "Pascal-Case"}

	for _, name := range validNames {
		if !checker.Check(name) {
			t.Errorf("expected %q to be valid for PascalCaseChecker", name)
		}
	}

	for _, name := range invalidNames {
		if checker.Check(name) {
			t.Errorf("expected %q to be invalid for PascalCaseChecker", name)
		}
	}
}

func TestAlphanumericChecker(t *testing.T) {
	checker := AlphanumericChecker{}
	validNames := []string{"Alphanumeric123", "AnotherExample456"}
	invalidNames := []string{"Alphanumeric_123", "Another-Example456", "Alphanumeric 123"}

	for _, name := range validNames {
		if !checker.Check(name) {
			t.Errorf("expected %q to be valid for AlphanumericChecker", name)
		}
	}

	for _, name := range invalidNames {
		if checker.Check(name) {
			t.Errorf("expected %q to be invalid for AlphanumericChecker", name)
		}
	}
}

func TestNoSQLReservedWordsChecker(t *testing.T) {
	checker := NoSQLReservedWordsChecker{}
	validNames := []string{"ValidName", "AnotherExample"}
	invalidNames := []string{"SELECT", "INSERT", "DELETE"}

	for _, name := range validNames {
		if !checker.Check(name) {
			t.Errorf("expected %q to be valid for NoSQLReservedWordsChecker", name)
		}
	}

	for _, name := range invalidNames {
		if checker.Check(name) {
			t.Errorf("expected %q to be invalid for NoSQLReservedWordsChecker", name)
		}
	}
}

func TestMinimumLengthChecker(t *testing.T) {
	checker := MinimumLengthChecker{MinLength: 3}
	validNames := []string{"abc", "abcd"}
	invalidNames := []string{"ab", "a"}

	for _, name := range validNames {
		if !checker.Check(name) {
			t.Errorf("expected %q to be valid for MinimumLengthChecker", name)
		}
	}

	for _, name := range invalidNames {
		if checker.Check(name) {
			t.Errorf("expected %q to be invalid for MinimumLengthChecker", name)
		}
	}
}

func TestMaximumLengthChecker(t *testing.T) {
	checker := MaximumLengthChecker{MaxLength: 5}
	validNames := []string{"abc", "abcd", "abcde"}
	invalidNames := []string{"abcdef", "abcdefg"}

	for _, name := range validNames {
		if !checker.Check(name) {
			t.Errorf("expected %q to be valid for MaximumLengthChecker", name)
		}
	}

	for _, name := range invalidNames {
		if checker.Check(name) {
			t.Errorf("expected %q to be invalid for MaximumLengthChecker", name)
		}
	}
}

func TestAllowedSpecialCharactersChecker(t *testing.T) {
	checker := AllowedSpecialCharactersChecker{AllowedChars: "_-"}
	validNames := []string{"valid_name", "valid-name"}
	invalidNames := []string{"invalid$name", "invalid@name"}

	for _, name := range validNames {
		if !checker.Check(name) {
			t.Errorf("expected %q to be valid for AllowedSpecialCharactersChecker", name)
		}
	}

	for _, name := range invalidNames {
		if checker.Check(name) {
			t.Errorf("expected %q to be invalid for AllowedSpecialCharactersChecker", name)
		}
	}
}

func TestRegexChecker(t *testing.T) {
	checker := RegexChecker{Pattern: `^[a-z]+$`}
	validNames := []string{"lowercase", "anotherexample"}
	invalidNames := []string{"LowerCase", "lower_case", "lower-case"}

	for _, name := range validNames {
		if !checker.Check(name) {
			t.Errorf("expected %q to be valid for RegexChecker", name)
		}
	}

	for _, name := range invalidNames {
		if checker.Check(name) {
			t.Errorf("expected %q to be invalid for RegexChecker", name)
		}
	}
}

func TestPrefixChecker(t *testing.T) {
	checker := PrefixChecker{Prefix: "ex"}
	validNames := []string{"example", "exemplary"}
	invalidNames := []string{"Example", "notexample"}

	for _, name := range validNames {
		if !checker.Check(name) {
			t.Errorf("expected %q to be valid for PrefixChecker", name)
		}
	}

	for _, name := range invalidNames {
		if checker.Check(name) {
			t.Errorf("expected %q to be invalid for PrefixChecker", name)
		}
	}
}

func TestValidateName(t *testing.T) {
	name := "exampleName"
	checkers := []Checker{
		CamelCaseChecker{},
		PrefixChecker{Prefix: "ex"},
	}

	valid := ValidateName(name, checkers...)
	if !valid {
		t.Errorf("expected %q to be valid according to the configured checkers", name)
	}

	name = "invalidName"
	valid = ValidateName(name, checkers...)
	if valid {
		t.Errorf("expected %q to be invalid according to the configured checkers", name)
	}
}
