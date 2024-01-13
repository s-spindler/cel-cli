package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
)

func main() {
	var jsonMap map[string]interface{}
	json.Unmarshal([]byte(`{"name": "horst"}`), &jsonMap)

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

	out, _, err := program.Eval(map[string]interface{}{"i": jsonMap})
	if err != nil {
		log.Fatalf("failed to evaluate program: %s", err)
	}

	fmt.Println(out)
}
