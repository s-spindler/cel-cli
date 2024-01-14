package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/spf13/pflag"
)

func eval(jsonIn string, expression string) (result bool, err error) {

	var jsonMap map[string]interface{}
	err = json.Unmarshal([]byte(jsonIn), &jsonMap)
	if err != nil {
		err = fmt.Errorf("failed to parse input: %w", err)
		return
	}

	declarations := cel.Declarations(
		decls.NewVar("i", decls.NewMapType(decls.String, decls.Dyn)),
	)

	env, err := cel.NewEnv(declarations)
	if err != nil {
		err = fmt.Errorf("failed to create environment: %w", err)
		return
	}

	ast, iss := env.Parse(expression)

	if iss.Err() != nil {
		err = fmt.Errorf("failed to parse: %w", iss.Err())
		return
	}

	checked, iss := env.Check(ast)
	if iss.Err() != nil {
		err = fmt.Errorf("failed to check AST: %w", iss.Err())
		return
	}

	program, err := env.Program(checked)
	if err != nil {
		err = fmt.Errorf("failed to create program: %w", err)
		return
	}

	out, _, err := program.Eval(map[string]interface{}{"i": jsonMap})
	if err != nil {
		err = fmt.Errorf("failed to evaluate program: %w", err)
		return
	}
	if out.Type() != cel.BoolType {
		err = fmt.Errorf("expression did not evaluate to boolean but was of type: %s", out.Type())
		return
	}

	result = out.Value().(bool)
	return result, nil
}

func main() {
	var (
		jsonIn     string
		expression string
	)

	flags := pflag.NewFlagSet("cel-cli", pflag.ExitOnError)
	flags.StringVarP(&jsonIn, "input-json", "i", "", "JSON input")
	flags.StringVarP(&expression, "expression", "e", "", "expression to evaluate")

	flags.Parse(os.Args[1:])

	result, err := eval(jsonIn, expression)
	if err != nil {
		log.Fatalf("failed to evaluate: %s", err)
	}

	if !result {
		os.Exit(1)
	}
	// else: fall through and return default exit code 0
}
