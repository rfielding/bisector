package main

import "fmt"

type Expr interface {
	isExpr()
}

/* All expressions can be qualified with a restriction */
type RExpr struct {
	expr Expr
	env  Expr
}

/* All numbers are represented rationally */
type R struct {
	num int
	den int
}

func (r R) isExpr() {}

func gcd(a int, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

func reduce(a R) (r R) {
	if a.num == 0 {
		if a.den != 0 {
			return R{0, 1}
		} else {
			//Actually... this is quite absurd
			return R{0, 0}
		}

	}
	if a.den == 0 {
		//At least normalize this absurd result
		if a.num < 0 {
			return R{-1, 0}
		} else {
			return R{1, 0}
		}
	}
	if a.den < 0 {
		a.num = -a.num
		a.den = -a.den
	}
	g := gcd(a.num, a.den)
	return R{a.num / g, a.den / g}
}

func rAdd(a R, b R) R {
	return reduce(R{a.num*b.den + b.num*a.den, a.den * b.den})
}

func rSub(a R, b R) R {
	return reduce(R{a.num*b.den - b.num*a.den, a.den * b.den})
}

func rMul(a R, b R) R {
	return reduce(R{a.num * b.num, a.den * b.den})
}

func rDiv(a R, b R) R {
	return reduce(R{a.num * b.den, a.den * b.num})
}

type B bool

func (b B) isExpr() {}

type O string

type V string

func (v V) isExpr() {}

type Binop struct {
	left  Expr
	op    O
	right Expr
}

func (b Binop) isExpr() {}

func bop(a Expr, op O, b Expr) (r Binop) {
	return Binop{a, op, b}
}

func XAnd(a RExpr, b RExpr) (r RExpr) {
	var oldEnv Expr = B(true)
	if a.env != B(true) || b.env != B(true) {
		oldEnv = bop(a.env, "and", b.env)
	}
	newExpr := bop(a.expr, "and", b.expr)

	return RExpr{newExpr, oldEnv}
}

func XOr(a RExpr, b RExpr) (r RExpr) {
	var oldEnv Expr = B(true)
	if a.env != B(true) || b.env != B(true) {
		oldEnv = bop(a.env, "and", b.env)
	}
	newExpr := bop(a.expr, "or", b.expr)

	return RExpr{newExpr, oldEnv}
}

//This is the and function intended for use with
//restriction expressions, to compound restrictions
func rAnd(a Expr, b Expr) Expr {
	if a == B(true) {
		return b
	}
	if b == B(true) {
		return a
	}
	return bop(a, "and", b)
}

//This is one of the operations that forces a condition to
//be added!
func XDiv(a RExpr, b RExpr) (r RExpr) {
	//When we do a division, calculate restrictions
	//if we can't rule them out
	newRestriction := bop(bop(b.expr, "eq", R{0, 1}), "eq", B(false))
	newEnv := rAnd(rAnd(a.env, b.env), newRestriction)
	newExpr := bop(a.expr, "div", b.expr)
	return RExpr{newExpr, newEnv}
}

func XAdd(a RExpr, b RExpr) (r RExpr) {
	oldEnv := rAnd(a.env, b.env)
	newExpr := bop(a.expr, "add", b.expr)
	return RExpr{newExpr, oldEnv}
}

func XMul(a RExpr, b RExpr) (r RExpr) {
	oldEnv := rAnd(a.env, b.env)
	newExpr := bop(a.expr, "mul", b.expr)
	return RExpr{newExpr, oldEnv}
}

func XSub(a RExpr, b RExpr) (r RExpr) {
	oldEnv := rAnd(a.env, b.env)
	newExpr := bop(a.expr, "sub", b.expr)
	return RExpr{newExpr, oldEnv}
}

func commute(b Binop) (r Binop) {
	return bop(b.right, b.op, b.left)
}

func ldistribute(b Binop) (r Binop) {
	r = b
	switch bi := b.right.(type) {
	case Binop:
		x := bop(b.left, b.op, bi.left)
		y := bop(b.left, b.op, bi.right)
		r = bop(x, bi.op, y)
	}
	return
}

func rdistribute(b Binop) (r Binop) {
	r = b
	switch bi := b.left.(type) {
	case Binop:
		x := bop(bi.left, b.op, b.right)
		y := bop(bi.right, b.op, b.right)
		r = bop(x, bi.op, y)
	}
	return
}

func lunDistribute(b Binop) (r Binop) {
	r = b
	switch bl := b.left.(type) {
	case Binop:
		switch br := b.right.(type) {
		case Binop:
			if bl.left == br.left {
				if bl.op == br.op {
					z := bop(bl.right, b.op, br.right)
					r = bop(bl.left, bl.op, z)
				}
			}
		}
	}
	return
}

func runDistribute(b Binop) (r Binop) {
	r = b
	switch br := b.right.(type) {
	case Binop:
		switch bl := b.left.(type) {
		case Binop:
			if bl.right == br.right {
				if br.op == bl.op {
					z := bop(bl.left, b.op, br.left)
					r = bop(z, br.op, br.right)
				}
			}
		}
	}
	return
}

func XCommute(a RExpr) (r RExpr) {
	oldEnv := a.env
	r = a
	switch e := a.expr.(type) {
	case Binop:
		e2 := commute(e)
		r = RExpr{e2, oldEnv}
	}
	return
}

//TODO: ldistribute<*,+> (a * b) where b is not a binop:
//    must create side condition and temp variables:
//     (a * b) ->  ((a * b_1) + (a * b_2)),b = b_1+b_2
func XLDistribute(a RExpr) (r RExpr) {
	oldEnv := a.env
	r = a
	switch e := a.expr.(type) {
	case Binop:
		e2 := ldistribute(e)
		r = RExpr{e2, oldEnv}
	}
	return
}

func XRDistribute(a RExpr) (r RExpr) {
	oldEnv := a.env
	r = a
	switch e := a.expr.(type) {
	case Binop:
		e2 := rdistribute(e)
		r = RExpr{e2, oldEnv}
	}
	return
}

func XLunDistribute(a RExpr) (r RExpr) {
	oldEnv := a.env
	r = a
	switch e := a.expr.(type) {
	case Binop:
		e2 := lunDistribute(e)
		r = RExpr{e2, oldEnv}
	}
	return
}

func XRunDistribute(a RExpr) (r RExpr) {
	oldEnv := a.env
	r = a
	switch e := a.expr.(type) {
	case Binop:
		e2 := runDistribute(e)
		r = RExpr{e2, oldEnv}
	}
	return
}

func XRational(a R) RExpr {
	return RExpr{a, B(true)}
}

func XVar(v V) RExpr {
	return RExpr{v, B(true)}
}

func bRationalEval(b Binop) (r Expr) {
	xe := RationalEval(b.left)
	ye := RationalEval(b.right)
	r = bop(xe, b.op, ye)
	switch x := xe.(type) {
	case R:
		switch y := ye.(type) {
		case R:
			//TODO: when a "/" binop is involved, do
			// rational evaluation
			switch b.op {
			case "add":
				r = rAdd(x, y)
			case "sub":
				r = rSub(x, y)
			case "mul":
				r = rMul(x, y)
			case "div":
				r = rDiv(x, y)
			}
		}
	}
	return
}

func RationalEval(expr Expr) (r Expr) {
	r = expr
	switch b := expr.(type) {
	case Binop:
		r = bRationalEval(b)
	}
	return
}

func XRationalEval(a RExpr) (r RExpr) {
	oldEnv := a.env
	r = a
	switch e := a.expr.(type) {
	case Binop:
		e2 := RationalEval(e)
		r = RExpr{e2, oldEnv}
	}
	return
}

func main() {
	//Rational arithmetic
	r0 := XRational(R{5, 3})
	r1 := XRational(R{7, 11})
	r2 := XRational(R{4, 5})
	r := XMul(r0, XAdd(r1, r2))
	fmt.Println(r)
	x := XVar(V("x"))
	y := XVar(V("y"))
	z := XVar(V("z"))
	e0 := XDiv(r, x)
	e1 := XMul(e0, y)
	fmt.Println(e1)
	//Algebra
	e2 := XCommute(e1)
	fmt.Println(e2)
	e3 := XDiv(e1, z)
	fmt.Println(e3)
	e4 := XRationalEval(e3)
	fmt.Println(e4)
}
