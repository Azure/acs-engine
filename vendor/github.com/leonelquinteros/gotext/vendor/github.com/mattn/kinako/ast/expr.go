package ast

type Token struct {
	Tok int
	Lit string
}

// Position provides interface to store code locations.
type Position struct {
	Line   int
	Column int
}

// Expr provides all of interfaces for expression.
type Expr interface {
	expr()
}

// ExprImpl provide commonly implementations for Expr.
type ExprImpl struct {
}

// expr provide restraint interface.
func (x *ExprImpl) expr() {}

// NumberExpr provide Number expression.
type NumberExpr struct {
	ExprImpl
	Lit string
}

// UnaryExpr provide unary minus expression. ex: -1, ^1, ~1.
type UnaryExpr struct {
	ExprImpl
	Operator string
	Expr     Expr
}

// IdentExpr provide identity expression.
type IdentExpr struct {
	ExprImpl
	Lit string
}

// Stmt provides all of interfaces for statement.
type Stmt interface {
	stmt()
}

// StmtImpl provide commonly implementations for Stmt..
type StmtImpl struct {
}

// stmt provide restraint interface.
func (x *StmtImpl) stmt() {}

// LetsStmt provide multiple statement of let.
type LetsStmt struct {
	StmtImpl
	Lhss     []Expr
	Operator string
	Rhss     []Expr
}

// StringExpr provide String expression.
type StringExpr struct {
	ExprImpl
	Lit string
}

type TernaryOpExpr struct {
	ExprImpl
	Expr Expr
	Lhs  Expr
	Rhs  Expr
}

// CallExpr provide calling expression.
type CallExpr struct {
	ExprImpl
	Func     interface{}
	Name     string
	SubExprs []Expr
}

// ParenExpr provide parent block expression.
type ParenExpr struct {
	ExprImpl
	SubExpr Expr
}

// BinOpExpr provide binary operator expression.
type BinOpExpr struct {
	ExprImpl
	Lhs      Expr
	Operator string
	Rhs      Expr
}

// ExprStmt provide expression statement.
type ExprStmt struct {
	StmtImpl
	Expr Expr
}

// LetStmt provide statement of let.
type LetStmt struct {
	StmtImpl
	Lhs      Expr
	Operator string
	Rhs      Expr
}
