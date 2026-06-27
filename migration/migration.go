package migration

import (
	"fmt"
	"strings"
	providers "valkyrie/dbProviders"
	"valkyrie/schema"
)

func GenerateUpMigrations(schemaDef *schema.Schema) (string, error) {

	dialect := GetDialect(schemaDef.Datasource.Provider)
	if dialect == nil {
		return "", fmt.Errorf("unknown provider: %s", schemaDef.Datasource.Provider)
	}

	provider := schemaDef.Datasource.Provider

	var sb strings.Builder

	if provider == providers.Sqlite {
		sb.WriteString("PRAGMA foreign_keys = ON;\n\n")
	}
	generateEnums(schemaDef.Enums, dialect, &sb)

	generateTables(schemaDef, dialect, &sb)

	return sb.String(), nil
}

func GenerateDownMigrations(schemaDef *schema.Schema) (string, error) {

	dialect := GetDialect(schemaDef.Datasource.Provider)
	if dialect == nil {
		return "", fmt.Errorf("unknown provider: %s", schemaDef.Datasource.Provider)
	}

	provider := schemaDef.Datasource.Provider

	var downBuilder strings.Builder

	if provider == providers.Sqlite {
		downBuilder.WriteString("PRAGMA foreign_keys = ON;\n\n")
	}

	dropTables(schemaDef, dialect, &downBuilder)
	dropEnums(schemaDef, dialect, &downBuilder)

	return downBuilder.String(), nil
}

func GenerateMigration(schemaDef *schema.Schema) (string, error) {

	upSQL, err := GenerateUpMigrations(schemaDef)
	if err != nil {
		return "", err
	}

	downSQL, err := GenerateDownMigrations(schemaDef)
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	sb.WriteString("-- +goose Up\n")
	sb.WriteString(upSQL)
	sb.WriteString("\n-- +goose Down\n")
	sb.WriteString(downSQL)
	return sb.String(), nil
}
