package migration

import (
	"database/sql"
	"strings"
	providers "valk/dbProviders"
	"valk/schema"

	"ariga.io/atlas/sql/migrate"
)

type Dialect interface {
	GetSQLType(sf *schema.ScalarField) string
	GetSQLDefault(dv *schema.DefaultValue, pslType string) string
	QuoteIdent(name string) string

	GenerateEnum(enum *schema.Enum) string
	FormatAutoIncrement(sqlType string) (typeOverride string, extraKeyword string)
	FormatSinglePK(tableName, colName string, isAutoInc bool) (inlineSQL string, tableConstraint string)
	FormatEnumConstraint(colName string, enum *schema.Enum) string
	SupportsNativeEnums() bool
	OpenConn(db *sql.DB) (migrate.Driver, error)
}

func GetDialect(provider providers.DbProvider) Dialect {
	switch provider {

	case providers.Mysql:
		return nil //TODO

	case providers.Postgres, providers.Postgresql:
		return &PostgresDialect{}

	case providers.Sqlite:
		return &SqliteDialect{}

	default:
		return nil

	}
}

func getSQLType(sf *schema.ScalarField, provider providers.DbProvider) string {
	dialect := GetDialect(provider)
	if dialect == nil {
		return strings.ToUpper(sf.SQLType)
	}

	return dialect.GetSQLType(sf)
}

func getSQLDefault(dv *schema.DefaultValue, pslType string, provider providers.DbProvider) string {
	dialect := GetDialect(provider)
	if dialect == nil {
		return ""
	}

	return dialect.GetSQLDefault(dv, pslType)
}
