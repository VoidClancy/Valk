package valkyrie

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"strings"
	"time"
)

var _ = time.Time{}
var _ = fmt.Sprintf
var _ = strings.Join
var _ = context.Background
var _ = sql.LevelDefault
var _ = slices.Contains[[]string, string]

// Profile represents the database model
type Profile struct {
	Id     string  `db:"id" json:"id"`
	Bio    *string `db:"bio" json:"bio"`
	UserId string  `db:"userId" json:"userId"`
	User   *User   `json:"user,omitempty"`
}

// ProfileCreateInput represents the input structure for creation
type ProfileCreateInput struct {
	Id     *string `json:"id"`
	Bio    *string `json:"bio"`
	UserId string  `json:"userId"`
}

// ProfileSelect specifies which fields to include
type ProfileSelect struct {
	Id     bool        `json:"id"`
	Bio    bool        `json:"bio"`
	UserId bool        `json:"userId"`
	User   *UserSelect `json:"user,omitempty"`
}

// ProfileOmit specifies which fields to exclude
type ProfileOmit struct {
	Id     bool      `json:"id"`
	Bio    bool      `json:"bio"`
	UserId bool      `json:"userId"`
	User   *UserOmit `json:"user,omitempty"`
}

type ProfileDelegate struct {
	client *Queries
}

func (m *Profile) ScanFields(cols []string) []any {
	targets := make([]any, len(cols))
	for i, col := range cols {
		switch col {
		case "id":
			targets[i] = &m.Id
		case "bio":
			targets[i] = &m.Bio
		case "userId":
			targets[i] = &m.UserId
		}
	}
	return targets
}

var profileDefaultCols = []string{
	"id",
	"bio",
	"userId",
}

func (q *Queries) selectProfileCols(selects *ProfileSelect, omits *ProfileOmit, forceCols ...string) []string {
	if selects == nil && omits == nil && len(forceCols) == 0 {
		return profileDefaultCols
	}

	var cols []string
	var anySelected bool
	if selects != nil {
		if selects.Id {
			anySelected = true
		}
		if selects.Bio {
			anySelected = true
		}
		if selects.UserId {
			anySelected = true
		}
		if selects.User != nil {
			anySelected = true
		}
	}
	{
		include := true
		if selects != nil {
			include = false
			if !anySelected {
				include = true
			} else if selects.Id {
				include = true
			}
		} else if omits != nil {
			if omits.Id {
				include = false
			}
		}
		if include {
			cols = append(cols, "id")
		}
	}
	{
		include := true
		if selects != nil {
			include = false
			if !anySelected {
				include = true
			} else if selects.Bio {
				include = true
			}
		} else if omits != nil {
			if omits.Bio {
				include = false
			}
		}
		if include {
			cols = append(cols, "bio")
		}
	}
	{
		include := true
		if selects != nil {
			include = false
			if !anySelected {
				include = true
			} else if selects.UserId {
				include = true
			}
			// Force-include FK when its relation is selected
			if selects.User != nil {
				include = true
			}
		} else if omits != nil {
			if omits.UserId {
				include = false
			}
		}
		if include {
			cols = append(cols, "userId")
		}
	}

	if len(cols) == 0 {
		cols = append(cols, "id")
		cols = append(cols, "bio")
		cols = append(cols, "userId")
	}

	// Force-include any requested columns
	for _, f := range forceCols {
		if !slices.Contains(cols, f) {
			cols = append(cols, f)
		}
	}

	return cols
}
func (d *ProfileDelegate) Create(input ProfileCreateInput) *CreateBuilder[Profile, ProfileCreateInput, ProfileSelect, ProfileOmit] {
	return &CreateBuilder[Profile, ProfileCreateInput, ProfileSelect, ProfileOmit]{
		client:   d.client,
		input:    input,
		execFunc: d.client.executeProfileCreate,
	}
}

func (q *Queries) executeProfileCreate(ctx context.Context, input ProfileCreateInput, selects *ProfileSelect, omits *ProfileOmit) (*Profile, error) {
	var cols []string
	var vals []any
	if input.Id != nil {
		cols = append(cols, q.dialect.Quote("id"))
		vals = append(vals, *input.Id)
	} else {
		cols = append(cols, q.dialect.Quote("id"))
		vals = append(vals, generateCUID())
	}
	if input.Bio != nil {
		cols = append(cols, q.dialect.Quote("bio"))
		vals = append(vals, *input.Bio)
	}
	cols = append(cols, q.dialect.Quote("userId"))
	vals = append(vals, input.UserId)

	returningCols := q.selectProfileCols(selects, omits)

	scanFunc := func(res *Profile, cols []string) []any {
		return res.ScanFields(cols)
	}

	idCol := "id"

	res, err := executeInsert(ctx, q, "Profile", cols, vals, returningCols, idCol, scanFunc)
	if err != nil {
		return nil, err
	}

	if err := q.loadProfileRelations(ctx, []*Profile{res}, selects); err != nil {
		return nil, err
	}

	return res, nil
}
func (q *Queries) loadProfileRelations(ctx context.Context, records []*Profile, selects *ProfileSelect) error {
	if selects == nil || len(records) == 0 {
		return nil
	}
	if selects.User != nil {
		returningCols := q.selectUserCols(selects.User, nil, "id")
		// Current model holds the FK: Profile.userId
		allChildren, err := loadRelation(
			ctx, q, records,
			func(p *Profile) (string, bool) {
				return fmt.Sprint(p.UserId), true
			},
			"User",
			"id",
			returningCols,
			func(rows *sql.Rows, child *User) error {
				return rows.Scan(child.ScanFields(returningCols)...)
			},
			func(child *User) (string, bool) {
				return fmt.Sprint(child.Id), true
			},
			func(p *Profile, children []*User) {
				if len(children) > 0 {
					p.User = children[0]
				}
			},
		)
		if err != nil {
			return fmt.Errorf("loading user: %w", err)
		}
		if err := q.loadUserRelations(ctx, allChildren, selects.User); err != nil {
			return err
		}
	}

	return nil
}
