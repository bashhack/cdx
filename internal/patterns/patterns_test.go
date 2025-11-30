package patterns

import (
	"testing"
)

func TestDetectLanguage(t *testing.T) {
	tests := []struct {
		ext  string
		want Language
	}{
		{".go", Go},
		{".ts", TypeScript},
		{".tsx", TypeScript},
		{".js", JavaScript},
		{".jsx", JavaScript},
		{".mjs", JavaScript},
		{".py", Python},
		{".rs", Rust},
		{".unknown", Unknown},
		{"", Unknown},
	}

	for _, tt := range tests {
		t.Run(tt.ext, func(t *testing.T) {
			if got := DetectLanguage(tt.ext); got != tt.want {
				t.Errorf("DetectLanguage(%q) = %q, want %q", tt.ext, got, tt.want)
			}
		})
	}
}

func TestForLanguage(t *testing.T) {
	tests := []struct {
		lang    Language
		wantNil bool
	}{
		{Go, false},
		{TypeScript, false},
		{JavaScript, false},
		{Python, false},
		{Rust, false},
		{Unknown, true},
		{Language("invalid"), true},
	}

	for _, tt := range tests {
		t.Run(string(tt.lang), func(t *testing.T) {
			got := ForLanguage(tt.lang)
			if (got == nil) != tt.wantNil {
				t.Errorf("ForLanguage(%q) nil = %v, want nil = %v", tt.lang, got == nil, tt.wantNil)
			}
		})
	}
}

func TestGoPatterns_MatchFunctions(t *testing.T) {
	lp := ForLanguage(Go)
	if lp == nil {
		t.Fatal("expected Go patterns")
	}

	tests := []struct {
		name     string
		line     string
		wantKind string
		wantName string
	}{
		{
			name:     "simple function",
			line:     "func GetUser(id int) *User {",
			wantKind: "function",
			wantName: "GetUser",
		},
		{
			name:     "method with pointer receiver",
			line:     "func (s *Server) Start() error {",
			wantKind: "method",
			wantName: "Start",
		},
		{
			name:     "method with value receiver",
			line:     "func (u User) String() string {",
			wantKind: "method",
			wantName: "String",
		},
		{
			name:     "struct type",
			line:     "type User struct {",
			wantKind: "type",
			wantName: "User",
		},
		{
			name:     "interface type",
			line:     "type Repository interface {",
			wantKind: "interface",
			wantName: "Repository",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var matched bool
			var matchedKind string

			for _, p := range lp.Definition {
				if matches := p.Regex.FindStringSubmatch(tt.line); len(matches) > 1 {
					if matches[1] == tt.wantName {
						matched = true
						matchedKind = p.Kind
						break
					}
				}
			}

			if !matched {
				t.Errorf("expected line to match pattern for %q", tt.wantName)
			}

			if matchedKind != tt.wantKind {
				t.Errorf("kind = %q, want %q", matchedKind, tt.wantKind)
			}
		})
	}
}

func TestDefinitionPatternFor(t *testing.T) {
	tests := []struct {
		name       string
		symbol     string
		lang       Language
		testLine   string
		shouldFind bool
	}{
		{
			name:       "Go function",
			symbol:     "GetUser",
			lang:       Go,
			testLine:   "func GetUser(id int) *User {",
			shouldFind: true,
		},
		{
			name:       "Go function no match",
			symbol:     "GetUser",
			lang:       Go,
			testLine:   "func GetUserByID(id int) *User {",
			shouldFind: false,
		},
		{
			name:       "Go method",
			symbol:     "Start",
			lang:       Go,
			testLine:   "func (s *Server) Start() error {",
			shouldFind: true,
		},
		{
			name:       "Go type",
			symbol:     "User",
			lang:       Go,
			testLine:   "type User struct {",
			shouldFind: true,
		},
		{
			name:       "TypeScript class",
			symbol:     "UserService",
			lang:       TypeScript,
			testLine:   "export class UserService {",
			shouldFind: true,
		},
		{
			name:       "TypeScript function",
			symbol:     "createUser",
			lang:       TypeScript,
			testLine:   "export function createUser(name: string): User {",
			shouldFind: true,
		},
		{
			name:       "Python function",
			symbol:     "get_user",
			lang:       Python,
			testLine:   "def get_user(user_id: int) -> User:",
			shouldFind: true,
		},
		{
			name:       "Python class",
			symbol:     "UserService",
			lang:       Python,
			testLine:   "class UserService:",
			shouldFind: true,
		},
		{
			name:       "Rust function",
			symbol:     "create_user",
			lang:       Rust,
			testLine:   "pub fn create_user(name: String) -> User {",
			shouldFind: true,
		},
		{
			name:       "Rust struct",
			symbol:     "User",
			lang:       Rust,
			testLine:   "pub struct User {",
			shouldFind: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			patterns := DefinitionPatternFor(tt.symbol, tt.lang)
			if len(patterns) == 0 {
				t.Fatal("expected patterns to be generated")
			}

			var found bool
			for _, p := range patterns {
				if p.MatchString(tt.testLine) {
					found = true
					break
				}
			}

			if found != tt.shouldFind {
				t.Errorf("pattern match = %v, want %v for line %q", found, tt.shouldFind, tt.testLine)
			}
		})
	}
}

func TestLanguagePatterns_TestFileDetection(t *testing.T) {
	tests := []struct {
		lang     Language
		filename string
		isTest   bool
	}{
		{Go, "user_test.go", true},
		{Go, "user.go", false},
		{TypeScript, "user.test.ts", true},
		{TypeScript, "user.spec.ts", true},
		{TypeScript, "user.ts", false},
		{JavaScript, "user.test.js", true},
		{JavaScript, "user.spec.js", true},
		{JavaScript, "user.js", false},
		{Python, "test_user.py", true},
		{Python, "user_test.py", true},
		{Python, "user.py", false},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			lp := ForLanguage(tt.lang)
			if lp == nil {
				t.Fatalf("no patterns for %q", tt.lang)
			}

			if lp.TestFile == nil {
				t.Fatalf("no TestFile pattern for %q", tt.lang)
			}

			got := lp.TestFile.MatchString(tt.filename)
			if got != tt.isTest {
				t.Errorf("TestFile.Match(%q) = %v, want %v", tt.filename, got, tt.isTest)
			}
		})
	}
}

func TestAllLanguages(t *testing.T) {
	langs := AllLanguages()

	if len(langs) < 5 {
		t.Errorf("expected at least 5 languages, got %d", len(langs))
	}

	// Verify all returned languages have patterns
	for _, lang := range langs {
		if ForLanguage(lang) == nil {
			t.Errorf("AllLanguages() returned %q but ForLanguage(%q) is nil", lang, lang)
		}
	}
}

func TestArrowFunctionPatterns(t *testing.T) {
	tests := []struct {
		name        string
		lang        Language
		line        string
		wantName    string
		shouldMatch bool
	}{
		// TypeScript - should match
		{
			name:        "TS basic arrow function",
			lang:        TypeScript,
			line:        "const fetchUser = () => {",
			shouldMatch: true,
			wantName:    "fetchUser",
		},
		{
			name:        "TS arrow with params",
			lang:        TypeScript,
			line:        "const add = (a: number, b: number) => a + b",
			shouldMatch: true,
			wantName:    "add",
		},
		{
			name:        "TS async arrow function",
			lang:        TypeScript,
			line:        "const fetchData = async () => {",
			shouldMatch: true,
			wantName:    "fetchData",
		},
		{
			name:        "TS exported arrow function",
			lang:        TypeScript,
			line:        "export const handler = () => {",
			shouldMatch: true,
			wantName:    "handler",
		},
		{
			name:        "TS nested function type in params",
			lang:        TypeScript,
			line:        "const withCallback = (cb: (x: number) => void) => {",
			shouldMatch: true,
			wantName:    "withCallback",
		},
		{
			name:        "TS single param without parens",
			lang:        TypeScript,
			line:        "const double = x => x * 2",
			shouldMatch: true,
			wantName:    "double",
		},
		{
			name:        "TS async single param",
			lang:        TypeScript,
			line:        "const processItem = async item => {",
			shouldMatch: true,
			wantName:    "processItem",
		},
		{
			name:        "TS arrow with return type annotation",
			lang:        TypeScript,
			line:        "export const fetchUser = async (id: number): Promise<User> => {",
			shouldMatch: true,
			wantName:    "fetchUser",
		},
		{
			name:        "TS arrow with simple return type",
			lang:        TypeScript,
			line:        "const getName = (user: User): string => {",
			shouldMatch: true,
			wantName:    "getName",
		},
		// TypeScript - should NOT match (false positives)
		{
			name:        "TS parenthesized expression",
			lang:        TypeScript,
			line:        "const result = (a + b) * 2",
			shouldMatch: false,
			wantName:    "",
		},
		{
			name:        "TS function call result",
			lang:        TypeScript,
			line:        "const value = (calculateSomething())",
			shouldMatch: false,
			wantName:    "",
		},
		{
			name:        "TS grouping in expression",
			lang:        TypeScript,
			line:        "const x = (y + z)",
			shouldMatch: false,
			wantName:    "",
		},
		// JavaScript - should match
		{
			name:        "JS basic arrow function",
			lang:        JavaScript,
			line:        "const greet = () => {",
			shouldMatch: true,
			wantName:    "greet",
		},
		{
			name:        "JS single param without parens",
			lang:        JavaScript,
			line:        "const square = n => n * n",
			shouldMatch: true,
			wantName:    "square",
		},
		{
			name:        "JS exported async arrow",
			lang:        JavaScript,
			line:        "export const fetchAPI = async () => {",
			shouldMatch: true,
			wantName:    "fetchAPI",
		},
		// JavaScript - should NOT match
		{
			name:        "JS parenthesized expression",
			lang:        JavaScript,
			line:        "const total = (price + tax)",
			shouldMatch: false,
			wantName:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lp := ForLanguage(tt.lang)
			if lp == nil {
				t.Fatalf("no patterns for %q", tt.lang)
			}

			var matched bool
			var matchedName string

			for _, p := range lp.Definition {
				if p.Kind != "function" {
					continue
				}
				if matches := p.Regex.FindStringSubmatch(tt.line); len(matches) > 1 {
					matched = true
					matchedName = matches[1]
					break
				}
			}

			if matched != tt.shouldMatch {
				t.Errorf("match = %v, want %v for line %q", matched, tt.shouldMatch, tt.line)
			}

			if tt.shouldMatch && matchedName != tt.wantName {
				t.Errorf("matched name = %q, want %q", matchedName, tt.wantName)
			}
		})
	}
}

func TestGoConstPatterns(t *testing.T) {
	tests := []struct {
		name        string
		line        string
		wantName    string
		shouldMatch bool
	}{
		// Should match
		{
			name:        "standalone const",
			line:        "const MaxRetries = 3",
			shouldMatch: true,
			wantName:    "MaxRetries",
		},
		{
			name:        "const with type",
			line:        "const DefaultTimeout time.Duration = 30",
			shouldMatch: true,
			wantName:    "DefaultTimeout",
		},
		{
			name:        "tab-indented const block member",
			line:        "\tMaxConnections = 100",
			shouldMatch: true,
			wantName:    "MaxConnections",
		},
		{
			name:        "tab-indented const with type",
			line:        "\tStatusOK Status = 200",
			shouldMatch: true,
			wantName:    "StatusOK",
		},
		// Should NOT match (false positives)
		{
			name:        "space-indented var (not const block)",
			line:        "    MAX_VALUE := 100",
			shouldMatch: false,
			wantName:    "",
		},
		{
			name:        "lowercase identifier",
			line:        "\tmaxValue = 100",
			shouldMatch: false,
			wantName:    "",
		},
	}

	lp := ForLanguage(Go)
	if lp == nil {
		t.Fatal("no patterns for Go")
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var matched bool
			var matchedName string

			for _, p := range lp.Definition {
				if p.Kind != "const" {
					continue
				}
				if matches := p.Regex.FindStringSubmatch(tt.line); len(matches) > 1 {
					matched = true
					matchedName = matches[1]
					break
				}
			}

			if matched != tt.shouldMatch {
				t.Errorf("match = %v, want %v for line %q", matched, tt.shouldMatch, tt.line)
			}

			if tt.shouldMatch && matchedName != tt.wantName {
				t.Errorf("matched name = %q, want %q", matchedName, tt.wantName)
			}
		})
	}
}
