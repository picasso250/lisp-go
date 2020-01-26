package main

import (
	"bufio"
	"io"
	"os"
	"unicode"
	// "fmt"
)

func parseStdin() ([]Ast, error) {
	return parse(bufio.NewReader(os.Stdin))
}
func parse(inStream *bufio.Reader) ([]Ast, error) {
	// only support S-expression and atom
	ret := make([]Ast, 0)
	var cur *Ast
	state := 0
	sb := make([]rune, 0)
	in := inStream

	for {
		c, _, err := in.ReadRune()
		if err == io.EOF {
			if len(sb) > 0 {
				// if only an atom
				ret = append(ret, *NewASTreeAtom(string(sb)))
			}
			return ret, nil
		}
		if err != nil {
			return nil, err
		}
		switch state {
		case 0: // initial state
			if unicode.IsSpace(c) {
				continue
			}
			if c == '(' {
				if cur == nil {
					cur = NewASTreeNil()
				} else {
					t := NewASTreeNilParent(cur)
					cur.add(t)
					cur = t
				}
			} else if c == ')' {
				if cur.parent == cur { // expression ends
					ret = append(ret, *cur)
					cur = nil
				} else {
					cur = cur.parent
				}
			} else {
				// word
				state = 1
				sb = append(sb, c)
			}
			break
		case 1: // in an atom
			if unicode.IsSpace(c) || c == '(' || c == ')' {
				e := NewASTreeAtom(string(sb))
				if cur == nil {
					// in root
					ret = append(ret, *NewASTreeAtom(string(sb)))
				} else {
					cur.add(e)
				}
				sb = make([]rune, 0)
				state = 0
				if c == '(' || c == ')' {
					err = in.UnreadByte()
					if err != nil {
						return ret, err
					}
				}
			} else {
				sb = append(sb, c)
			}
			break
		}
	}
}
