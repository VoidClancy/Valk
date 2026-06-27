package migration

import (
	"strings"
	"valkyrie/schema"
)

type Dialect interface {
	GetSQLType(sf *schema.ScalarField) string
	GetSQLDefault(dv *schema.DefaultValue, pslType string) string
	QuoteIdent(name string) string
}

func GetDialect(provider schema.DbProvider) Dialect {
	switch provider {
	case schema.Mysql:
		return nil //TODO
	case schema.Postgres, schema.Postgresql:
		return &PostgresDialect{}
	case schema.Sqlite:
		return &SqliteDialect{}
	default:
		return nil
	}
}

func getSQLType(sf *schema.ScalarField, provider schema.DbProvider) string {
	dialect := GetDialect(provider)
	if dialect == nil {
		return strings.ToUpper(sf.SQLType)
	}
	return dialect.GetSQLType(sf)
}

func getSQLDefault(dv *schema.DefaultValue, pslType string, provider schema.DbProvider) string {
	dialect := GetDialect(provider)
	if dialect == nil {
		return ""
	}
	return dialect.GetSQLDefault(dv, pslType)
}
