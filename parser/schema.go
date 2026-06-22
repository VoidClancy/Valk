package parser

import (
	"fmt"
)

type Schema struct {
	Models []Model
}

type Model struct {
	Name       string
	Fields     []Field
	Attributes []Attribute
}

type Field struct {
	Name       string
	Type       string
	IsArray    bool
	IsOptional bool
	Attributes []Attribute
}

type Attribute struct {
	Name string
	Args []Argument
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
	return []byte(fmt.Sprintf("%q", vt.String())), nil
}

type Value struct {
	Type   ValueType
	Scalar string     `json:",omitempty"`
	Array  []Value    `json:",omitempty"`
	Args   []Argument `json:",omitempty"`
	Left   *Value     `json:",omitempty"`
	Right  *Value     `json:",omitempty"`
}
