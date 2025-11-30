// Package patterns provides language-specific regex patterns for code search.
package patterns

import "regexp"

// Language represents a programming language.
type Language string

const (
	Go         Language = "go"
	TypeScript Language = "ts"
	JavaScript Language = "js"
	Python     Language = "py"
	Rust       Language = "rust"
	Unknown    Language = ""
)

// Pattern holds a compiled regex and metadata about what it matches.
type Pattern struct {
	Regex *regexp.Regexp
	Kind  string // "function", "type", "method", "interface", "const", "var"
}

// LanguagePatterns holds all definition patterns for a language.
type LanguagePatterns struct {
	Language   Language
	TestFile   *regexp.Regexp // Pattern to identify test files
	Definition []Pattern
	Extensions []string
}

// registry maps languages to their patterns.
var registry = map[Language]*LanguagePatterns{
	Go:         goPatterns(),
	TypeScript: tsPatterns(),
	JavaScript: jsPatterns(),
	Python:     pythonPatterns(),
	Rust:       rustPatterns(),
}

// ForLanguage returns patterns for the given language.
func ForLanguage(lang Language) *LanguagePatterns {
	if p, ok := registry[lang]; ok {
		return p
	}
	return nil
}

// DetectLanguage determines language from file extension.
func DetectLanguage(ext string) Language {
	switch ext {
	case ".go":
		return Go
	case ".ts", ".tsx":
		return TypeScript
	case ".js", ".jsx", ".mjs":
		return JavaScript
	case ".py":
		return Python
	case ".rs":
		return Rust
	default:
		return Unknown
	}
}

// AllLanguages returns all supported languages.
func AllLanguages() []Language {
	return []Language{Go, TypeScript, JavaScript, Python, Rust}
}

// goPatterns returns Go-specific patterns.
func goPatterns() *LanguagePatterns {
	return &LanguagePatterns{
		Language:   Go,
		Extensions: []string{".go"},
		Definition: []Pattern{
			// func FunctionName(
			{
				Regex: regexp.MustCompile(`^func\s+([A-Za-z_][A-Za-z0-9_]*)\s*\(`),
				Kind:  "function",
			},
			// func (receiver) MethodName(
			{
				Regex: regexp.MustCompile(`^func\s+\([^)]+\)\s+([A-Za-z_][A-Za-z0-9_]*)\s*\(`),
				Kind:  "method",
			},
			// type TypeName struct/interface
			{
				Regex: regexp.MustCompile(`^type\s+([A-Za-z_][A-Za-z0-9_]*)\s+struct\b`),
				Kind:  "type",
			},
			{
				Regex: regexp.MustCompile(`^type\s+([A-Za-z_][A-Za-z0-9_]*)\s+interface\b`),
				Kind:  "interface",
			},
			// type TypeName = ... (type alias)
			{
				Regex: regexp.MustCompile(`^type\s+([A-Za-z_][A-Za-z0-9_]*)\s+=`),
				Kind:  "type",
			},
			// type TypeName SomeOtherType (type definition)
			{
				Regex: regexp.MustCompile(`^type\s+([A-Za-z_][A-Za-z0-9_]*)\s+[A-Za-z]`),
				Kind:  "type",
			},
			// const ConstName = or const ( ConstName =
			{
				Regex: regexp.MustCompile(`^\s*([A-Z_][A-Za-z0-9_]*)\s*(?:=|[A-Za-z])`),
				Kind:  "const",
			},
			// var VarName = or var VarName Type
			{
				Regex: regexp.MustCompile(`^var\s+([A-Za-z_][A-Za-z0-9_]*)\s*(?:=|[A-Za-z\[])`),
				Kind:  "var",
			},
		},
		TestFile: regexp.MustCompile(`_test\.go$`),
	}
}

// tsPatterns returns TypeScript-specific patterns.
func tsPatterns() *LanguagePatterns {
	return &LanguagePatterns{
		Language:   TypeScript,
		Extensions: []string{".ts", ".tsx"},
		Definition: []Pattern{
			// function functionName(
			{
				Regex: regexp.MustCompile(`^(?:export\s+)?(?:async\s+)?function\s+([A-Za-z_$][A-Za-z0-9_$]*)\s*[<(]`),
				Kind:  "function",
			},
			// const functionName = (arrow function)
			{
				Regex: regexp.MustCompile(`^(?:export\s+)?const\s+([A-Za-z_$][A-Za-z0-9_$]*)\s*=\s*(?:async\s*)?\(`),
				Kind:  "function",
			},
			// class ClassName
			{
				Regex: regexp.MustCompile(`^(?:export\s+)?(?:abstract\s+)?class\s+([A-Za-z_$][A-Za-z0-9_$]*)`),
				Kind:  "type",
			},
			// interface InterfaceName
			{
				Regex: regexp.MustCompile(`^(?:export\s+)?interface\s+([A-Za-z_$][A-Za-z0-9_$]*)`),
				Kind:  "interface",
			},
			// type TypeName =
			{
				Regex: regexp.MustCompile(`^(?:export\s+)?type\s+([A-Za-z_$][A-Za-z0-9_$]*)\s*[<=]`),
				Kind:  "type",
			},
			// enum EnumName
			{
				Regex: regexp.MustCompile(`^(?:export\s+)?enum\s+([A-Za-z_$][A-Za-z0-9_$]*)`),
				Kind:  "type",
			},
		},
		TestFile: regexp.MustCompile(`\.(test|spec)\.tsx?$`),
	}
}

// jsPatterns returns JavaScript-specific patterns.
func jsPatterns() *LanguagePatterns {
	return &LanguagePatterns{
		Language:   JavaScript,
		Extensions: []string{".js", ".jsx", ".mjs"},
		Definition: []Pattern{
			// function functionName(
			{
				Regex: regexp.MustCompile(`^(?:export\s+)?(?:async\s+)?function\s+([A-Za-z_$][A-Za-z0-9_$]*)\s*\(`),
				Kind:  "function",
			},
			// const functionName = (arrow function)
			{
				Regex: regexp.MustCompile(`^(?:export\s+)?const\s+([A-Za-z_$][A-Za-z0-9_$]*)\s*=\s*(?:async\s*)?\(`),
				Kind:  "function",
			},
			// class ClassName
			{
				Regex: regexp.MustCompile(`^(?:export\s+)?class\s+([A-Za-z_$][A-Za-z0-9_$]*)`),
				Kind:  "type",
			},
		},
		TestFile: regexp.MustCompile(`\.(test|spec)\.(js|jsx|mjs)$`),
	}
}

// pythonPatterns returns Python-specific patterns.
func pythonPatterns() *LanguagePatterns {
	return &LanguagePatterns{
		Language:   Python,
		Extensions: []string{".py"},
		Definition: []Pattern{
			// def function_name(
			{
				Regex: regexp.MustCompile(`^(?:async\s+)?def\s+([A-Za-z_][A-Za-z0-9_]*)\s*\(`),
				Kind:  "function",
			},
			// class ClassName
			{
				Regex: regexp.MustCompile(`^class\s+([A-Za-z_][A-Za-z0-9_]*)`),
				Kind:  "type",
			},
		},
		TestFile: regexp.MustCompile(`(^test_|_test\.py$)`),
	}
}

// rustPatterns returns Rust-specific patterns.
func rustPatterns() *LanguagePatterns {
	return &LanguagePatterns{
		Language:   Rust,
		Extensions: []string{".rs"},
		Definition: []Pattern{
			// fn function_name(
			{
				Regex: regexp.MustCompile(`^(?:pub\s+)?(?:async\s+)?fn\s+([A-Za-z_][A-Za-z0-9_]*)\s*[<(]`),
				Kind:  "function",
			},
			// struct StructName
			{
				Regex: regexp.MustCompile(`^(?:pub\s+)?struct\s+([A-Za-z_][A-Za-z0-9_]*)`),
				Kind:  "type",
			},
			// enum EnumName
			{
				Regex: regexp.MustCompile(`^(?:pub\s+)?enum\s+([A-Za-z_][A-Za-z0-9_]*)`),
				Kind:  "type",
			},
			// trait TraitName
			{
				Regex: regexp.MustCompile(`^(?:pub\s+)?trait\s+([A-Za-z_][A-Za-z0-9_]*)`),
				Kind:  "interface",
			},
			// impl TraitName for or impl StructName
			{
				Regex: regexp.MustCompile(`^impl\s+(?:<[^>]+>\s+)?([A-Za-z_][A-Za-z0-9_]*)`),
				Kind:  "type",
			},
		},
		TestFile: regexp.MustCompile(`(^test_|_test\.rs$|/tests/)`),
	}
}

// DefinitionPatternFor builds a regex pattern to find definitions of a specific symbol.
func DefinitionPatternFor(symbol string, lang Language) []*regexp.Regexp {
	lp := ForLanguage(lang)
	if lp == nil {
		return nil
	}

	var patterns []*regexp.Regexp
	for _, p := range lp.Definition {
		// Build a pattern that specifically matches this symbol
		var patStr string
		switch lang {
		case Go:
			switch p.Kind {
			case "function":
				patStr = `^func\s+` + regexp.QuoteMeta(symbol) + `\s*\(`
			case "method":
				patStr = `^func\s+\([^)]+\)\s+` + regexp.QuoteMeta(symbol) + `\s*\(`
			case "type", "interface":
				patStr = `^type\s+` + regexp.QuoteMeta(symbol) + `\s+`
			case "const":
				patStr = `(?:^const\s+|^\s+)` + regexp.QuoteMeta(symbol) + `\s*(?:=|[A-Za-z])`
			case "var":
				patStr = `^var\s+` + regexp.QuoteMeta(symbol) + `\s*`
			}
		case TypeScript, JavaScript:
			switch p.Kind {
			case "function":
				patStr = `(?:^(?:export\s+)?(?:async\s+)?function\s+` + regexp.QuoteMeta(symbol) + `|^(?:export\s+)?const\s+` + regexp.QuoteMeta(symbol) + `\s*=)`
			case "type", "interface":
				patStr = `^(?:export\s+)?(?:class|interface|type|enum)\s+` + regexp.QuoteMeta(symbol)
			}
		case Python:
			switch p.Kind {
			case "function":
				patStr = `^(?:async\s+)?def\s+` + regexp.QuoteMeta(symbol) + `\s*\(`
			case "type":
				patStr = `^class\s+` + regexp.QuoteMeta(symbol)
			}
		case Rust:
			switch p.Kind {
			case "function":
				patStr = `^(?:pub\s+)?(?:async\s+)?fn\s+` + regexp.QuoteMeta(symbol) + `\s*[<(]`
			case "type":
				patStr = `^(?:pub\s+)?(?:struct|enum)\s+` + regexp.QuoteMeta(symbol)
			case "interface":
				patStr = `^(?:pub\s+)?trait\s+` + regexp.QuoteMeta(symbol)
			}
		}

		if patStr != "" {
			if re, err := regexp.Compile(patStr); err == nil {
				patterns = append(patterns, re)
			}
		}
	}

	return patterns
}
