package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
)


func main() {
	env, err := cel.NewEnv(
		cel.Declarations(
			decls.NewVar("data", decls.NewMapType(decls.String, decls.Dyn)),
		),
	)

	if len(os.Args) != 2 {
		log.Fatal("This program supports one arg: the path to the cel program")
	}

	src, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("could not read program: %v", err)
	}

	ast, issues := env.Compile(string(src))
	if issues != nil && issues.Err() != nil {
		log.Fatalf("type-check error: %s", issues.Err())
	}

	prg, err := env.Program(ast)
	if err != nil {
		log.Fatalf("program construction error: %s", err)
	}

	stdinScanner := bufio.NewScanner(os.Stdin)

	for stdinScanner.Scan() {
		data := make(map[string]any)

		err := json.Unmarshal(stdinScanner.Bytes(), &data)
		if err != nil {
			log.Fatalf("could not parse json: %v", err)
		}

		out, _, err := prg.Eval(map[string]any {
			"data": data,
		})

		if err != nil {
			log.Fatalf("eval error", err)
		}

		fmt.Println(out)
	}
}
