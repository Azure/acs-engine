//line parser.go.y:2
package parser

import __yyfmt__ "fmt"

//line parser.go.y:2
import (
	"github.com/mattn/kinako/ast"
)

//line parser.go.y:16
type yySymType struct {
	yys       int
	compstmt  []ast.Stmt
	stmts     []ast.Stmt
	stmt      ast.Stmt
	expr      ast.Expr
	exprs     []ast.Expr
	tok       ast.Token
	term      ast.Token
	terms     ast.Token
	opt_terms ast.Token
}

const IDENT = 57346
const NUMBER = 57347
const STRING = 57348
const EQEQ = 57349
const NEQ = 57350
const GE = 57351
const LE = 57352
const OROR = 57353
const ANDAND = 57354
const POW = 57355
const SHIFTLEFT = 57356
const SHIFTRIGHT = 57357
const PLUSPLUS = 57358
const MINUSMINUS = 57359
const UNARY = 57360

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"IDENT",
	"NUMBER",
	"STRING",
	"EQEQ",
	"NEQ",
	"GE",
	"LE",
	"OROR",
	"ANDAND",
	"POW",
	"'='",
	"'?'",
	"':'",
	"','",
	"'>'",
	"'<'",
	"SHIFTLEFT",
	"SHIFTRIGHT",
	"'+'",
	"'-'",
	"PLUSPLUS",
	"MINUSMINUS",
	"'*'",
	"'/'",
	"'%'",
	"UNARY",
	"'!'",
	"'^'",
	"'('",
	"')'",
	"'|'",
	"'&'",
	"';'",
	"'\\n'",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line parser.go.y:213

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
	-1, 50,
	7, 0,
	8, 0,
	-2, 20,
	-1, 51,
	7, 0,
	8, 0,
	-2, 21,
}

const yyNprod = 40
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 251

var yyAct = [...]int{

	9, 6, 7, 33, 35, 37, 22, 23, 60, 3,
	24, 25, 26, 38, 39, 40, 1, 41, 33, 35,
	8, 43, 44, 45, 46, 47, 48, 49, 50, 51,
	52, 53, 54, 55, 56, 57, 58, 59, 61, 42,
	27, 28, 30, 32, 34, 36, 65, 0, 21, 63,
	4, 29, 31, 2, 18, 22, 23, 17, 0, 24,
	25, 26, 64, 0, 66, 0, 67, 33, 35, 27,
	28, 30, 32, 34, 36, 0, 0, 21, 0, 0,
	29, 31, 0, 0, 22, 23, 0, 0, 24, 25,
	26, 0, 0, 0, 0, 62, 33, 35, 27, 28,
	30, 32, 34, 36, 0, 20, 21, 0, 0, 29,
	31, 0, 0, 22, 23, 5, 0, 24, 25, 26,
	19, 0, 0, 0, 0, 33, 35, 27, 28, 30,
	32, 34, 36, 0, 19, 21, 0, 0, 29, 31,
	0, 0, 22, 23, 0, 0, 24, 25, 26, 0,
	0, 0, 0, 0, 33, 35, 27, 28, 30, 32,
	0, 36, 0, 0, 0, 0, 0, 29, 31, 0,
	0, 22, 23, 0, 0, 24, 25, 26, 27, 28,
	30, 32, 0, 33, 35, 0, 0, 0, 0, 29,
	31, 0, 0, 22, 23, 0, 0, 24, 25, 26,
	30, 32, 10, 11, 15, 33, 35, 0, 0, 29,
	31, 0, 0, 22, 23, 0, 0, 24, 25, 26,
	0, 12, 10, 11, 15, 33, 35, 0, 13, 14,
	16, 24, 25, 26, 6, 7, 0, 0, 0, 33,
	35, 12, 0, 0, 0, 0, 0, 0, 13, 14,
	16,
}
var yyPact = [...]int{

	-35, -1000, 218, -35, -35, -1000, -1000, -1000, -1000, 91,
	-27, -1000, 218, 218, 218, -1000, 218, -1000, 198, -1000,
	218, 218, 218, 218, 218, 218, 218, 218, 218, 218,
	218, 218, 218, 218, 218, 218, 218, 218, -31, -31,
	-31, 62, -1000, 120, 33, 205, 205, -31, -31, -31,
	191, 191, -16, -16, -16, -16, 120, 149, 120, 171,
	29, 120, -1000, 218, -1000, 218, 120, 120,
}
var yyPgo = [...]int{

	0, 16, 9, 20, 0, 8, 53, 50, 115,
}
var yyR1 = [...]int{

	0, 1, 1, 2, 2, 3, 3, 4, 4, 4,
	4, 4, 4, 4, 4, 4, 4, 4, 4, 4,
	4, 4, 4, 4, 4, 4, 4, 4, 4, 4,
	4, 5, 5, 5, 6, 6, 7, 7, 8, 8,
}
var yyR2 = [...]int{

	0, 1, 2, 2, 3, 3, 1, 1, 1, 2,
	2, 2, 1, 5, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	4, 0, 1, 3, 0, 1, 1, 2, 1, 1,
}
var yyChk = [...]int{

	-1000, -1, -6, -2, -7, -8, 36, 37, -3, -4,
	4, 5, 23, 30, 31, 6, 32, -6, -7, -8,
	14, 15, 22, 23, 26, 27, 28, 7, 8, 18,
	9, 19, 10, 34, 11, 35, 12, 32, -4, -4,
	-4, -4, -3, -4, -4, -4, -4, -4, -4, -4,
	-4, -4, -4, -4, -4, -4, -4, -4, -4, -4,
	-5, -4, 33, 16, 33, 17, -4, -4,
}
var yyDef = [...]int{

	34, -2, 1, 34, 35, 36, 38, 39, 3, 6,
	7, 8, 0, 0, 0, 12, 0, 2, 35, 37,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 31, 9, 10,
	11, 0, 4, 5, 0, 15, 16, 17, 18, 19,
	-2, -2, 22, 23, 24, 25, 26, 27, 28, 29,
	0, 32, 14, 0, 30, 0, 13, 33,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	37, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 30, 3, 3, 3, 28, 35, 3,
	32, 33, 26, 22, 17, 23, 3, 27, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 16, 36,
	19, 14, 18, 15, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 31, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 34,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 20, 21, 24, 25, 29,
}
var yyTok3 = [...]int{
	0,
}

var yyErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	yyDebug        = 0
	yyErrorVerbose = false
)

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

type yyParser interface {
	Parse(yyLexer) int
	Lookahead() int
}

type yyParserImpl struct {
	lval  yySymType
	stack [yyInitialStackSize]yySymType
	char  int
}

func (p *yyParserImpl) Lookahead() int {
	return p.char
}

func yyNewParser() yyParser {
	return &yyParserImpl{}
}

const yyFlag = -1000

func yyTokname(c int) string {
	if c >= 1 && c-1 < len(yyToknames) {
		if yyToknames[c-1] != "" {
			return yyToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func yyErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !yyErrorVerbose {
		return "syntax error"
	}

	for _, e := range yyErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + yyTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := yyPact[state]
	for tok := TOKSTART; tok-1 < len(yyToknames); tok++ {
		if n := base + tok; n >= 0 && n < yyLast && yyChk[yyAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if yyDef[state] == -2 {
		i := 0
		for yyExca[i] != -1 || yyExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; yyExca[i] >= 0; i += 2 {
			tok := yyExca[i]
			if tok < TOKSTART || yyExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if yyExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += yyTokname(tok)
	}
	return res
}

func yylex1(lex yyLexer, lval *yySymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		token = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			token = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		token = yyTok3[i+0]
		if token == char {
			token = yyTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(token), uint(char))
	}
	return char, token
}

func yyParse(yylex yyLexer) int {
	return yyNewParser().Parse(yylex)
}

func (yyrcvr *yyParserImpl) Parse(yylex yyLexer) int {
	var yyn int
	var yyVAL yySymType
	var yyDollar []yySymType
	_ = yyDollar // silence set and not used
	yyS := yyrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yyrcvr.char = -1
	yytoken := -1 // yyrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		yystate = -1
		yyrcvr.char = -1
		yytoken = -1
	}()
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yytoken), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yyrcvr.char < 0 {
		yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
	}
	yyn += yytoken
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yytoken { /* valid shift */
		yyrcvr.char = -1
		yytoken = -1
		yyVAL = yyrcvr.lval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yyrcvr.char < 0 {
			yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yytoken {
				break
			}
		}
		yyn = yyExca[xi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error(yyErrorMessage(yystate, yytoken))
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yytoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yytoken))
			}
			if yytoken == yyEofCode {
				goto ret1
			}
			yyrcvr.char = -1
			yytoken = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	// yyp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if yyp+1 >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 1:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:45
		{
			yyVAL.compstmt = nil
		}
	case 2:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:49
		{
			yyVAL.compstmt = yyDollar[1].stmts
		}
	case 3:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:55
		{
			yyVAL.stmts = []ast.Stmt{yyDollar[2].stmt}
			if l, ok := yylex.(*Lexer); ok {
				l.stmts = yyVAL.stmts
			}
		}
	case 4:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:62
		{
			if yyDollar[3].stmt != nil {
				yyVAL.stmts = append(yyDollar[1].stmts, yyDollar[3].stmt)
				if l, ok := yylex.(*Lexer); ok {
					l.stmts = yyVAL.stmts
				}
			}
		}
	case 5:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:73
		{
			yyVAL.stmt = &ast.LetStmt{Lhs: yyDollar[1].expr, Operator: "=", Rhs: yyDollar[3].expr}
		}
	case 6:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:77
		{
			yyVAL.stmt = &ast.ExprStmt{Expr: yyDollar[1].expr}
		}
	case 7:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:83
		{
			yyVAL.expr = &ast.IdentExpr{Lit: yyDollar[1].tok.Lit}
		}
	case 8:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:87
		{
			yyVAL.expr = &ast.NumberExpr{Lit: yyDollar[1].tok.Lit}
		}
	case 9:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:91
		{
			yyVAL.expr = &ast.UnaryExpr{Operator: "-", Expr: yyDollar[2].expr}
		}
	case 10:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:95
		{
			yyVAL.expr = &ast.UnaryExpr{Operator: "!", Expr: yyDollar[2].expr}
		}
	case 11:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:99
		{
			yyVAL.expr = &ast.UnaryExpr{Operator: "^", Expr: yyDollar[2].expr}
		}
	case 12:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:103
		{
			yyVAL.expr = &ast.StringExpr{Lit: yyDollar[1].tok.Lit}
		}
	case 13:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:107
		{
			yyVAL.expr = &ast.TernaryOpExpr{Expr: yyDollar[1].expr, Lhs: yyDollar[3].expr, Rhs: yyDollar[5].expr}
		}
	case 14:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:111
		{
			yyVAL.expr = &ast.ParenExpr{SubExpr: yyDollar[2].expr}
		}
	case 15:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:115
		{
			yyVAL.expr = &ast.BinOpExpr{Lhs: yyDollar[1].expr, Operator: "+", Rhs: yyDollar[3].expr}
		}
	case 16:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:119
		{
			yyVAL.expr = &ast.BinOpExpr{Lhs: yyDollar[1].expr, Operator: "-", Rhs: yyDollar[3].expr}
		}
	case 17:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:123
		{
			yyVAL.expr = &ast.BinOpExpr{Lhs: yyDollar[1].expr, Operator: "*", Rhs: yyDollar[3].expr}
		}
	case 18:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:127
		{
			yyVAL.expr = &ast.BinOpExpr{Lhs: yyDollar[1].expr, Operator: "/", Rhs: yyDollar[3].expr}
		}
	case 19:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:131
		{
			yyVAL.expr = &ast.BinOpExpr{Lhs: yyDollar[1].expr, Operator: "%", Rhs: yyDollar[3].expr}
		}
	case 20:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:135
		{
			yyVAL.expr = &ast.BinOpExpr{Lhs: yyDollar[1].expr, Operator: "==", Rhs: yyDollar[3].expr}
		}
	case 21:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:139
		{
			yyVAL.expr = &ast.BinOpExpr{Lhs: yyDollar[1].expr, Operator: "!=", Rhs: yyDollar[3].expr}
		}
	case 22:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:143
		{
			yyVAL.expr = &ast.BinOpExpr{Lhs: yyDollar[1].expr, Operator: ">", Rhs: yyDollar[3].expr}
		}
	case 23:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:147
		{
			yyVAL.expr = &ast.BinOpExpr{Lhs: yyDollar[1].expr, Operator: ">=", Rhs: yyDollar[3].expr}
		}
	case 24:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:151
		{
			yyVAL.expr = &ast.BinOpExpr{Lhs: yyDollar[1].expr, Operator: "<", Rhs: yyDollar[3].expr}
		}
	case 25:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:155
		{
			yyVAL.expr = &ast.BinOpExpr{Lhs: yyDollar[1].expr, Operator: "<=", Rhs: yyDollar[3].expr}
		}
	case 26:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:159
		{
			yyVAL.expr = &ast.BinOpExpr{Lhs: yyDollar[1].expr, Operator: "|", Rhs: yyDollar[3].expr}
		}
	case 27:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:163
		{
			yyVAL.expr = &ast.BinOpExpr{Lhs: yyDollar[1].expr, Operator: "||", Rhs: yyDollar[3].expr}
		}
	case 28:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:167
		{
			yyVAL.expr = &ast.BinOpExpr{Lhs: yyDollar[1].expr, Operator: "&", Rhs: yyDollar[3].expr}
		}
	case 29:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:171
		{
			yyVAL.expr = &ast.BinOpExpr{Lhs: yyDollar[1].expr, Operator: "&&", Rhs: yyDollar[3].expr}
		}
	case 30:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:175
		{
			yyVAL.expr = &ast.CallExpr{Name: yyDollar[1].tok.Lit, SubExprs: yyDollar[3].exprs}
		}
	case 31:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.go.y:180
		{
			yyVAL.exprs = nil
		}
	case 32:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:184
		{
			yyVAL.exprs = []ast.Expr{yyDollar[1].expr}
		}
	case 33:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:188
		{
			yyVAL.exprs = append(yyDollar[1].exprs, yyDollar[3].expr)
		}
	case 36:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:198
		{
		}
	case 37:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:201
		{
		}
	case 38:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:206
		{
		}
	case 39:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:209
		{
		}
	}
	goto yystack /* stack new state and value */
}
