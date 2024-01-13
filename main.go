package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
)

func eval(jsonIn map[string]interface{}) (bool, error) {
	declarations := cel.Declarations(
		decls.NewVar("i", decls.NewMapType(decls.String, decls.Dyn)),
	)

	env, err := cel.NewEnv(declarations)
	if err != nil {
		log.Fatalf("failed to create environment: %s", err)
	}

	ast, iss := env.Parse("i.name == 'horst'")

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
	var jsonIn map[string]interface{}
	json.Unmarshal([]byte(`{"name": "horst"}`), &jsonIn)

	out, err := eval(jsonIn)
	if err != nil {
		log.Fatalf("failed to evaluate: %s", err)
	}

	fmt.Println(out)
}
