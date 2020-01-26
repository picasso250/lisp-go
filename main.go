package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	// "strings"
)

func main() {
	astList, err := parse(bufio.NewReader(os.Stdin))
	if err != nil {
		log.Fatal(err)
	}
	env := NewEnv()
	var v *Ast
	for _, e := range astList {
		v, err = eval(&e, env)
		if err != nil {
			log.Fatal(err)
		}
	}
	if v != nil {
		fmt.Printf(v.String())
	}
}
