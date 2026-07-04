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

// User represents the database model
type User struct {
	Id           string       `db:"id" json:"id"`
	Email        string       `db:"email" json:"email"`
	PhoneNum     string       `db:"phoneNum" json:"phoneNum"`
	Role         UserRoleType `db:"role" json:"role"`
	ReferredById *string      `db:"referredById" json:"referredById"`
	Profile      *Profile     `json:"profile,omitempty"`
	Posts        []*Post      `json:"posts,omitempty"`
	Comments     []*Comment   `json:"comments,omitempty"`
	ReferredBy   *User        `json:"referredBy,omitempty"`
	Referrals    []*User      `json:"referrals,omitempty"`
}

// UserCreateInput represents the input structure for creation
type UserCreateInput struct {
	Id           *string       `json:"id"`
	Email        string        `json:"email"`
	PhoneNum     string        `json:"phoneNum"`
	Role         *UserRoleType `json:"role"`
	ReferredById *string       `json:"referredById"`
}

// UserSelect specifies which fields to include
type UserSelect struct {
	Id           bool           `json:"id"`
	Email        bool           `json:"email"`
	PhoneNum     bool           `json:"phoneNum"`
	Role         bool           `json:"role"`
	ReferredById bool           `json:"referredById"`
	Profile      *ProfileSelect `json:"profile,omitempty"`
	Posts        *PostSelect    `json:"posts,omitempty"`
	Comments     *CommentSelect `json:"comments,omitempty"`
	ReferredBy   *UserSelect    `json:"referredBy,omitempty"`
	Referrals    *UserSelect    `json:"referrals,omitempty"`
}

// UserOmit specifies which fields to exclude
type UserOmit struct {
	Id           bool         `json:"id"`
	Email        bool         `json:"email"`
	PhoneNum     bool         `json:"phoneNum"`
	Role         bool         `json:"role"`
	ReferredById bool         `json:"referredById"`
	Profile      *ProfileOmit `json:"profile,omitempty"`
	Posts        *PostOmit    `json:"posts,omitempty"`
	Comments     *CommentOmit `json:"comments,omitempty"`
	ReferredBy   *UserOmit    `json:"referredBy,omitempty"`
	Referrals    *UserOmit    `json:"referrals,omitempty"`
}

type UserDelegate struct {
	client *Queries
}

func (m *User) ScanFields(cols []string) []any {
	targets := make([]any, len(cols))
	for i, col := range cols {
		switch col {
		case "id":
			targets[i] = &m.Id
		case "email":
			targets[i] = &m.Email
		case "phoneNum":
			targets[i] = &m.PhoneNum
		case "role":
			targets[i] = &m.Role
		case "referredById":
			targets[i] = &m.ReferredById
		}
	}
	return targets
}

var userDefaultCols = []string{
	"id",
	"email",
	"phoneNum",
	"role",
	"referredById",
}

func (q *Queries) selectUserCols(selects *UserSelect, omits *UserOmit, forceCols ...string) []string {
	if selects == nil && omits == nil && len(forceCols) == 0 {
		return userDefaultCols
	}

	anySelected := selects != nil && (selects.Id || selects.Email || selects.PhoneNum || selects.Role || selects.ReferredById || selects.Profile != nil || selects.Posts != nil || selects.Comments != nil || selects.ReferredBy != nil || selects.Referrals != nil)

	specs := []colSpec{
		{"id", selects != nil && selects.Id, omits != nil && omits.Id, false},
		{"email", selects != nil && selects.Email, omits != nil && omits.Email, false},
		{"phoneNum", selects != nil && selects.PhoneNum, omits != nil && omits.PhoneNum, false},
		{"role", selects != nil && selects.Role, omits != nil && omits.Role, false},
		{"referredById", selects != nil && selects.ReferredById, omits != nil && omits.ReferredById, selects != nil && selects.ReferredBy != nil},
	}

	cols := computeCols(specs, selects != nil, anySelected)

	for _, f := range forceCols {
		if !slices.Contains(cols, f) {
			cols = append(cols, f)
		}
	}

	return cols
}
func (d *UserDelegate) Create(input UserCreateInput) *CreateBuilder[User, UserCreateInput, UserSelect, UserOmit] {
	return &CreateBuilder[User, UserCreateInput, UserSelect, UserOmit]{
		client:   d.client,
		input:    input,
		execFunc: d.client.executeUserCreate,
	}
}

func (q *Queries) executeUserCreate(ctx context.Context, input UserCreateInput, selects *UserSelect, omits *UserOmit) (*User, error) {
	var cols []string
	var vals []any
	if input.Id != nil {
		cols = append(cols, q.dialect.Quote("id"))
		vals = append(vals, *input.Id)
	} else {
		cols = append(cols, q.dialect.Quote("id"))
		vals = append(vals, generateCUID())
	}
	cols = append(cols, q.dialect.Quote("email"))
	vals = append(vals, input.Email)
	cols = append(cols, q.dialect.Quote("phoneNum"))
	vals = append(vals, input.PhoneNum)
	if input.Role != nil {
		cols = append(cols, q.dialect.Quote("role"))
		vals = append(vals, *input.Role)
	}
	if input.ReferredById != nil {
		cols = append(cols, q.dialect.Quote("referredById"))
		vals = append(vals, *input.ReferredById)
	}

	returningCols := q.selectUserCols(selects, omits)

	scanFunc := func(res *User, cols []string) []any {
		return res.ScanFields(cols)
	}

	idCol := "id"

	res, err := executeInsert(ctx, q, "User", cols, vals, returningCols, idCol, scanFunc)
	if err != nil {
		return nil, err
	}

	if err := q.loadUserRelations(ctx, []*User{res}, selects); err != nil {
		return nil, err
	}

	return res, nil
}
func (q *Queries) loadUserRelations(ctx context.Context, records []*User, selects *UserSelect) error {
	if selects == nil || len(records) == 0 {
		return nil
	}
	if selects.Profile != nil {
		returningCols := q.selectProfileCols(selects.Profile, nil, "userId")
		// Inverse holds the FK: Profile.userId
		allChildren, err := loadRelation(
			ctx, q, records,
			directKey(func(p *User) string { return p.Id }),
			"Profile",
			"userId",
			returningCols,
			scanInto(returningCols, (*Profile).ScanFields),
			directKey(func(c *Profile) string { return c.UserId }),
			setOne(func(p *User, c *Profile) { p.Profile = c }),
		)
		if err != nil {
			return fmt.Errorf("loading profile: %w", err)
		}
		if err := q.loadProfileRelations(ctx, allChildren, selects.Profile); err != nil {
			return err
		}
	}
	if selects.Posts != nil {
		returningCols := q.selectPostCols(selects.Posts, nil, "authorId")
		// Inverse holds the FK: Post.authorId
		allChildren, err := loadRelation(
			ctx, q, records,
			directKey(func(p *User) string { return p.Id }),
			"Post",
			"authorId",
			returningCols,
			scanInto(returningCols, (*Post).ScanFields),
			directKey(func(c *Post) string { return c.AuthorId }),
			appendMany(func(p *User) *[]*Post { return &p.Posts }),
		)
		if err != nil {
			return fmt.Errorf("loading posts: %w", err)
		}
		if err := q.loadPostRelations(ctx, allChildren, selects.Posts); err != nil {
			return err
		}
	}
	if selects.Comments != nil {
		returningCols := q.selectCommentCols(selects.Comments, nil, "authorId")
		// Inverse holds the FK: Comment.authorId
		allChildren, err := loadRelation(
			ctx, q, records,
			directKey(func(p *User) string { return p.Id }),
			"Comment",
			"authorId",
			returningCols,
			scanInto(returningCols, (*Comment).ScanFields),
			directKey(func(c *Comment) string { return c.AuthorId }),
			appendMany(func(p *User) *[]*Comment { return &p.Comments }),
		)
		if err != nil {
			return fmt.Errorf("loading comments: %w", err)
		}
		if err := q.loadCommentRelations(ctx, allChildren, selects.Comments); err != nil {
			return err
		}
	}
	if selects.ReferredBy != nil {
		returningCols := q.selectUserCols(selects.ReferredBy, nil, "id")
		// Current model holds the FK: User.referredById
		allChildren, err := loadRelation(
			ctx, q, records,
			optionalKey(func(p *User) *string { return p.ReferredById }),
			"User",
			"id",
			returningCols,
			scanInto(returningCols, (*User).ScanFields),
			directKey(func(c *User) string { return c.Id }),
			setOne(func(p *User, c *User) { p.ReferredBy = c }),
		)
		if err != nil {
			return fmt.Errorf("loading referredBy: %w", err)
		}
		if err := q.loadUserRelations(ctx, allChildren, selects.ReferredBy); err != nil {
			return err
		}
	}
	if selects.Referrals != nil {
		returningCols := q.selectUserCols(selects.Referrals, nil, "referredById")
		// Inverse holds the FK: User.referredById
		allChildren, err := loadRelation(
			ctx, q, records,
			directKey(func(p *User) string { return p.Id }),
			"User",
			"referredById",
			returningCols,
			scanInto(returningCols, (*User).ScanFields),
			optionalKey(func(c *User) *string { return c.ReferredById }),
			appendMany(func(p *User) *[]*User { return &p.Referrals }),
		)
		if err != nil {
			return fmt.Errorf("loading referrals: %w", err)
		}
		if err := q.loadUserRelations(ctx, allChildren, selects.Referrals); err != nil {
			return err
		}
	}

	return nil
}
