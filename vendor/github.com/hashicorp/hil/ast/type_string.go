// Code generated by "stringer -type=Type"; DO NOT EDIT

package ast

import "fmt"

const (
	_Type_name_0 = "TypeInvalid"
	_Type_name_1 = "TypeAny"
	_Type_name_2 = "TypeString"
	_Type_name_3 = "TypeInt"
	_Type_name_4 = "TypeFloat"
	_Type_name_5 = "TypeList"
)

var (
	_Type_index_0 = [...]uint8{0, 11}
	_Type_index_1 = [...]uint8{0, 7}
	_Type_index_2 = [...]uint8{0, 10}
	_Type_index_3 = [...]uint8{0, 7}
	_Type_index_4 = [...]uint8{0, 9}
	_Type_index_5 = [...]uint8{0, 8}
)

func (i Type) String() string {
	switch {
	case i == 0:
		return _Type_name_0
	case i == 2:
		return _Type_name_1
	case i == 4:
		return _Type_name_2
	case i == 8:
		return _Type_name_3
	case i == 16:
		return _Type_name_4
	case i == 32:
		return _Type_name_5
	default:
		return fmt.Sprintf("Type(%d)", i)
	}
}
