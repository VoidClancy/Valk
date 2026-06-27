package providers

import "fmt"

type DbProvider string

const (
	Postgres   DbProvider = "postgres"
	Postgresql DbProvider = "postgresql"
	Sqlite     DbProvider = "sqlite"
	Mysql      DbProvider = "mysql"
)

func ParseDbProvider(s string) (DbProvider, error) {
	switch s {
	case "postgres", "postgresql":
		return Postgres, nil

	case "mysql":
		return "", fmt.Errorf("Mysql is not supported yet") //TODO

	case "sqlite":
		return Sqlite, nil

	default:
		return "", fmt.Errorf("unknown provider: %s", s)
	}
}
