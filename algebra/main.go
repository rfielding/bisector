package main

import "fmt"

type Expr interface {
	isExpr()
}

type O string

type V string

func (v V) isExpr() {}

type I int

func (i I) isExpr() {}

type B bool

func (b B) isExpr() {}

type Binop struct {
	left  Expr
	op    O
	right Expr
}

func (b Binop) isExpr() {}

func (b Binop) Commute() (r Expr) {
	return Binop{b.right, b.op, b.left}
}

func (b Binop) LDistribute() (r Expr) {
	r = b
	switch bi := b.right.(type) {
	case Binop:
		x := Binop{b.left, b.op, bi.left}
		y := Binop{b.left, b.op, bi.right}
		r = Binop{x, bi.op, y}
	}
	return
}

func (b Binop) RDistribute() (r Expr) {
	r = b
	switch bi := b.left.(type) {
	case Binop:
		x := Binop{bi.left, b.op, b.right}
		y := Binop{bi.right, b.op, b.right}
		r = Binop{x, bi.op, y}
	}
	return
}

func (b Binop) LUnDistribute() (r Expr) {
	r = b
	switch bl := b.left.(type) {
	case Binop:
		switch br := b.right.(type) {
		case Binop:
			if bl.left == br.left {
				if bl.op == br.op {
					z := Binop{bl.right, b.op, br.right}
					r = Binop{bl.left, bl.op, z}
				}
			}
		}
	}
	return
}

func (b Binop) RUnDistribute() (r Expr) {
	r = b
	switch br := b.right.(type) {
	case Binop:
		switch bl := b.left.(type) {
		case Binop:
			if bl.right == br.right {
				if br.op == bl.op {
					z := Binop{bl.left, b.op, br.left}
					r = Binop{z, br.op, br.right}
				}
			}
		}
	}
	return
}

func main() {
	r := Binop{I(5), "/", I(4)}
	e := Binop{V("z"), "+", r}
	e2 := e.Commute()
	e3 := Binop{e2, "*", I(6)}
	e4 := e3.RDistribute()
	e5 := Binop{V("x"), "=", e4}
	fmt.Println(e5)
}
