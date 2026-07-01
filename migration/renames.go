package migration

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"ariga.io/atlas/sql/schema"
)

func typeCategory(t schema.Type) string {
	switch t.(type) {
	case *schema.IntegerType:
		return "integer"
	case *schema.StringType:
		return "string"
	case *schema.UUIDType:
		return "uuid"
	case *schema.BoolType:
		return "bool"
	case *schema.TimeType:
		return "time"
	case *schema.FloatType, *schema.DecimalType:
		return "float"
	case *schema.BinaryType:
		return "binary"
	case *schema.JSONType:
		return "json"
	case *schema.EnumType:
		return "enum"
	case *schema.SpatialType:
		return "spatial"
	default:
		return "other"
	}
}

func typesCompatible(c1, c2 *schema.Column) bool {
	if c1.Type == nil || c2.Type == nil {
		return false
	}
	if c1.Type.Type != nil && c2.Type.Type != nil {
		return typeCategory(c1.Type.Type) == typeCategory(c2.Type.Type)
	}
	return c1.Type.Raw == c2.Type.Raw
}

// renameMatchFn determines whether a dropped column should be paired with an added column as a rename.
// Returns the matched AddColumn, or nil if no match.
type renameMatchFn func(tableName string, drop *schema.DropColumn, adds []*schema.AddColumn, matchedAdds map[*schema.AddColumn]bool) *schema.AddColumn

func processRenames(changes []schema.Change, matchFn renameMatchFn, onMatch func(tableName, newName, oldName string)) {
	for _, change := range changes {
		modifyTable, ok := change.(*schema.ModifyTable)
		if !ok {
			continue
		}

		var drops []*schema.DropColumn
		var adds []*schema.AddColumn
		for _, tblChange := range modifyTable.Changes {
			if drop, ok := tblChange.(*schema.DropColumn); ok {
				drops = append(drops, drop)
			} else if add, ok := tblChange.(*schema.AddColumn); ok {
				adds = append(adds, add)
			}
		}

		matchedDrops := make(map[*schema.DropColumn]bool)
		matchedAdds := make(map[*schema.AddColumn]bool)
		var newChanges []schema.Change

		for _, drop := range drops {
			if renamedTo := matchFn(modifyTable.T.Name, drop, adds, matchedAdds); renamedTo != nil {
				matchedAdds[renamedTo] = true
				matchedDrops[drop] = true
				newChanges = append(newChanges, &schema.RenameColumn{
					From: drop.C,
					To:   renamedTo.C,
				})
				if onMatch != nil {
					onMatch(modifyTable.T.Name, renamedTo.C.Name, drop.C.Name)
				}
			}
		}

		for _, tblChange := range modifyTable.Changes {
			if drop, ok := tblChange.(*schema.DropColumn); ok {
				if !matchedDrops[drop] {
					newChanges = append(newChanges, drop)
				}
			} else if add, ok := tblChange.(*schema.AddColumn); ok {
				if !matchedAdds[add] {
					newChanges = append(newChanges, add)
				}
			} else {
				newChanges = append(newChanges, tblChange)
			}
		}
		modifyTable.Changes = newChanges
	}
}

func (p *diffPlanner) detectUpRenames(changes []schema.Change) {
	processRenames(changes, func(tableName string, drop *schema.DropColumn, adds []*schema.AddColumn, matchedAdds map[*schema.AddColumn]bool) *schema.AddColumn {
		if p.isInteractive {
			reader := bufio.NewReader(os.Stdin)
			for _, add := range adds {
				if matchedAdds[add] {
					continue
				}
				if !typesCompatible(drop.C, add.C) {
					continue
				}
				fmt.Printf("[WARNING]: Destructive change: dropping column %q and adding column %q in table %q.\n",
					drop.C.Name, add.C.Name, tableName)
				fmt.Printf("   Are you renaming column %q to %q? [y/N]: ", drop.C.Name, add.C.Name)

				response, err := reader.ReadString('\n')
				if err == nil {
					response = strings.TrimSpace(strings.ToLower(response))
					if response == "y" || response == "yes" {
						return add
					}
				}
			}
		} else {
			for _, add := range adds {
				if !typesCompatible(drop.C, add.C) {
					continue
				}
				fmt.Printf("[WARNING]: Destructive change: dropping column %q and adding column %q in table %q (skipped rename detection in non-interactive mode).\n",
					drop.C.Name, add.C.Name, tableName)
			}
		}
		return nil
	}, func(tableName, newName, oldName string) {
		if p.renamedMap[tableName] == nil {
			p.renamedMap[tableName] = make(map[string]string)
		}
		p.renamedMap[tableName][newName] = oldName
	})
}

func (p *diffPlanner) applyDownRenames(changes []schema.Change) {
	processRenames(changes, func(tableName string, drop *schema.DropColumn, adds []*schema.AddColumn, matchedAdds map[*schema.AddColumn]bool) *schema.AddColumn {
		renameTargetMap, ok := p.renamedMap[tableName]
		if !ok {
			return nil
		}
		// in down, we dropping the NEW name and adding the OLD name
		oldName, ok := renameTargetMap[drop.C.Name]
		if !ok {
			return nil
		}
		for _, add := range adds {
			if add.C.Name == oldName && !matchedAdds[add] {
				return add
			}
		}
		return nil
	}, nil)
}
