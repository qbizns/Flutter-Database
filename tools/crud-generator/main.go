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
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip CREATE TABLE line, constraints, and closing
		if strings.HasPrefix(strings.ToUpper(line), "CREATE TABLE") ||
			strings.HasPrefix(strings.ToUpper(line), "CONSTRAINT") ||
			strings.HasPrefix(strings.ToUpper(line), "PRIMARY KEY") ||
			strings.HasPrefix(strings.ToUpper(line), "FOREIGN KEY") ||
			strings.HasPrefix(strings.ToUpper(line), "UNIQUE") ||
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

func main() {
	fmt.Println("=== POS Backend CRUD Generator ===")
	fmt.Println("Production-Grade Code Generator for 170+ Tables")
	fmt.Println()

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

	fmt.Println("Schema analysis complete!")
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("1. Generate repository layer")
	fmt.Println("2. Generate service layer")
	fmt.Println("3. Generate handler layer")
	fmt.Println("4. Generate DTOs")
	fmt.Println("5. Generate routes")
	fmt.Println("6. Generate validators")
	fmt.Println("7. Generate tests")
}
