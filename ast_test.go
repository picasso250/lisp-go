package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEval(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"1", "1"},
		{"(+ 1 2)", "3"},
		{"(atom (' ()))", "nil"},
		{"(atom (' a))", "#t"},
		{"(car ('(a b)))", "a"},
		{"(cdr ('(a b)))", "(b)"},
		{"(cons 'a ('(b)))", "(a b)"},
		{"(eq (' a) (' b))", "nil"},
		{"(eq (' a) (' a))", "#t"},
		{"(eq (' a) (' a))", "#t"},
	}
	for _, c := range cases {
		ast_list, err := parse(bufio.NewReader(strings.NewReader(c.in)))
		if err != nil {
			t.Error(err)
		}
		env := &Env{}
		var v *Ast
		for _, e := range ast_list {
			v, err = eval(&e, env)
			if err != nil {
				t.Error(err)
			}
		}
		got := v.String()
		if got != c.want {
			t.Errorf("%q == %q, want %q", c.in, got, c.want)
		}
	}
}
func testFile(t *testing.T, file_name string) {
	file, err := os.Open(file_name)
	if err != nil {
		t.Error(err)
	}
	defer file.Close()
	rd := bufio.NewReader(file)
	line, err := rd.ReadString('\n')
	if err != nil {
		t.Error(err)
	}
	want := strings.TrimSpace(line)
	ast_list, err := parse(rd)
	if err != nil {
		t.Error(err)
	}
	env := &Env{}
	var v *Ast
	for _, e := range ast_list {
		// fmt.Printf("eval %s\n", e.SimpleString())
		v, err = eval(&e, env)
		if err != nil {
			t.Error(err)
		}
	}
	got := v.String()
	if got != want {
		t.Errorf("%q == %q, want %q", "", got, want)
	}
}

func TestCaseFile(t *testing.T) {
	err := filepath.Walk("test_case", func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			fmt.Printf("Test:\t%s\n", path)
			testFile(t, path)
		}
		return nil
	})
  if err != nil {
    t.Errorf("error: %s", err.Error())
  }
}
