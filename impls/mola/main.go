package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"mola/language"
	"mola/syntax"
	"os"
	"os/signal"
)

/*
LISP FUNCTIONS
*/
func READ(s string) (*language.MalValue, error) {
	mv, err := syntax.ReadStr(s)
	if err != nil {
		return nil, err
	}
	return mv, nil
}

func EVAL(ast language.MalValue, env MalEnv) (*language.MalValue, error) {

	if env == nil {
		return nil, fmt.Errorf("nil env provided")
	}

	switch ast.TypeId {
	case language.Symbol:
		f, ok := env[ast.Symbol]
		if !ok {
			return nil, fmt.Errorf("symbol \"%s\" definition not found in environment", ast.Symbol)
		}
		return f, nil
	case language.List:
		// ok I'd want to do this in-place though
		// but for now we do it badly so we can get it to work
		var i_list []language.MalValue
		for _, mv := range ast.List {
			eval_mv, err := EVAL(mv, env)
			if err != nil {
				return nil, err
			}
			i_list = append(i_list, *eval_mv)
		}
		// check list has length
		if len(i_list) > 0 {
			// list doesn't start with a function
			if i_list[0].TypeId != language.Function {
				rv := language.PackList(i_list)
				return &rv, nil
			}

			// super sick usage of golang's unpack operator (...)
			rv, err := i_list[0].Function(i_list[1:]...)
			if err != nil {
				return nil, err
			}
			return rv, nil
		}
	}

	// default case--return atomic value
	return &ast, nil
}
func PRINT(mv language.MalValue) string {
	s, err := syntax.Pr_Str(mv)
	if err != nil {
		return fmt.Sprintf("String representation error: %s\n", err.Error())
	}
	return s + "\n"
}

/* */
type MalEnv map[string]*language.MalValue

func rep(s string, env MalEnv) string {
	mv, err := READ(s)
	if err != nil {
		return fmt.Sprintf("Syntax error: %s\n", err.Error())
	}
	// if mv is nil--return empty string.
	// this is in the scenario where we encounter a comment immediately
	if mv == nil {
		return ""
	}

	r, err := EVAL(*mv, env)
	if err != nil {
		return fmt.Sprintf("Eval error: %s\n", err.Error())
	}

	return PRINT(*r)
}

// OS Level stuff
func readline(reader bufio.Reader) (string, error) {
	content, err := reader.ReadBytes('\n')
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// TODO: add command historyy & line editing
func run(ctx context.Context, reader bufio.Reader, writer io.Writer) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	add_fn := language.NewFunction(language.I_Add)
	sub_fn := language.NewFunction(language.I_Sub)
	mul_fn := language.NewFunction(language.I_Mul)
	div_fn := language.NewFunction(language.I_Div)

	mal_env := MalEnv{
		"+": &add_fn,
		"-": &sub_fn,
		"*": &mul_fn,
		"/": &div_fn,
	}

	for {
		// REPL loop
		fmt.Fprintf(writer, "user> ")
		content, err := readline(reader)
		if err != nil {
			return err
		}

		mal_value := rep(content, mal_env)

		fmt.Fprintf(writer, "%s", mal_value)

		// close our loop if ctrl+c was hit
		select {
		case <-ctx.Done():
			return nil
		default:
		}

	}

}

func main() {
	// I think it's possible in theory for NewReader to return nil which would cause a nil-pointer reference exception
	if err := run(context.Background(), *bufio.NewReader(os.Stdin), os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
