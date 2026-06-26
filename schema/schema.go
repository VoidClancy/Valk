package schema

import (
	"encoding/json"
	"fmt"
)

type Severity int

const (
	SevError Severity = iota
	SevWarning
)

func (s Severity) String() string {
	if s == SevWarning {
		return "warning"
	}
	return "error"
}

type Position struct {
	Line, Col, Offset int
}

func (p Position) String() string {
	return fmt.Sprintf("%d:%d", p.Line, p.Col)
}

type Diagnostic struct {
	Severity Severity
	Message  string
	Pos      Position
	Source   string // "lexer" | "parser" | "resolver"
}

func (d Diagnostic) String() string {
	return fmt.Sprintf("%s: %s: %s", d.Pos, d.Severity, d.Message)
}

func (d Diagnostic) Error() string {
	return d.String()
}

type DiagnosticList []Diagnostic

func (dl DiagnosticList) HasErrors() bool {
	for _, d := range dl {
		if d.Severity == SevError {
			return true
		}
	}
	return false
}

type Schema struct {
	Datasource Datasource
	Models     []*Model
	Enums      []*Enum
	Errors     DiagnosticList
}

type Datasource struct {
	Name     string
	Provider string
}

type Enum struct {
	Name         string
	Values       []string
	ValueMap     []EnumValue
	TableMapName string
}

type EnumValue struct {
	Name   string
	DBName string
}

type Model struct {
	Name            string
	TableName       string
	ScalarFields    []*ScalarField
	RelationFields  []*RelationField
	Attributes      []Attribute
	CompositePK     []string
	CompositeUnique []UniqueConstraint
	Indexes         []Index
}

type UniqueConstraint struct {
	Fields []string
	Name   string
}

type Index struct {
	Fields []string
	Name   string
}

type DefaultKind int

const (
	DefaultLiteral DefaultKind = iota
	DefaultFunc
	DefaultDBGenerated
	DefaultEnumValue
)

func (k DefaultKind) String() string {
	switch k {
	case DefaultLiteral:
		return "Literal"
	case DefaultFunc:
		return "Func"
	case DefaultDBGenerated:
		return "DBGenerated"
	case DefaultEnumValue:
		return "EnumValue"
	default:
		return "Unknown"
	}
}

type DefaultValue struct {
	Value        string // "autoincrement()", ENUM VAL, now(), etc
	Type         string // Literal, Func, etc
	Kind         DefaultKind
	Literal      string
	FuncName     string
	DBExpression string
	EnumValue    string
}

type ScalarField struct {
	Name        string
	Type        string // raw PSL type name, Int, String, etc
	ColName     string
	GoType      string
	SQLType     string
	IsArray     bool
	Optional    bool
	IsID        bool
	IsUnique    bool
	IsUpdatedAt bool
	Default     *DefaultValue
	EnumRef     *Enum
	NativeType  *NativeType
	Attributes  []Attribute
}

type NativeType struct {
	Name string
	Args []string
}

type RelationKind int

const (
	RelOneToOne RelationKind = iota
	RelOneToMany
	RelManyToOne
	RelManyToMany
)

func (k RelationKind) String() string {
	switch k {
	case RelOneToOne:
		return "OneToOne"
	case RelOneToMany:
		return "OneToMany"
	case RelManyToOne:
		return "ManyToOne"
	case RelManyToMany:
		return "ManyToMany"
	default:
		return "Unknown"
	}
}

type RelationField struct {
	Name            string
	Type            string // target model name
	Kind            RelationKind
	TargetModelName string
	TargetModel     *Model
	FKFields        []*ScalarField
	RefFields       []*ScalarField
	OnDelete        string
	OnUpdate        string
	Optional        bool
	IsArray         bool
	Inverse         *RelationField
	RelationName    string
	JoinTableName   string
}

type Attribute struct {
	Name string
	Args []Argument
	Line int
	Col  int
}

type Argument struct {
	Name  string
	Value Value
}

type ValueType int

const (
	ValLiteral ValueType = iota
	ValIdent
	ValArray
	ValFunc
	ValBinary
)

func (vt ValueType) String() string {
	switch vt {
	case ValLiteral:
		return "Literal"
	case ValIdent:
		return "Ident"
	case ValArray:
		return "Array"
	case ValFunc:
		return "Func"
	case ValBinary:
		return "Binary"
	default:
		return "Unknown"
	}
}

func (vt ValueType) MarshalJSON() ([]byte, error) {
	return fmt.Appendf(nil, "%q", vt.String()), nil
}

type Value struct {
	Type   ValueType
	Scalar string     `json:",omitempty"`
	Array  []Value    `json:",omitempty"`
	Args   []Argument `json:",omitempty"`
	Left   *Value     `json:",omitempty"`
	Right  *Value     `json:",omitempty"`
}

func (rf *RelationField) MarshalJSON() ([]byte, error) {
	type Alias RelationField
	var inverseName string
	if rf.Inverse != nil {
		inverseName = rf.Inverse.Name
	}
	var targetModelName string
	if rf.TargetModel != nil {
		targetModelName = rf.TargetModel.Name
	}
	return json.Marshal(&struct {
		*Alias
		TargetModel string `json:"TargetModel,omitempty"`
		Inverse     string `json:"Inverse,omitempty"`
	}{
		Alias:       (*Alias)(rf),
		TargetModel: targetModelName,
		Inverse:     inverseName,
	})
}

func (sf *ScalarField) MarshalJSON() ([]byte, error) {
	type Alias ScalarField
	var enumName string
	if sf.EnumRef != nil {
		enumName = sf.EnumRef.Name
	}
	return json.Marshal(&struct {
		*Alias
		EnumRef string `json:"EnumRef,omitempty"`
	}{
		Alias:   (*Alias)(sf),
		EnumRef: enumName,
	})
}
