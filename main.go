package main

import (
	"fmt"
	"log"

	"github.com/google/cel-go/cel"
)

func main() {
	env, err := cel.NewEnv()
	if err != nil {
		log.Fatalf("failed to create environment: %s", err)
	}

	ast, iss := env.Parse(`{'name': 'horst'}`)

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

	out, _, err := program.Eval(cel.NoVars())
	if err != nil {
		log.Fatalf("failed to evaluate program: %s", err)
	}

	fmt.Println(out)

}
