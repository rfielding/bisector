package main

import "fmt"

type Expr interface {
	isExpr()
}

type RExpr struct {
	expr       Expr
	constraint Expr
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

func Bop(a Expr, op O, b Expr) (r Binop) {
	return Binop{a, op, b}
}

func Commute(b Binop) (r Binop) {
	return Bop(b.right, b.op, b.left)
}

func LDistribute(b Binop) (r Binop) {
	r = b
	switch bi := b.right.(type) {
	case Binop:
		x := Bop(b.left, b.op, bi.left)
		y := Bop(b.left, b.op, bi.right)
		r = Bop(x, bi.op, y)
	}
	return
}

func RDistribute(b Binop) (r Binop) {
	r = b
	switch bi := b.left.(type) {
	case Binop:
		x := Bop(bi.left, b.op, b.right)
		y := Bop(bi.right, b.op, b.right)
		r = Bop(x, bi.op, y)
	}
	return
}

func LUnDistribute(b Binop) (r Binop) {
	r = b
	switch bl := b.left.(type) {
	case Binop:
		switch br := b.right.(type) {
		case Binop:
			if bl.left == br.left {
				if bl.op == br.op {
					z := Bop(bl.right, b.op, br.right)
					r = Bop(bl.left, bl.op, z)
				}
			}
		}
	}
	return
}

func RUnDistribute(b Binop) (r Binop) {
	r = b
	switch br := b.right.(type) {
	case Binop:
		switch bl := b.left.(type) {
		case Binop:
			if bl.right == br.right {
				if br.op == bl.op {
					z := Bop(bl.left, b.op, br.left)
					r = Bop(z, br.op, br.right)
				}
			}
		}
	}
	return
}

func And(a Expr, b Expr) (r Binop) {
	return Bop(a, "and", b)
}

func Or(a Expr, b Expr) (r Binop) {
	return Bop(a, "or", b)
}

func Div(a Expr, b Expr) (r Binop) {
	//constraintNew := Binop{b, "!=", I(0)}
	return Bop(a, "/", b)
}

func Add(a Expr, b Expr) (r Binop) {
	return Bop(a, "+", b)
}

func Mul(a Expr, b Expr) (r Binop) {
	return Bop(a, "*", b)
}

func Eq(a Expr, b Expr) (r Binop) {
	return Bop(a, "=", b)
}

func main() {
	r := Div(I(5), I(4))
	e := Add(V("z"), r)
	e2 := Commute(e)
	e3 := Mul(e2, I(6))
	e4 := RDistribute(e3)
	e5 := Eq(V("x"), e4)
	fmt.Println(e5)
}
