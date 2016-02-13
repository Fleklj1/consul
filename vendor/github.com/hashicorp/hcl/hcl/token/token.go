// Package token defines constants representing the lexical tokens for HCL
// (HashiCorp Configuration Language)
package token

import (
	"fmt"
	"strconv"
	"strings"

	hclstrconv "github.com/hashicorp/hcl/hcl/strconv"
)

// Token defines a single HCL token which can be obtained via the Scanner
type Token struct {
	Type Type
	Pos  Pos
	Text string
	JSON bool
}

// Type is the set of lexical tokens of the HCL (HashiCorp Configuration Language)
type Type int

const (
	// Special tokens
	ILLEGAL Type = iota
	EOF
	COMMENT

	identifier_beg
	IDENT // literals
	literal_beg
	NUMBER  // 12345
	FLOAT   // 123.45
	BOOL    // true,false
	STRING  // "abc"
	HEREDOC // <<FOO\nbar\nFOO
	literal_end
	identifier_end

	operator_beg
	LBRACK // [
	LBRACE // {
	COMMA  // ,
	PERIOD // .

	RBRACK // ]
	RBRACE // }

	ASSIGN // =
	ADD    // +
	SUB    // -
	operator_end
)

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",

	EOF:     "EOF",
	COMMENT: "COMMENT",

	IDENT:  "IDENT",
	NUMBER: "NUMBER",
	FLOAT:  "FLOAT",
	BOOL:   "BOOL",
	STRING: "STRING",

	LBRACK:  "LBRACK",
	LBRACE:  "LBRACE",
	COMMA:   "COMMA",
	PERIOD:  "PERIOD",
	HEREDOC: "HEREDOC",

	RBRACK: "RBRACK",
	RBRACE: "RBRACE",

	ASSIGN: "ASSIGN",
	ADD:    "ADD",
	SUB:    "SUB",
}

// String returns the string corresponding to the token tok.
func (t Type) String() string {
	s := ""
	if 0 <= t && t < Type(len(tokens)) {
		s = tokens[t]
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(t)) + ")"
	}
	return s
}

// IsIdentifier returns true for tokens corresponding to identifiers and basic
// type literals; it returns false otherwise.
func (t Type) IsIdentifier() bool { return identifier_beg < t && t < identifier_end }

// IsLiteral returns true for tokens corresponding to basic type literals; it
// returns false otherwise.
func (t Type) IsLiteral() bool { return literal_beg < t && t < literal_end }

// IsOperator returns true for tokens corresponding to operators and
// delimiters; it returns false otherwise.
func (t Type) IsOperator() bool { return operator_beg < t && t < operator_end }

// String returns the token's literal text. Note that this is only
// applicable for certain token types, such as token.IDENT,
// token.STRING, etc..
func (t Token) String() string {
	return fmt.Sprintf("%s %s %s", t.Pos.String(), t.Type.String(), t.Text)
}

// Value returns the properly typed value for this token. The type of
// the returned interface{} is guaranteed based on the Type field.
//
// This can only be called for literal types. If it is called for any other
// type, this will panic.
func (t Token) Value() interface{} {
	switch t.Type {
	case BOOL:
		if t.Text == "true" {
			return true
		} else if t.Text == "false" {
			return false
		}

		panic("unknown bool value: " + t.Text)
	case FLOAT:
		v, err := strconv.ParseFloat(t.Text, 64)
		if err != nil {
			panic(err)
		}

		return float64(v)
	case NUMBER:
		v, err := strconv.ParseInt(t.Text, 0, 64)
		if err != nil {
			panic(err)
		}

		return int64(v)
	case IDENT:
		return t.Text
	case HEREDOC:
		// We need to find the end of the marker
		idx := strings.IndexByte(t.Text, '\n')
		if idx == -1 {
			panic("heredoc doesn't contain newline")
		}

		return string(t.Text[idx+1 : len(t.Text)-idx+1])
	case STRING:
		// Determine the Unquote method to use. If it came from JSON,
		// then we need to use the built-in unquote since we have to
		// escape interpolations there.
		f := hclstrconv.Unquote
		if t.JSON {
			f = strconv.Unquote
		}

		v, err := f(t.Text)
		if err != nil {
			panic(fmt.Sprintf("unquote %s err: %s", t.Text, err))
		}

		return v
	default:
		panic(fmt.Sprintf("unimplemented Value for type: %s", t.Type))
	}
}
