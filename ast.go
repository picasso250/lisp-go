package main

import (
	"errors"
	"fmt"

	// "fmt"
	"strconv"
	"strings"
)

// http://www.jianshu.com/p/509505d3bd50
func eval(expr *Ast, env *Env) (*Ast, error) {
	// #t
	if expr.isAtom() && expr.equalsStr("#t") {
		return expr, nil
	}
	// look up variable
	if expr.isAtom() {
		v := env.get(expr.atom)
		if v == nil {
			return nil, errors.New("undefined variable " + expr.atom)
		}
		return v, nil
	}
	// preserve int
	if expr.isInt() {
		return expr, nil
	}

	// only the root of lisp ( 7 function)
	if expr.isList() {
		if expr.hasHead() {
			head := expr.getHead()
			tail := expr.getTail()
			return evalList(head, tail, env)
		} else {
			return expr, nil
		}
	}
	return NewASTreeNil(), nil
}

func evalList(head *Ast, tail *Ast, env *Env) (*Ast, error) {
	if head.isAtom() {
		// Quote
		if head.equalsStr("quote") || head.equalsStr("'") {
			return tail.get(0), nil
		}
		// Atom
		if head.equalsStr("atom") {
			v, err := eval(tail.get(0), env)
			if err != nil {
				return nil, err
			}
			return NewAstBool(v.isAtom()), nil
		}
		// Eq
		if head.equalsStr("eq") {
			v1, err := eval(tail.get(0), env)
			if err != nil {
				return nil, err
			}
			v2, err := eval(tail.get(1), env)
			if err != nil {
				return nil, err
			}
			if v1.isAtom() && v2.isAtom() {
				return NewAstBool(v1.equals(v2)), nil
			}
			if v1.isInt() && v2.isInt() {
				return NewAstBool(v1.ival == v2.ival), nil
			}
			return NewAstBool(v1.isNil() && v2.isNil()), nil
		}
		// Car
		if head.equalsStr("car") {
			v, err := eval(tail.get(0), env)
			if err != nil {
				return nil, err
			}
			return v.get(0), nil
		}
		// Cdr
		if head.equalsStr("cdr") {
			v, err := eval(tail.get(0), env)
			if err != nil {
				return nil, err
			}
			return v.getTail(), nil
		}
		// Cons
		if head.equalsStr("cons") {
			v1, err := eval(tail.get(0), env)
			if err != nil {
				return nil, err
			}
			v2, err := eval(tail.get(1), env)
			if err != nil {
				return nil, err
			}
			return cons(v1, v2), nil
		}
		// Cond
		if head.equalsStr("cond") {
			return cond(tail, env)
		}

		if head.equalsStr("let") { // (let (name val) expr)
			name_val := tail.get(0)
			val, err := eval(name_val.get(1), env)
			if err != nil {
				return nil, err
			}
			ne := NewEnvNameValue(name_val.get(0).atom, val)
			env_new := extend(env, ne)
			return eval(tail.get(1), env_new)
		}

		// defun
		if head.equalsStr("defun") { // (defun name value) or (defun name (arg...) e)
			// defun a var or a function
			err := doDef(tail, env)
			if err != nil {
				return nil, err
			}
			return tail, nil
		}

		// lambda
		if head.equalsStr("lambda") { // (lambda (arg...) e)
			_param := tail.get(0)
			_body := tail.get(1)
			c := Closure{_param.toStringList(), _body, env, ""}
			return NewASTreeClosure(c), nil
		}

		// + - *
		if head.atom == "+" || head.atom == "-" || head.atom == "*" {
			v1, err := eval(tail.get(0), env)
			if err != nil {
				return nil, err
			}
			v2, err := eval(tail.get(1), env)
			if err != nil {
				return nil, err
			}
			if head.equalsStr("+") {
				return add(tail, env)
			} else if head.equalsStr("-") {
				return sub(v1, v2), nil
			} else if head.equalsStr("*") {
				return mul(tail, env)
			}
		}

		if head.equalsStr("list") {
			ls := make([]*Ast, 0, len(tail.lst))
			for _, v := range tail.lst {
				vv, err := eval(v, env)
				if err != nil {
					return nil, err
				}
				ls = append(ls, vv)
			}
			return NewASTreeList(ls), nil
		}
	}
	// apply user defined function
	hval, err := eval(head, env)
	if err != nil {
		return nil, err
	}
	if hval.isLambda() {
		return apply(hval, tail.lst, env)
	} else {
		return nil, errors.New("not a function")
	}
}
func apply(fun *Ast, lst []*Ast, env *Env) (*Ast, error) {
	c := fun.closure
	e := &Env{}
	fmt.Printf("apply ")
	for i, name := range c.params {
		v, err := eval(lst[i], env)
		fmt.Printf("%s ", v.String())
		if err != nil {
			return nil, err
		}
		e.add(name, v)
	}
	fmt.Printf("\n")
	newEnv := extend(c.env, e)
	return eval(c.body, newEnv)
}
func doDef(tree *Ast, env *Env) error {
	name := tree.get(0)
	var value *Ast
	var err error
	if tree.len() == 2 {
		// variable
		value, err = eval(tree.get(1), env)
		if err != nil {
			return err
		}
	} else {
		// function
		_param := tree.get(1)
		_body := tree.get(2)
		closure := Closure{_param.toStringList(), _body, env, name.atom}
		value = NewASTreeClosure(closure)
	}
	env.add(name.atom, value)
	return nil
}

func cond(tree *Ast, env *Env) (*Ast, error) {
	for tree.len() > 0 {
		entry := tree.getHead()
		cond := entry.get(0)
		fmt.Printf("cond %s\n", cond.String())
		condRes, err := eval(cond, env)
		if err != nil {
			return nil, err
		}
		if condRes.isTrue() {
			return eval(entry.get(1), env)
		}
		tree = tree.getTail()
	}
	return NewASTreeNil(), nil
}

func cons(t1 *Ast, t2 *Ast) *Ast {
	lst := append([]*Ast(nil), t1)
	lst = append(lst, t2.lst...)
	return NewASTreeList(lst)
}

func add(v *Ast, env *Env) (*Ast, error) {
	s := 0
	for _, vv := range v.lst {
		vvv, err := eval(vv, env)
		if err != nil {
			return NewASTreeAtom(strconv.Itoa(s)), err
		}
		s += vvv.ival
	}
	return NewASTreeAtom(strconv.Itoa(s)), nil
}
func mul(v *Ast, env *Env) (*Ast, error) {
	s := 1
	for _, vv := range v.lst {
		vvv, err := eval(vv, env)
		if err != nil {
			return NewASTreeAtom(strconv.Itoa(s)), err
		}
		s *= vvv.ival
	}
	return NewASTreeAtom(strconv.Itoa(s)), nil
}
func sub(v1 *Ast, v2 *Ast) *Ast {
	i := v1.ival - v2.ival
	return NewASTreeAtom(strconv.Itoa(i))
}

const (
	TATOM  = 1
	TSEXPR = 2
	TLMD   = 3
	TINT   = 4
)

type Ast struct {
	parent  *Ast
	type_   int
	lst     []*Ast
	atom    string
	ival    int
	closure Closure
}

func NewASTreeAtom(atom string) *Ast {
	ast := Ast{}
	if atom == "nil" {
		ast.type_ = TSEXPR
	} else if IsDigit(atom[0]) {
		ast.ival, _ = strconv.Atoi(atom)
		ast.type_ = TINT
	} else if atom[0] == '\'' && len(atom) >= 2 {
		// 'a is short for (' a)
		ast.lst = append(ast.lst, NewASTreeAtom("quote"))
		ast.lst = append(ast.lst, NewASTreeAtom(atom[1:]))
		ast.type_ = (TSEXPR)
	} else {
		ast.atom = atom
		ast.type_ = TATOM
	}
	return &ast
}

func NewASTreeNilParent(parent *Ast) *Ast {
	ast := Ast{}
	ast.parent = parent
	ast.type_ = (TSEXPR)
	return &ast
}

func NewASTreeNil() *Ast {
	ast := Ast{}
	ast.parent = &ast
	ast.type_ = (TSEXPR)
	return &ast
}
func NewASTreeList(lst []*Ast) *Ast {
	ast := Ast{}
	ast.lst = lst
	ast.parent = &ast
	ast.type_ = (TSEXPR)
	return &ast
}

func NewASTreeClosure(closure Closure) *Ast {
	ast := Ast{}
	// lambda
	ast.closure = closure
	ast.type_ = (TLMD)
	return &ast
}
func NewAstBool(b bool) *Ast {
	if b {
		return NewASTreeAtom("#t")
	}
	return NewASTreeNil()
}
func (ast *Ast) String() string {
	if ast.type_ == TATOM || ast.type_ == TINT || ast.type_ == TSEXPR {
		return ast.SimpleString()
	}
	if ast.type_ == TLMD {
		return ast.closure.String()
	}
	return ""
}
func (ast *Ast) SimpleString() string {
	if ast.type_ == TATOM {
		return ast.atom
	}
	if ast.type_ == TSEXPR && ast.len() == 0 {
		return "nil"
	}
	if ast.type_ == TSEXPR {
		sm := MapAst(ast.lst, func(a *Ast) string {
			return a.SimpleString()
		})
		str := strings.Join(sm, " ")
		return "(" + str + ")"
	}
	if ast.type_ == TLMD {
		return ast.closure.SimpleString()
	}
	if ast.type_ == TINT {
		return strconv.Itoa(ast.ival)
	}
	return ""
}

func (ast *Ast) add(e *Ast) {
	e.parent = ast
	ast.lst = append(ast.lst, e)
}

func (ast *Ast) len() int {
	return len(ast.lst)
}

func (ast *Ast) isList() bool {
	return ast.type_ == TSEXPR
}

func (ast *Ast) getTail() *Ast {
	lst := append([]*Ast(nil), ast.lst[1:]...)
	return NewASTreeList(lst)
}

func (ast *Ast) isInt() bool {
	return ast.type_ == TINT
}

func (ast *Ast) isDef() bool {
	if ast.hasHead() && ast.getHead().isAtom() && ast.getHead().equalsStr("defun") {
		return true
	}
	return false
}

func (ast *Ast) isLambda() bool {
	return ast.type_ == TLMD
}

func (ast *Ast) isNil() bool {
	return ast.type_ == TSEXPR && ast.len() == 0
}

func (ast *Ast) isTrue() bool {
	if ast.type_ == TATOM {
		return ast.equalsStr("#t")
	}
	return ast.len() != 0
}
func (ast *Ast) isAtom() bool {
	return ast.type_ == TATOM
}
func (ast *Ast) equalsStr(str string) bool {
	return ast.atom == (str)
}
func (ast *Ast) equals(atom *Ast) bool {
	return ast.atom == (atom.atom)
}
func (ast *Ast) hasHead() bool {
	return ast.len() > 0
}
func (ast *Ast) getHead() *Ast {
	return ast.get(0)
}
func (ast *Ast) get(i int) *Ast {
	return ast.lst[i]
}

func (ast *Ast) isRoot() bool {
	return ast == ast.parent
}

func (ast *Ast) toStringList() []string {
	ls := make([]string, 0)
	for _, t := range ast.lst {
		ls = append(ls, t.atom)
	}
	return ls
}

// ==== closure ====

type Closure struct {
	params []string
	body   *Ast
	env    *Env
	name   string
}

func (c Closure) String() string {
	s := strings.Join(c.params, ",")
	return "Closure[env=" + c.env.String() + ", param=(" + s + "), body=" + c.body.String() + "]"
}
func (c Closure) SimpleString() string {
	s := strings.Join(c.params, ",")
	return "Closure(" + s + ")"
}

// ==== closure end ====

// ==== Env ====
type Env struct {
	map_   map[string]*Ast
	parent *Env
}

func NewEnvFromMap(map_ map[string]*Ast) *Env {
	env := Env{}
	env.map_ = map_
	return &env
}

func NewEnv() *Env {
	env := Env{map[string]*Ast{}, nil}
	return &env
}

func NewEnvNameValue(name string, value *Ast) *Env {
	env := NewEnv()
	env.map_[name] = value
	return env
}

func NewEnvFromParams(params []string, arg_list []*Ast) *Env {
	env := NewEnv()
	for i, p := range params {
		env.map_[p] = arg_list[i]
	}
	return env
}

func extend(old *Env, new_ *Env) *Env {
	e := NewEnv()
	for k, v := range new_.map_ {
		e.map_[k] = v
	}
	e.parent = old
	return e
}

func (e *Env) get(name string) *Ast {
	v, ok := e.map_[name]
	if ok {
		return v
	}
	if e.parent != nil {
		return e.parent.get(name)
	}
	return nil
}
func (e *Env) add(name string, value *Ast) {
	if e.map_ == nil {
		e.map_ = make(map[string]*Ast)
	}
	e.map_[name] = value
}
func (e *Env) String() string {
	c := make([]string, 0)
	for k, v := range e.map_ {
		c = append(c, k+": "+v.SimpleString())
	}
	p := ""
	if e.parent != nil {
		p = e.parent.String()
	}
	return "{" + strings.Join(c, ", ") + p + "}"
}

// ==== Env end ====

// ==== helper function ====
func MapAst(lst []*Ast, f func(*Ast) string) []string {
	ret := make([]string, 0, len(lst))
	for _, v := range lst {
		ret = append(ret, f(v))
	}
	return ret
}
func IsDigit(b byte) bool {
	return '0' <= b && b <= '9'
}
