package syntax

import (
	"fmt"
	"mola/language"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"unicode"
)

// this is not efficient, but just working on getting it to "work" in the mal exercise
type TokenReader struct {
	tokens *TokenElement
	len    int
}

type TokenElement struct {
	v    string
	next *TokenElement
}

func NewReader() *TokenReader {
	return &TokenReader{tokens: nil, len: 0}
}

// return current token in incoming stream
func (rdr *TokenReader) peek() (string, error) {
	if rdr.tokens == nil {
		return "", fmt.Errorf("reader token list not initialized")
	}
	if rdr.len == 0 {
		return "", fmt.Errorf("reader token stream empty")
	}

	// sanity check
	if len(rdr.tokens.v) == 0 {
		return "", fmt.Errorf("reader token empty")
	}
	return rdr.tokens.v, nil
}

// return current token AND increment stream position
func (rdr *TokenReader) next() (string, error) {
	token, err := rdr.peek()
	if err != nil {
		return "", err
	}
	rdr.tokens = rdr.tokens.next
	rdr.len -= 1
	return token, nil
}

func ReadStr(s string) (*language.MalValue, error) {
	// should call tokenize
	tok_reader, err := Tokenize(s)
	if err != nil {
		return nil, err
	}
	// then call read_form
	mal, err := read_form(&tok_reader)
	if err != nil {
		return nil, err
	}
	return mal, nil
}

func read_list(r *TokenReader) (*language.MalValue, error) {
	// initializing lists with a capacity of 6 for now
	// first thing we do, is consume the "(" token
	left_paren, err := r.next()
	if err != nil {
		return nil, err
	}
	if left_paren[0] != '(' {
		return nil, fmt.Errorf("read_list called on non-list starting token: \"%s\"", left_paren)
	}

	list := language.NewList()

	// for debugging, will disable or remove in final build
	deadlock := 0
	for {
		deadlock += 1
		tok, err := r.peek()

		// if we get an error here, that means we reached the end of our token stream
		// so we're missing a closing parentheses ")"
		if err != nil {
			return nil, fmt.Errorf("\")\" missing in list declaration; unexpectedly found end of token stream")
		}

		if tok[0] == ')' {
			// let's also pop this token off the reader
			_, err := r.next()
			if err != nil {
				return nil, err
			}
			return &list, nil
		}

		mv, err := read_form(r)
		if err != nil {
			return nil, err
		}

		if mv == nil {
			return nil, fmt.Errorf("unexpceted end of token stream found.  Expected \")\" ")
		}
		list.List = append(list.List, *mv)

		if deadlock > 100 {
			return nil, fmt.Errorf("DEADLOCK FOUND")
		}

	}
}

func is_integer(s string) bool {
	r := true
	for idx, c := range s {
		if !unicode.IsDigit(c) {
			// allow negative sign
			if idx == 0 && c == '-' && len(s) > 1 {
				continue
			}
			return false
		}
	}
	return r
}

func is_symbol(s string) bool {
	r := true
	symbolic_chars := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!+-=/*"

	for _, c := range s {
		if !strings.ContainsRune(symbolic_chars, c) {
			return false
		}
	}
	return r
}

func parse_string(s string) (string, error) {
	return s, nil
}

func read_atom(r *TokenReader) (*language.MalValue, error) {
	tok, err := r.next()
	if err != nil {
		return nil, err
	}

	if is_integer(tok) {
		num, err := strconv.Atoi(tok)
		if err != nil {
			return nil, err
		}
		mv := language.NewInt(num, &tok)
		return &mv, nil
	}

	// string parsing
	if tok[0] == '"' || tok[0] == '\'' {
		mv := language.NewString(tok)
		return &mv, nil
	}

	if tok == "nil" {
		mv := language.NewNil()
		return &mv, nil
	}

	if is_symbol(tok) {
		mv := language.NewSymbol(tok)
		return &mv, nil
	}

	return nil, fmt.Errorf("unknown atomic type token: \"%s\"", tok)
}

// builds AST
func read_form(r *TokenReader) (*language.MalValue, error) {
	tok, err := r.peek()
	if err != nil {
		return nil, err
	}
	var mal *language.MalValue
	switch tok[0] {
	// list handle
	case '(':
		mal, err = read_list(r)
		if err != nil {
			return nil, err
		}
	// comment handle
	case ';':
		return nil, nil
	default:
		mal, err = read_atom(r)
		if err != nil {
			return nil, err
		}
	}
	return mal, nil
}

func Tokenize(s string) (TokenReader, error) {
	r := NewReader()
	// \x60 = backticks
	reg_expression, err := regexp.Compile(`[\s,]*(~@|[\[\]{}()'\x60~^@]|"(?:\\.|[^\\"])*"?|;.*|[^\s\[\]{}('"\x60,;)]*)`)
	if err != nil {
		return *r, err
	}

	// second parameter is the length of the provided value?  Defaults to len(s)+1 if we provided a
	// below zero value based on the google source code: https://cs.opensource.google/go/go/+/refs/tags/go1.25.4:src/regexp/regexp.go;l=1114
	tokens := reg_expression.FindAllString(s, -1)
	if tokens == nil {
		return *r, fmt.Errorf("no tokens found")
	}

	// Debug
	// fmt.Println(strings.Join(tokens, ","))

	// I need to flip tokens here
	slices.Reverse(tokens)
	for _, tok := range tokens {
		r.len += 1
		elem := TokenElement{v: strings.TrimSpace(tok), next: r.tokens}
		r.tokens = &elem
	}

	return *r, nil
}
