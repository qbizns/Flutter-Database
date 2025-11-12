package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

// TableInfo represents a database table structure
type TableInfo struct {
	Name            string
	GoName          string        // PascalCase name for Go structs
	CamelName       string        // camelCase for variables
	Package         string        // package name
	Columns         []ColumnInfo
	PrimaryKey      string
	HasOrganization bool // true if table has organization_id
	HasSoftDelete   bool // true if table has deleted_at
	HasTimestamps   bool // true if has created_at/updated_at
	ForeignKeys     []ForeignKey
	Indexes         []IndexInfo
}

// ColumnInfo represents a table column
type ColumnInfo struct {
	Name       string
	GoName     string // PascalCase for struct field
	Type       string // SQL type
	GoType     string // Go type
	Nullable   bool
	IsPK       bool
	IsFK       bool
	FKTable    string
	DefaultVal string
	Validation string // validation tags
}

// ForeignKey represents a foreign key relationship
type ForeignKey struct {
	Column      string
	RefTable    string
	RefColumn   string
	OnDelete    string
	OnUpdate    string
}

// IndexInfo represents a table index
type IndexInfo struct {
	Name    string
	Columns []string
	Unique  bool
}

// SchemaAnalyzer parses SQL migrations to extract table structures
type SchemaAnalyzer struct {
	tables map[string]*TableInfo
}

// NewSchemaAnalyzer creates a new schema analyzer
func NewSchemaAnalyzer() *SchemaAnalyzer {
	return &SchemaAnalyzer{
		tables: make(map[string]*TableInfo),
	}
}

// AnalyzeMigrations reads all SQL migration files and extracts table structures
func (sa *SchemaAnalyzer) AnalyzeMigrations(migrationDirs []string) error {
	for _, dir := range migrationDirs {
		files, err := filepath.Glob(filepath.Join(dir, "*.sql"))
		if err != nil {
			return fmt.Errorf("failed to list migrations in %s: %w", dir, err)
		}

		for _, file := range files {
			if err := sa.parseFile(file); err != nil {
				fmt.Printf("Warning: failed to parse %s: %v\n", file, err)
				continue
			}
		}
	}

	return nil
}

// parseFile parses a single SQL migration file
func (sa *SchemaAnalyzer) parseFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var currentSQL strings.Builder
	inCreateTable := false

	for scanner.Scan() {
		line := scanner.Text()

		// Skip comments
		if strings.HasPrefix(strings.TrimSpace(line), "--") {
			continue
		}

		// Check for CREATE TABLE
		if strings.Contains(strings.ToUpper(line), "CREATE TABLE") {
			inCreateTable = true
			currentSQL.Reset()
		}

		if inCreateTable {
			currentSQL.WriteString(line)
			currentSQL.WriteString("\n")

			// Check if statement is complete
			if strings.Contains(line, ");") {
				sa.parseCreateTable(currentSQL.String())
				inCreateTable = false
				currentSQL.Reset()
			}
		}
	}

	return scanner.Err()
}

// parseCreateTable extracts table information from CREATE TABLE statement
func (sa *SchemaAnalyzer) parseCreateTable(sql string) {
	// Extract table name
	re := regexp.MustCompile(`CREATE TABLE\s+(?:IF NOT EXISTS\s+)?([a-zA-Z0-9_]+)\s*\(`)
	matches := re.FindStringSubmatch(sql)
	if len(matches) < 2 {
		return
	}

	tableName := matches[1]

	table := &TableInfo{
		Name:      tableName,
		GoName:    toGoName(tableName),
		CamelName: toCamelCase(tableName),
		Package:   toPackageName(tableName),
		Columns:   []ColumnInfo{},
	}

	// Parse columns
	lines := strings.Split(sql, "\n")
	inGeneratedColumn := false

	for _, line := range lines {
		line = strings.TrimSpace(line)
		upperLine := strings.ToUpper(line)

		// Check if entering a GENERATED column definition
		if strings.Contains(upperLine, "GENERATED ALWAYS AS") {
			inGeneratedColumn = true
		}

		// Check if exiting a GENERATED column definition
		if inGeneratedColumn && (strings.Contains(line, ") STORED") || strings.Contains(line, ") VIRTUAL")) {
			inGeneratedColumn = false
			continue
		}

		// Skip lines inside GENERATED column definitions
		if inGeneratedColumn {
			continue
		}

		// Skip CREATE TABLE line, constraints, and closing
		if strings.HasPrefix(upperLine, "CREATE TABLE") ||
			strings.HasPrefix(upperLine, "CONSTRAINT") ||
			strings.HasPrefix(upperLine, "PRIMARY KEY") ||
			strings.HasPrefix(upperLine, "FOREIGN KEY") ||
			strings.HasPrefix(upperLine, "UNIQUE") ||
			strings.HasPrefix(upperLine, "CHECK") ||
			strings.HasPrefix(line, ");") ||
			line == "" {
			continue
		}

		col := sa.parseColumn(line, table)
		if col.Name != "" {
			table.Columns = append(table.Columns, col)

			// Check for special columns
			if col.Name == "organization_id" {
				table.HasOrganization = true
			}
			if col.Name == "deleted_at" {
				table.HasSoftDelete = true
			}
			if col.Name == "created_at" || col.Name == "updated_at" {
				table.HasTimestamps = true
			}
			if col.IsPK {
				table.PrimaryKey = col.Name
			}
		}
	}

	// Only add tables with valid columns
	if len(table.Columns) > 0 {
		sa.tables[tableName] = table
	}
}

// parseColumn extracts column information from a line
func (sa *SchemaAnalyzer) parseColumn(line string, table *TableInfo) ColumnInfo {
	// Remove trailing comma
	line = strings.TrimSuffix(line, ",")

	parts := strings.Fields(line)
	if len(parts) < 2 {
		return ColumnInfo{}
	}

	col := ColumnInfo{
		Name:   parts[0],
		GoName: toGoName(parts[0]),
		Type:   parts[1],
	}

	// Parse type and constraints
	upperLine := strings.ToUpper(line)

	col.Nullable = !strings.Contains(upperLine, "NOT NULL")
	col.IsPK = strings.Contains(upperLine, "PRIMARY KEY")

	// Determine Go type
	col.GoType = sqlTypeToGoType(col.Type, col.Nullable)

	// Add validation tags
	col.Validation = generateValidationTags(col)

	// Check if it's a foreign key (ends with _id and not the primary key)
	if strings.HasSuffix(col.Name, "_id") && col.Name != "id" && !col.IsPK {
		col.IsFK = true
		col.FKTable = strings.TrimSuffix(col.Name, "_id")
	}

	return col
}

// sqlTypeToGoType converts SQL types to Go types
func sqlTypeToGoType(sqlType string, nullable bool) string {
	sqlType = strings.ToUpper(sqlType)

	var baseType string

	switch {
	case strings.Contains(sqlType, "UUID"):
		baseType = "uuid.UUID"
	case strings.Contains(sqlType, "VARCHAR"), strings.Contains(sqlType, "TEXT"):
		baseType = "string"
	case strings.Contains(sqlType, "INTEGER"), strings.Contains(sqlType, "INT"):
		baseType = "int64"
	case strings.Contains(sqlType, "DECIMAL"), strings.Contains(sqlType, "NUMERIC"):
		baseType = "float64"
	case strings.Contains(sqlType, "BOOLEAN"), strings.Contains(sqlType, "BOOL"):
		baseType = "bool"
	case strings.Contains(sqlType, "TIMESTAMP"), strings.Contains(sqlType, "DATE"):
		baseType = "time.Time"
	case strings.Contains(sqlType, "JSONB"), strings.Contains(sqlType, "JSON"):
		baseType = "json.RawMessage"
	default:
		baseType = "string"
	}

	if nullable && baseType != "json.RawMessage" {
		return "*" + baseType
	}

	return baseType
}

// generateValidationTags generates validation tags for a column
func generateValidationTags(col ColumnInfo) string {
	var tags []string

	if !col.Nullable && col.Name != "id" && !strings.HasSuffix(col.Name, "_at") {
		tags = append(tags, "required")
	}

	if col.GoType == "string" || col.GoType == "*string" {
		if strings.Contains(col.Name, "email") {
			tags = append(tags, "email")
		}
		if strings.Contains(col.Name, "url") {
			tags = append(tags, "url")
		}
		if col.Name == "phone" || strings.Contains(col.Name, "phone") {
			tags = append(tags, "e164")
		}
	}

	if col.IsFK {
		tags = append(tags, "uuid")
	}

	if len(tags) > 0 {
		return strings.Join(tags, ",")
	}

	return ""
}

// toGoName converts snake_case to PascalCase
func toGoName(name string) string {
	parts := strings.Split(name, "_")
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + part[1:]
		}
	}
	return strings.Join(parts, "")
}

// toCamelCase converts snake_case to camelCase
func toCamelCase(name string) string {
	goName := toGoName(name)
	if len(goName) > 0 {
		return strings.ToLower(goName[:1]) + goName[1:]
	}
	return goName
}

// toPackageName converts table name to package name
func toPackageName(tableName string) string {
	// Remove plural 's' and convert to singular
	name := tableName
	if strings.HasSuffix(name, "ies") {
		name = strings.TrimSuffix(name, "ies") + "y"
	} else if strings.HasSuffix(name, "ses") || strings.HasSuffix(name, "ches") || strings.HasSuffix(name, "xes") {
		name = strings.TrimSuffix(name, "es")
	} else if strings.HasSuffix(name, "s") && !strings.HasSuffix(name, "ss") {
		name = strings.TrimSuffix(name, "s")
	}

	return strings.ToLower(name)
}

// GetTables returns all analyzed tables
func (sa *SchemaAnalyzer) GetTables() map[string]*TableInfo {
	return sa.tables
}

// PrintSummary prints a summary of analyzed tables
func (sa *SchemaAnalyzer) PrintSummary() {
	fmt.Printf("\n=== Schema Analysis Summary ===\n")
	fmt.Printf("Total tables analyzed: %d\n\n", len(sa.tables))

	for name, table := range sa.tables {
		fmt.Printf("Table: %s\n", name)
		fmt.Printf("  Go Name: %s\n", table.GoName)
		fmt.Printf("  Package: %s\n", table.Package)
		fmt.Printf("  Columns: %d\n", len(table.Columns))
		fmt.Printf("  Primary Key: %s\n", table.PrimaryKey)
		fmt.Printf("  Has Organization: %v\n", table.HasOrganization)
		fmt.Printf("  Has Soft Delete: %v\n", table.HasSoftDelete)
		fmt.Printf("  Has Timestamps: %v\n", table.HasTimestamps)
		fmt.Println()
	}
}

// CodeGenerator generates production-grade CRUD code
type CodeGenerator struct {
	analyzer  *SchemaAnalyzer
	templates map[string]*template.Template
	outputDir string
}

// NewCodeGenerator creates a new code generator
func NewCodeGenerator(analyzer *SchemaAnalyzer, outputDir string) (*CodeGenerator, error) {
	gen := &CodeGenerator{
		analyzer:  analyzer,
		templates: make(map[string]*template.Template),
		outputDir: outputDir,
	}

	// Load templates
	if err := gen.loadTemplates(); err != nil {
		return nil, err
	}

	return gen, nil
}

// loadTemplates loads all code generation templates
func (g *CodeGenerator) loadTemplates() error {
	templateDir := "tools/crud-generator/templates"

	// Define template files
	templateFiles := map[string]string{
		"repository": filepath.Join(templateDir, "repository.tmpl"),
		"service":    filepath.Join(templateDir, "service.tmpl"),
		"handler":    filepath.Join(templateDir, "handler.tmpl"),
		"dto":        filepath.Join(templateDir, "dto.tmpl"),
		"routes":     filepath.Join(templateDir, "routes.tmpl"),
		"validator":  filepath.Join(templateDir, "validator.tmpl"),
		"test":       filepath.Join(templateDir, "test.tmpl"),
	}

	// Load each template
	for name, file := range templateFiles {
		tmpl, err := template.New(filepath.Base(file)).Funcs(templateFuncs()).ParseFiles(file)
		if err != nil {
			return fmt.Errorf("failed to load template %s: %w", name, err)
		}
		g.templates[name] = tmpl
	}

	return nil
}

// templateFuncs returns custom template functions
func templateFuncs() template.FuncMap {
	return template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"toGoName": toGoName,
		"contains": strings.Contains,
		"trimPrefix": strings.TrimPrefix,
	}
}

// Generate generates code for all tables
func (g *CodeGenerator) Generate() error {
	tables := g.analyzer.GetTables()

	fmt.Printf("\n=== Generating Code ===\n")
	fmt.Printf("Total tables: %d\n", len(tables))
	fmt.Printf("Output directory: %s\n\n", g.outputDir)

	generated := 0
	failed := 0

	for tableName, table := range tables {
		fmt.Printf("Generating %s... ", tableName)

		if err := g.generateTable(table); err != nil {
			fmt.Printf("❌ FAILED: %v\n", err)
			failed++
			continue
		}

		fmt.Printf("✅ SUCCESS\n")
		generated++
	}

	fmt.Printf("\n=== Generation Summary ===\n")
	fmt.Printf("✅ Success: %d tables\n", generated)
	fmt.Printf("❌ Failed: %d tables\n", failed)
	fmt.Printf("📁 Total files: %d\n", generated*7)

	return nil
}

// generateTable generates all files for a single table
func (g *CodeGenerator) generateTable(table *TableInfo) error {
	// Create package directory
	pkgDir := filepath.Join(g.outputDir, "internal", table.Package)
	if err := os.MkdirAll(pkgDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Generate each file
	files := map[string]string{
		"repository": "repository.go",
		"service":    "service.go",
		"handler":    "handler.go",
		"dto":        "dto.go",
		"routes":     "routes.go",
		"validator":  "validator.go",
		"test":       "repository_test.go",
	}

	for templateName, filename := range files {
		outputPath := filepath.Join(pkgDir, filename)
		if err := g.generateFile(templateName, table, outputPath); err != nil {
			return fmt.Errorf("failed to generate %s: %w", filename, err)
		}
	}

	return nil
}

// generateFile generates a single file from a template
func (g *CodeGenerator) generateFile(templateName string, table *TableInfo, outputPath string) error {
	tmpl, ok := g.templates[templateName]
	if !ok {
		return fmt.Errorf("template %s not found", templateName)
	}

	// Create output file
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Execute template
	if err := tmpl.Execute(file, table); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Format code
	if err := g.formatFile(outputPath); err != nil {
		return fmt.Errorf("failed to format file: %w", err)
	}

	return nil
}

// formatFile formats a Go source file using gofmt
func (g *CodeGenerator) formatFile(path string) error {
	// Skip formatting for now - will be done in bulk with `go fmt ./...`
	return nil
}

func main() {
	fmt.Println("=== POS Backend CRUD Generator ===")
	fmt.Println("Production-Grade Code Generator for 170+ Tables")
	fmt.Println()

	// Parse command line flags
	outputDir := "."
	if len(os.Args) > 1 {
		outputDir = os.Args[1]
	}

	analyzer := NewSchemaAnalyzer()

	// Analyze migrations
	migrationDirs := []string{
		"postgres/migrations",
		"accounting/migrations",
	}

	fmt.Println("Analyzing database schema...")
	if err := analyzer.AnalyzeMigrations(migrationDirs); err != nil {
		fmt.Fprintf(os.Stderr, "Error analyzing migrations: %v\n", err)
		os.Exit(1)
	}

	analyzer.PrintSummary()

	// Ask user to proceed
	fmt.Println("\nReady to generate code for all tables.")
	fmt.Print("Proceed? (y/n): ")

	var response string
	fmt.Scanln(&response)

	if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
		fmt.Println("Generation cancelled.")
		os.Exit(0)
	}

	// Create code generator
	generator, err := NewCodeGenerator(analyzer, outputDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating generator: %v\n", err)
		os.Exit(1)
	}

	// Generate all code
	if err := generator.Generate(); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating code: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n=== Code Generation Complete! ===")
	fmt.Println("\nNext steps:")
	fmt.Println("1. Run: go fmt ./...")
	fmt.Println("2. Run: go mod tidy")
	fmt.Println("3. Run: go build ./...")
	fmt.Println("4. Run: go test ./...")
	fmt.Println("5. Review and customize generated code as needed")
}
