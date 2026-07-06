package migration

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	providers "github.com/voidclancy/valk/dbProviders"
	vs "github.com/voidclancy/valk/schema"

	"ariga.io/atlas/sql/migrate"
	"ariga.io/atlas/sql/schema"
)

type diffPlanner struct {
	db            *sql.DB
	provider      providers.DbProvider
	targetSchema  *vs.Schema
	isInteractive bool

	ctx           context.Context
	driver        migrate.Driver
	currentSchema *schema.Schema
	targetAtlas   *schema.Schema

	renamedMap map[string]map[string]string
}

// DiffAndPlan compares the current database schema state to the parsed target schema definition,
// generating the corresponding Goose Up and Down migration SQL strings.
func DiffAndPlan(db *sql.DB, provider providers.DbProvider, targetSchema *vs.Schema, isInteractive bool) (upSQL, downSQL string, err error) {
	planner := &diffPlanner{
		db:            db,
		provider:      provider,
		targetSchema:  targetSchema,
		isInteractive: isInteractive,
		ctx:           context.Background(),
		renamedMap:    make(map[string]map[string]string),
	}

	if err := planner.initDriver(); err != nil {
		return "", "", err
	}

	if err := planner.inspectCurrentSchema(); err != nil {
		return "", "", err
	}

	if err := planner.buildTargetSchema(); err != nil {
		return "", "", err
	}

	upSQL, err = planner.planDirection(planner.currentSchema, planner.targetAtlas, "up", planner.detectUpRenames)
	if err != nil {
		return "", "", err
	}

	downSQL, err = planner.planDirection(planner.targetAtlas, planner.currentSchema, "down", planner.applyDownRenames)
	if err != nil {
		return "", "", err
	}

	return upSQL, downSQL, nil
}

func (p *diffPlanner) initDriver() error {
	dialect := GetDialect(p.provider)
	driver, err := dialect.OpenConn(p.db)
	if err != nil {
		return fmt.Errorf("failed to open atlas driver: %w", err)
	}
	p.driver = driver
	return nil
}

func (p *diffPlanner) inspectCurrentSchema() error {
	current, err := p.driver.InspectSchema(p.ctx, "", nil)
	if err != nil {
		// Fallback to empty schema if inspection fails or database is brand new
		current = &schema.Schema{
			Name: "public",
		}
		if p.provider == providers.Sqlite {
			current.Name = ""
		}
	} else {
		// Filter out internal migration tracking tables like goose_db_version
		var filteredTables []*schema.Table
		for _, tbl := range current.Tables {
			if tbl.Name != "goose_db_version" {
				filteredTables = append(filteredTables, tbl)
			}
		}
		current.Tables = filteredTables
	}
	p.currentSchema = current
	return nil
}

func (p *diffPlanner) buildTargetSchema() error {
	target, err := ConvertToAtlasSchema(p.targetSchema, p.provider, p.currentSchema.Name)
	if err != nil {
		return fmt.Errorf("failed to build target schema: %w", err)
	}
	p.targetAtlas = target
	return nil
}

func (p *diffPlanner) planDirection(from, to *schema.Schema, direction string, applyRenames func([]schema.Change)) (string, error) {
	changes, err := p.driver.SchemaDiff(from, to)
	if err != nil {
		return "", fmt.Errorf("failed to diff schemas (%s): %w", direction, err)
	}

	applyRenames(changes)

	if len(changes) == 0 {
		return "", nil
	}

	plan, err := p.driver.PlanChanges(p.ctx, "migration", changes)
	if err != nil {
		return "", fmt.Errorf("failed to plan %s changes: %w", direction, err)
	}

	var stmts []string
	for _, c := range plan.Changes {
		stmt := c.Cmd
		if !strings.HasSuffix(stmt, ";") {
			stmt += ";"
		}
		stmts = append(stmts, formatSQL(stmt))
	}
	if p.provider != providers.Sqlite {
		stmts = SortMigrationStatements(stmts)
	}
	return strings.Join(stmts, "\n"), nil
}
