package valk

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
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

// AllFieldsSoFarSelect specifies which fields to include
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

// AllFieldsSoFarOmit specifies which fields to exclude
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

type AllFieldsSoFarDelegate struct {
	client          *Queries
	beforeCreate    func(context.Context, *AllFieldsSoFarCreate) error
	afterCreate     func(context.Context, []*AllFieldsSoFar) error
	afterCreateMany func(context.Context, []AllFieldsSoFarCreate, int64) error
}

func (d *AllFieldsSoFarDelegate) BeforeCreate(hook func(context.Context, *AllFieldsSoFarCreate) error) {
	d.beforeCreate = hook
}

func (d *AllFieldsSoFarDelegate) AfterCreate(hook func(context.Context, []*AllFieldsSoFar) error) {
	d.afterCreate = hook
}

func (d *AllFieldsSoFarDelegate) AfterCreateMany(hook func(context.Context, []AllFieldsSoFarCreate, int64) error) {
	d.afterCreateMany = hook
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

func (q *Queries) selectAllFieldsSoFarCols(selects *AllFieldsSoFarSelect, omits *AllFieldsSoFarOmit, forceCols ...string) []string {
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

var AllFieldsSoFarColOrder = []string{
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

func (s *AllFieldsSoFarSelect) hasAnyRelation() bool {
	if s == nil {
		return false
	}
	return false
}

func (d *AllFieldsSoFarDelegate) Create(assignments ...FieldAssignment) *CreateBuilder[AllFieldsSoFar, AllFieldsSoFarSelect, AllFieldsSoFarOmit] {
	return &CreateBuilder[AllFieldsSoFar, AllFieldsSoFarSelect, AllFieldsSoFarOmit]{
		client:      d.client,
		assignments: assignments,
		execFunc:    d.client.executeAllFieldsSoFarCreate,
	}
}

func validateAllFieldsSoFarCreate(assignments []FieldAssignment) error {
	errs := &ValidationError{}

	provided := make(map[string]bool)
	for _, a := range assignments {
		provided[a.Col] = true
		switch a.Col {
		case "id":
			if v, ok := a.Val.(int32); ok {
				ValidateInt32(errs, "id", v, "")
			} else {
				errs.Add("id", a.Val, "type", "field id must be of type int32")
			}
		case "stringReq":
			if v, ok := a.Val.(string); ok {
				ValidateString(errs, "stringReq", v, true, 0, false, false)
			} else {
				errs.Add("stringReq", a.Val, "type", "field stringReq must be of type string")
			}
		case "stringOpt":
			if v, ok := a.Val.(string); ok {
				ValidateString(errs, "stringOpt", v, false, 0, false, false)
			} else {
				errs.Add("stringOpt", a.Val, "type", "field stringOpt must be of type string")
			}
		case "stringDefault":
			if v, ok := a.Val.(string); ok {
				ValidateString(errs, "stringDefault", v, false, 0, false, false)
			} else {
				errs.Add("stringDefault", a.Val, "type", "field stringDefault must be of type string")
			}
		case "stringVarchar":
			if v, ok := a.Val.(string); ok {
				ValidateString(errs, "stringVarchar", v, true, 255, false, false)
			} else {
				errs.Add("stringVarchar", a.Val, "type", "field stringVarchar must be of type string")
			}
		case "stringChar":
			if v, ok := a.Val.(string); ok {
				ValidateString(errs, "stringChar", v, true, 10, false, false)
			} else {
				errs.Add("stringChar", a.Val, "type", "field stringChar must be of type string")
			}
		case "bitVal":
			if v, ok := a.Val.(string); ok {
				ValidateString(errs, "bitVal", v, true, 0, true, false)
			} else {
				errs.Add("bitVal", a.Val, "type", "field bitVal must be of type string")
			}
		case "varBitVal":
			if v, ok := a.Val.(string); ok {
				ValidateString(errs, "varBitVal", v, true, 0, true, false)
			} else {
				errs.Add("varBitVal", a.Val, "type", "field varBitVal must be of type string")
			}
		case "inetVal":
			if v, ok := a.Val.(string); ok {
				ValidateString(errs, "inetVal", v, true, 0, false, true)
			} else {
				errs.Add("inetVal", a.Val, "type", "field inetVal must be of type string")
			}
		case "xmlVal":
			if v, ok := a.Val.(string); ok {
				ValidateString(errs, "xmlVal", v, true, 0, false, false)
			} else {
				errs.Add("xmlVal", a.Val, "type", "field xmlVal must be of type string")
			}
		case "cuidDefault":
			if v, ok := a.Val.(string); ok {
				ValidateString(errs, "cuidDefault", v, false, 0, false, false)
			} else {
				errs.Add("cuidDefault", a.Val, "type", "field cuidDefault must be of type string")
			}
		case "cuid1Default":
			if v, ok := a.Val.(string); ok {
				ValidateString(errs, "cuid1Default", v, false, 0, false, false)
			} else {
				errs.Add("cuid1Default", a.Val, "type", "field cuid1Default must be of type string")
			}
		case "cuid2Default":
			if v, ok := a.Val.(string); ok {
				ValidateString(errs, "cuid2Default", v, false, 0, false, false)
			} else {
				errs.Add("cuid2Default", a.Val, "type", "field cuid2Default must be of type string")
			}
		case "uuidDefault":
			if v, ok := a.Val.(string); ok {
				ValidateString(errs, "uuidDefault", v, false, 0, false, false)
			} else {
				errs.Add("uuidDefault", a.Val, "type", "field uuidDefault must be of type string")
			}
		case "uuid4Default":
			if v, ok := a.Val.(string); ok {
				ValidateString(errs, "uuid4Default", v, false, 0, false, false)
			} else {
				errs.Add("uuid4Default", a.Val, "type", "field uuid4Default must be of type string")
			}
		case "uuid7Default":
			if v, ok := a.Val.(string); ok {
				ValidateString(errs, "uuid7Default", v, false, 0, false, false)
			} else {
				errs.Add("uuid7Default", a.Val, "type", "field uuid7Default must be of type string")
			}
		case "ulidDefault":
			if v, ok := a.Val.(string); ok {
				ValidateString(errs, "ulidDefault", v, false, 0, false, false)
			} else {
				errs.Add("ulidDefault", a.Val, "type", "field ulidDefault must be of type string")
			}
		case "nanoidDefault":
			if v, ok := a.Val.(string); ok {
				ValidateString(errs, "nanoidDefault", v, false, 0, false, false)
			} else {
				errs.Add("nanoidDefault", a.Val, "type", "field nanoidDefault must be of type string")
			}
		case "uuidDb":
			if v, ok := a.Val.(string); ok {
				ValidateString(errs, "uuidDb", v, true, 0, false, false)
			} else {
				errs.Add("uuidDb", a.Val, "type", "field uuidDb must be of type string")
			}
		case "intReq":
			if v, ok := a.Val.(int32); ok {
				ValidateInt32(errs, "intReq", v, "")
			} else {
				errs.Add("intReq", a.Val, "type", "field intReq must be of type int32")
			}
		case "intOpt":
			if v, ok := a.Val.(int32); ok {
				ValidateInt32(errs, "intOpt", v, "")
			} else {
				errs.Add("intOpt", a.Val, "type", "field intOpt must be of type int32")
			}
		case "intDefault":
			if v, ok := a.Val.(int32); ok {
				ValidateInt32(errs, "intDefault", v, "")
			} else {
				errs.Add("intDefault", a.Val, "type", "field intDefault must be of type int32")
			}
		case "integerVal":
			if v, ok := a.Val.(int32); ok {
				ValidateInt32(errs, "integerVal", v, "")
			} else {
				errs.Add("integerVal", a.Val, "type", "field integerVal must be of type int32")
			}
		case "smallInt":
			if v, ok := a.Val.(int32); ok {
				ValidateInt32(errs, "smallInt", v, "SmallInt")
			} else {
				errs.Add("smallInt", a.Val, "type", "field smallInt must be of type int32")
			}
		case "tinyInt":
			if v, ok := a.Val.(int32); ok {
				ValidateInt32(errs, "tinyInt", v, "")
			} else {
				errs.Add("tinyInt", a.Val, "type", "field tinyInt must be of type int32")
			}
		case "oidVal":
			if v, ok := a.Val.(int32); ok {
				ValidateInt32(errs, "oidVal", v, "Oid")
			} else {
				errs.Add("oidVal", a.Val, "type", "field oidVal must be of type int32")
			}
		case "bigIntReq":
			if v, ok := a.Val.(int64); ok {
				ValidateInt64(errs, "bigIntReq", v, "")
			} else {
				errs.Add("bigIntReq", a.Val, "type", "field bigIntReq must be of type int64")
			}
		case "bigIntOpt":
			if v, ok := a.Val.(int64); ok {
				ValidateInt64(errs, "bigIntOpt", v, "")
			} else {
				errs.Add("bigIntOpt", a.Val, "type", "field bigIntOpt must be of type int64")
			}
		case "floatReq":
			if _, ok := a.Val.(float64); !ok {
				errs.Add("floatReq", a.Val, "type", "field floatReq must be of type float64")
			}
		case "floatOpt":
			if _, ok := a.Val.(float64); !ok {
				errs.Add("floatOpt", a.Val, "type", "field floatOpt must be of type float64")
			}
		case "realVal":
			if _, ok := a.Val.(float64); !ok {
				errs.Add("realVal", a.Val, "type", "field realVal must be of type float64")
			}
		case "decimalReq":
			if v, ok := a.Val.(string); ok {
				ValidateString(errs, "decimalReq", v, true, 0, false, false)
			} else {
				errs.Add("decimalReq", a.Val, "type", "field decimalReq must be of type string")
			}
		case "decimalOpt":
			if v, ok := a.Val.(string); ok {
				ValidateString(errs, "decimalOpt", v, false, 0, false, false)
			} else {
				errs.Add("decimalOpt", a.Val, "type", "field decimalOpt must be of type string")
			}
		case "decimalPrecise":
			if v, ok := a.Val.(string); ok {
				ValidateString(errs, "decimalPrecise", v, true, 0, false, false)
			} else {
				errs.Add("decimalPrecise", a.Val, "type", "field decimalPrecise must be of type string")
			}
		case "moneyVal":
			if v, ok := a.Val.(string); ok {
				ValidateString(errs, "moneyVal", v, true, 0, false, false)
			} else {
				errs.Add("moneyVal", a.Val, "type", "field moneyVal must be of type string")
			}
		case "boolReq":
			if _, ok := a.Val.(bool); !ok {
				errs.Add("boolReq", a.Val, "type", "field boolReq must be of type bool")
			}
		case "boolOpt":
			if _, ok := a.Val.(bool); !ok {
				errs.Add("boolOpt", a.Val, "type", "field boolOpt must be of type bool")
			}
		case "boolDefault":
			if _, ok := a.Val.(bool); !ok {
				errs.Add("boolDefault", a.Val, "type", "field boolDefault must be of type bool")
			}
		case "dateTimeReq":
			if _, ok := a.Val.(time.Time); !ok {
				errs.Add("dateTimeReq", a.Val, "type", "field dateTimeReq must be of type time.Time")
			}
		case "dateTimeOpt":
			if _, ok := a.Val.(time.Time); !ok {
				errs.Add("dateTimeOpt", a.Val, "type", "field dateTimeOpt must be of type time.Time")
			}
		case "dateTimeDefault":
			if _, ok := a.Val.(time.Time); !ok {
				errs.Add("dateTimeDefault", a.Val, "type", "field dateTimeDefault must be of type time.Time")
			}
		case "updatedAt":
			if _, ok := a.Val.(time.Time); !ok {
				errs.Add("updatedAt", a.Val, "type", "field updatedAt must be of type time.Time")
			}
		case "dateTimeTz":
			if _, ok := a.Val.(time.Time); !ok {
				errs.Add("dateTimeTz", a.Val, "type", "field dateTimeTz must be of type time.Time")
			}
		case "timestampVal":
			if _, ok := a.Val.(time.Time); !ok {
				errs.Add("timestampVal", a.Val, "type", "field timestampVal must be of type time.Time")
			}
		case "timeVal":
			if _, ok := a.Val.(time.Time); !ok {
				errs.Add("timeVal", a.Val, "type", "field timeVal must be of type time.Time")
			}
		case "timetzVal":
			if _, ok := a.Val.(time.Time); !ok {
				errs.Add("timetzVal", a.Val, "type", "field timetzVal must be of type time.Time")
			}
		case "jsonReq":
			if _, ok := a.Val.(json.RawMessage); !ok {
				errs.Add("jsonReq", a.Val, "type", "field jsonReq must be of type json.RawMessage")
			}
		case "jsonOpt":
			if _, ok := a.Val.(json.RawMessage); !ok {
				errs.Add("jsonOpt", a.Val, "type", "field jsonOpt must be of type json.RawMessage")
			}
		case "jsonVal":
			if _, ok := a.Val.(json.RawMessage); !ok {
				errs.Add("jsonVal", a.Val, "type", "field jsonVal must be of type json.RawMessage")
			}
		case "bytesReq":
			if _, ok := a.Val.([]byte); !ok {
				errs.Add("bytesReq", a.Val, "type", "field bytesReq must be of type []byte")
			}
		case "bytesOpt":
			if _, ok := a.Val.([]byte); !ok {
				errs.Add("bytesOpt", a.Val, "type", "field bytesOpt must be of type []byte")
			}
		case "hstoreField":
			if _, ok := a.Val.(map[string]*string); !ok {
				errs.Add("hstoreField", a.Val, "type", "field hstoreField must be of type map[string]*string")
			}
		case "ltreeField":
			if v, ok := a.Val.(string); ok {
				ValidateString(errs, "ltreeField", v, true, 0, false, false)
			} else {
				errs.Add("ltreeField", a.Val, "type", "field ltreeField must be of type string")
			}
		case "citextField":
			if v, ok := a.Val.(string); ok {
				ValidateString(errs, "citextField", v, false, 0, false, false)
			} else {
				errs.Add("citextField", a.Val, "type", "field citextField must be of type string")
			}
		}
	}
	if !provided["stringReq"] {
		errs.Add("stringReq", "", "required", "field StringReq is required")
	}
	if !provided["stringVarchar"] {
		errs.Add("stringVarchar", "", "required", "field StringVarchar is required")
	}
	if !provided["stringChar"] {
		errs.Add("stringChar", "", "required", "field StringChar is required")
	}
	if !provided["bitVal"] {
		errs.Add("bitVal", "", "required", "field BitVal is required")
	}
	if !provided["varBitVal"] {
		errs.Add("varBitVal", "", "required", "field VarBitVal is required")
	}
	if !provided["inetVal"] {
		errs.Add("inetVal", "", "required", "field InetVal is required")
	}
	if !provided["xmlVal"] {
		errs.Add("xmlVal", "", "required", "field XmlVal is required")
	}
	if !provided["uuidDb"] {
		errs.Add("uuidDb", "", "required", "field UuidDb is required")
	}
	if !provided["intReq"] {
		errs.Add("intReq", nil, "required", "field IntReq is required")
	}
	if !provided["integerVal"] {
		errs.Add("integerVal", nil, "required", "field IntegerVal is required")
	}
	if !provided["smallInt"] {
		errs.Add("smallInt", nil, "required", "field SmallInt is required")
	}
	if !provided["tinyInt"] {
		errs.Add("tinyInt", nil, "required", "field TinyInt is required")
	}
	if !provided["oidVal"] {
		errs.Add("oidVal", nil, "required", "field OidVal is required")
	}
	if !provided["bigIntReq"] {
		errs.Add("bigIntReq", nil, "required", "field BigIntReq is required")
	}
	if !provided["floatReq"] {
		errs.Add("floatReq", nil, "required", "field FloatReq is required")
	}
	if !provided["realVal"] {
		errs.Add("realVal", nil, "required", "field RealVal is required")
	}
	if !provided["decimalReq"] {
		errs.Add("decimalReq", "", "required", "field DecimalReq is required")
	}
	if !provided["decimalPrecise"] {
		errs.Add("decimalPrecise", "", "required", "field DecimalPrecise is required")
	}
	if !provided["moneyVal"] {
		errs.Add("moneyVal", "", "required", "field MoneyVal is required")
	}
	if !provided["boolReq"] {
		errs.Add("boolReq", nil, "required", "field BoolReq is required")
	}
	if !provided["dateTimeReq"] {
		errs.Add("dateTimeReq", nil, "required", "field DateTimeReq is required")
	}
	if !provided["updatedAt"] {
		errs.Add("updatedAt", nil, "required", "field UpdatedAt is required")
	}
	if !provided["dateTimeTz"] {
		errs.Add("dateTimeTz", nil, "required", "field DateTimeTz is required")
	}
	if !provided["timestampVal"] {
		errs.Add("timestampVal", nil, "required", "field TimestampVal is required")
	}
	if !provided["timeVal"] {
		errs.Add("timeVal", nil, "required", "field TimeVal is required")
	}
	if !provided["timetzVal"] {
		errs.Add("timetzVal", nil, "required", "field TimetzVal is required")
	}
	if !provided["jsonReq"] {
		errs.Add("jsonReq", nil, "required", "field JsonReq is required")
	}
	if !provided["jsonVal"] {
		errs.Add("jsonVal", nil, "required", "field JsonVal is required")
	}
	if !provided["bytesReq"] {
		errs.Add("bytesReq", nil, "required", "field BytesReq is required")
	}
	if !provided["ltreeField"] {
		errs.Add("ltreeField", "", "required", "field LtreeField is required")
	}

	if errs.HasErrors() {
		return *errs
	}
	return nil
}

func assignmentsToAllFieldsSoFarCreate(assignments []FieldAssignment) AllFieldsSoFarCreate {
	var input AllFieldsSoFarCreate
	for _, a := range assignments {
		switch a.Col {
		case "id":
			if v, ok := a.Val.(int32); ok {
				input.Id = &v
			}
		case "stringReq":
			if v, ok := a.Val.(string); ok {
				input.StringReq = v
			}
		case "stringOpt":
			if v, ok := a.Val.(string); ok {
				input.StringOpt = &v
			}
		case "stringDefault":
			if v, ok := a.Val.(string); ok {
				input.StringDefault = &v
			}
		case "stringVarchar":
			if v, ok := a.Val.(string); ok {
				input.StringVarchar = v
			}
		case "stringChar":
			if v, ok := a.Val.(string); ok {
				input.StringChar = v
			}
		case "bitVal":
			if v, ok := a.Val.(string); ok {
				input.BitVal = v
			}
		case "varBitVal":
			if v, ok := a.Val.(string); ok {
				input.VarBitVal = v
			}
		case "inetVal":
			if v, ok := a.Val.(string); ok {
				input.InetVal = v
			}
		case "xmlVal":
			if v, ok := a.Val.(string); ok {
				input.XmlVal = v
			}
		case "cuidDefault":
			if v, ok := a.Val.(string); ok {
				input.CuidDefault = &v
			}
		case "cuid1Default":
			if v, ok := a.Val.(string); ok {
				input.Cuid1Default = &v
			}
		case "cuid2Default":
			if v, ok := a.Val.(string); ok {
				input.Cuid2Default = &v
			}
		case "uuidDefault":
			if v, ok := a.Val.(string); ok {
				input.UuidDefault = &v
			}
		case "uuid4Default":
			if v, ok := a.Val.(string); ok {
				input.Uuid4Default = &v
			}
		case "uuid7Default":
			if v, ok := a.Val.(string); ok {
				input.Uuid7Default = &v
			}
		case "ulidDefault":
			if v, ok := a.Val.(string); ok {
				input.UlidDefault = &v
			}
		case "nanoidDefault":
			if v, ok := a.Val.(string); ok {
				input.NanoidDefault = &v
			}
		case "uuidDb":
			if v, ok := a.Val.(string); ok {
				input.UuidDb = v
			}
		case "intReq":
			if v, ok := a.Val.(int32); ok {
				input.IntReq = v
			}
		case "intOpt":
			if v, ok := a.Val.(int32); ok {
				input.IntOpt = &v
			}
		case "intDefault":
			if v, ok := a.Val.(int32); ok {
				input.IntDefault = &v
			}
		case "integerVal":
			if v, ok := a.Val.(int32); ok {
				input.IntegerVal = v
			}
		case "smallInt":
			if v, ok := a.Val.(int32); ok {
				input.SmallInt = v
			}
		case "tinyInt":
			if v, ok := a.Val.(int32); ok {
				input.TinyInt = v
			}
		case "oidVal":
			if v, ok := a.Val.(int32); ok {
				input.OidVal = v
			}
		case "bigIntReq":
			if v, ok := a.Val.(int64); ok {
				input.BigIntReq = v
			}
		case "bigIntOpt":
			if v, ok := a.Val.(int64); ok {
				input.BigIntOpt = &v
			}
		case "floatReq":
			if v, ok := a.Val.(float64); ok {
				input.FloatReq = v
			}
		case "floatOpt":
			if v, ok := a.Val.(float64); ok {
				input.FloatOpt = &v
			}
		case "realVal":
			if v, ok := a.Val.(float64); ok {
				input.RealVal = v
			}
		case "decimalReq":
			if v, ok := a.Val.(string); ok {
				input.DecimalReq = v
			}
		case "decimalOpt":
			if v, ok := a.Val.(string); ok {
				input.DecimalOpt = &v
			}
		case "decimalPrecise":
			if v, ok := a.Val.(string); ok {
				input.DecimalPrecise = v
			}
		case "moneyVal":
			if v, ok := a.Val.(string); ok {
				input.MoneyVal = v
			}
		case "boolReq":
			if v, ok := a.Val.(bool); ok {
				input.BoolReq = v
			}
		case "boolOpt":
			if v, ok := a.Val.(bool); ok {
				input.BoolOpt = &v
			}
		case "boolDefault":
			if v, ok := a.Val.(bool); ok {
				input.BoolDefault = &v
			}
		case "dateTimeReq":
			if v, ok := a.Val.(time.Time); ok {
				input.DateTimeReq = v
			}
		case "dateTimeOpt":
			if v, ok := a.Val.(time.Time); ok {
				input.DateTimeOpt = &v
			}
		case "dateTimeDefault":
			if v, ok := a.Val.(time.Time); ok {
				input.DateTimeDefault = &v
			}
		case "updatedAt":
			if v, ok := a.Val.(time.Time); ok {
				input.UpdatedAt = v
			}
		case "dateTimeTz":
			if v, ok := a.Val.(time.Time); ok {
				input.DateTimeTz = v
			}
		case "timestampVal":
			if v, ok := a.Val.(time.Time); ok {
				input.TimestampVal = v
			}
		case "timeVal":
			if v, ok := a.Val.(time.Time); ok {
				input.TimeVal = v
			}
		case "timetzVal":
			if v, ok := a.Val.(time.Time); ok {
				input.TimetzVal = v
			}
		case "jsonReq":
			if v, ok := a.Val.(json.RawMessage); ok {
				input.JsonReq = v
			}
		case "jsonOpt":
			if v, ok := a.Val.(json.RawMessage); ok {
				input.JsonOpt = &v
			}
		case "jsonVal":
			if v, ok := a.Val.(json.RawMessage); ok {
				input.JsonVal = v
			}
		case "bytesReq":
			if v, ok := a.Val.([]byte); ok {
				input.BytesReq = v
			}
		case "bytesOpt":
			if v, ok := a.Val.([]byte); ok {
				input.BytesOpt = &v
			}
		case "hstoreField":
			if v, ok := a.Val.(map[string]*string); ok {
				input.HstoreField = &v
			}
		case "ltreeField":
			if v, ok := a.Val.(string); ok {
				input.LtreeField = v
			}
		case "citextField":
			if v, ok := a.Val.(string); ok {
				input.CitextField = &v
			}
		}
	}
	return input
}

func (s *AllFieldsSoFarCreate) ToRowMap() map[string]any {
	m := make(map[string]any, 54)
	if s.Id != nil {
		m["id"] = *s.Id
	}
	m["stringReq"] = s.StringReq
	if s.StringOpt != nil {
		m["stringOpt"] = *s.StringOpt
	}
	if s.StringDefault != nil {
		m["stringDefault"] = *s.StringDefault
	}
	m["stringVarchar"] = s.StringVarchar
	m["stringChar"] = s.StringChar
	m["bitVal"] = s.BitVal
	m["varBitVal"] = s.VarBitVal
	m["inetVal"] = s.InetVal
	m["xmlVal"] = s.XmlVal
	if s.CuidDefault != nil {
		m["cuidDefault"] = *s.CuidDefault
	} else {
		m["cuidDefault"] = generateCUID()
	}
	if s.Cuid1Default != nil {
		m["cuid1Default"] = *s.Cuid1Default
	} else {
		m["cuid1Default"] = generateCUID()
	}
	if s.Cuid2Default != nil {
		m["cuid2Default"] = *s.Cuid2Default
	} else {
		m["cuid2Default"] = generateCUID2()
	}
	if s.UuidDefault != nil {
		m["uuidDefault"] = *s.UuidDefault
	} else {
		m["uuidDefault"] = generateUUID()
	}
	if s.Uuid4Default != nil {
		m["uuid4Default"] = *s.Uuid4Default
	} else {
		m["uuid4Default"] = generateUUID()
	}
	if s.Uuid7Default != nil {
		m["uuid7Default"] = *s.Uuid7Default
	} else {
		m["uuid7Default"] = generateUUID7()
	}
	if s.UlidDefault != nil {
		m["ulidDefault"] = *s.UlidDefault
	} else {
		m["ulidDefault"] = generateULID()
	}
	if s.NanoidDefault != nil {
		m["nanoidDefault"] = *s.NanoidDefault
	} else {
		m["nanoidDefault"] = generateNanoID()
	}
	m["uuidDb"] = s.UuidDb
	m["intReq"] = s.IntReq
	if s.IntOpt != nil {
		m["intOpt"] = *s.IntOpt
	}
	if s.IntDefault != nil {
		m["intDefault"] = *s.IntDefault
	}
	m["integerVal"] = s.IntegerVal
	m["smallInt"] = s.SmallInt
	m["tinyInt"] = s.TinyInt
	m["oidVal"] = s.OidVal
	m["bigIntReq"] = s.BigIntReq
	if s.BigIntOpt != nil {
		m["bigIntOpt"] = *s.BigIntOpt
	}
	m["floatReq"] = s.FloatReq
	if s.FloatOpt != nil {
		m["floatOpt"] = *s.FloatOpt
	}
	m["realVal"] = s.RealVal
	m["decimalReq"] = s.DecimalReq
	if s.DecimalOpt != nil {
		m["decimalOpt"] = *s.DecimalOpt
	}
	m["decimalPrecise"] = s.DecimalPrecise
	m["moneyVal"] = s.MoneyVal
	m["boolReq"] = s.BoolReq
	if s.BoolOpt != nil {
		m["boolOpt"] = *s.BoolOpt
	}
	if s.BoolDefault != nil {
		m["boolDefault"] = *s.BoolDefault
	}
	m["dateTimeReq"] = s.DateTimeReq
	if s.DateTimeOpt != nil {
		m["dateTimeOpt"] = *s.DateTimeOpt
	}
	if s.DateTimeDefault != nil {
		m["dateTimeDefault"] = *s.DateTimeDefault
	} else {
		m["dateTimeDefault"] = time.Now()
	}
	m["updatedAt"] = s.UpdatedAt
	m["dateTimeTz"] = s.DateTimeTz
	m["timestampVal"] = s.TimestampVal
	m["timeVal"] = s.TimeVal
	m["timetzVal"] = s.TimetzVal
	m["jsonReq"] = s.JsonReq
	if s.JsonOpt != nil {
		m["jsonOpt"] = *s.JsonOpt
	}
	m["jsonVal"] = s.JsonVal
	m["bytesReq"] = s.BytesReq
	if s.BytesOpt != nil {
		m["bytesOpt"] = *s.BytesOpt
	}
	if s.HstoreField != nil {
		m["hstoreField"] = ToHstore(*s.HstoreField)
	}
	m["ltreeField"] = s.LtreeField
	if s.CitextField != nil {
		m["citextField"] = *s.CitextField
	}
	return m
}

func (q *Queries) executeAllFieldsSoFarCreate(ctx context.Context, assignments []FieldAssignment, selects *AllFieldsSoFarSelect, omits *AllFieldsSoFarOmit) (*AllFieldsSoFar, error) {
	input := assignmentsToAllFieldsSoFarCreate(assignments)

	if q.AllFieldsSoFar.beforeCreate != nil {
		if err := q.AllFieldsSoFar.beforeCreate(ctx, &input); err != nil {
			return nil, err
		}
	}

	if err := validateAllFieldsSoFarCreate(assignments); err != nil {
		return nil, err
	}

	rowMap := input.ToRowMap()
	var cols []string
	var vals []any
	for _, col := range AllFieldsSoFarColOrder {
		if val, ok := rowMap[col]; ok {
			cols = append(cols, col)
			vals = append(vals, val)
		}
	}

	returningCols := q.selectAllFieldsSoFarCols(selects, omits)

	scanFunc := func(res *AllFieldsSoFar, cols []string) []any {
		return res.ScanFields(cols)
	}

	idCol := "id"

	hasRelations := selects.hasAnyRelation()

	var res *AllFieldsSoFar
	var err error
	if hasRelations {
		err = q.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = executeInsert(ctx, txQ, "AllFieldsSoFar", cols, vals, returningCols, idCol, scanFunc)
			if err != nil {
				return err
			}
			return txQ.loadAllFieldsSoFarRelations(ctx, []*AllFieldsSoFar{res}, selects)
		})
	} else {
		res, err = executeInsert(ctx, q, "AllFieldsSoFar", cols, vals, returningCols, idCol, scanFunc)
	}
	if err != nil {
		return nil, err
	}

	if q.AllFieldsSoFar.afterCreate != nil {
		if err := q.AllFieldsSoFar.afterCreate(ctx, []*AllFieldsSoFar{res}); err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (d *AllFieldsSoFarDelegate) CreateMany(records ...RecordInput) *CreateManyBuilder[AllFieldsSoFar] {
	return &CreateManyBuilder[AllFieldsSoFar]{
		client:   d.client,
		records:  records,
		execFunc: d.client.executeAllFieldsSoFarCreateMany,
	}
}

func (d *AllFieldsSoFarDelegate) CreateManyAndReturn(records ...RecordInput) *CreateManyAndReturnBuilder[AllFieldsSoFar, AllFieldsSoFarSelect, AllFieldsSoFarOmit] {
	return &CreateManyAndReturnBuilder[AllFieldsSoFar, AllFieldsSoFarSelect, AllFieldsSoFarOmit]{
		client:   d.client,
		records:  records,
		execFunc: d.client.executeAllFieldsSoFarCreateManyAndReturn,
	}
}

func (q *Queries) executeAllFieldsSoFarCreateMany(ctx context.Context, records []RecordInput) (int64, error) {
	rowMaps := make([]map[string]any, len(records))
	inputs := make([]AllFieldsSoFarCreate, len(records))
	for i, rec := range records {
		if err := validateAllFieldsSoFarCreate(rec.Assignments); err != nil {
			return 0, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
		input := assignmentsToAllFieldsSoFarCreate(rec.Assignments)
		if q.AllFieldsSoFar.beforeCreate != nil {
			if err := q.AllFieldsSoFar.beforeCreate(ctx, &input); err != nil {
				return 0, err
			}
		}
		rowMaps[i] = input.ToRowMap()
		inputs[i] = input
	}
	count, err := executeCreateMany(ctx, q, rowMaps, "AllFieldsSoFar", AllFieldsSoFarColOrder)
	if err != nil {
		return 0, err
	}
	if q.AllFieldsSoFar.afterCreateMany != nil {
		if err := q.AllFieldsSoFar.afterCreateMany(ctx, inputs, count); err != nil {
			return 0, err
		}
	}
	return count, nil
}

func (q *Queries) executeAllFieldsSoFarCreateManyAndReturn(ctx context.Context, records []RecordInput, selects *AllFieldsSoFarSelect, omits *AllFieldsSoFarOmit) ([]*AllFieldsSoFar, error) {
	rowMaps := make([]map[string]any, len(records))
	idCol := "id"
	for i, rec := range records {
		if err := validateAllFieldsSoFarCreate(rec.Assignments); err != nil {
			return nil, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
		input := assignmentsToAllFieldsSoFarCreate(rec.Assignments)
		if q.AllFieldsSoFar.beforeCreate != nil {
			if err := q.AllFieldsSoFar.beforeCreate(ctx, &input); err != nil {
				return nil, err
			}
		}
		rowMaps[i] = input.ToRowMap()
	}
	results, err := executeCreateManyAndReturn(ctx, q, rowMaps, "AllFieldsSoFar", AllFieldsSoFarColOrder, selects, omits,
		q.selectAllFieldsSoFarCols,
		q.loadAllFieldsSoFarRelations,
		(*AllFieldsSoFar).ScanFields,
		(*AllFieldsSoFarSelect).hasAnyRelation,
		idCol,
	)
	if err != nil {
		return nil, err
	}
	if q.AllFieldsSoFar.afterCreate != nil {
		if err := q.AllFieldsSoFar.afterCreate(ctx, results); err != nil {
			return nil, err
		}
	}
	return results, nil
}
func (d *AllFieldsSoFarDelegate) FindUnique(where UniquePredicate) *FindUniqueBuilder[AllFieldsSoFar, AllFieldsSoFarSelect, AllFieldsSoFarOmit] {
	return &FindUniqueBuilder[AllFieldsSoFar, AllFieldsSoFarSelect, AllFieldsSoFarOmit]{
		client:   d.client,
		where:    where,
		execFunc: d.client.executeAllFieldsSoFarFindUnique,
	}
}

func (d *AllFieldsSoFarDelegate) FindFirst(preds ...Predicate) *FindFirstBuilder[AllFieldsSoFar, AllFieldsSoFarSelect, AllFieldsSoFarOmit] {
	return &FindFirstBuilder[AllFieldsSoFar, AllFieldsSoFarSelect, AllFieldsSoFarOmit]{
		client:   d.client,
		where:    preds,
		execFunc: d.client.executeAllFieldsSoFarFindFirst,
	}
}

func (d *AllFieldsSoFarDelegate) FindMany(preds ...Predicate) *FindManyBuilder[AllFieldsSoFar, AllFieldsSoFarSelect, AllFieldsSoFarOmit] {
	return &FindManyBuilder[AllFieldsSoFar, AllFieldsSoFarSelect, AllFieldsSoFarOmit]{
		client:   d.client,
		where:    preds,
		execFunc: d.client.executeAllFieldsSoFarFindMany,
	}
}

func (q *Queries) executeAllFieldsSoFarFindUnique(ctx context.Context, where UniquePredicate, selects *AllFieldsSoFarSelect, omits *AllFieldsSoFarOmit) (*AllFieldsSoFar, error) {
	if where == nil {
		return nil, fmt.Errorf("at least one unique field must be set for FindUnique")
	}
	if err := where.Validate(); err != nil {
		return nil, err
	}
	whereClause, vals := CompilePredicates(q.dialect, []Predicate{where})
	if whereClause != "" {
		whereClause = " WHERE " + whereClause
	}
	returningCols := q.selectAllFieldsSoFarCols(selects, omits)
	return executeSingleWithRelations(ctx, q, "AllFieldsSoFar", whereClause, vals, returningCols,
		func(res *AllFieldsSoFar, cols []string) []any { return res.ScanFields(cols) },
		selects.hasAnyRelation(),
		func(ctx context.Context, txQ *Queries, results []*AllFieldsSoFar) error {
			return txQ.loadAllFieldsSoFarRelations(ctx, results, selects)
		},
		nil,
	)
}

func (q *Queries) executeAllFieldsSoFarFindFirst(
	ctx context.Context,
	params QueryParams,
	selects *AllFieldsSoFarSelect,
	omits *AllFieldsSoFarOmit,
) (*AllFieldsSoFar, error) {
	for _, p := range params.Where {
		if p != nil {
			if err := p.Validate(); err != nil {
				return nil, err
			}
		}
	}
	whereClause, vals := CompilePredicates(q.dialect, params.Where)
	if whereClause != "" {
		whereClause = " WHERE " + whereClause
	}
	returningCols := q.selectAllFieldsSoFarCols(selects, omits)
	return executeSingleWithRelations(ctx, q, "AllFieldsSoFar", whereClause, vals, returningCols,
		func(res *AllFieldsSoFar, cols []string) []any { return res.ScanFields(cols) },
		selects.hasAnyRelation(),
		func(ctx context.Context, txQ *Queries, results []*AllFieldsSoFar) error {
			return txQ.loadAllFieldsSoFarRelations(ctx, results, selects)
		},
		params.Skip,
	)
}

func (q *Queries) executeAllFieldsSoFarFindMany(
	ctx context.Context,
	params QueryParams,
	selects *AllFieldsSoFarSelect,
	omits *AllFieldsSoFarOmit,
) ([]*AllFieldsSoFar, error) {
	for _, p := range params.Where {
		if p != nil {
			if err := p.Validate(); err != nil {
				return nil, err
			}
		}
	}
	whereClause, vals := CompilePredicates(q.dialect, params.Where)
	if whereClause != "" {
		whereClause = " WHERE " + whereClause
	}
	returningCols := q.selectAllFieldsSoFarCols(selects, omits)
	return executeManyWithRelations(ctx, q, "AllFieldsSoFar", whereClause, vals, returningCols,
		func(res *AllFieldsSoFar, cols []string) []any { return res.ScanFields(cols) },
		selects.hasAnyRelation(),
		func(ctx context.Context, txQ *Queries, results []*AllFieldsSoFar) error {
			return txQ.loadAllFieldsSoFarRelations(ctx, results, selects)
		},
		params.Take,
		params.Skip,
	)
}
func (q *Queries) loadAllFieldsSoFarRelations(ctx context.Context, records []*AllFieldsSoFar, selects *AllFieldsSoFarSelect) error {
	if selects == nil || len(records) == 0 {
		return nil
	}

	return nil
}
