package valk

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"time"
)

// AllFieldsSoFar represents the database model
type AllFieldsSoFar struct {
	Id              int32               `db:"id" json:"id"`
	StringReq       string              `db:"stringReq" json:"stringReq"`
	StringOpt       *string             `db:"stringOpt" json:"stringOpt,omitempty"`
	StringDefault   string              `db:"stringDefault" json:"stringDefault"`
	StringVarchar   string              `db:"stringVarchar" json:"stringVarchar"`
	StringChar      string              `db:"stringChar" json:"stringChar"`
	BitVal          string              `db:"bitVal" json:"bitVal"`
	VarBitVal       string              `db:"varBitVal" json:"varBitVal"`
	InetVal         string              `db:"inetVal" json:"inetVal"`
	XmlVal          string              `db:"xmlVal" json:"xmlVal"`
	CuidDefault     string              `db:"cuidDefault" json:"cuidDefault"`
	Cuid1Default    string              `db:"cuid1Default" json:"cuid1Default"`
	Cuid2Default    string              `db:"cuid2Default" json:"cuid2Default"`
	UuidDefault     string              `db:"uuidDefault" json:"uuidDefault"`
	Uuid4Default    string              `db:"uuid4Default" json:"uuid4Default"`
	Uuid7Default    string              `db:"uuid7Default" json:"uuid7Default"`
	UlidDefault     string              `db:"ulidDefault" json:"ulidDefault"`
	NanoidDefault   string              `db:"nanoidDefault" json:"nanoidDefault"`
	UuidDb          string              `db:"uuidDb" json:"uuidDb"`
	IntReq          int32               `db:"intReq" json:"intReq"`
	IntOpt          *int32              `db:"intOpt" json:"intOpt,omitempty"`
	IntDefault      int32               `db:"intDefault" json:"intDefault"`
	IntegerVal      int32               `db:"integerVal" json:"integerVal"`
	SmallInt        int32               `db:"smallInt" json:"smallInt"`
	TinyInt         int32               `db:"tinyInt" json:"tinyInt"`
	OidVal          int32               `db:"oidVal" json:"oidVal"`
	BigIntReq       int64               `db:"bigIntReq" json:"bigIntReq"`
	BigIntOpt       *int64              `db:"bigIntOpt" json:"bigIntOpt,omitempty"`
	FloatReq        float64             `db:"floatReq" json:"floatReq"`
	FloatOpt        *float64            `db:"floatOpt" json:"floatOpt,omitempty"`
	RealVal         float64             `db:"realVal" json:"realVal"`
	DecimalReq      string              `db:"decimalReq" json:"decimalReq"`
	DecimalOpt      *string             `db:"decimalOpt" json:"decimalOpt,omitempty"`
	DecimalPrecise  string              `db:"decimalPrecise" json:"decimalPrecise"`
	MoneyVal        string              `db:"moneyVal" json:"moneyVal"`
	BoolReq         bool                `db:"boolReq" json:"boolReq"`
	BoolOpt         *bool               `db:"boolOpt" json:"boolOpt,omitempty"`
	BoolDefault     bool                `db:"boolDefault" json:"boolDefault"`
	DateTimeReq     time.Time           `db:"dateTimeReq" json:"dateTimeReq"`
	DateTimeOpt     *time.Time          `db:"dateTimeOpt" json:"dateTimeOpt,omitempty"`
	DateTimeDefault time.Time           `db:"dateTimeDefault" json:"dateTimeDefault"`
	UpdatedAt       time.Time           `db:"updatedAt" json:"updatedAt"`
	DateTimeTz      time.Time           `db:"dateTimeTz" json:"dateTimeTz"`
	TimestampVal    time.Time           `db:"timestampVal" json:"timestampVal"`
	TimeVal         time.Time           `db:"timeVal" json:"timeVal"`
	TimetzVal       time.Time           `db:"timetzVal" json:"timetzVal"`
	JsonReq         json.RawMessage     `db:"jsonReq" json:"jsonReq"`
	JsonOpt         *json.RawMessage    `db:"jsonOpt" json:"jsonOpt,omitempty"`
	JsonVal         json.RawMessage     `db:"jsonVal" json:"jsonVal"`
	BytesReq        []byte              `db:"bytesReq" json:"bytesReq"`
	BytesOpt        *[]byte             `db:"bytesOpt" json:"bytesOpt,omitempty"`
	HstoreField     *map[string]*string `db:"hstoreField" json:"hstoreField,omitempty"`
	LtreeField      string              `db:"ltreeField" json:"ltreeField"`
	CitextField     *string             `db:"citextField" json:"citextField,omitempty"`
}

// AllFieldsSoFarCreate is used for hooks only — the Create API uses FieldAssignment
type AllFieldsSoFarCreate struct {
	Id              *int32              `json:"id"`
	StringReq       string              `json:"stringReq"`
	StringOpt       *string             `json:"stringOpt"`
	StringDefault   *string             `json:"stringDefault"`
	StringVarchar   string              `json:"stringVarchar"`
	StringChar      string              `json:"stringChar"`
	BitVal          string              `json:"bitVal"`
	VarBitVal       string              `json:"varBitVal"`
	InetVal         string              `json:"inetVal"`
	XmlVal          string              `json:"xmlVal"`
	CuidDefault     *string             `json:"cuidDefault"`
	Cuid1Default    *string             `json:"cuid1Default"`
	Cuid2Default    *string             `json:"cuid2Default"`
	UuidDefault     *string             `json:"uuidDefault"`
	Uuid4Default    *string             `json:"uuid4Default"`
	Uuid7Default    *string             `json:"uuid7Default"`
	UlidDefault     *string             `json:"ulidDefault"`
	NanoidDefault   *string             `json:"nanoidDefault"`
	UuidDb          string              `json:"uuidDb"`
	IntReq          int32               `json:"intReq"`
	IntOpt          *int32              `json:"intOpt"`
	IntDefault      *int32              `json:"intDefault"`
	IntegerVal      int32               `json:"integerVal"`
	SmallInt        int32               `json:"smallInt"`
	TinyInt         int32               `json:"tinyInt"`
	OidVal          int32               `json:"oidVal"`
	BigIntReq       int64               `json:"bigIntReq"`
	BigIntOpt       *int64              `json:"bigIntOpt"`
	FloatReq        float64             `json:"floatReq"`
	FloatOpt        *float64            `json:"floatOpt"`
	RealVal         float64             `json:"realVal"`
	DecimalReq      string              `json:"decimalReq"`
	DecimalOpt      *string             `json:"decimalOpt"`
	DecimalPrecise  string              `json:"decimalPrecise"`
	MoneyVal        string              `json:"moneyVal"`
	BoolReq         bool                `json:"boolReq"`
	BoolOpt         *bool               `json:"boolOpt"`
	BoolDefault     *bool               `json:"boolDefault"`
	DateTimeReq     time.Time           `json:"dateTimeReq"`
	DateTimeOpt     *time.Time          `json:"dateTimeOpt"`
	DateTimeDefault *time.Time          `json:"dateTimeDefault"`
	UpdatedAt       time.Time           `json:"updatedAt"`
	DateTimeTz      time.Time           `json:"dateTimeTz"`
	TimestampVal    time.Time           `json:"timestampVal"`
	TimeVal         time.Time           `json:"timeVal"`
	TimetzVal       time.Time           `json:"timetzVal"`
	JsonReq         json.RawMessage     `json:"jsonReq"`
	JsonOpt         *json.RawMessage    `json:"jsonOpt"`
	JsonVal         json.RawMessage     `json:"jsonVal"`
	BytesReq        []byte              `json:"bytesReq"`
	BytesOpt        *[]byte             `json:"bytesOpt"`
	HstoreField     *map[string]*string `json:"hstoreField"`
	LtreeField      string              `json:"ltreeField"`
	CitextField     *string             `json:"citextField"`
}

// colMask returns a bit mask of columns that are set
func (s *AllFieldsSoFarCreate) colMask() uint64 {
	var mask uint64
	if s.Id != nil {
		mask |= 1 << 0
	}
	mask |= 1 << 1
	if s.StringOpt != nil {
		mask |= 1 << 2
	}
	if s.StringDefault != nil {
		mask |= 1 << 3
	}
	mask |= 1 << 4
	mask |= 1 << 5
	mask |= 1 << 6
	mask |= 1 << 7
	mask |= 1 << 8
	mask |= 1 << 9
	mask |= 1 << 10
	mask |= 1 << 11
	mask |= 1 << 12
	mask |= 1 << 13
	mask |= 1 << 14
	mask |= 1 << 15
	mask |= 1 << 16
	mask |= 1 << 17
	mask |= 1 << 18
	mask |= 1 << 19
	if s.IntOpt != nil {
		mask |= 1 << 20
	}
	if s.IntDefault != nil {
		mask |= 1 << 21
	}
	mask |= 1 << 22
	mask |= 1 << 23
	mask |= 1 << 24
	mask |= 1 << 25
	mask |= 1 << 26
	if s.BigIntOpt != nil {
		mask |= 1 << 27
	}
	mask |= 1 << 28
	if s.FloatOpt != nil {
		mask |= 1 << 29
	}
	mask |= 1 << 30
	mask |= 1 << 31
	if s.DecimalOpt != nil {
		mask |= 1 << 32
	}
	mask |= 1 << 33
	mask |= 1 << 34
	mask |= 1 << 35
	if s.BoolOpt != nil {
		mask |= 1 << 36
	}
	if s.BoolDefault != nil {
		mask |= 1 << 37
	}
	mask |= 1 << 38
	if s.DateTimeOpt != nil {
		mask |= 1 << 39
	}
	mask |= 1 << 40
	mask |= 1 << 41
	mask |= 1 << 42
	mask |= 1 << 43
	mask |= 1 << 44
	mask |= 1 << 45
	mask |= 1 << 46
	if s.JsonOpt != nil {
		mask |= 1 << 47
	}
	mask |= 1 << 48
	mask |= 1 << 49
	if s.BytesOpt != nil {
		mask |= 1 << 50
	}
	if s.HstoreField != nil {
		mask |= 1 << 51
	}
	mask |= 1 << 52
	if s.CitextField != nil {
		mask |= 1 << 53
	}
	return mask
}

type AllFieldsSoFarSelect struct {
	Id              bool `json:"id"`
	StringReq       bool `json:"stringReq"`
	StringOpt       bool `json:"stringOpt"`
	StringDefault   bool `json:"stringDefault"`
	StringVarchar   bool `json:"stringVarchar"`
	StringChar      bool `json:"stringChar"`
	BitVal          bool `json:"bitVal"`
	VarBitVal       bool `json:"varBitVal"`
	InetVal         bool `json:"inetVal"`
	XmlVal          bool `json:"xmlVal"`
	CuidDefault     bool `json:"cuidDefault"`
	Cuid1Default    bool `json:"cuid1Default"`
	Cuid2Default    bool `json:"cuid2Default"`
	UuidDefault     bool `json:"uuidDefault"`
	Uuid4Default    bool `json:"uuid4Default"`
	Uuid7Default    bool `json:"uuid7Default"`
	UlidDefault     bool `json:"ulidDefault"`
	NanoidDefault   bool `json:"nanoidDefault"`
	UuidDb          bool `json:"uuidDb"`
	IntReq          bool `json:"intReq"`
	IntOpt          bool `json:"intOpt"`
	IntDefault      bool `json:"intDefault"`
	IntegerVal      bool `json:"integerVal"`
	SmallInt        bool `json:"smallInt"`
	TinyInt         bool `json:"tinyInt"`
	OidVal          bool `json:"oidVal"`
	BigIntReq       bool `json:"bigIntReq"`
	BigIntOpt       bool `json:"bigIntOpt"`
	FloatReq        bool `json:"floatReq"`
	FloatOpt        bool `json:"floatOpt"`
	RealVal         bool `json:"realVal"`
	DecimalReq      bool `json:"decimalReq"`
	DecimalOpt      bool `json:"decimalOpt"`
	DecimalPrecise  bool `json:"decimalPrecise"`
	MoneyVal        bool `json:"moneyVal"`
	BoolReq         bool `json:"boolReq"`
	BoolOpt         bool `json:"boolOpt"`
	BoolDefault     bool `json:"boolDefault"`
	DateTimeReq     bool `json:"dateTimeReq"`
	DateTimeOpt     bool `json:"dateTimeOpt"`
	DateTimeDefault bool `json:"dateTimeDefault"`
	UpdatedAt       bool `json:"updatedAt"`
	DateTimeTz      bool `json:"dateTimeTz"`
	TimestampVal    bool `json:"timestampVal"`
	TimeVal         bool `json:"timeVal"`
	TimetzVal       bool `json:"timetzVal"`
	JsonReq         bool `json:"jsonReq"`
	JsonOpt         bool `json:"jsonOpt"`
	JsonVal         bool `json:"jsonVal"`
	BytesReq        bool `json:"bytesReq"`
	BytesOpt        bool `json:"bytesOpt"`
	HstoreField     bool `json:"hstoreField"`
	LtreeField      bool `json:"ltreeField"`
	CitextField     bool `json:"citextField"`
}

type AllFieldsSoFarOmit struct {
	Id              bool `json:"id"`
	StringReq       bool `json:"stringReq"`
	StringOpt       bool `json:"stringOpt"`
	StringDefault   bool `json:"stringDefault"`
	StringVarchar   bool `json:"stringVarchar"`
	StringChar      bool `json:"stringChar"`
	BitVal          bool `json:"bitVal"`
	VarBitVal       bool `json:"varBitVal"`
	InetVal         bool `json:"inetVal"`
	XmlVal          bool `json:"xmlVal"`
	CuidDefault     bool `json:"cuidDefault"`
	Cuid1Default    bool `json:"cuid1Default"`
	Cuid2Default    bool `json:"cuid2Default"`
	UuidDefault     bool `json:"uuidDefault"`
	Uuid4Default    bool `json:"uuid4Default"`
	Uuid7Default    bool `json:"uuid7Default"`
	UlidDefault     bool `json:"ulidDefault"`
	NanoidDefault   bool `json:"nanoidDefault"`
	UuidDb          bool `json:"uuidDb"`
	IntReq          bool `json:"intReq"`
	IntOpt          bool `json:"intOpt"`
	IntDefault      bool `json:"intDefault"`
	IntegerVal      bool `json:"integerVal"`
	SmallInt        bool `json:"smallInt"`
	TinyInt         bool `json:"tinyInt"`
	OidVal          bool `json:"oidVal"`
	BigIntReq       bool `json:"bigIntReq"`
	BigIntOpt       bool `json:"bigIntOpt"`
	FloatReq        bool `json:"floatReq"`
	FloatOpt        bool `json:"floatOpt"`
	RealVal         bool `json:"realVal"`
	DecimalReq      bool `json:"decimalReq"`
	DecimalOpt      bool `json:"decimalOpt"`
	DecimalPrecise  bool `json:"decimalPrecise"`
	MoneyVal        bool `json:"moneyVal"`
	BoolReq         bool `json:"boolReq"`
	BoolOpt         bool `json:"boolOpt"`
	BoolDefault     bool `json:"boolDefault"`
	DateTimeReq     bool `json:"dateTimeReq"`
	DateTimeOpt     bool `json:"dateTimeOpt"`
	DateTimeDefault bool `json:"dateTimeDefault"`
	UpdatedAt       bool `json:"updatedAt"`
	DateTimeTz      bool `json:"dateTimeTz"`
	TimestampVal    bool `json:"timestampVal"`
	TimeVal         bool `json:"timeVal"`
	TimetzVal       bool `json:"timetzVal"`
	JsonReq         bool `json:"jsonReq"`
	JsonOpt         bool `json:"jsonOpt"`
	JsonVal         bool `json:"jsonVal"`
	BytesReq        bool `json:"bytesReq"`
	BytesOpt        bool `json:"bytesOpt"`
	HstoreField     bool `json:"hstoreField"`
	LtreeField      bool `json:"ltreeField"`
	CitextField     bool `json:"citextField"`
}

type AllFieldsSoFarSelectQuery interface {
	GetRelationParams() (*AllFieldsSoFarSelect, *AllFieldsSoFarOmit, QueryParams[AllFieldsSoFar])
}

func (s *AllFieldsSoFarSelect) GetRelationParams() (*AllFieldsSoFarSelect, *AllFieldsSoFarOmit, QueryParams[AllFieldsSoFar]) {
	return s, nil, QueryParams[AllFieldsSoFar]{}
}

type AllFieldsSoFarQueryBuilder struct {
	selects *AllFieldsSoFarSelect
	omits   *AllFieldsSoFarOmit
	where   []PredicateOf[AllFieldsSoFar]
	take    *int
	skip    *int
	orderBy []OrderBy[AllFieldsSoFar]
	cursor  UniquePredicate[AllFieldsSoFar]
}

func (b *AllFieldsSoFarQueryBuilder) Where(preds ...PredicateOf[AllFieldsSoFar]) *AllFieldsSoFarQueryBuilder {
	b.where = append(b.where, preds...)
	return b
}

func (b *AllFieldsSoFarQueryBuilder) Take(limit int) *AllFieldsSoFarQueryBuilder {
	b.take = &limit
	return b
}

func (b *AllFieldsSoFarQueryBuilder) Skip(offset int) *AllFieldsSoFarQueryBuilder {
	b.skip = &offset
	return b
}

func (b *AllFieldsSoFarQueryBuilder) OrderBy(orders ...OrderBy[AllFieldsSoFar]) *AllFieldsSoFarQueryBuilder {
	b.orderBy = append(b.orderBy, orders...)
	return b
}

func (b *AllFieldsSoFarQueryBuilder) Cursor(where UniquePredicate[AllFieldsSoFar]) *AllFieldsSoFarQueryBuilder {
	b.cursor = where
	return b
}

func (b *AllFieldsSoFarQueryBuilder) Select(s AllFieldsSoFarSelect) *AllFieldsSoFarQueryBuilder {
	b.selects = &s
	return b
}

func (b *AllFieldsSoFarQueryBuilder) Omit(o AllFieldsSoFarOmit) *AllFieldsSoFarQueryBuilder {
	b.omits = &o
	return b
}

func (b *AllFieldsSoFarQueryBuilder) GetRelationParams() (*AllFieldsSoFarSelect, *AllFieldsSoFarOmit, QueryParams[AllFieldsSoFar]) {
	if b == nil {
		return nil, nil, QueryParams[AllFieldsSoFar]{}
	}
	return b.selects, b.omits, QueryParams[AllFieldsSoFar]{
		Where:   b.where,
		Take:    b.take,
		Skip:    b.skip,
		OrderBy: b.orderBy,
		Cursor:  b.cursor,
	}
}

type AllFieldsSoFarCreateQuery = func(ctx context.Context, args *AllFieldsSoFarCreate) (*AllFieldsSoFar, error)
type AllFieldsSoFarCreateManyQuery = func(ctx context.Context, args []*AllFieldsSoFarCreate) (int64, error)
type AllFieldsSoFarCreateManyAndReturnQuery = func(ctx context.Context, args []*AllFieldsSoFarCreate) ([]*AllFieldsSoFar, error)
type AllFieldsSoFarFindUniqueQuery = func(ctx context.Context, where UniquePredicate[AllFieldsSoFar], additional []PredicateOf[AllFieldsSoFar], selects *AllFieldsSoFarSelect, omits *AllFieldsSoFarOmit) (*AllFieldsSoFar, error)
type AllFieldsSoFarFindFirstQuery = func(ctx context.Context, params QueryParams[AllFieldsSoFar], selects *AllFieldsSoFarSelect, omits *AllFieldsSoFarOmit) (*AllFieldsSoFar, error)
type AllFieldsSoFarFindManyQuery = func(ctx context.Context, params QueryParams[AllFieldsSoFar], selects *AllFieldsSoFarSelect, omits *AllFieldsSoFarOmit) ([]*AllFieldsSoFar, error)
type AllFieldsSoFarDeleteManyQuery = func(ctx context.Context, preds []PredicateOf[AllFieldsSoFar]) (int64, error)
type AllFieldsSoFarDeleteQuery = func(ctx context.Context, where UniquePredicate[AllFieldsSoFar], selects *AllFieldsSoFarSelect, omits *AllFieldsSoFarOmit) (*AllFieldsSoFar, error)
type AllFieldsSoFarCountQuery = func(ctx context.Context, params QueryParams[AllFieldsSoFar]) (int64, error)

type AllFieldsSoFarExtension struct {
	Create              func(ctx context.Context, input *AllFieldsSoFarCreate, next AllFieldsSoFarCreateQuery) (*AllFieldsSoFar, error)
	CreateMany          func(ctx context.Context, inputs []*AllFieldsSoFarCreate, next AllFieldsSoFarCreateManyQuery) (int64, error)
	CreateManyAndReturn func(ctx context.Context, inputs []*AllFieldsSoFarCreate, next AllFieldsSoFarCreateManyAndReturnQuery) ([]*AllFieldsSoFar, error)
	FindUnique          func(ctx context.Context, where UniquePredicate[AllFieldsSoFar], additional []PredicateOf[AllFieldsSoFar], selects *AllFieldsSoFarSelect, omits *AllFieldsSoFarOmit, next AllFieldsSoFarFindUniqueQuery) (*AllFieldsSoFar, error)
	FindFirst           func(ctx context.Context, params QueryParams[AllFieldsSoFar], selects *AllFieldsSoFarSelect, omits *AllFieldsSoFarOmit, next AllFieldsSoFarFindFirstQuery) (*AllFieldsSoFar, error)
	FindMany            func(ctx context.Context, params QueryParams[AllFieldsSoFar], selects *AllFieldsSoFarSelect, omits *AllFieldsSoFarOmit, next AllFieldsSoFarFindManyQuery) ([]*AllFieldsSoFar, error)
	DeleteMany          func(ctx context.Context, preds []PredicateOf[AllFieldsSoFar], next AllFieldsSoFarDeleteManyQuery) (int64, error)
	Delete              func(ctx context.Context, where UniquePredicate[AllFieldsSoFar], selects *AllFieldsSoFarSelect, omits *AllFieldsSoFarOmit, next AllFieldsSoFarDeleteQuery) (*AllFieldsSoFar, error)
	Count               func(ctx context.Context, params QueryParams[AllFieldsSoFar], next AllFieldsSoFarCountQuery) (int64, error)
}

type AllFieldsSoFarDelegate struct {
	client     *Queries
	extensions []AllFieldsSoFarExtension
}

func (d *AllFieldsSoFarDelegate) Use(exts ...AllFieldsSoFarExtension) {
	d.extensions = append(d.extensions, exts...)
}

func (m *AllFieldsSoFar) ScanFields(cols []string) []any {
	targets := make([]any, len(cols))
	for i, col := range cols {
		switch col {
		case "id":
			targets[i] = &m.Id
		case "stringReq":
			targets[i] = &m.StringReq
		case "stringOpt":
			targets[i] = &m.StringOpt
		case "stringDefault":
			targets[i] = &m.StringDefault
		case "stringVarchar":
			targets[i] = &m.StringVarchar
		case "stringChar":
			targets[i] = &m.StringChar
		case "bitVal":
			targets[i] = &m.BitVal
		case "varBitVal":
			targets[i] = &m.VarBitVal
		case "inetVal":
			targets[i] = &m.InetVal
		case "xmlVal":
			targets[i] = &m.XmlVal
		case "cuidDefault":
			targets[i] = &m.CuidDefault
		case "cuid1Default":
			targets[i] = &m.Cuid1Default
		case "cuid2Default":
			targets[i] = &m.Cuid2Default
		case "uuidDefault":
			targets[i] = &m.UuidDefault
		case "uuid4Default":
			targets[i] = &m.Uuid4Default
		case "uuid7Default":
			targets[i] = &m.Uuid7Default
		case "ulidDefault":
			targets[i] = &m.UlidDefault
		case "nanoidDefault":
			targets[i] = &m.NanoidDefault
		case "uuidDb":
			targets[i] = &m.UuidDb
		case "intReq":
			targets[i] = &m.IntReq
		case "intOpt":
			targets[i] = &m.IntOpt
		case "intDefault":
			targets[i] = &m.IntDefault
		case "integerVal":
			targets[i] = &m.IntegerVal
		case "smallInt":
			targets[i] = &m.SmallInt
		case "tinyInt":
			targets[i] = &m.TinyInt
		case "oidVal":
			targets[i] = &m.OidVal
		case "bigIntReq":
			targets[i] = &m.BigIntReq
		case "bigIntOpt":
			targets[i] = &m.BigIntOpt
		case "floatReq":
			targets[i] = &m.FloatReq
		case "floatOpt":
			targets[i] = &m.FloatOpt
		case "realVal":
			targets[i] = &m.RealVal
		case "decimalReq":
			targets[i] = &m.DecimalReq
		case "decimalOpt":
			targets[i] = &m.DecimalOpt
		case "decimalPrecise":
			targets[i] = &m.DecimalPrecise
		case "moneyVal":
			targets[i] = &m.MoneyVal
		case "boolReq":
			targets[i] = &m.BoolReq
		case "boolOpt":
			targets[i] = &m.BoolOpt
		case "boolDefault":
			targets[i] = &m.BoolDefault
		case "dateTimeReq":
			targets[i] = &m.DateTimeReq
		case "dateTimeOpt":
			targets[i] = &m.DateTimeOpt
		case "dateTimeDefault":
			targets[i] = &m.DateTimeDefault
		case "updatedAt":
			targets[i] = &m.UpdatedAt
		case "dateTimeTz":
			targets[i] = &m.DateTimeTz
		case "timestampVal":
			targets[i] = &m.TimestampVal
		case "timeVal":
			targets[i] = &m.TimeVal
		case "timetzVal":
			targets[i] = &m.TimetzVal
		case "jsonReq":
			targets[i] = &m.JsonReq
		case "jsonOpt":
			targets[i] = &m.JsonOpt
		case "jsonVal":
			targets[i] = &m.JsonVal
		case "bytesReq":
			targets[i] = &m.BytesReq
		case "bytesOpt":
			targets[i] = &m.BytesOpt
		case "hstoreField":
			targets[i] = HstoreScan{P: &m.HstoreField}
		case "ltreeField":
			targets[i] = &m.LtreeField
		case "citextField":
			targets[i] = &m.CitextField
		}
	}
	return targets
}

var allFieldsSoFarDefaultCols = []string{
	"id",
	"stringReq",
	"stringOpt",
	"stringDefault",
	"stringVarchar",
	"stringChar",
	"bitVal",
	"varBitVal",
	"inetVal",
	"xmlVal",
	"cuidDefault",
	"cuid1Default",
	"cuid2Default",
	"uuidDefault",
	"uuid4Default",
	"uuid7Default",
	"ulidDefault",
	"nanoidDefault",
	"uuidDb",
	"intReq",
	"intOpt",
	"intDefault",
	"integerVal",
	"smallInt",
	"tinyInt",
	"oidVal",
	"bigIntReq",
	"bigIntOpt",
	"floatReq",
	"floatOpt",
	"realVal",
	"decimalReq",
	"decimalOpt",
	"decimalPrecise",
	"moneyVal",
	"boolReq",
	"boolOpt",
	"boolDefault",
	"dateTimeReq",
	"dateTimeOpt",
	"dateTimeDefault",
	"updatedAt",
	"dateTimeTz",
	"timestampVal",
	"timeVal",
	"timetzVal",
	"jsonReq",
	"jsonOpt",
	"jsonVal",
	"bytesReq",
	"bytesOpt",
	"hstoreField",
	"ltreeField",
	"citextField",
}

var allFieldsSoFarPKCols = []string{
	"id",
}

var allFieldsSoFarUniqueCols = []string{
	"id",
}

func selectAllFieldsSoFarCols(selects *AllFieldsSoFarSelect, omits *AllFieldsSoFarOmit, forceCols ...string) []string {
	if selects == nil && omits == nil && len(forceCols) == 0 {
		return allFieldsSoFarDefaultCols
	}

	anySelected := selects != nil && (selects.Id || selects.StringReq || selects.StringOpt || selects.StringDefault || selects.StringVarchar || selects.StringChar || selects.BitVal || selects.VarBitVal || selects.InetVal || selects.XmlVal || selects.CuidDefault || selects.Cuid1Default || selects.Cuid2Default || selects.UuidDefault || selects.Uuid4Default || selects.Uuid7Default || selects.UlidDefault || selects.NanoidDefault || selects.UuidDb || selects.IntReq || selects.IntOpt || selects.IntDefault || selects.IntegerVal || selects.SmallInt || selects.TinyInt || selects.OidVal || selects.BigIntReq || selects.BigIntOpt || selects.FloatReq || selects.FloatOpt || selects.RealVal || selects.DecimalReq || selects.DecimalOpt || selects.DecimalPrecise || selects.MoneyVal || selects.BoolReq || selects.BoolOpt || selects.BoolDefault || selects.DateTimeReq || selects.DateTimeOpt || selects.DateTimeDefault || selects.UpdatedAt || selects.DateTimeTz || selects.TimestampVal || selects.TimeVal || selects.TimetzVal || selects.JsonReq || selects.JsonOpt || selects.JsonVal || selects.BytesReq || selects.BytesOpt || selects.HstoreField || selects.LtreeField || selects.CitextField)

	specs := []colSpec{
		{"id", selects != nil && selects.Id, omits != nil && omits.Id, selects != nil && selects.hasAnyRelation()},
		{"stringReq", selects != nil && selects.StringReq, omits != nil && omits.StringReq, false},
		{"stringOpt", selects != nil && selects.StringOpt, omits != nil && omits.StringOpt, false},
		{"stringDefault", selects != nil && selects.StringDefault, omits != nil && omits.StringDefault, false},
		{"stringVarchar", selects != nil && selects.StringVarchar, omits != nil && omits.StringVarchar, false},
		{"stringChar", selects != nil && selects.StringChar, omits != nil && omits.StringChar, false},
		{"bitVal", selects != nil && selects.BitVal, omits != nil && omits.BitVal, false},
		{"varBitVal", selects != nil && selects.VarBitVal, omits != nil && omits.VarBitVal, false},
		{"inetVal", selects != nil && selects.InetVal, omits != nil && omits.InetVal, false},
		{"xmlVal", selects != nil && selects.XmlVal, omits != nil && omits.XmlVal, false},
		{"cuidDefault", selects != nil && selects.CuidDefault, omits != nil && omits.CuidDefault, false},
		{"cuid1Default", selects != nil && selects.Cuid1Default, omits != nil && omits.Cuid1Default, false},
		{"cuid2Default", selects != nil && selects.Cuid2Default, omits != nil && omits.Cuid2Default, false},
		{"uuidDefault", selects != nil && selects.UuidDefault, omits != nil && omits.UuidDefault, false},
		{"uuid4Default", selects != nil && selects.Uuid4Default, omits != nil && omits.Uuid4Default, false},
		{"uuid7Default", selects != nil && selects.Uuid7Default, omits != nil && omits.Uuid7Default, false},
		{"ulidDefault", selects != nil && selects.UlidDefault, omits != nil && omits.UlidDefault, false},
		{"nanoidDefault", selects != nil && selects.NanoidDefault, omits != nil && omits.NanoidDefault, false},
		{"uuidDb", selects != nil && selects.UuidDb, omits != nil && omits.UuidDb, false},
		{"intReq", selects != nil && selects.IntReq, omits != nil && omits.IntReq, false},
		{"intOpt", selects != nil && selects.IntOpt, omits != nil && omits.IntOpt, false},
		{"intDefault", selects != nil && selects.IntDefault, omits != nil && omits.IntDefault, false},
		{"integerVal", selects != nil && selects.IntegerVal, omits != nil && omits.IntegerVal, false},
		{"smallInt", selects != nil && selects.SmallInt, omits != nil && omits.SmallInt, false},
		{"tinyInt", selects != nil && selects.TinyInt, omits != nil && omits.TinyInt, false},
		{"oidVal", selects != nil && selects.OidVal, omits != nil && omits.OidVal, false},
		{"bigIntReq", selects != nil && selects.BigIntReq, omits != nil && omits.BigIntReq, false},
		{"bigIntOpt", selects != nil && selects.BigIntOpt, omits != nil && omits.BigIntOpt, false},
		{"floatReq", selects != nil && selects.FloatReq, omits != nil && omits.FloatReq, false},
		{"floatOpt", selects != nil && selects.FloatOpt, omits != nil && omits.FloatOpt, false},
		{"realVal", selects != nil && selects.RealVal, omits != nil && omits.RealVal, false},
		{"decimalReq", selects != nil && selects.DecimalReq, omits != nil && omits.DecimalReq, false},
		{"decimalOpt", selects != nil && selects.DecimalOpt, omits != nil && omits.DecimalOpt, false},
		{"decimalPrecise", selects != nil && selects.DecimalPrecise, omits != nil && omits.DecimalPrecise, false},
		{"moneyVal", selects != nil && selects.MoneyVal, omits != nil && omits.MoneyVal, false},
		{"boolReq", selects != nil && selects.BoolReq, omits != nil && omits.BoolReq, false},
		{"boolOpt", selects != nil && selects.BoolOpt, omits != nil && omits.BoolOpt, false},
		{"boolDefault", selects != nil && selects.BoolDefault, omits != nil && omits.BoolDefault, false},
		{"dateTimeReq", selects != nil && selects.DateTimeReq, omits != nil && omits.DateTimeReq, false},
		{"dateTimeOpt", selects != nil && selects.DateTimeOpt, omits != nil && omits.DateTimeOpt, false},
		{"dateTimeDefault", selects != nil && selects.DateTimeDefault, omits != nil && omits.DateTimeDefault, false},
		{"updatedAt", selects != nil && selects.UpdatedAt, omits != nil && omits.UpdatedAt, false},
		{"dateTimeTz", selects != nil && selects.DateTimeTz, omits != nil && omits.DateTimeTz, false},
		{"timestampVal", selects != nil && selects.TimestampVal, omits != nil && omits.TimestampVal, false},
		{"timeVal", selects != nil && selects.TimeVal, omits != nil && omits.TimeVal, false},
		{"timetzVal", selects != nil && selects.TimetzVal, omits != nil && omits.TimetzVal, false},
		{"jsonReq", selects != nil && selects.JsonReq, omits != nil && omits.JsonReq, false},
		{"jsonOpt", selects != nil && selects.JsonOpt, omits != nil && omits.JsonOpt, false},
		{"jsonVal", selects != nil && selects.JsonVal, omits != nil && omits.JsonVal, false},
		{"bytesReq", selects != nil && selects.BytesReq, omits != nil && omits.BytesReq, false},
		{"bytesOpt", selects != nil && selects.BytesOpt, omits != nil && omits.BytesOpt, false},
		{"hstoreField", selects != nil && selects.HstoreField, omits != nil && omits.HstoreField, false},
		{"ltreeField", selects != nil && selects.LtreeField, omits != nil && omits.LtreeField, false},
		{"citextField", selects != nil && selects.CitextField, omits != nil && omits.CitextField, false},
	}

	cols := computeCols(specs, selects != nil, anySelected)

	for _, f := range forceCols {
		if !slices.Contains(cols, f) {
			cols = append(cols, f)
		}
	}

	return cols
}

func (s *AllFieldsSoFarSelect) hasAnyRelation() bool {
	if s == nil {
		return false
	}
	return false
}

type AllFieldsSoFarCreateBuilder struct {
	*CreateBuilder[AllFieldsSoFar, AllFieldsSoFarSelect, AllFieldsSoFarOmit]
}

func (b *AllFieldsSoFarCreateBuilder) OnConflict(target UniqueConstraintTarget) *AllFieldsSoFarConflictBuilder[AllFieldsSoFarCreateBuilder] {
	return &AllFieldsSoFarConflictBuilder[AllFieldsSoFarCreateBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

func (b *AllFieldsSoFarCreateBuilder) SetId(v int32) *AllFieldsSoFarCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "id", Val: v})
	return b
}
func (b *AllFieldsSoFarCreateBuilder) SetStringReq(v string) *AllFieldsSoFarCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "stringReq", Val: v})
	return b
}
func (b *AllFieldsSoFarCreateBuilder) SetStringOpt(v string) *AllFieldsSoFarCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "stringOpt", Val: v})
	return b
}
func (b *AllFieldsSoFarCreateBuilder) SetStringDefault(v string) *AllFieldsSoFarCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "stringDefault", Val: v})
	return b
}
func (b *AllFieldsSoFarCreateBuilder) SetStringVarchar(v string) *AllFieldsSoFarCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "stringVarchar", Val: v})
	return b
}
func (b *AllFieldsSoFarCreateBuilder) SetStringChar(v string) *AllFieldsSoFarCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "stringChar", Val: v})
	return b
}
func (b *AllFieldsSoFarCreateBuilder) SetBitVal(v string) *AllFieldsSoFarCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "bitVal", Val: v})
	return b
}
func (b *AllFieldsSoFarCreateBuilder) SetVarBitVal(v string) *AllFieldsSoFarCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "varBitVal", Val: v})
	return b
}
func (b *AllFieldsSoFarCreateBuilder) SetInetVal(v string) *AllFieldsSoFarCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "inetVal", Val: v})
	return b
}
func (b *AllFieldsSoFarCreateBuilder) SetXmlVal(v string) *AllFieldsSoFarCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "xmlVal", Val: v})
	return b
}
func (b *AllFieldsSoFarCreateBuilder) SetCuidDefault(v string) *AllFieldsSoFarCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "cuidDefault", Val: v})
	return b
}
func (b *AllFieldsSoFarCreateBuilder) SetCuid1Default(v string) *AllFieldsSoFarCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "cuid1Default", Val: v})
	return b
}
func (b *AllFieldsSoFarCreateBuilder) SetCuid2Default(v string) *AllFieldsSoFarCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "cuid2Default", Val: v})
	return b
}
func (b *AllFieldsSoFarCreateBuilder) SetUuidDefault(v string) *AllFieldsSoFarCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "uuidDefault", Val: v})
	return b
}
func (b *AllFieldsSoFarCreateBuilder) SetUuid4Default(v string) *AllFieldsSoFarCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "uuid4Default", Val: v})
	return b
}
func (b *AllFieldsSoFarCreateBuilder) SetUuid7Default(v string) *AllFieldsSoFarCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "uuid7Default", Val: v})
	return b
}
func (b *AllFieldsSoFarCreateBuilder) SetUlidDefault(v string) *AllFieldsSoFarCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "ulidDefault", Val: v})
	return b
}
func (b *AllFieldsSoFarCreateBuilder) SetNanoidDefault(v string) *AllFieldsSoFarCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "nanoidDefault", Val: v})
	return b
}
func (b *AllFieldsSoFarCreateBuilder) SetUuidDb(v string) *AllFieldsSoFarCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "uuidDb", Val: v})
	return b
}
func (b *AllFieldsSoFarCreateBuilder) SetIntReq(v int32) *AllFieldsSoFarCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "intReq", Val: v})
	return b
}
func (b *AllFieldsSoFarCreateBuilder) SetIntOpt(v int32) *AllFieldsSoFarCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "intOpt", Val: v})
	return b
}
func (b *AllFieldsSoFarCreateBuilder) SetIntDefault(v int32) *AllFieldsSoFarCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "intDefault", Val: v})
	return b
}
func (b *AllFieldsSoFarCreateBuilder) SetIntegerVal(v int32) *AllFieldsSoFarCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "integerVal", Val: v})
	return b
}
func (b *AllFieldsSoFarCreateBuilder) SetSmallInt(v int32) *AllFieldsSoFarCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "smallInt", Val: v})
	return b
}
func (b *AllFieldsSoFarCreateBuilder) SetTinyInt(v int32) *AllFieldsSoFarCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "tinyInt", Val: v})
	return b
}
func (b *AllFieldsSoFarCreateBuilder) SetOidVal(v int32) *AllFieldsSoFarCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "oidVal", Val: v})
	return b
}
func (b *AllFieldsSoFarCreateBuilder) SetBigIntReq(v int64) *AllFieldsSoFarCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "bigIntReq", Val: v})
	return b
}
func (b *AllFieldsSoFarCreateBuilder) SetBigIntOpt(v int64) *AllFieldsSoFarCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "bigIntOpt", Val: v})
	return b
}
func (b *AllFieldsSoFarCreateBuilder) SetFloatReq(v float64) *AllFieldsSoFarCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "floatReq", Val: v})
	return b
}
func (b *AllFieldsSoFarCreateBuilder) SetFloatOpt(v float64) *AllFieldsSoFarCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "floatOpt", Val: v})
	return b
}
func (b *AllFieldsSoFarCreateBuilder) SetRealVal(v float64) *AllFieldsSoFarCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "realVal", Val: v})
	return b
}
func (b *AllFieldsSoFarCreateBuilder) SetDecimalReq(v string) *AllFieldsSoFarCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "decimalReq", Val: v})
	return b
}
func (b *AllFieldsSoFarCreateBuilder) SetDecimalOpt(v string) *AllFieldsSoFarCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "decimalOpt", Val: v})
	return b
}
func (b *AllFieldsSoFarCreateBuilder) SetDecimalPrecise(v string) *AllFieldsSoFarCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "decimalPrecise", Val: v})
	return b
}
func (b *AllFieldsSoFarCreateBuilder) SetMoneyVal(v string) *AllFieldsSoFarCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "moneyVal", Val: v})
	return b
}
func (b *AllFieldsSoFarCreateBuilder) SetBoolReq(v bool) *AllFieldsSoFarCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "boolReq", Val: v})
	return b
}
func (b *AllFieldsSoFarCreateBuilder) SetBoolOpt(v bool) *AllFieldsSoFarCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "boolOpt", Val: v})
	return b
}
func (b *AllFieldsSoFarCreateBuilder) SetBoolDefault(v bool) *AllFieldsSoFarCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "boolDefault", Val: v})
	return b
}
func (b *AllFieldsSoFarCreateBuilder) SetDateTimeReq(v time.Time) *AllFieldsSoFarCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "dateTimeReq", Val: v})
	return b
}
func (b *AllFieldsSoFarCreateBuilder) SetDateTimeOpt(v time.Time) *AllFieldsSoFarCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "dateTimeOpt", Val: v})
	return b
}
func (b *AllFieldsSoFarCreateBuilder) SetDateTimeDefault(v time.Time) *AllFieldsSoFarCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "dateTimeDefault", Val: v})
	return b
}
func (b *AllFieldsSoFarCreateBuilder) SetUpdatedAt(v time.Time) *AllFieldsSoFarCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "updatedAt", Val: v})
	return b
}
func (b *AllFieldsSoFarCreateBuilder) SetDateTimeTz(v time.Time) *AllFieldsSoFarCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "dateTimeTz", Val: v})
	return b
}
func (b *AllFieldsSoFarCreateBuilder) SetTimestampVal(v time.Time) *AllFieldsSoFarCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "timestampVal", Val: v})
	return b
}
func (b *AllFieldsSoFarCreateBuilder) SetTimeVal(v time.Time) *AllFieldsSoFarCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "timeVal", Val: v})
	return b
}
func (b *AllFieldsSoFarCreateBuilder) SetTimetzVal(v time.Time) *AllFieldsSoFarCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "timetzVal", Val: v})
	return b
}
func (b *AllFieldsSoFarCreateBuilder) SetJsonReq(v json.RawMessage) *AllFieldsSoFarCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "jsonReq", Val: v})
	return b
}
func (b *AllFieldsSoFarCreateBuilder) SetJsonOpt(v json.RawMessage) *AllFieldsSoFarCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "jsonOpt", Val: v})
	return b
}
func (b *AllFieldsSoFarCreateBuilder) SetJsonVal(v json.RawMessage) *AllFieldsSoFarCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "jsonVal", Val: v})
	return b
}
func (b *AllFieldsSoFarCreateBuilder) SetBytesReq(v []byte) *AllFieldsSoFarCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "bytesReq", Val: v})
	return b
}
func (b *AllFieldsSoFarCreateBuilder) SetBytesOpt(v []byte) *AllFieldsSoFarCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "bytesOpt", Val: v})
	return b
}
func (b *AllFieldsSoFarCreateBuilder) SetHstoreField(v map[string]*string) *AllFieldsSoFarCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "hstoreField", Val: v})
	return b
}
func (b *AllFieldsSoFarCreateBuilder) SetLtreeField(v string) *AllFieldsSoFarCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "ltreeField", Val: v})
	return b
}
func (b *AllFieldsSoFarCreateBuilder) SetCitextField(v string) *AllFieldsSoFarCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "citextField", Val: v})
	return b
}

func (d *AllFieldsSoFarDelegate) Create(assignments ...FieldAssignment) *AllFieldsSoFarCreateBuilder {
	return &AllFieldsSoFarCreateBuilder{
		CreateBuilder: &CreateBuilder[AllFieldsSoFar, AllFieldsSoFarSelect, AllFieldsSoFarOmit]{
			assignments: assignments,
			execFunc:    d.executeCreate,
		},
	}
}

const (
	providedAllFieldsSoFarId              uint64 = 1 << 0
	providedAllFieldsSoFarStringReq       uint64 = 1 << 1
	providedAllFieldsSoFarStringOpt       uint64 = 1 << 2
	providedAllFieldsSoFarStringDefault   uint64 = 1 << 3
	providedAllFieldsSoFarStringVarchar   uint64 = 1 << 4
	providedAllFieldsSoFarStringChar      uint64 = 1 << 5
	providedAllFieldsSoFarBitVal          uint64 = 1 << 6
	providedAllFieldsSoFarVarBitVal       uint64 = 1 << 7
	providedAllFieldsSoFarInetVal         uint64 = 1 << 8
	providedAllFieldsSoFarXmlVal          uint64 = 1 << 9
	providedAllFieldsSoFarCuidDefault     uint64 = 1 << 10
	providedAllFieldsSoFarCuid1Default    uint64 = 1 << 11
	providedAllFieldsSoFarCuid2Default    uint64 = 1 << 12
	providedAllFieldsSoFarUuidDefault     uint64 = 1 << 13
	providedAllFieldsSoFarUuid4Default    uint64 = 1 << 14
	providedAllFieldsSoFarUuid7Default    uint64 = 1 << 15
	providedAllFieldsSoFarUlidDefault     uint64 = 1 << 16
	providedAllFieldsSoFarNanoidDefault   uint64 = 1 << 17
	providedAllFieldsSoFarUuidDb          uint64 = 1 << 18
	providedAllFieldsSoFarIntReq          uint64 = 1 << 19
	providedAllFieldsSoFarIntOpt          uint64 = 1 << 20
	providedAllFieldsSoFarIntDefault      uint64 = 1 << 21
	providedAllFieldsSoFarIntegerVal      uint64 = 1 << 22
	providedAllFieldsSoFarSmallInt        uint64 = 1 << 23
	providedAllFieldsSoFarTinyInt         uint64 = 1 << 24
	providedAllFieldsSoFarOidVal          uint64 = 1 << 25
	providedAllFieldsSoFarBigIntReq       uint64 = 1 << 26
	providedAllFieldsSoFarBigIntOpt       uint64 = 1 << 27
	providedAllFieldsSoFarFloatReq        uint64 = 1 << 28
	providedAllFieldsSoFarFloatOpt        uint64 = 1 << 29
	providedAllFieldsSoFarRealVal         uint64 = 1 << 30
	providedAllFieldsSoFarDecimalReq      uint64 = 1 << 31
	providedAllFieldsSoFarDecimalOpt      uint64 = 1 << 32
	providedAllFieldsSoFarDecimalPrecise  uint64 = 1 << 33
	providedAllFieldsSoFarMoneyVal        uint64 = 1 << 34
	providedAllFieldsSoFarBoolReq         uint64 = 1 << 35
	providedAllFieldsSoFarBoolOpt         uint64 = 1 << 36
	providedAllFieldsSoFarBoolDefault     uint64 = 1 << 37
	providedAllFieldsSoFarDateTimeReq     uint64 = 1 << 38
	providedAllFieldsSoFarDateTimeOpt     uint64 = 1 << 39
	providedAllFieldsSoFarDateTimeDefault uint64 = 1 << 40
	providedAllFieldsSoFarUpdatedAt       uint64 = 1 << 41
	providedAllFieldsSoFarDateTimeTz      uint64 = 1 << 42
	providedAllFieldsSoFarTimestampVal    uint64 = 1 << 43
	providedAllFieldsSoFarTimeVal         uint64 = 1 << 44
	providedAllFieldsSoFarTimetzVal       uint64 = 1 << 45
	providedAllFieldsSoFarJsonReq         uint64 = 1 << 46
	providedAllFieldsSoFarJsonOpt         uint64 = 1 << 47
	providedAllFieldsSoFarJsonVal         uint64 = 1 << 48
	providedAllFieldsSoFarBytesReq        uint64 = 1 << 49
	providedAllFieldsSoFarBytesOpt        uint64 = 1 << 50
	providedAllFieldsSoFarHstoreField     uint64 = 1 << 51
	providedAllFieldsSoFarLtreeField      uint64 = 1 << 52
	providedAllFieldsSoFarCitextField     uint64 = 1 << 53
)

func assignmentsToAllFieldsSoFarCreate(assignments []FieldAssignment) (AllFieldsSoFarCreate, error) {
	var input AllFieldsSoFarCreate
	var errs ValidationError
	var provided uint64

	for _, a := range assignments {
		switch a.Col {
		case "id":
			provided |= providedAllFieldsSoFarId
			if v, ok := a.Val.(int32); ok {
				input.Id = &v
				ValidateInt32(&errs, "id", v, "")
			} else {
				errs.Add("id", a.Val, "type", "field id must be of type int32")
			}
		case "stringReq":
			provided |= providedAllFieldsSoFarStringReq
			if v, ok := a.Val.(string); ok {
				input.StringReq = v
				ValidateString(&errs, "stringReq", v, true, 0, false, false)
			} else {
				errs.Add("stringReq", a.Val, "type", "field stringReq must be of type string")
			}
		case "stringOpt":
			provided |= providedAllFieldsSoFarStringOpt
			if v, ok := a.Val.(string); ok {
				input.StringOpt = &v
				ValidateString(&errs, "stringOpt", v, false, 0, false, false)
			} else {
				errs.Add("stringOpt", a.Val, "type", "field stringOpt must be of type string")
			}
		case "stringDefault":
			provided |= providedAllFieldsSoFarStringDefault
			if v, ok := a.Val.(string); ok {
				input.StringDefault = &v
				ValidateString(&errs, "stringDefault", v, false, 0, false, false)
			} else {
				errs.Add("stringDefault", a.Val, "type", "field stringDefault must be of type string")
			}
		case "stringVarchar":
			provided |= providedAllFieldsSoFarStringVarchar
			if v, ok := a.Val.(string); ok {
				input.StringVarchar = v
				ValidateString(&errs, "stringVarchar", v, true, 255, false, false)
			} else {
				errs.Add("stringVarchar", a.Val, "type", "field stringVarchar must be of type string")
			}
		case "stringChar":
			provided |= providedAllFieldsSoFarStringChar
			if v, ok := a.Val.(string); ok {
				input.StringChar = v
				ValidateString(&errs, "stringChar", v, true, 10, false, false)
			} else {
				errs.Add("stringChar", a.Val, "type", "field stringChar must be of type string")
			}
		case "bitVal":
			provided |= providedAllFieldsSoFarBitVal
			if v, ok := a.Val.(string); ok {
				input.BitVal = v
				ValidateString(&errs, "bitVal", v, true, 0, true, false)
			} else {
				errs.Add("bitVal", a.Val, "type", "field bitVal must be of type string")
			}
		case "varBitVal":
			provided |= providedAllFieldsSoFarVarBitVal
			if v, ok := a.Val.(string); ok {
				input.VarBitVal = v
				ValidateString(&errs, "varBitVal", v, true, 0, true, false)
			} else {
				errs.Add("varBitVal", a.Val, "type", "field varBitVal must be of type string")
			}
		case "inetVal":
			provided |= providedAllFieldsSoFarInetVal
			if v, ok := a.Val.(string); ok {
				input.InetVal = v
				ValidateString(&errs, "inetVal", v, true, 0, false, true)
			} else {
				errs.Add("inetVal", a.Val, "type", "field inetVal must be of type string")
			}
		case "xmlVal":
			provided |= providedAllFieldsSoFarXmlVal
			if v, ok := a.Val.(string); ok {
				input.XmlVal = v
				ValidateString(&errs, "xmlVal", v, true, 0, false, false)
			} else {
				errs.Add("xmlVal", a.Val, "type", "field xmlVal must be of type string")
			}
		case "cuidDefault":
			provided |= providedAllFieldsSoFarCuidDefault
			if v, ok := a.Val.(string); ok {
				input.CuidDefault = &v
				ValidateString(&errs, "cuidDefault", v, false, 0, false, false)
			} else {
				errs.Add("cuidDefault", a.Val, "type", "field cuidDefault must be of type string")
			}
		case "cuid1Default":
			provided |= providedAllFieldsSoFarCuid1Default
			if v, ok := a.Val.(string); ok {
				input.Cuid1Default = &v
				ValidateString(&errs, "cuid1Default", v, false, 0, false, false)
			} else {
				errs.Add("cuid1Default", a.Val, "type", "field cuid1Default must be of type string")
			}
		case "cuid2Default":
			provided |= providedAllFieldsSoFarCuid2Default
			if v, ok := a.Val.(string); ok {
				input.Cuid2Default = &v
				ValidateString(&errs, "cuid2Default", v, false, 0, false, false)
			} else {
				errs.Add("cuid2Default", a.Val, "type", "field cuid2Default must be of type string")
			}
		case "uuidDefault":
			provided |= providedAllFieldsSoFarUuidDefault
			if v, ok := a.Val.(string); ok {
				input.UuidDefault = &v
				ValidateString(&errs, "uuidDefault", v, false, 0, false, false)
			} else {
				errs.Add("uuidDefault", a.Val, "type", "field uuidDefault must be of type string")
			}
		case "uuid4Default":
			provided |= providedAllFieldsSoFarUuid4Default
			if v, ok := a.Val.(string); ok {
				input.Uuid4Default = &v
				ValidateString(&errs, "uuid4Default", v, false, 0, false, false)
			} else {
				errs.Add("uuid4Default", a.Val, "type", "field uuid4Default must be of type string")
			}
		case "uuid7Default":
			provided |= providedAllFieldsSoFarUuid7Default
			if v, ok := a.Val.(string); ok {
				input.Uuid7Default = &v
				ValidateString(&errs, "uuid7Default", v, false, 0, false, false)
			} else {
				errs.Add("uuid7Default", a.Val, "type", "field uuid7Default must be of type string")
			}
		case "ulidDefault":
			provided |= providedAllFieldsSoFarUlidDefault
			if v, ok := a.Val.(string); ok {
				input.UlidDefault = &v
				ValidateString(&errs, "ulidDefault", v, false, 0, false, false)
			} else {
				errs.Add("ulidDefault", a.Val, "type", "field ulidDefault must be of type string")
			}
		case "nanoidDefault":
			provided |= providedAllFieldsSoFarNanoidDefault
			if v, ok := a.Val.(string); ok {
				input.NanoidDefault = &v
				ValidateString(&errs, "nanoidDefault", v, false, 0, false, false)
			} else {
				errs.Add("nanoidDefault", a.Val, "type", "field nanoidDefault must be of type string")
			}
		case "uuidDb":
			provided |= providedAllFieldsSoFarUuidDb
			if v, ok := a.Val.(string); ok {
				input.UuidDb = v
				ValidateString(&errs, "uuidDb", v, true, 0, false, false)
			} else {
				errs.Add("uuidDb", a.Val, "type", "field uuidDb must be of type string")
			}
		case "intReq":
			provided |= providedAllFieldsSoFarIntReq
			if v, ok := a.Val.(int32); ok {
				input.IntReq = v
				ValidateInt32(&errs, "intReq", v, "")
			} else {
				errs.Add("intReq", a.Val, "type", "field intReq must be of type int32")
			}
		case "intOpt":
			provided |= providedAllFieldsSoFarIntOpt
			if v, ok := a.Val.(int32); ok {
				input.IntOpt = &v
				ValidateInt32(&errs, "intOpt", v, "")
			} else {
				errs.Add("intOpt", a.Val, "type", "field intOpt must be of type int32")
			}
		case "intDefault":
			provided |= providedAllFieldsSoFarIntDefault
			if v, ok := a.Val.(int32); ok {
				input.IntDefault = &v
				ValidateInt32(&errs, "intDefault", v, "")
			} else {
				errs.Add("intDefault", a.Val, "type", "field intDefault must be of type int32")
			}
		case "integerVal":
			provided |= providedAllFieldsSoFarIntegerVal
			if v, ok := a.Val.(int32); ok {
				input.IntegerVal = v
				ValidateInt32(&errs, "integerVal", v, "")
			} else {
				errs.Add("integerVal", a.Val, "type", "field integerVal must be of type int32")
			}
		case "smallInt":
			provided |= providedAllFieldsSoFarSmallInt
			if v, ok := a.Val.(int32); ok {
				input.SmallInt = v
				ValidateInt32(&errs, "smallInt", v, "SmallInt")
			} else {
				errs.Add("smallInt", a.Val, "type", "field smallInt must be of type int32")
			}
		case "tinyInt":
			provided |= providedAllFieldsSoFarTinyInt
			if v, ok := a.Val.(int32); ok {
				input.TinyInt = v
				ValidateInt32(&errs, "tinyInt", v, "")
			} else {
				errs.Add("tinyInt", a.Val, "type", "field tinyInt must be of type int32")
			}
		case "oidVal":
			provided |= providedAllFieldsSoFarOidVal
			if v, ok := a.Val.(int32); ok {
				input.OidVal = v
				ValidateInt32(&errs, "oidVal", v, "Oid")
			} else {
				errs.Add("oidVal", a.Val, "type", "field oidVal must be of type int32")
			}
		case "bigIntReq":
			provided |= providedAllFieldsSoFarBigIntReq
			if v, ok := a.Val.(int64); ok {
				input.BigIntReq = v
				ValidateInt64(&errs, "bigIntReq", v, "")
			} else {
				errs.Add("bigIntReq", a.Val, "type", "field bigIntReq must be of type int64")
			}
		case "bigIntOpt":
			provided |= providedAllFieldsSoFarBigIntOpt
			if v, ok := a.Val.(int64); ok {
				input.BigIntOpt = &v
				ValidateInt64(&errs, "bigIntOpt", v, "")
			} else {
				errs.Add("bigIntOpt", a.Val, "type", "field bigIntOpt must be of type int64")
			}
		case "floatReq":
			provided |= providedAllFieldsSoFarFloatReq
			if v, ok := a.Val.(float64); ok {
				input.FloatReq = v
			} else {
				errs.Add("floatReq", a.Val, "type", "field floatReq must be of type float64")
			}
		case "floatOpt":
			provided |= providedAllFieldsSoFarFloatOpt
			if v, ok := a.Val.(float64); ok {
				input.FloatOpt = &v
			} else {
				errs.Add("floatOpt", a.Val, "type", "field floatOpt must be of type float64")
			}
		case "realVal":
			provided |= providedAllFieldsSoFarRealVal
			if v, ok := a.Val.(float64); ok {
				input.RealVal = v
			} else {
				errs.Add("realVal", a.Val, "type", "field realVal must be of type float64")
			}
		case "decimalReq":
			provided |= providedAllFieldsSoFarDecimalReq
			if v, ok := a.Val.(string); ok {
				input.DecimalReq = v
				ValidateString(&errs, "decimalReq", v, true, 0, false, false)
			} else {
				errs.Add("decimalReq", a.Val, "type", "field decimalReq must be of type string")
			}
		case "decimalOpt":
			provided |= providedAllFieldsSoFarDecimalOpt
			if v, ok := a.Val.(string); ok {
				input.DecimalOpt = &v
				ValidateString(&errs, "decimalOpt", v, false, 0, false, false)
			} else {
				errs.Add("decimalOpt", a.Val, "type", "field decimalOpt must be of type string")
			}
		case "decimalPrecise":
			provided |= providedAllFieldsSoFarDecimalPrecise
			if v, ok := a.Val.(string); ok {
				input.DecimalPrecise = v
				ValidateString(&errs, "decimalPrecise", v, true, 0, false, false)
			} else {
				errs.Add("decimalPrecise", a.Val, "type", "field decimalPrecise must be of type string")
			}
		case "moneyVal":
			provided |= providedAllFieldsSoFarMoneyVal
			if v, ok := a.Val.(string); ok {
				input.MoneyVal = v
				ValidateString(&errs, "moneyVal", v, true, 0, false, false)
			} else {
				errs.Add("moneyVal", a.Val, "type", "field moneyVal must be of type string")
			}
		case "boolReq":
			provided |= providedAllFieldsSoFarBoolReq
			if v, ok := a.Val.(bool); ok {
				input.BoolReq = v
			} else {
				errs.Add("boolReq", a.Val, "type", "field boolReq must be of type bool")
			}
		case "boolOpt":
			provided |= providedAllFieldsSoFarBoolOpt
			if v, ok := a.Val.(bool); ok {
				input.BoolOpt = &v
			} else {
				errs.Add("boolOpt", a.Val, "type", "field boolOpt must be of type bool")
			}
		case "boolDefault":
			provided |= providedAllFieldsSoFarBoolDefault
			if v, ok := a.Val.(bool); ok {
				input.BoolDefault = &v
			} else {
				errs.Add("boolDefault", a.Val, "type", "field boolDefault must be of type bool")
			}
		case "dateTimeReq":
			provided |= providedAllFieldsSoFarDateTimeReq
			if v, ok := a.Val.(time.Time); ok {
				input.DateTimeReq = v
			} else {
				errs.Add("dateTimeReq", a.Val, "type", "field dateTimeReq must be of type time.Time")
			}
		case "dateTimeOpt":
			provided |= providedAllFieldsSoFarDateTimeOpt
			if v, ok := a.Val.(time.Time); ok {
				input.DateTimeOpt = &v
			} else {
				errs.Add("dateTimeOpt", a.Val, "type", "field dateTimeOpt must be of type time.Time")
			}
		case "dateTimeDefault":
			provided |= providedAllFieldsSoFarDateTimeDefault
			if v, ok := a.Val.(time.Time); ok {
				input.DateTimeDefault = &v
			} else {
				errs.Add("dateTimeDefault", a.Val, "type", "field dateTimeDefault must be of type time.Time")
			}
		case "updatedAt":
			provided |= providedAllFieldsSoFarUpdatedAt
			if v, ok := a.Val.(time.Time); ok {
				input.UpdatedAt = v
			} else {
				errs.Add("updatedAt", a.Val, "type", "field updatedAt must be of type time.Time")
			}
		case "dateTimeTz":
			provided |= providedAllFieldsSoFarDateTimeTz
			if v, ok := a.Val.(time.Time); ok {
				input.DateTimeTz = v
			} else {
				errs.Add("dateTimeTz", a.Val, "type", "field dateTimeTz must be of type time.Time")
			}
		case "timestampVal":
			provided |= providedAllFieldsSoFarTimestampVal
			if v, ok := a.Val.(time.Time); ok {
				input.TimestampVal = v
			} else {
				errs.Add("timestampVal", a.Val, "type", "field timestampVal must be of type time.Time")
			}
		case "timeVal":
			provided |= providedAllFieldsSoFarTimeVal
			if v, ok := a.Val.(time.Time); ok {
				input.TimeVal = v
			} else {
				errs.Add("timeVal", a.Val, "type", "field timeVal must be of type time.Time")
			}
		case "timetzVal":
			provided |= providedAllFieldsSoFarTimetzVal
			if v, ok := a.Val.(time.Time); ok {
				input.TimetzVal = v
			} else {
				errs.Add("timetzVal", a.Val, "type", "field timetzVal must be of type time.Time")
			}
		case "jsonReq":
			provided |= providedAllFieldsSoFarJsonReq
			if v, ok := a.Val.(json.RawMessage); ok {
				input.JsonReq = v
			} else {
				errs.Add("jsonReq", a.Val, "type", "field jsonReq must be of type json.RawMessage")
			}
		case "jsonOpt":
			provided |= providedAllFieldsSoFarJsonOpt
			if v, ok := a.Val.(json.RawMessage); ok {
				input.JsonOpt = &v
			} else {
				errs.Add("jsonOpt", a.Val, "type", "field jsonOpt must be of type json.RawMessage")
			}
		case "jsonVal":
			provided |= providedAllFieldsSoFarJsonVal
			if v, ok := a.Val.(json.RawMessage); ok {
				input.JsonVal = v
			} else {
				errs.Add("jsonVal", a.Val, "type", "field jsonVal must be of type json.RawMessage")
			}
		case "bytesReq":
			provided |= providedAllFieldsSoFarBytesReq
			if v, ok := a.Val.([]byte); ok {
				input.BytesReq = v
			} else {
				errs.Add("bytesReq", a.Val, "type", "field bytesReq must be of type []byte")
			}
		case "bytesOpt":
			provided |= providedAllFieldsSoFarBytesOpt
			if v, ok := a.Val.([]byte); ok {
				input.BytesOpt = &v
			} else {
				errs.Add("bytesOpt", a.Val, "type", "field bytesOpt must be of type []byte")
			}
		case "hstoreField":
			provided |= providedAllFieldsSoFarHstoreField
			if v, ok := a.Val.(map[string]*string); ok {
				input.HstoreField = &v
			} else {
				errs.Add("hstoreField", a.Val, "type", "field hstoreField must be of type map[string]*string")
			}
		case "ltreeField":
			provided |= providedAllFieldsSoFarLtreeField
			if v, ok := a.Val.(string); ok {
				input.LtreeField = v
				ValidateString(&errs, "ltreeField", v, true, 0, false, false)
			} else {
				errs.Add("ltreeField", a.Val, "type", "field ltreeField must be of type string")
			}
		case "citextField":
			provided |= providedAllFieldsSoFarCitextField
			if v, ok := a.Val.(string); ok {
				input.CitextField = &v
				ValidateString(&errs, "citextField", v, false, 0, false, false)
			} else {
				errs.Add("citextField", a.Val, "type", "field citextField must be of type string")
			}
		}
	}
	if provided&providedAllFieldsSoFarStringReq == 0 {
		errs.Add("stringReq", "", "required", "field StringReq is required")
	}
	if provided&providedAllFieldsSoFarStringVarchar == 0 {
		errs.Add("stringVarchar", "", "required", "field StringVarchar is required")
	}
	if provided&providedAllFieldsSoFarStringChar == 0 {
		errs.Add("stringChar", "", "required", "field StringChar is required")
	}
	if provided&providedAllFieldsSoFarBitVal == 0 {
		errs.Add("bitVal", "", "required", "field BitVal is required")
	}
	if provided&providedAllFieldsSoFarVarBitVal == 0 {
		errs.Add("varBitVal", "", "required", "field VarBitVal is required")
	}
	if provided&providedAllFieldsSoFarInetVal == 0 {
		errs.Add("inetVal", "", "required", "field InetVal is required")
	}
	if provided&providedAllFieldsSoFarXmlVal == 0 {
		errs.Add("xmlVal", "", "required", "field XmlVal is required")
	}
	if provided&providedAllFieldsSoFarUuidDb == 0 {
		errs.Add("uuidDb", "", "required", "field UuidDb is required")
	}
	if provided&providedAllFieldsSoFarIntReq == 0 {
		errs.Add("intReq", nil, "required", "field IntReq is required")
	}
	if provided&providedAllFieldsSoFarIntegerVal == 0 {
		errs.Add("integerVal", nil, "required", "field IntegerVal is required")
	}
	if provided&providedAllFieldsSoFarSmallInt == 0 {
		errs.Add("smallInt", nil, "required", "field SmallInt is required")
	}
	if provided&providedAllFieldsSoFarTinyInt == 0 {
		errs.Add("tinyInt", nil, "required", "field TinyInt is required")
	}
	if provided&providedAllFieldsSoFarOidVal == 0 {
		errs.Add("oidVal", nil, "required", "field OidVal is required")
	}
	if provided&providedAllFieldsSoFarBigIntReq == 0 {
		errs.Add("bigIntReq", nil, "required", "field BigIntReq is required")
	}
	if provided&providedAllFieldsSoFarFloatReq == 0 {
		errs.Add("floatReq", nil, "required", "field FloatReq is required")
	}
	if provided&providedAllFieldsSoFarRealVal == 0 {
		errs.Add("realVal", nil, "required", "field RealVal is required")
	}
	if provided&providedAllFieldsSoFarDecimalReq == 0 {
		errs.Add("decimalReq", "", "required", "field DecimalReq is required")
	}
	if provided&providedAllFieldsSoFarDecimalPrecise == 0 {
		errs.Add("decimalPrecise", "", "required", "field DecimalPrecise is required")
	}
	if provided&providedAllFieldsSoFarMoneyVal == 0 {
		errs.Add("moneyVal", "", "required", "field MoneyVal is required")
	}
	if provided&providedAllFieldsSoFarBoolReq == 0 {
		errs.Add("boolReq", nil, "required", "field BoolReq is required")
	}
	if provided&providedAllFieldsSoFarDateTimeReq == 0 {
		errs.Add("dateTimeReq", nil, "required", "field DateTimeReq is required")
	}
	if provided&providedAllFieldsSoFarUpdatedAt == 0 {
		errs.Add("updatedAt", nil, "required", "field UpdatedAt is required")
	}
	if provided&providedAllFieldsSoFarDateTimeTz == 0 {
		errs.Add("dateTimeTz", nil, "required", "field DateTimeTz is required")
	}
	if provided&providedAllFieldsSoFarTimestampVal == 0 {
		errs.Add("timestampVal", nil, "required", "field TimestampVal is required")
	}
	if provided&providedAllFieldsSoFarTimeVal == 0 {
		errs.Add("timeVal", nil, "required", "field TimeVal is required")
	}
	if provided&providedAllFieldsSoFarTimetzVal == 0 {
		errs.Add("timetzVal", nil, "required", "field TimetzVal is required")
	}
	if provided&providedAllFieldsSoFarJsonReq == 0 {
		errs.Add("jsonReq", nil, "required", "field JsonReq is required")
	}
	if provided&providedAllFieldsSoFarJsonVal == 0 {
		errs.Add("jsonVal", nil, "required", "field JsonVal is required")
	}
	if provided&providedAllFieldsSoFarBytesReq == 0 {
		errs.Add("bytesReq", nil, "required", "field BytesReq is required")
	}
	if provided&providedAllFieldsSoFarLtreeField == 0 {
		errs.Add("ltreeField", "", "required", "field LtreeField is required")
	}

	if errs.HasErrors() {
		return input, errs
	}
	return input, nil
}

func (s *AllFieldsSoFarCreate) ToColsVals() (cols []string, vals []any) {
	cols = make([]string, 0, 54)
	vals = make([]any, 0, 54)
	if s.Id != nil {
		cols = append(cols, "id")
		vals = append(vals, *s.Id)
	}
	cols = append(cols, "stringReq")
	vals = append(vals, s.StringReq)
	if s.StringOpt != nil {
		cols = append(cols, "stringOpt")
		vals = append(vals, *s.StringOpt)
	}
	if s.StringDefault != nil {
		cols = append(cols, "stringDefault")
		vals = append(vals, *s.StringDefault)
	}
	cols = append(cols, "stringVarchar")
	vals = append(vals, s.StringVarchar)
	cols = append(cols, "stringChar")
	vals = append(vals, s.StringChar)
	cols = append(cols, "bitVal")
	vals = append(vals, s.BitVal)
	cols = append(cols, "varBitVal")
	vals = append(vals, s.VarBitVal)
	cols = append(cols, "inetVal")
	vals = append(vals, s.InetVal)
	cols = append(cols, "xmlVal")
	vals = append(vals, s.XmlVal)
	cols = append(cols, "cuidDefault")
	if s.CuidDefault != nil {
		vals = append(vals, *s.CuidDefault)
	} else {
		vals = append(vals, generateCUID())
	}
	cols = append(cols, "cuid1Default")
	if s.Cuid1Default != nil {
		vals = append(vals, *s.Cuid1Default)
	} else {
		vals = append(vals, generateCUID())
	}
	cols = append(cols, "cuid2Default")
	if s.Cuid2Default != nil {
		vals = append(vals, *s.Cuid2Default)
	} else {
		vals = append(vals, generateCUID2())
	}
	cols = append(cols, "uuidDefault")
	if s.UuidDefault != nil {
		vals = append(vals, *s.UuidDefault)
	} else {
		vals = append(vals, generateUUID())
	}
	cols = append(cols, "uuid4Default")
	if s.Uuid4Default != nil {
		vals = append(vals, *s.Uuid4Default)
	} else {
		vals = append(vals, generateUUID())
	}
	cols = append(cols, "uuid7Default")
	if s.Uuid7Default != nil {
		vals = append(vals, *s.Uuid7Default)
	} else {
		vals = append(vals, generateUUID7())
	}
	cols = append(cols, "ulidDefault")
	if s.UlidDefault != nil {
		vals = append(vals, *s.UlidDefault)
	} else {
		vals = append(vals, generateULID())
	}
	cols = append(cols, "nanoidDefault")
	if s.NanoidDefault != nil {
		vals = append(vals, *s.NanoidDefault)
	} else {
		vals = append(vals, generateNanoID())
	}
	cols = append(cols, "uuidDb")
	vals = append(vals, s.UuidDb)
	cols = append(cols, "intReq")
	vals = append(vals, s.IntReq)
	if s.IntOpt != nil {
		cols = append(cols, "intOpt")
		vals = append(vals, *s.IntOpt)
	}
	if s.IntDefault != nil {
		cols = append(cols, "intDefault")
		vals = append(vals, *s.IntDefault)
	}
	cols = append(cols, "integerVal")
	vals = append(vals, s.IntegerVal)
	cols = append(cols, "smallInt")
	vals = append(vals, s.SmallInt)
	cols = append(cols, "tinyInt")
	vals = append(vals, s.TinyInt)
	cols = append(cols, "oidVal")
	vals = append(vals, s.OidVal)
	cols = append(cols, "bigIntReq")
	vals = append(vals, s.BigIntReq)
	if s.BigIntOpt != nil {
		cols = append(cols, "bigIntOpt")
		vals = append(vals, *s.BigIntOpt)
	}
	cols = append(cols, "floatReq")
	vals = append(vals, s.FloatReq)
	if s.FloatOpt != nil {
		cols = append(cols, "floatOpt")
		vals = append(vals, *s.FloatOpt)
	}
	cols = append(cols, "realVal")
	vals = append(vals, s.RealVal)
	cols = append(cols, "decimalReq")
	vals = append(vals, s.DecimalReq)
	if s.DecimalOpt != nil {
		cols = append(cols, "decimalOpt")
		vals = append(vals, *s.DecimalOpt)
	}
	cols = append(cols, "decimalPrecise")
	vals = append(vals, s.DecimalPrecise)
	cols = append(cols, "moneyVal")
	vals = append(vals, s.MoneyVal)
	cols = append(cols, "boolReq")
	vals = append(vals, s.BoolReq)
	if s.BoolOpt != nil {
		cols = append(cols, "boolOpt")
		vals = append(vals, *s.BoolOpt)
	}
	if s.BoolDefault != nil {
		cols = append(cols, "boolDefault")
		vals = append(vals, *s.BoolDefault)
	}
	cols = append(cols, "dateTimeReq")
	vals = append(vals, s.DateTimeReq)
	if s.DateTimeOpt != nil {
		cols = append(cols, "dateTimeOpt")
		vals = append(vals, *s.DateTimeOpt)
	}
	cols = append(cols, "dateTimeDefault")
	if s.DateTimeDefault != nil {
		vals = append(vals, *s.DateTimeDefault)
	} else {
		vals = append(vals, time.Now())
	}
	cols = append(cols, "updatedAt")
	vals = append(vals, s.UpdatedAt)
	cols = append(cols, "dateTimeTz")
	vals = append(vals, s.DateTimeTz)
	cols = append(cols, "timestampVal")
	vals = append(vals, s.TimestampVal)
	cols = append(cols, "timeVal")
	vals = append(vals, s.TimeVal)
	cols = append(cols, "timetzVal")
	vals = append(vals, s.TimetzVal)
	cols = append(cols, "jsonReq")
	vals = append(vals, s.JsonReq)
	if s.JsonOpt != nil {
		cols = append(cols, "jsonOpt")
		vals = append(vals, *s.JsonOpt)
	}
	cols = append(cols, "jsonVal")
	vals = append(vals, s.JsonVal)
	cols = append(cols, "bytesReq")
	vals = append(vals, s.BytesReq)
	if s.BytesOpt != nil {
		cols = append(cols, "bytesOpt")
		vals = append(vals, *s.BytesOpt)
	}
	if s.HstoreField != nil {
		cols = append(cols, "hstoreField")
		vals = append(vals, ToHstore(*s.HstoreField))
	}
	cols = append(cols, "ltreeField")
	vals = append(vals, s.LtreeField)
	if s.CitextField != nil {
		cols = append(cols, "citextField")
		vals = append(vals, *s.CitextField)
	}
	return
}

func partitionAllFieldsSoFarInputs(dialect Dialect, inputs []*AllFieldsSoFarCreate) [][]*AllFieldsSoFarCreate {
	if !dialect.SupportsBulkInsert {
		result := make([][]*AllFieldsSoFarCreate, len(inputs))
		for i, input := range inputs {
			result[i] = []*AllFieldsSoFarCreate{input}
		}
		return result
	}

	if !dialect.SupportsDefaultKeyword {
		groups := make(map[uint64][]*AllFieldsSoFarCreate)
		var masks []uint64
		for _, input := range inputs {
			mask := input.colMask()
			if _, exists := groups[mask]; !exists {
				masks = append(masks, mask)
			}
			groups[mask] = append(groups[mask], input)
		}
		result := make([][]*AllFieldsSoFarCreate, len(masks))
		for i, mask := range masks {
			result[i] = groups[mask]
		}
		return result
	}

	return [][]*AllFieldsSoFarCreate{inputs}
}

func (d *AllFieldsSoFarDelegate) executeCreate(ctx context.Context, assignments []FieldAssignment, selects *AllFieldsSoFarSelect, omits *AllFieldsSoFarOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (*AllFieldsSoFar, error) {
	input, err := assignmentsToAllFieldsSoFarCreate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()
	returningCols := selectAllFieldsSoFarCols(selects, omits)

	if len(d.extensions) == 0 {
		hasRelations := selects.hasAnyRelation()
		if hasRelations {
			var res *AllFieldsSoFar
			err = d.client.transaction(ctx, func(txQ *Queries) error {
				var err error
				res, err = txQ.AllFieldsSoFar.runCreate(ctx, cols, vals, returningCols, allFieldsSoFarPKCols, conflictTarget, conflictAction)
				if err != nil {
					return err
				}
				return txQ.AllFieldsSoFar.loadRelations(ctx, []*AllFieldsSoFar{res}, selects)
			})
			return res, err
		}
		return d.runCreate(ctx, cols, vals, returningCols, allFieldsSoFarPKCols, conflictTarget, conflictAction)
	}

	curr := func(c context.Context, args *AllFieldsSoFarCreate) (*AllFieldsSoFar, error) {
		cols, vals := args.ToColsVals()
		returningCols := selectAllFieldsSoFarCols(selects, omits)

		hasRelations := selects.hasAnyRelation()
		var res *AllFieldsSoFar
		var err error
		if hasRelations {
			err = d.client.transaction(c, func(txQ *Queries) error {
				var err error
				res, err = txQ.AllFieldsSoFar.runCreate(c, cols, vals, returningCols, allFieldsSoFarPKCols, conflictTarget, conflictAction)
				if err != nil {
					return err
				}
				return txQ.AllFieldsSoFar.loadRelations(c, []*AllFieldsSoFar{res}, selects)
			})
		} else {
			res, err = d.runCreate(c, cols, vals, returningCols, allFieldsSoFarPKCols, conflictTarget, conflictAction)
		}
		if err != nil {
			return nil, err
		}
		return res, nil
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.Create != nil {
			next, hook := curr, ext.Create
			curr = func(c context.Context, input *AllFieldsSoFarCreate) (*AllFieldsSoFar, error) {
				return hook(c, input, next)
			}
		}
	}

	return curr(ctx, &input)
}

type AllFieldsSoFarCreateManyBuilder struct {
	*CreateManyBuilder[AllFieldsSoFar]
}

func (b *AllFieldsSoFarCreateManyBuilder) OnConflict(target UniqueConstraintTarget) *AllFieldsSoFarConflictBuilder[AllFieldsSoFarCreateManyBuilder] {
	return &AllFieldsSoFarConflictBuilder[AllFieldsSoFarCreateManyBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

type AllFieldsSoFarCreateManyAndReturnBuilder struct {
	*CreateManyAndReturnBuilder[AllFieldsSoFar, AllFieldsSoFarSelect, AllFieldsSoFarOmit]
}

func (b *AllFieldsSoFarCreateManyAndReturnBuilder) OnConflict(target UniqueConstraintTarget) *AllFieldsSoFarConflictBuilder[AllFieldsSoFarCreateManyAndReturnBuilder] {
	return &AllFieldsSoFarConflictBuilder[AllFieldsSoFarCreateManyAndReturnBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

func (d *AllFieldsSoFarDelegate) CreateMany(builders ...*AllFieldsSoFarCreateBuilder) *AllFieldsSoFarCreateManyBuilder {
	records := make([]RecordInput, len(builders))
	for i, b := range builders {
		records[i] = RecordInput{Assignments: b.assignments}
	}
	return &AllFieldsSoFarCreateManyBuilder{
		CreateManyBuilder: &CreateManyBuilder[AllFieldsSoFar]{
			records:  records,
			execFunc: d.executeCreateMany,
		},
	}
}

func (d *AllFieldsSoFarDelegate) CreateManyAndReturn(builders ...*AllFieldsSoFarCreateBuilder) *AllFieldsSoFarCreateManyAndReturnBuilder {
	records := make([]RecordInput, len(builders))
	for i, b := range builders {
		records[i] = RecordInput{Assignments: b.assignments}
	}
	return &AllFieldsSoFarCreateManyAndReturnBuilder{
		CreateManyAndReturnBuilder: &CreateManyAndReturnBuilder[AllFieldsSoFar, AllFieldsSoFarSelect, AllFieldsSoFarOmit]{
			records:  records,
			execFunc: d.executeCreateManyAndReturn,
		},
	}
}

func (d *AllFieldsSoFarDelegate) executeCreateMany(ctx context.Context, records []RecordInput, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (int64, error) {
	inputs := make([]*AllFieldsSoFarCreate, len(records))
	for i, rec := range records {
		input, err := assignmentsToAllFieldsSoFarCreate(rec.Assignments)
		if err != nil {
			return 0, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
		inputs[i] = &input
	}

	if len(d.extensions) == 0 {
		return d.runCreateMany(ctx, inputs, conflictTarget, conflictAction)
	}

	curr := func(c context.Context, args []*AllFieldsSoFarCreate) (int64, error) {
		return d.runCreateMany(c, args, conflictTarget, conflictAction)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.CreateMany != nil {
			next, hook := curr, ext.CreateMany
			curr = func(c context.Context, inputs []*AllFieldsSoFarCreate) (int64, error) {
				return hook(c, inputs, next)
			}
		}
	}

	return curr(ctx, inputs)
}

func (d *AllFieldsSoFarDelegate) executeCreateManyAndReturn(ctx context.Context, records []RecordInput, selects *AllFieldsSoFarSelect, omits *AllFieldsSoFarOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) ([]*AllFieldsSoFar, error) {
	inputs := make([]*AllFieldsSoFarCreate, len(records))
	for i, rec := range records {
		input, err := assignmentsToAllFieldsSoFarCreate(rec.Assignments)
		if err != nil {
			return nil, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
		inputs[i] = &input
	}

	if len(d.extensions) == 0 {
		hasRelations := selects != nil && selects.hasAnyRelation()
		if hasRelations {
			var res []*AllFieldsSoFar
			err := d.client.transaction(ctx, func(txQ *Queries) error {
				var err error
				res, err = txQ.AllFieldsSoFar.runCreateManyAndReturn(ctx, inputs, selects, omits, conflictTarget, conflictAction)
				if err != nil {
					return err
				}
				return txQ.AllFieldsSoFar.loadRelations(ctx, res, selects)
			})
			return res, err
		}
		return d.runCreateManyAndReturn(ctx, inputs, selects, omits, conflictTarget, conflictAction)
	}

	curr := func(c context.Context, args []*AllFieldsSoFarCreate) ([]*AllFieldsSoFar, error) {
		hasRelations := selects != nil && selects.hasAnyRelation()
		if hasRelations {
			var res []*AllFieldsSoFar
			err := d.client.transaction(c, func(txQ *Queries) error {
				var err error
				res, err = txQ.AllFieldsSoFar.runCreateManyAndReturn(c, args, selects, omits, conflictTarget, conflictAction)
				if err != nil {
					return err
				}
				return txQ.AllFieldsSoFar.loadRelations(c, res, selects)
			})
			return res, err
		}
		return d.runCreateManyAndReturn(c, args, selects, omits, conflictTarget, conflictAction)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.CreateManyAndReturn != nil {
			next, hook := curr, ext.CreateManyAndReturn
			curr = func(c context.Context, inputs []*AllFieldsSoFarCreate) ([]*AllFieldsSoFar, error) {
				return hook(c, inputs, next)
			}
		}
	}

	return curr(ctx, inputs)
}

func (d *AllFieldsSoFarDelegate) runCreate(
	ctx context.Context,
	cols []string,
	vals []any,
	returningCols []string,
	pkCols []string,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) (*AllFieldsSoFar, error) {
	query, clauseArgs := buildSingleInsertSQL(d.client, "AllFieldsSoFar", cols, returningCols, pkCols, conflictTarget, conflictAction, len(vals))
	if len(clauseArgs) > 0 {
		vals = append(vals, clauseArgs...)
	}

	var res AllFieldsSoFar
	if d.client.dialect.SupportsInsertReturning {
		rows, err := d.client.query(ctx, query, vals...)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		if rows.Next() {
			if err := rows.Scan(res.ScanFields(returningCols)...); err != nil {
				return nil, err
			}
			return &res, nil
		}
		return nil, rows.Err()
	}

	return d.runCreateFallback(ctx, query, vals, cols, returningCols, pkCols)
}

func (d *AllFieldsSoFarDelegate) runCreateFallback(ctx context.Context, query string, vals []any, cols []string, returningCols []string, pkCols []string) (*AllFieldsSoFar, error) {
	result, err := d.client.exec(ctx, query, vals...)
	if err != nil {
		return nil, err
	}

	var pkVals []any
	for _, pkCol := range pkCols {
		var val any
		for i, c := range cols {
			if c == pkCol {
				val = vals[i]
				break
			}
		}
		if val == nil && len(pkCols) == 1 {
			lastID, err := result.LastInsertId()
			if err != nil {
				return nil, err
			}
			val = lastID
		}
		pkVals = append(pkVals, val)
	}

	var selectSb strings.Builder
	selectSb.Grow(64 + len(returningCols)*15 + len("AllFieldsSoFar") + len(pkCols)*15)
	selectSb.WriteString("SELECT ")
	for i, col := range returningCols {
		if i > 0 {
			selectSb.WriteString(", ")
		}
		d.client.dialect.WriteQuotedIdent(&selectSb, col)
	}
	selectSb.WriteString(" FROM ")
	d.client.dialect.WriteQuotedIdent(&selectSb, "AllFieldsSoFar")
	selectSb.WriteString(" WHERE ")
	for i, pkCol := range pkCols {
		if i > 0 {
			selectSb.WriteString(" AND ")
		}
		d.client.dialect.WriteQuotedIdent(&selectSb, pkCol)
		selectSb.WriteString(" = ")
		d.client.dialect.WritePlaceholder(&selectSb, i+1)
	}

	rows, err := d.client.query(ctx, selectSb.String(), pkVals...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res AllFieldsSoFar
	if rows.Next() {
		if err := rows.Scan(res.ScanFields(returningCols)...); err != nil {
			return nil, err
		}
		return &res, nil
	}
	return nil, rows.Err()
}

func (d *AllFieldsSoFarDelegate) buildBulkInsertSQL(q *Queries, batch []*AllFieldsSoFarCreate, paramStartIdx int) (cols []string, vals []any, queryStr string) {
	var colMask uint64
	for _, input := range batch {
		colMask |= input.colMask()
	}

	cols = make([]string, 0, 54)
	for i, c := range allFieldsSoFarDefaultCols {
		if colMask&(1<<i) != 0 {
			cols = append(cols, c)
		}
	}

	vals = make([]any, 0, len(batch)*len(cols))
	var sb strings.Builder
	sb.Grow(128 + len(batch)*len(cols)*10)
	sb.WriteString("INSERT INTO ")
	q.dialect.WriteQuotedIdent(&sb, "AllFieldsSoFar")
	sb.WriteString(" (")
	for i, col := range cols {
		if i > 0 {
			sb.WriteString(", ")
		}
		q.dialect.WriteQuotedIdent(&sb, col)
	}
	sb.WriteString(") VALUES ")

	paramIdx := paramStartIdx
	for ri, input := range batch {
		if ri > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString("(")
		for j, col := range cols {
			if j > 0 {
				sb.WriteString(", ")
			}
			writeDefault := false
			switch col {
			case "id":
				if input.Id != nil {
					vals = append(vals, *input.Id)
				} else {
					writeDefault = true
				}
			case "stringReq":
				vals = append(vals, input.StringReq)
			case "stringOpt":
				if input.StringOpt != nil {
					vals = append(vals, *input.StringOpt)
				} else {
					writeDefault = true
				}
			case "stringDefault":
				if input.StringDefault != nil {
					vals = append(vals, *input.StringDefault)
				} else {
					writeDefault = true
				}
			case "stringVarchar":
				vals = append(vals, input.StringVarchar)
			case "stringChar":
				vals = append(vals, input.StringChar)
			case "bitVal":
				vals = append(vals, input.BitVal)
			case "varBitVal":
				vals = append(vals, input.VarBitVal)
			case "inetVal":
				vals = append(vals, input.InetVal)
			case "xmlVal":
				vals = append(vals, input.XmlVal)
			case "cuidDefault":
				if input.CuidDefault != nil {
					vals = append(vals, *input.CuidDefault)
				} else {
					vals = append(vals, generateCUID())
				}
			case "cuid1Default":
				if input.Cuid1Default != nil {
					vals = append(vals, *input.Cuid1Default)
				} else {
					vals = append(vals, generateCUID())
				}
			case "cuid2Default":
				if input.Cuid2Default != nil {
					vals = append(vals, *input.Cuid2Default)
				} else {
					vals = append(vals, generateCUID2())
				}
			case "uuidDefault":
				if input.UuidDefault != nil {
					vals = append(vals, *input.UuidDefault)
				} else {
					vals = append(vals, generateUUID())
				}
			case "uuid4Default":
				if input.Uuid4Default != nil {
					vals = append(vals, *input.Uuid4Default)
				} else {
					vals = append(vals, generateUUID())
				}
			case "uuid7Default":
				if input.Uuid7Default != nil {
					vals = append(vals, *input.Uuid7Default)
				} else {
					vals = append(vals, generateUUID7())
				}
			case "ulidDefault":
				if input.UlidDefault != nil {
					vals = append(vals, *input.UlidDefault)
				} else {
					vals = append(vals, generateULID())
				}
			case "nanoidDefault":
				if input.NanoidDefault != nil {
					vals = append(vals, *input.NanoidDefault)
				} else {
					vals = append(vals, generateNanoID())
				}
			case "uuidDb":
				vals = append(vals, input.UuidDb)
			case "intReq":
				vals = append(vals, input.IntReq)
			case "intOpt":
				if input.IntOpt != nil {
					vals = append(vals, *input.IntOpt)
				} else {
					writeDefault = true
				}
			case "intDefault":
				if input.IntDefault != nil {
					vals = append(vals, *input.IntDefault)
				} else {
					writeDefault = true
				}
			case "integerVal":
				vals = append(vals, input.IntegerVal)
			case "smallInt":
				vals = append(vals, input.SmallInt)
			case "tinyInt":
				vals = append(vals, input.TinyInt)
			case "oidVal":
				vals = append(vals, input.OidVal)
			case "bigIntReq":
				vals = append(vals, input.BigIntReq)
			case "bigIntOpt":
				if input.BigIntOpt != nil {
					vals = append(vals, *input.BigIntOpt)
				} else {
					writeDefault = true
				}
			case "floatReq":
				vals = append(vals, input.FloatReq)
			case "floatOpt":
				if input.FloatOpt != nil {
					vals = append(vals, *input.FloatOpt)
				} else {
					writeDefault = true
				}
			case "realVal":
				vals = append(vals, input.RealVal)
			case "decimalReq":
				vals = append(vals, input.DecimalReq)
			case "decimalOpt":
				if input.DecimalOpt != nil {
					vals = append(vals, *input.DecimalOpt)
				} else {
					writeDefault = true
				}
			case "decimalPrecise":
				vals = append(vals, input.DecimalPrecise)
			case "moneyVal":
				vals = append(vals, input.MoneyVal)
			case "boolReq":
				vals = append(vals, input.BoolReq)
			case "boolOpt":
				if input.BoolOpt != nil {
					vals = append(vals, *input.BoolOpt)
				} else {
					writeDefault = true
				}
			case "boolDefault":
				if input.BoolDefault != nil {
					vals = append(vals, *input.BoolDefault)
				} else {
					writeDefault = true
				}
			case "dateTimeReq":
				vals = append(vals, input.DateTimeReq)
			case "dateTimeOpt":
				if input.DateTimeOpt != nil {
					vals = append(vals, *input.DateTimeOpt)
				} else {
					writeDefault = true
				}
			case "dateTimeDefault":
				if input.DateTimeDefault != nil {
					vals = append(vals, *input.DateTimeDefault)
				} else {
					vals = append(vals, time.Now())
				}
			case "updatedAt":
				vals = append(vals, input.UpdatedAt)
			case "dateTimeTz":
				vals = append(vals, input.DateTimeTz)
			case "timestampVal":
				vals = append(vals, input.TimestampVal)
			case "timeVal":
				vals = append(vals, input.TimeVal)
			case "timetzVal":
				vals = append(vals, input.TimetzVal)
			case "jsonReq":
				vals = append(vals, input.JsonReq)
			case "jsonOpt":
				if input.JsonOpt != nil {
					vals = append(vals, *input.JsonOpt)
				} else {
					writeDefault = true
				}
			case "jsonVal":
				vals = append(vals, input.JsonVal)
			case "bytesReq":
				vals = append(vals, input.BytesReq)
			case "bytesOpt":
				if input.BytesOpt != nil {
					vals = append(vals, *input.BytesOpt)
				} else {
					writeDefault = true
				}
			case "hstoreField":
				if input.HstoreField != nil {
					vals = append(vals, ToHstore(*input.HstoreField))
				} else {
					writeDefault = true
				}
			case "ltreeField":
				vals = append(vals, input.LtreeField)
			case "citextField":
				if input.CitextField != nil {
					vals = append(vals, *input.CitextField)
				} else {
					writeDefault = true
				}
			}
			if writeDefault {
				sb.WriteString("DEFAULT")
			} else {
				q.dialect.WritePlaceholder(&sb, paramIdx)
				paramIdx++
			}
		}
		sb.WriteString(")")
	}
	queryStr = sb.String()
	return cols, vals, queryStr
}

func (d *AllFieldsSoFarDelegate) runCreateMany(ctx context.Context, inputs []*AllFieldsSoFarCreate, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (int64, error) {
	if len(inputs) == 0 {
		return 0, nil
	}

	batches := partitionAllFieldsSoFarInputs(d.client.dialect, inputs)

	var count int64
	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)

		var conflictCols []string
		if conflictTarget != nil {
			conflictCols = conflictTarget.UniqueColumns()
		}
		var nonConflictCols []string
		if conflictAction != nil && conflictAction.Type == ConflictActionUpdateNewValues {
			nonConflictCols = computeNonConflictCols(cols, conflictCols, allFieldsSoFarPKCols)
		}
		clause, clauseArgs := d.client.dialect.BuildConflictClause(conflictCols, conflictAction, nonConflictCols, len(vals)+1)
		queryStr += clause
		vals = append(vals, clauseArgs...)

		result, err := d.client.exec(ctx, queryStr, vals...)
		if err != nil {
			return 0, err
		}
		affected, err := result.RowsAffected()
		if err != nil {
			return 0, err
		}
		count += affected
	}
	return count, nil
}

func (d *AllFieldsSoFarDelegate) runCreateManyAndReturn(
	ctx context.Context,
	inputs []*AllFieldsSoFarCreate,
	selects *AllFieldsSoFarSelect,
	omits *AllFieldsSoFarOmit,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) ([]*AllFieldsSoFar, error) {
	if len(inputs) == 0 {
		return nil, nil
	}

	batches := partitionAllFieldsSoFarInputs(d.client.dialect, inputs)
	returningCols := selectAllFieldsSoFarCols(selects, omits)
	hasRelations := selects != nil && selects.hasAnyRelation()

	recordsOut := make([]*AllFieldsSoFar, 0, len(inputs))

	runBatch := func(txQ *Queries, batch []*AllFieldsSoFarCreate) error {
		cols, vals, queryStr := d.buildBulkInsertSQL(txQ, batch, 1)

		var conflictCols []string
		if conflictTarget != nil {
			conflictCols = conflictTarget.UniqueColumns()
		}
		var nonConflictCols []string
		if conflictAction != nil && conflictAction.Type == ConflictActionUpdateNewValues {
			nonConflictCols = computeNonConflictCols(cols, conflictCols, allFieldsSoFarPKCols)
		}
		clause, clauseArgs := txQ.dialect.BuildConflictClause(conflictCols, conflictAction, nonConflictCols, len(vals)+1)
		queryStr += clause
		vals = append(vals, clauseArgs...)

		if txQ.dialect.SupportsInsertReturning && len(returningCols) > 0 {
			var retSb strings.Builder
			retSb.Grow(12 + len(returningCols)*15)
			retSb.WriteString(" RETURNING ")
			for i, col := range returningCols {
				if i > 0 {
					retSb.WriteString(", ")
				}
				txQ.dialect.WriteQuotedIdent(&retSb, col)
			}
			queryStr += retSb.String()
			rows, err := txQ.query(ctx, queryStr, vals...)
			if err != nil {
				return err
			}
			defer rows.Close()

			for rows.Next() {
				var res AllFieldsSoFar
				if err := rows.Scan(res.ScanFields(returningCols)...); err != nil {
					return err
				}
				recordsOut = append(recordsOut, &res)
			}
			return rows.Err()
		}

		// Fallback for dialects without RETURNING (MySQL)
		result, err := txQ.exec(ctx, queryStr, vals...)
		if err != nil {
			return err
		}

		// We need to fetch the inserted records for this batch
		// Note: MySQL bulk inserts only return the ID of the FIRST inserted row
		lastID, err := result.LastInsertId()
		if err != nil {
			return err
		}

		// Query back the rows by IDs (assuming autoincrement ID and single PK)
		// If composite PK, it's more complex, but this is a standard fallback
		var selectSb strings.Builder
		selectSb.Grow(64 + len(returningCols)*15 + len("AllFieldsSoFar") + len(batch)*15)
		selectSb.WriteString("SELECT ")
		for i, col := range returningCols {
			if i > 0 {
				selectSb.WriteString(", ")
			}
			txQ.dialect.WriteQuotedIdent(&selectSb, col)
		}
		selectSb.WriteString(" FROM ")
		txQ.dialect.WriteQuotedIdent(&selectSb, "AllFieldsSoFar")
		selectSb.WriteString(" WHERE ")
		txQ.dialect.WriteQuotedIdent(&selectSb, allFieldsSoFarPKCols[0])
		selectSb.WriteString(" >= ")
		txQ.dialect.WritePlaceholder(&selectSb, 1)
		selectSb.WriteString(" AND ")
		txQ.dialect.WriteQuotedIdent(&selectSb, allFieldsSoFarPKCols[0])
		selectSb.WriteString(" < ")
		txQ.dialect.WritePlaceholder(&selectSb, 2)

		rows, err := txQ.query(ctx, selectSb.String(), lastID, lastID+int64(len(batch)))
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var res AllFieldsSoFar
			if err := rows.Scan(res.ScanFields(returningCols)...); err != nil {
				return err
			}
			recordsOut = append(recordsOut, &res)
		}
		return rows.Err()
	}

	// Always wrap in transaction if we have multiple batches OR if we need to load relations
	if len(batches) > 1 || hasRelations || !d.client.dialect.SupportsInsertReturning {
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			for _, batch := range batches {
				if err := runBatch(txQ, batch); err != nil {
					return err
				}
			}
			if hasRelations {
				return txQ.AllFieldsSoFar.loadRelations(ctx, recordsOut, selects)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	} else {
		if err := runBatch(d.client, batches[0]); err != nil {
			return nil, err
		}
	}

	return recordsOut, nil
}

type AllFieldsSoFarConflictBuilder[B any] struct {
	builder        *B
	setAction      func(ConflictAction, UniqueConstraintTarget)
	conflictTarget UniqueConstraintTarget
}

func (cb *AllFieldsSoFarConflictBuilder[B]) Ignore() *B {
	cb.setAction(ConflictAction{Type: ConflictActionIgnore}, cb.conflictTarget)
	return cb.builder
}

func (cb *AllFieldsSoFarConflictBuilder[B]) UpdateNewValues() *B {
	cb.setAction(ConflictAction{Type: ConflictActionUpdateNewValues}, cb.conflictTarget)
	return cb.builder
}

func (cb *AllFieldsSoFarConflictBuilder[B]) Update(fn func(u *AllFieldsSoFarUpsert)) *B {
	var up ConflictUpdate
	u := newAllFieldsSoFarUpsert(&up)
	fn(u)
	cb.setAction(ConflictAction{
		Type:        ConflictActionUpdateCustom,
		Assignments: up.assignments,
		Args:        up.args,
	}, cb.conflictTarget)
	return cb.builder
}

type AllFieldsSoFarUpsert struct {
	Id              numericFieldUpsert[int32]
	StringReq       fieldUpsert[string]
	StringOpt       fieldUpsert[*string]
	StringDefault   fieldUpsert[string]
	StringVarchar   fieldUpsert[string]
	StringChar      fieldUpsert[string]
	BitVal          fieldUpsert[string]
	VarBitVal       fieldUpsert[string]
	InetVal         fieldUpsert[string]
	XmlVal          fieldUpsert[string]
	CuidDefault     fieldUpsert[string]
	Cuid1Default    fieldUpsert[string]
	Cuid2Default    fieldUpsert[string]
	UuidDefault     fieldUpsert[string]
	Uuid4Default    fieldUpsert[string]
	Uuid7Default    fieldUpsert[string]
	UlidDefault     fieldUpsert[string]
	NanoidDefault   fieldUpsert[string]
	UuidDb          fieldUpsert[string]
	IntReq          numericFieldUpsert[int32]
	IntOpt          numericFieldUpsert[*int32]
	IntDefault      numericFieldUpsert[int32]
	IntegerVal      numericFieldUpsert[int32]
	SmallInt        numericFieldUpsert[int32]
	TinyInt         numericFieldUpsert[int32]
	OidVal          numericFieldUpsert[int32]
	BigIntReq       numericFieldUpsert[int64]
	BigIntOpt       numericFieldUpsert[*int64]
	FloatReq        numericFieldUpsert[float64]
	FloatOpt        numericFieldUpsert[*float64]
	RealVal         numericFieldUpsert[float64]
	DecimalReq      numericFieldUpsert[string]
	DecimalOpt      numericFieldUpsert[*string]
	DecimalPrecise  numericFieldUpsert[string]
	MoneyVal        numericFieldUpsert[string]
	BoolReq         fieldUpsert[bool]
	BoolOpt         fieldUpsert[*bool]
	BoolDefault     fieldUpsert[bool]
	DateTimeReq     fieldUpsert[time.Time]
	DateTimeOpt     fieldUpsert[*time.Time]
	DateTimeDefault fieldUpsert[time.Time]
	UpdatedAt       fieldUpsert[time.Time]
	DateTimeTz      fieldUpsert[time.Time]
	TimestampVal    fieldUpsert[time.Time]
	TimeVal         fieldUpsert[time.Time]
	TimetzVal       fieldUpsert[time.Time]
	JsonReq         fieldUpsert[json.RawMessage]
	JsonOpt         fieldUpsert[*json.RawMessage]
	JsonVal         fieldUpsert[json.RawMessage]
	BytesReq        fieldUpsert[[]byte]
	BytesOpt        fieldUpsert[*[]byte]
	HstoreField     fieldUpsert[*map[string]*string]
	LtreeField      fieldUpsert[string]
	CitextField     fieldUpsert[*string]
}

func newAllFieldsSoFarUpsert(up *ConflictUpdate) *AllFieldsSoFarUpsert {
	return &AllFieldsSoFarUpsert{
		Id: numericFieldUpsert[int32]{
			fieldUpsert: fieldUpsert[int32]{column: "id", update: up},
			tableName:   "AllFieldsSoFar",
		},
		StringReq:     fieldUpsert[string]{column: "stringReq", update: up},
		StringOpt:     fieldUpsert[*string]{column: "stringOpt", update: up},
		StringDefault: fieldUpsert[string]{column: "stringDefault", update: up},
		StringVarchar: fieldUpsert[string]{column: "stringVarchar", update: up},
		StringChar:    fieldUpsert[string]{column: "stringChar", update: up},
		BitVal:        fieldUpsert[string]{column: "bitVal", update: up},
		VarBitVal:     fieldUpsert[string]{column: "varBitVal", update: up},
		InetVal:       fieldUpsert[string]{column: "inetVal", update: up},
		XmlVal:        fieldUpsert[string]{column: "xmlVal", update: up},
		CuidDefault:   fieldUpsert[string]{column: "cuidDefault", update: up},
		Cuid1Default:  fieldUpsert[string]{column: "cuid1Default", update: up},
		Cuid2Default:  fieldUpsert[string]{column: "cuid2Default", update: up},
		UuidDefault:   fieldUpsert[string]{column: "uuidDefault", update: up},
		Uuid4Default:  fieldUpsert[string]{column: "uuid4Default", update: up},
		Uuid7Default:  fieldUpsert[string]{column: "uuid7Default", update: up},
		UlidDefault:   fieldUpsert[string]{column: "ulidDefault", update: up},
		NanoidDefault: fieldUpsert[string]{column: "nanoidDefault", update: up},
		UuidDb:        fieldUpsert[string]{column: "uuidDb", update: up},
		IntReq: numericFieldUpsert[int32]{
			fieldUpsert: fieldUpsert[int32]{column: "intReq", update: up},
			tableName:   "AllFieldsSoFar",
		},
		IntOpt: numericFieldUpsert[*int32]{
			fieldUpsert: fieldUpsert[*int32]{column: "intOpt", update: up},
			tableName:   "AllFieldsSoFar",
		},
		IntDefault: numericFieldUpsert[int32]{
			fieldUpsert: fieldUpsert[int32]{column: "intDefault", update: up},
			tableName:   "AllFieldsSoFar",
		},
		IntegerVal: numericFieldUpsert[int32]{
			fieldUpsert: fieldUpsert[int32]{column: "integerVal", update: up},
			tableName:   "AllFieldsSoFar",
		},
		SmallInt: numericFieldUpsert[int32]{
			fieldUpsert: fieldUpsert[int32]{column: "smallInt", update: up},
			tableName:   "AllFieldsSoFar",
		},
		TinyInt: numericFieldUpsert[int32]{
			fieldUpsert: fieldUpsert[int32]{column: "tinyInt", update: up},
			tableName:   "AllFieldsSoFar",
		},
		OidVal: numericFieldUpsert[int32]{
			fieldUpsert: fieldUpsert[int32]{column: "oidVal", update: up},
			tableName:   "AllFieldsSoFar",
		},
		BigIntReq: numericFieldUpsert[int64]{
			fieldUpsert: fieldUpsert[int64]{column: "bigIntReq", update: up},
			tableName:   "AllFieldsSoFar",
		},
		BigIntOpt: numericFieldUpsert[*int64]{
			fieldUpsert: fieldUpsert[*int64]{column: "bigIntOpt", update: up},
			tableName:   "AllFieldsSoFar",
		},
		FloatReq: numericFieldUpsert[float64]{
			fieldUpsert: fieldUpsert[float64]{column: "floatReq", update: up},
			tableName:   "AllFieldsSoFar",
		},
		FloatOpt: numericFieldUpsert[*float64]{
			fieldUpsert: fieldUpsert[*float64]{column: "floatOpt", update: up},
			tableName:   "AllFieldsSoFar",
		},
		RealVal: numericFieldUpsert[float64]{
			fieldUpsert: fieldUpsert[float64]{column: "realVal", update: up},
			tableName:   "AllFieldsSoFar",
		},
		DecimalReq: numericFieldUpsert[string]{
			fieldUpsert: fieldUpsert[string]{column: "decimalReq", update: up},
			tableName:   "AllFieldsSoFar",
		},
		DecimalOpt: numericFieldUpsert[*string]{
			fieldUpsert: fieldUpsert[*string]{column: "decimalOpt", update: up},
			tableName:   "AllFieldsSoFar",
		},
		DecimalPrecise: numericFieldUpsert[string]{
			fieldUpsert: fieldUpsert[string]{column: "decimalPrecise", update: up},
			tableName:   "AllFieldsSoFar",
		},
		MoneyVal: numericFieldUpsert[string]{
			fieldUpsert: fieldUpsert[string]{column: "moneyVal", update: up},
			tableName:   "AllFieldsSoFar",
		},
		BoolReq:         fieldUpsert[bool]{column: "boolReq", update: up},
		BoolOpt:         fieldUpsert[*bool]{column: "boolOpt", update: up},
		BoolDefault:     fieldUpsert[bool]{column: "boolDefault", update: up},
		DateTimeReq:     fieldUpsert[time.Time]{column: "dateTimeReq", update: up},
		DateTimeOpt:     fieldUpsert[*time.Time]{column: "dateTimeOpt", update: up},
		DateTimeDefault: fieldUpsert[time.Time]{column: "dateTimeDefault", update: up},
		UpdatedAt:       fieldUpsert[time.Time]{column: "updatedAt", update: up},
		DateTimeTz:      fieldUpsert[time.Time]{column: "dateTimeTz", update: up},
		TimestampVal:    fieldUpsert[time.Time]{column: "timestampVal", update: up},
		TimeVal:         fieldUpsert[time.Time]{column: "timeVal", update: up},
		TimetzVal:       fieldUpsert[time.Time]{column: "timetzVal", update: up},
		JsonReq:         fieldUpsert[json.RawMessage]{column: "jsonReq", update: up},
		JsonOpt:         fieldUpsert[*json.RawMessage]{column: "jsonOpt", update: up},
		JsonVal:         fieldUpsert[json.RawMessage]{column: "jsonVal", update: up},
		BytesReq:        fieldUpsert[[]byte]{column: "bytesReq", update: up},
		BytesOpt:        fieldUpsert[*[]byte]{column: "bytesOpt", update: up},
		HstoreField:     fieldUpsert[*map[string]*string]{column: "hstoreField", update: up},
		LtreeField:      fieldUpsert[string]{column: "ltreeField", update: up},
		CitextField:     fieldUpsert[*string]{column: "citextField", update: up},
	}
}
func (d *AllFieldsSoFarDelegate) FindUnique(where UniquePredicate[AllFieldsSoFar], additional ...PredicateOf[AllFieldsSoFar]) *FindUniqueBuilder[AllFieldsSoFar, AllFieldsSoFarSelect, AllFieldsSoFarOmit] {
	return &FindUniqueBuilder[AllFieldsSoFar, AllFieldsSoFarSelect, AllFieldsSoFarOmit]{
		where:      where,
		additional: additional,
		execFunc:   d.executeFindUnique,
	}
}

func (d *AllFieldsSoFarDelegate) FindFirst(preds ...PredicateOf[AllFieldsSoFar]) *FindFirstBuilder[AllFieldsSoFar, AllFieldsSoFarSelect, AllFieldsSoFarOmit] {
	return &FindFirstBuilder[AllFieldsSoFar, AllFieldsSoFarSelect, AllFieldsSoFarOmit]{
		where:    preds,
		execFunc: d.executeFindFirst,
	}
}

func (d *AllFieldsSoFarDelegate) FindMany(preds ...PredicateOf[AllFieldsSoFar]) *FindManyBuilder[AllFieldsSoFar, AllFieldsSoFarSelect, AllFieldsSoFarOmit] {
	return &FindManyBuilder[AllFieldsSoFar, AllFieldsSoFarSelect, AllFieldsSoFarOmit]{
		where:    preds,
		execFunc: d.executeFindMany,
	}
}

func (d *AllFieldsSoFarDelegate) executeFindUnique(ctx context.Context, where UniquePredicate[AllFieldsSoFar], additional []PredicateOf[AllFieldsSoFar], selects *AllFieldsSoFarSelect, omits *AllFieldsSoFarOmit) (*AllFieldsSoFar, error) {
	if len(d.extensions) == 0 {
		return d.runFindUnique(ctx, where, additional, selects, omits)
	}

	curr := func(c context.Context, w UniquePredicate[AllFieldsSoFar], add []PredicateOf[AllFieldsSoFar], sel *AllFieldsSoFarSelect, o *AllFieldsSoFarOmit) (*AllFieldsSoFar, error) {
		return d.runFindUnique(c, w, add, sel, o)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.FindUnique != nil {
			next, hook := curr, ext.FindUnique
			curr = func(c context.Context, w UniquePredicate[AllFieldsSoFar], add []PredicateOf[AllFieldsSoFar], sel *AllFieldsSoFarSelect, o *AllFieldsSoFarOmit) (*AllFieldsSoFar, error) {
				return hook(c, w, add, sel, o, next)
			}
		}
	}

	return curr(ctx, where, additional, selects, omits)
}

func (d *AllFieldsSoFarDelegate) executeFindFirst(
	ctx context.Context,
	params QueryParams[AllFieldsSoFar],
	selects *AllFieldsSoFarSelect,
	omits *AllFieldsSoFarOmit,
) (*AllFieldsSoFar, error) {
	if len(d.extensions) == 0 {
		return d.runFindFirst(ctx, params, selects, omits)
	}

	curr := func(c context.Context, p QueryParams[AllFieldsSoFar], sel *AllFieldsSoFarSelect, o *AllFieldsSoFarOmit) (*AllFieldsSoFar, error) {
		return d.runFindFirst(c, p, sel, o)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.FindFirst != nil {
			next, hook := curr, ext.FindFirst
			curr = func(c context.Context, p QueryParams[AllFieldsSoFar], sel *AllFieldsSoFarSelect, o *AllFieldsSoFarOmit) (*AllFieldsSoFar, error) {
				return hook(c, p, sel, o, next)
			}
		}
	}

	return curr(ctx, params, selects, omits)
}

func (d *AllFieldsSoFarDelegate) executeFindMany(
	ctx context.Context,
	params QueryParams[AllFieldsSoFar],
	selects *AllFieldsSoFarSelect,
	omits *AllFieldsSoFarOmit,
) ([]*AllFieldsSoFar, error) {
	if len(d.extensions) == 0 {
		return d.runFindMany(ctx, params, selects, omits)
	}

	curr := func(c context.Context, p QueryParams[AllFieldsSoFar], sel *AllFieldsSoFarSelect, o *AllFieldsSoFarOmit) ([]*AllFieldsSoFar, error) {
		return d.runFindMany(c, p, sel, o)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.FindMany != nil {
			next, hook := curr, ext.FindMany
			curr = func(c context.Context, p QueryParams[AllFieldsSoFar], sel *AllFieldsSoFarSelect, o *AllFieldsSoFarOmit) ([]*AllFieldsSoFar, error) {
				return hook(c, p, sel, o, next)
			}
		}
	}

	return curr(ctx, params, selects, omits)
}

func (d *AllFieldsSoFarDelegate) runFindUnique(ctx context.Context, where UniquePredicate[AllFieldsSoFar], additional []PredicateOf[AllFieldsSoFar], selects *AllFieldsSoFarSelect, omits *AllFieldsSoFarOmit) (*AllFieldsSoFar, error) {
	if err := where.Validate(); err != nil {
		return nil, err
	}
	for _, p := range additional {
		if p != nil {
			if err := p.Validate(); err != nil {
				return nil, err
			}
		}
	}
	allPreds := append([]PredicateOf[AllFieldsSoFar]{where}, additional...)
	whereClause, vals, _ := CompilePredicates(d.client.dialect, allPreds)
	if whereClause != "" {
		whereClause = " WHERE " + whereClause
	}
	returningCols := selectAllFieldsSoFarCols(selects, omits)

	var res *AllFieldsSoFar
	var err error
	if selects.hasAnyRelation() {
		err = d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = txQ.AllFieldsSoFar.queryOne(ctx, whereClause, "", vals, returningCols, nil)
			if err != nil {
				return err
			}
			if res == nil {
				return nil
			}
			return txQ.AllFieldsSoFar.loadRelations(ctx, []*AllFieldsSoFar{res}, selects)
		})
	} else {
		res, err = d.queryOne(ctx, whereClause, "", vals, returningCols, nil)
	}
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (d *AllFieldsSoFarDelegate) runFindFirst(
	ctx context.Context,
	params QueryParams[AllFieldsSoFar],
	selects *AllFieldsSoFarSelect,
	omits *AllFieldsSoFarOmit,
) (*AllFieldsSoFar, error) {
	for _, pr := range params.Where {
		if pr != nil {
			if err := pr.Validate(); err != nil {
				return nil, err
			}
		}
	}
	whereClause, vals, nextIdx := CompilePredicates(d.client.dialect, params.Where)
	isCursorQuery := (params.Cursor.Data.Column != "" || len(params.Cursor.Data.Children) > 0)
	if isCursorQuery {
		cClause, cVals, err := compileCursorClause(d.client.dialect, params.Cursor, params.OrderBy, allFieldsSoFarPKCols, allFieldsSoFarUniqueCols, "AllFieldsSoFar", nextIdx, params.Take)
		if err != nil {
			return nil, err
		}
		if cClause != "" {
			if whereClause == "" {
				whereClause = cClause
			} else {
				whereClause = "(" + whereClause + ") AND " + cClause
			}
			vals = append(vals, cVals...)
		}
	}
	if whereClause != "" {
		whereClause = " WHERE " + whereClause
	}
	orderByClause := formatOrderBySQL(d.client.dialect, params.OrderBy, allFieldsSoFarPKCols, allFieldsSoFarUniqueCols, isCursorQuery, params.Take)
	returningCols := selectAllFieldsSoFarCols(selects, omits)

	var res *AllFieldsSoFar
	var err error
	if selects.hasAnyRelation() {
		err = d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = txQ.AllFieldsSoFar.queryOne(ctx, whereClause, orderByClause, vals, returningCols, params.Skip)
			if err != nil {
				return err
			}
			if res == nil {
				return nil
			}
			return txQ.AllFieldsSoFar.loadRelations(ctx, []*AllFieldsSoFar{res}, selects)
		})
	} else {
		res, err = d.queryOne(ctx, whereClause, orderByClause, vals, returningCols, params.Skip)
	}
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (d *AllFieldsSoFarDelegate) runFindMany(
	ctx context.Context,
	params QueryParams[AllFieldsSoFar],
	selects *AllFieldsSoFarSelect,
	omits *AllFieldsSoFarOmit,
) ([]*AllFieldsSoFar, error) {
	for _, pr := range params.Where {
		if pr != nil {
			if err := pr.Validate(); err != nil {
				return nil, err
			}
		}
	}
	whereClause, vals, nextIdx := CompilePredicates(d.client.dialect, params.Where)
	isCursorQuery := (params.Cursor.Data.Column != "" || len(params.Cursor.Data.Children) > 0)
	if isCursorQuery {
		cClause, cVals, err := compileCursorClause(d.client.dialect, params.Cursor, params.OrderBy, allFieldsSoFarPKCols, allFieldsSoFarUniqueCols, "AllFieldsSoFar", nextIdx, params.Take)
		if err != nil {
			return nil, err
		}
		if cClause != "" {
			if whereClause == "" {
				whereClause = cClause
			} else {
				whereClause = "(" + whereClause + ") AND " + cClause
			}
			vals = append(vals, cVals...)
		}
	}
	if whereClause != "" {
		whereClause = " WHERE " + whereClause
	}
	orderByClause := formatOrderBySQL(d.client.dialect, params.OrderBy, allFieldsSoFarPKCols, allFieldsSoFarUniqueCols, isCursorQuery, params.Take)
	returningCols := selectAllFieldsSoFarCols(selects, omits)

	var results []*AllFieldsSoFar
	var err error
	if selects.hasAnyRelation() {
		err = d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			results, err = txQ.AllFieldsSoFar.queryMany(ctx, whereClause, orderByClause, vals, returningCols, params.Take, params.Skip)
			if err != nil {
				return err
			}
			if len(results) == 0 {
				return nil
			}
			return txQ.AllFieldsSoFar.loadRelations(ctx, results, selects)
		})
	} else {
		results, err = d.queryMany(ctx, whereClause, orderByClause, vals, returningCols, params.Take, params.Skip)
	}
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (d *AllFieldsSoFarDelegate) queryOne(ctx context.Context, whereClause string, orderByClause string, whereVals []any, returningCols []string, skip *int) (*AllFieldsSoFar, error) {
	limitOne := 1
	query := buildSelectSQL(d.client, "AllFieldsSoFar", returningCols, whereClause, orderByClause, &limitOne, skip)
	rows, err := d.client.query(ctx, query, whereVals...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			return nil, err
		}
		return nil, nil
	}

	var res AllFieldsSoFar
	if err := rows.Scan(res.ScanFields(returningCols)...); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &res, nil
}

func (d *AllFieldsSoFarDelegate) queryMany(ctx context.Context, whereClause string, orderByClause string, whereVals []any, returningCols []string, take *int, skip *int) ([]*AllFieldsSoFar, error) {
	query := buildSelectSQL(d.client, "AllFieldsSoFar", returningCols, whereClause, orderByClause, take, skip)
	rows, err := d.client.query(ctx, query, whereVals...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]*AllFieldsSoFar, 0)
	for rows.Next() {
		var res AllFieldsSoFar
		if err := rows.Scan(res.ScanFields(returningCols)...); err != nil {
			return nil, err
		}
		results = append(results, &res)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if take != nil && *take < 0 {
		reverseSlice(results)
	}
	return results, nil
}
func (d *AllFieldsSoFarDelegate) DeleteMany(preds ...PredicateOf[AllFieldsSoFar]) *DeleteManyBuilder[AllFieldsSoFar] {
	return &DeleteManyBuilder[AllFieldsSoFar]{
		where:    preds,
		execFunc: d.executeDeleteMany,
	}
}

func (d *AllFieldsSoFarDelegate) executeDeleteMany(ctx context.Context, preds []PredicateOf[AllFieldsSoFar]) (int64, error) {
	if len(d.extensions) == 0 {
		return d.runDeleteMany(ctx, preds)
	}

	curr := func(c context.Context, p []PredicateOf[AllFieldsSoFar]) (int64, error) {
		return d.runDeleteMany(c, p)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.DeleteMany != nil {
			next, hook := curr, ext.DeleteMany
			curr = func(c context.Context, p []PredicateOf[AllFieldsSoFar]) (int64, error) {
				return hook(c, p, next)
			}
		}
	}

	return curr(ctx, preds)
}

func (d *AllFieldsSoFarDelegate) runDeleteMany(ctx context.Context, preds []PredicateOf[AllFieldsSoFar]) (int64, error) {
	for _, pr := range preds {
		if pr != nil {
			if err := pr.Validate(); err != nil {
				return 0, err
			}
		}
	}

	whereClause, vals, _ := CompilePredicates(d.client.dialect, preds)

	var sb strings.Builder
	sb.WriteString("DELETE FROM ")
	d.client.dialect.WriteQuotedIdent(&sb, "AllFieldsSoFar")
	if whereClause != "" {
		sb.WriteString(" WHERE ")
		sb.WriteString(whereClause)
	}

	result, err := d.client.exec(ctx, sb.String(), vals...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (d *AllFieldsSoFarDelegate) Delete(where UniquePredicate[AllFieldsSoFar]) *DeleteBuilder[AllFieldsSoFar, AllFieldsSoFarSelect, AllFieldsSoFarOmit] {
	return &DeleteBuilder[AllFieldsSoFar, AllFieldsSoFarSelect, AllFieldsSoFarOmit]{
		where:    where,
		execFunc: d.executeDelete,
	}
}

func (d *AllFieldsSoFarDelegate) executeDelete(ctx context.Context, where UniquePredicate[AllFieldsSoFar], selects *AllFieldsSoFarSelect, omits *AllFieldsSoFarOmit) (*AllFieldsSoFar, error) {
	if len(d.extensions) == 0 {
		return d.runDelete(ctx, where, selects, omits)
	}

	curr := func(c context.Context, w UniquePredicate[AllFieldsSoFar], s *AllFieldsSoFarSelect, o *AllFieldsSoFarOmit) (*AllFieldsSoFar, error) {
		return d.runDelete(c, w, s, o)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.Delete != nil {
			next, hook := curr, ext.Delete
			curr = func(c context.Context, w UniquePredicate[AllFieldsSoFar], s *AllFieldsSoFarSelect, o *AllFieldsSoFarOmit) (*AllFieldsSoFar, error) {
				return hook(c, w, s, o, next)
			}
		}
	}

	return curr(ctx, where, selects, omits)
}

func (d *AllFieldsSoFarDelegate) runDelete(ctx context.Context, where UniquePredicate[AllFieldsSoFar], selects *AllFieldsSoFarSelect, omits *AllFieldsSoFarOmit) (*AllFieldsSoFar, error) {
	if err := where.Validate(); err != nil {
		return nil, err
	}

	returningCols := selectAllFieldsSoFarCols(selects, omits, allFieldsSoFarPKCols...)

	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := !d.client.dialect.SupportsDeleteReturning || hasRelations

	if useTx {
		var res *AllFieldsSoFar
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = txQ.AllFieldsSoFar.executeFindUnique(ctx, where, nil, selects, omits)
			if err != nil {
				return err
			}
			if res == nil {
				return sql.ErrNoRows
			}

			// Build DELETE statement by PK
			var deleteSb strings.Builder
			deleteSb.WriteString("DELETE FROM ")
			txQ.dialect.WriteQuotedIdent(&deleteSb, "AllFieldsSoFar")
			deleteSb.WriteString(" WHERE ")

			var pkPreds []PredicateOf[AllFieldsSoFar]
			pkPreds = append(pkPreds, Predicate[AllFieldsSoFar]{
				Data: PredicateData{
					Column:   "id",
					Operator: "=",
					Value:    res.Id,
				},
			})

			whereClause, vals, _ := CompilePredicates(txQ.dialect, pkPreds)
			deleteSb.WriteString(whereClause)

			_, err = txQ.exec(ctx, deleteSb.String(), vals...)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
		return res, nil
	}

	// Dialect supports RETURNING, and no relations need loading: run direct DELETE ... RETURNING
	var sb strings.Builder
	sb.WriteString("DELETE FROM ")
	d.client.dialect.WriteQuotedIdent(&sb, "AllFieldsSoFar")

	whereClause, vals, _ := CompilePredicates(d.client.dialect, []PredicateOf[AllFieldsSoFar]{where})
	if whereClause != "" {
		sb.WriteString(" WHERE ")
		sb.WriteString(whereClause)
	}

	sb.WriteString(" RETURNING ")
	for i, col := range returningCols {
		if i > 0 {
			sb.WriteString(", ")
		}
		d.client.dialect.WriteQuotedIdent(&sb, col)
	}

	rows, err := d.client.query(ctx, sb.String(), vals...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, err
		}
		return nil, sql.ErrNoRows
	}

	var row AllFieldsSoFar
	if err := rows.Scan(row.ScanFields(returningCols)...); err != nil {
		return nil, err
	}
	return &row, nil
}
func (d *AllFieldsSoFarDelegate) Count(preds ...PredicateOf[AllFieldsSoFar]) *CountBuilder[AllFieldsSoFar] {
	return &CountBuilder[AllFieldsSoFar]{
		where:    preds,
		execFunc: d.executeCount,
	}
}

func (d *AllFieldsSoFarDelegate) executeCount(ctx context.Context, params QueryParams[AllFieldsSoFar]) (int64, error) {
	if len(d.extensions) == 0 {
		return d.runCount(ctx, params)
	}

	curr := func(c context.Context, p QueryParams[AllFieldsSoFar]) (int64, error) {
		return d.runCount(c, p)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.Count != nil {
			next, hook := curr, ext.Count
			curr = func(c context.Context, p QueryParams[AllFieldsSoFar]) (int64, error) {
				return hook(c, p, next)
			}
		}
	}

	return curr(ctx, params)
}

func (d *AllFieldsSoFarDelegate) runCount(ctx context.Context, params QueryParams[AllFieldsSoFar]) (int64, error) {
	for _, pr := range params.Where {
		if pr != nil {
			if err := pr.Validate(); err != nil {
				return 0, err
			}
		}
	}

	whereClause, vals, _ := CompilePredicates(d.client.dialect, params.Where)
	if whereClause != "" {
		whereClause = " WHERE " + whereClause
	}

	var query string
	if params.Take != nil || params.Skip != nil {
		var subQuery strings.Builder
		subQuery.WriteString("SELECT 1 FROM ")
		d.client.dialect.WriteQuotedIdent(&subQuery, "AllFieldsSoFar")
		if whereClause != "" {
			subQuery.WriteString(whereClause)
		}
		subQuery.WriteString(d.client.dialect.FormatLimitOffset(params.Take, params.Skip))
		query = "SELECT COUNT(*) FROM (" + subQuery.String() + ") as sub"
	} else {
		var sb strings.Builder
		sb.WriteString("SELECT COUNT(*) FROM ")
		d.client.dialect.WriteQuotedIdent(&sb, "AllFieldsSoFar")
		if whereClause != "" {
			sb.WriteString(whereClause)
		}
		query = sb.String()
	}

	rows, err := d.client.query(ctx, query, vals...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var count int64
	if rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return 0, err
		}
	}
	if err := rows.Err(); err != nil {
		return 0, err
	}
	return count, nil
}
func (d *AllFieldsSoFarDelegate) loadRelations(ctx context.Context, records []*AllFieldsSoFar, selects *AllFieldsSoFarSelect) error {
	_ = ctx
	if selects == nil || len(records) == 0 {
		return nil
	}

	return nil
}
