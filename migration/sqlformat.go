package migration

import (
	"strings"
)

// SortMigrationStatements reorders SQL statements so that drops occur first,
// followed by creates/alters, and indexes/constraints last.
// 1. Drops (Indexes, constraints, tables)
// 2. Alters & Creates (Tables, columns, types)
// 3. Additions (Indexes, constraints, foreign keys)
func SortMigrationStatements(stmts []string) []string {
	var drops []string
	var alters []string
	var indexes []string

	for _, stmt := range stmts {
		trimmed := strings.ToUpper(strings.TrimSpace(stmt))
		// Identify Drops
		isDrop := strings.HasPrefix(trimmed, "DROP ") ||
			(strings.Contains(trimmed, "ALTER TABLE ") && strings.Contains(trimmed, " DROP "))

		// Identify Indexes/Add Constraints/Foreign Keys
		isAddIndexOrConstraint := strings.HasPrefix(trimmed, "CREATE INDEX") ||
			strings.HasPrefix(trimmed, "CREATE UNIQUE INDEX") ||
			(strings.Contains(trimmed, "ALTER TABLE ") && strings.Contains(trimmed, " ADD CONSTRAINT")) ||
			(strings.Contains(trimmed, "ALTER TABLE ") && strings.Contains(trimmed, " ADD FOREIGN KEY"))

		if isDrop {
			drops = append(drops, stmt)
		} else if isAddIndexOrConstraint {
			indexes = append(indexes, stmt)
		} else {
			alters = append(alters, stmt)
		}
	}

	result := append([]string{}, drops...)
	result = append(result, alters...)
	result = append(result, indexes...)
	return result
}

// formatSQL formats a CREATE TABLE statement to put column definitions on separate lines.
func formatSQL(query string) string {
	query = strings.TrimSpace(query)
	if !strings.HasPrefix(strings.ToUpper(query), "CREATE TABLE") {
		return query
	}

	firstParen := strings.Index(query, "(")
	if firstParen == -1 {
		return query
	}

	lastParen := strings.LastIndex(query, ")")
	if lastParen == -1 || lastParen < firstParen {
		return query
	}

	prefix := query[:firstParen+1]
	body := query[firstParen+1 : lastParen]
	suffix := query[lastParen:]

	var parts []string
	var current strings.Builder
	depth := 0
	for i := 0; i < len(body); i++ {
		char := body[i]
		switch char {
		case '(':
			depth++
		case ')':
			depth--
		}
		if char == ',' && depth == 0 {
			parts = append(parts, strings.TrimSpace(current.String()))
			current.Reset()
		} else {
			current.WriteByte(char)
		}
	}
	if current.Len() > 0 {
		parts = append(parts, strings.TrimSpace(current.String()))
	}

	var sb strings.Builder
	sb.WriteString(prefix)
	sb.WriteString("\n")
	for i, part := range parts {
		sb.WriteString("  ")
		sb.WriteString(part)
		if i < len(parts)-1 {
			sb.WriteString(",")
		}
		sb.WriteString("\n")
	}
	sb.WriteString(suffix)
	return sb.String()
}
