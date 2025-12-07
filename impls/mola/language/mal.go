package language

import (
	"fmt"
	"strconv"
)

type MalTypeId int

const (
	List MalTypeId = iota
	Function
	Integer
	Float
	String
	Symbol
	Bool
	Nil
)

var TYPE_PRINTABLE = map[MalTypeId]string{
	List:     "List",
	Function: "Function",
	Integer:  "Integer",
	Float:    "Float",
	Symbol:   "Symbol",
	Bool:     "Bool",
	Nil:      "Nil",
}

// go type for function
type InternalFunc func(...MalValue) (*MalValue, error)

type MalValue struct {
	List     []MalValue
	TypeId   MalTypeId
	Integer  int
	Float    float64
	Symbol   string
	Function InternalFunc
}

func NewFunction(fn InternalFunc) MalValue {
	return MalValue{
		Function: fn,
		TypeId:   Function,
	}
}

/*Built in functions*/

func I_Sub(vars ...MalValue) (*MalValue, error) {

	if len(vars) == 0 {
		return nil, fmt.Errorf("internal sub called with no-args")
	}

	id := vars[0].TypeId
	var acc MalValue = vars[0]
	for _, v := range vars[1:] {
		if v.TypeId != id {
			return nil, TypeMismatch("sub", id, v.TypeId)
		}
		switch id {
		case Integer:
			acc.Integer -= v.Integer
		case Float:
			acc.Float -= v.Float
		default:
			return nil, fmt.Errorf("unsupported operation on type \"%s\"", TypeString(id))
		}

	}

	return &acc, nil
}

func TypeMismatch(m_name string, a MalTypeId, b MalTypeId) error {
	a_str, a_ok := TYPE_PRINTABLE[a]
	b_str, b_ok := TYPE_PRINTABLE[b]
	if !a_ok || !b_ok {
		return fmt.Errorf("internal operator \"%s\" called with mismatched types: type_id=\"%d\" and type_id=\"%d\"", m_name, a, b)
	}
	return fmt.Errorf("internal operator \"%s\" called with mismatched types: \"%s\" and \"%s\"", m_name, a_str, b_str)
}

func TypeString(a MalTypeId) string {
	str, ok := TYPE_PRINTABLE[a]
	if !ok {
		return fmt.Sprint(a)
	}
	return str
}

func I_Add(vars ...MalValue) (*MalValue, error) {

	if len(vars) == 0 {
		return nil, fmt.Errorf("internal add called with no-args")
	}

	id := vars[0].TypeId
	var acc MalValue = vars[0]
	for _, v := range vars[1:] {
		if v.TypeId != id {
			return nil, TypeMismatch("add", id, v.TypeId)
		}
		switch id {
		case Integer:
			acc.Integer += v.Integer
		case Float:
			acc.Float += v.Float
		case String:
			// drop quotes so "hello" +"world" = "hello world" and not "hello""world"
			acc.Symbol = acc.Symbol[:len(acc.Symbol)-1] + v.Symbol[1:]
		default:
			return nil, fmt.Errorf("unsupported operation on type \"%s\"", TypeString(id))
		}

	}

	return &acc, nil
}

func I_Mul(vars ...MalValue) (*MalValue, error) {

	if len(vars) == 0 {
		return nil, fmt.Errorf("internal mul called with no-args")
	}

	id := vars[0].TypeId
	var acc MalValue = vars[0]
	for _, v := range vars[1:] {
		if v.TypeId != id {
			return nil, TypeMismatch("mul", id, v.TypeId)
		}
		switch id {
		case Integer:
			acc.Integer *= v.Integer
		case Float:
			acc.Float *= v.Float
		default:
			return nil, fmt.Errorf("unsupported operation on type \"%s\"", TypeString(id))
		}

	}

	return &acc, nil
}

func I_Div(vars ...MalValue) (*MalValue, error) {

	if len(vars) == 0 {
		return nil, fmt.Errorf("internal div called with no-args")
	}

	id := vars[0].TypeId
	var acc MalValue = vars[0]
	for _, v := range vars[1:] {
		if v.TypeId != id {
			return nil, TypeMismatch("div", id, v.TypeId)
		}
		switch id {
		case Integer:
			if v.Integer == 0 {
				return nil, fmt.Errorf("division by zero error")
			}
			acc.Integer /= v.Integer
		case Float:
			if v.Float == 0 {
				return nil, fmt.Errorf("division by zero error")
			}
			acc.Float /= v.Float
		default:
			return nil, fmt.Errorf("unsupported operation on type \"%s\"", TypeString(id))
		}

	}

	return &acc, nil
}

func NewInt(v int, s *string) MalValue {
	if s == nil {
		return MalValue{List: nil, TypeId: Integer, Symbol: strconv.Itoa(v), Integer: v, Float: 0.0}
	}
	return MalValue{List: nil, TypeId: Integer, Integer: v, Symbol: *s, Float: 0.0}
}

func NewNil() MalValue {
	return MalValue{List: nil, TypeId: Nil}
}

// for booleans I'm just going to store them in the integer field
func NewBool(v bool) MalValue {
	if v {
		return MalValue{TypeId: Bool, Integer: 1}
	}
	return MalValue{TypeId: Bool, Integer: 0}
}

func NewString(s string) MalValue {
	return MalValue{List: nil, TypeId: String, Symbol: s, Integer: 0, Float: 0.0}
}

func NewSymbol(s string) MalValue {
	return MalValue{List: nil, TypeId: Symbol, Symbol: s, Integer: 0, Float: 0.0}
}

func NewFloat(f float64) MalValue {
	return MalValue{List: nil, TypeId: Float, Symbol: "<float>", Integer: 0, Float: f}
}

func NewList() MalValue {
	return MalValue{
		List:    make([]MalValue, 0, 6),
		Symbol:  "[...]",
		Integer: 0,
		TypeId:  List,
		Float:   0.0,
	}
}

// package []MalValue as a List MalValue type
func PackList(arr []MalValue) MalValue {
	return MalValue{
		List:   arr,
		Symbol: "[...]",
		TypeId: List,
	}
}
