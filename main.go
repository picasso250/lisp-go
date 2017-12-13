package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	// "strings"
)

func main() {
	ast_list, err := parse(bufio.NewReader(os.Stdin))
	if err != nil {
		log.Fatal(err)
	}
	env := NewEnv()
	var v *Ast
	for _, e := range ast_list {
		v, err = eval(&e, env)
		if err != nil {
			log.Fatal(err)
		}
	}
	if v != nil {
		fmt.Printf(v.String())
	}
}
