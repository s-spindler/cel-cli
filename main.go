package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/spf13/pflag"
)

func eval(jsonIn map[string]interface{}, expression string) (bool, error) {
	declarations := cel.Declarations(
		decls.NewVar("i", decls.NewMapType(decls.String, decls.Dyn)),
	)

	env, err := cel.NewEnv(declarations)
	if err != nil {
		log.Fatalf("failed to create environment: %s", err)
	}

	ast, iss := env.Parse(expression)

	if iss.Err() != nil {
		log.Fatalf("failed to parse: %s", iss.Err())
	}

	checked, iss := env.Check(ast)
	if iss.Err() != nil {
		log.Fatalf("failed to check AST: %s", iss.Err())
	}

	program, err := env.Program(checked)
	if err != nil {
		log.Fatalf("failed to create program: %s", err)
	}

	out, _, err := program.Eval(map[string]interface{}{"i": jsonIn})
	if err != nil {
		log.Fatalf("failed to evaluate program: %s", err)
	}
	if out.Type() != cel.BoolType {
		log.Fatalf("expression did not evaluate to boolean but was of type: %s", out.Type())
	}

	return out.Value().(bool), nil
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

	var jsonMap map[string]interface{}
	json.Unmarshal([]byte(jsonIn), &jsonMap)

	result, err := eval(jsonMap, expression)
	if err != nil {
		log.Fatalf("failed to evaluate: %s", err)
	}

	if !result {
		os.Exit(1)
	}
	// else: fall through and return default exit code 0
}
