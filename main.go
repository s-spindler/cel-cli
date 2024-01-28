package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types/ref"
	"github.com/spf13/pflag"
)

func eval(jsonIn string, expression string) (result ref.Val, err error) {

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

	result, _, err = program.Eval(map[string]interface{}{"i": jsonMap})
	if err != nil {
		err = fmt.Errorf("failed to evaluate program: %w", err)
		return
	}

	return
}

func main() {
	var (
		jsonIn         string
		expression     string
		boolean_result bool
		quiet_mode     bool
	)

	flags := pflag.NewFlagSet("cel-cli", pflag.ExitOnError)
	flags.StringVarP(&jsonIn, "input-json", "i", "", "JSON input to the expression.")
	flags.StringVarP(&expression, "expression", "e", "", "The expression to evaluate.")
	flags.BoolVarP(&boolean_result, "force-bool", "b", false,
		"Forces the expression to evaluate to a boolean value, terminating with a non-zero "+
			"status code otherwise. _Note_: Only a true expression will exit the program with "+
			"0 while false will be non-zero.")
	flags.BoolVarP(&quiet_mode, "quiet", "q", false,
		"Omits printing of the expression's result. Error messages are still printed.")

	flags.Parse(os.Args[1:])

	result, err := eval(jsonIn, expression)
	if err != nil {
		log.Fatalf("failed to evaluate: %s", err)
	}

	out := result.Value()

	exitCode := 0
	print := !quiet_mode

	if boolean_result {
		if result.Type() != cel.BoolType {
			exitCode = 2
			print = false
		} else if !out.(bool) {
			exitCode = 1
		}
	}

	if print {
		fmt.Println(out)
	}

	os.Exit(exitCode)
}
