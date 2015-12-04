//line query.y:2
package query

import __yyfmt__ "fmt"

//line query.y:2
import (
	"bufio"
	"fmt"
	"github.com/taylorchu/toki"
	"strconv"
	_time "time"
)

//line query.y:13
type QuerySymType struct {
	yys            int
	str            string
	selectTerm     SelectTerm
	selectTermList []SelectTerm
	whereTerm      WhereTerm
	whereClause    WhereClause
	time           _time.Time
	timediff       _time.Duration
}

const SELECT = 57346
const DISTINCT = 57347
const WHERE = 57348
const LVALUE = 57349
const QSTRING = 57350
const LIKE = 57351
const HAS = 57352
const NOW = 57353
const SET = 57354
const IBEFORE = 57355
const BEFORE = 57356
const IAFTER = 57357
const AFTER = 57358
const AND = 57359
const AS = 57360
const TO = 57361
const OR = 57362
const IN = 57363
const NOT = 57364
const FOR = 57365
const LPAREN = 57366
const RPAREN = 57367
const NEWLINE = 57368
const NUMBER = 57369
const SEMICOLON = 57370
const EQ = 57371
const NEQ = 57372
const COMMA = 57373
const ALL = 57374

var QueryToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"SELECT",
	"DISTINCT",
	"WHERE",
	"LVALUE",
	"QSTRING",
	"LIKE",
	"HAS",
	"NOW",
	"SET",
	"IBEFORE",
	"BEFORE",
	"IAFTER",
	"AFTER",
	"AND",
	"AS",
	"TO",
	"OR",
	"IN",
	"NOT",
	"FOR",
	"LPAREN",
	"RPAREN",
	"NEWLINE",
	"NUMBER",
	"SEMICOLON",
	"EQ",
	"NEQ",
	"COMMA",
	"ALL",
}
var QueryStatenames = [...]string{}

const QueryEofCode = 1
const QueryErrCode = 2
const QueryMaxDepth = 200

//line query.y:288
const eof = 0

var supported_formats = []string{"1/2/2006",
	"1-2-2006",
	"1/2/2006 03:04:05 PM MST",
	"1-2-2006 03:04:05 PM MST",
	"1/2/2006 15:04:05 MST",
	"1-2-2006 15:04:05 MST",
	"2006-1-2 15:04:05 MST"}

type QueryLex struct {
	Query       *Query
	querystring string
	scanner     *toki.Scanner
	lasttoken   string
	tokens      []string
	innertable  int
	Err         error
}

func (ql *QueryLex) NextLetter() string {
	var alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	ql.innertable += 1
	return string(alphabet[ql.innertable-1])
}

var termTemplate = `
	(
		select distinct data.uuid
		from data
		inner join
		(
			select distinct uuid, dkey, max(timestamp) as maxtime from data
			group by dkey, uuid order by timestamp desc
		) sorted
		on data.uuid = sorted.uuid and data.dkey = sorted.dkey and data.timestamp = sorted.maxtime
		where data.dval is not null
			%s
	) as %s
`

var termTemplateUnion = `
	(
		select distinct data.uuid
		from data
		inner join
		(
			select distinct uuid, dkey, max(timestamp) as maxtime from data
			group by dkey, uuid order by timestamp desc
		) sorted
		on data.uuid = sorted.uuid and data.dkey = sorted.dkey and data.timestamp = sorted.maxtime
		where data.dval is not null
			%s
	)
`

var timePredicateSingle = `
	select distinct uuid, dkey, max(timestamp) as maxtime from data
	where timestamp %s "%s"
	group by dkey, uuid order by timestamp desc
`

var timePredicateRange = `
	select distinct uuid, dkey, max(timestamp) as maxtime from data
	where timestamp %s "%s" and timestamp %s "%s"
	group by dkey, uuid order by timestamp desc
`

func NewQueryLexer(s string) *QueryLex {
	scanner := toki.NewScanner(
		[]toki.Def{
			{Token: WHERE, Pattern: "where"},
			{Token: SELECT, Pattern: "select"},
			{Token: DISTINCT, Pattern: "distinct"},
			{Token: ALL, Pattern: "\\*"},
			{Token: NOW, Pattern: "now"},
			{Token: SET, Pattern: "set"},
			{Token: BEFORE, Pattern: "before"},
			{Token: IBEFORE, Pattern: "ibefore"},
			{Token: AFTER, Pattern: "after"},
			{Token: IAFTER, Pattern: "iafter"},
			{Token: COMMA, Pattern: ","},
			{Token: AND, Pattern: "and"},
			{Token: AS, Pattern: "as"},
			{Token: TO, Pattern: "to"},
			{Token: FOR, Pattern: "for"},
			{Token: OR, Pattern: "or"},
			{Token: IN, Pattern: "in"},
			{Token: HAS, Pattern: "has"},
			{Token: NOT, Pattern: "not"},
			{Token: NEQ, Pattern: "!="},
			{Token: EQ, Pattern: "="},
			{Token: LPAREN, Pattern: "\\("},
			{Token: RPAREN, Pattern: "\\)"},
			{Token: SEMICOLON, Pattern: ";"},
			{Token: NEWLINE, Pattern: "\n"},
			{Token: LIKE, Pattern: "(like)|~"},
			{Token: NUMBER, Pattern: "([+-]?([0-9]*\\.)?[0-9]+)"},
			{Token: LVALUE, Pattern: "[a-zA-Z\\~\\$\\_][a-zA-Z0-9\\/\\%_\\-]*"},
			{Token: QSTRING, Pattern: "(\"[^\"\\\\]*?(\\.[^\"\\\\]*?)*?\")|('[^'\\\\]*?(\\.[^'\\\\]*?)*?')"},
		})
	scanner.SetInput(s)
	return &QueryLex{Query: &Query{}, querystring: s, scanner: scanner, Err: nil, lasttoken: "", tokens: []string{}}
}

func (lex *QueryLex) Lex(lval *QuerySymType) int {
	r := lex.scanner.Next()
	lex.lasttoken = r.String()
	if r.Pos.Line == 2 || len(r.Value) == 0 {
		return eof
	}
	lval.str = string(r.Value)
	lex.tokens = append(lex.tokens, lval.str)
	return int(r.Token)
}

func (lex *QueryLex) Error(s string) {
	lex.Err = fmt.Errorf(s)
}

func readline(fi *bufio.Reader) (string, bool) {
	fmt.Printf("aronnax> ")
	s, err := fi.ReadString('\n')
	if err != nil {
		return "", false
	}
	return s, true
}

//line yacctab:1
var QueryExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const QueryNprod = 33
const QueryPrivate = 57344

var QueryTokenNames []string
var QueryStates []string

const QueryLast = 57

var QueryAct = [...]int{

	52, 6, 30, 7, 11, 39, 42, 9, 20, 43,
	53, 51, 26, 25, 28, 27, 23, 50, 38, 22,
	24, 16, 31, 32, 17, 41, 13, 4, 8, 10,
	49, 48, 44, 45, 46, 47, 15, 55, 18, 19,
	54, 33, 29, 12, 2, 34, 37, 1, 21, 35,
	36, 40, 14, 5, 3, 0, 56,
}
var QueryPact = [...]int{

	40, -1000, -4, 1, -1000, -27, 36, -1000, -1000, 14,
	-1000, -4, -1000, -20, -1, 14, -7, 34, 14, -1000,
	-1000, -1000, 14, 14, -6, -2, -2, -2, -2, -1000,
	27, 23, 22, -1000, -8, -1000, -1000, -1000, -14, -1000,
	-17, 33, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, 30, -1000, -17, -1000,
}
var QueryPgo = [...]int{

	0, 27, 54, 53, 26, 52, 5, 51, 0, 48,
	47, 46,
}
var QueryR1 = [...]int{

	0, 10, 10, 2, 1, 1, 1, 3, 3, 4,
	4, 4, 4, 4, 5, 5, 5, 5, 5, 9,
	9, 9, 9, 9, 11, 6, 6, 7, 7, 7,
	7, 8, 8,
}
var QueryR2 = [...]int{

	0, 5, 3, 1, 1, 3, 2, 1, 1, 1,
	2, 3, 3, 2, 3, 3, 3, 2, 3, 2,
	2, 2, 2, 2, 2, 1, 2, 2, 1, 1,
	1, 2, 3,
}
var QueryChk = [...]int{

	-1000, -10, 4, -2, -1, -3, 5, 7, 32, 6,
	28, 31, 7, -4, -5, 22, 7, 10, 24, -1,
	28, -9, 20, 17, 21, 14, 13, 16, 15, -4,
	9, 29, 30, 7, -4, -4, -4, -11, 24, -6,
	-7, 27, 8, 11, -6, -6, -6, 8, 8, 8,
	25, 25, -8, 27, 7, 7, -8,
}
var QueryDef = [...]int{

	0, -2, 0, 0, 3, 4, 0, 7, 8, 0,
	2, 0, 6, 0, 9, 0, 0, 0, 0, 5,
	1, 10, 0, 0, 0, 0, 0, 0, 0, 13,
	0, 0, 0, 17, 0, 11, 12, 19, 0, 20,
	25, 28, 29, 30, 21, 22, 23, 14, 15, 16,
	18, 24, 26, 0, 27, 31, 32,
}
var QueryTok1 = [...]int{

	1,
}
var QueryTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32,
}
var QueryTok3 = [...]int{
	0,
}

var QueryErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	QueryDebug        = 0
	QueryErrorVerbose = false
)

type QueryLexer interface {
	Lex(lval *QuerySymType) int
	Error(s string)
}

type QueryParser interface {
	Parse(QueryLexer) int
	Lookahead() int
}

type QueryParserImpl struct {
	lookahead func() int
}

func (p *QueryParserImpl) Lookahead() int {
	return p.lookahead()
}

func QueryNewParser() QueryParser {
	p := &QueryParserImpl{
		lookahead: func() int { return -1 },
	}
	return p
}

const QueryFlag = -1000

func QueryTokname(c int) string {
	if c >= 1 && c-1 < len(QueryToknames) {
		if QueryToknames[c-1] != "" {
			return QueryToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func QueryStatname(s int) string {
	if s >= 0 && s < len(QueryStatenames) {
		if QueryStatenames[s] != "" {
			return QueryStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func QueryErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !QueryErrorVerbose {
		return "syntax error"
	}

	for _, e := range QueryErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + QueryTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := QueryPact[state]
	for tok := TOKSTART; tok-1 < len(QueryToknames); tok++ {
		if n := base + tok; n >= 0 && n < QueryLast && QueryChk[QueryAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if QueryDef[state] == -2 {
		i := 0
		for QueryExca[i] != -1 || QueryExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; QueryExca[i] >= 0; i += 2 {
			tok := QueryExca[i]
			if tok < TOKSTART || QueryExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if QueryExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += QueryTokname(tok)
	}
	return res
}

func Querylex1(lex QueryLexer, lval *QuerySymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = QueryTok1[0]
		goto out
	}
	if char < len(QueryTok1) {
		token = QueryTok1[char]
		goto out
	}
	if char >= QueryPrivate {
		if char < QueryPrivate+len(QueryTok2) {
			token = QueryTok2[char-QueryPrivate]
			goto out
		}
	}
	for i := 0; i < len(QueryTok3); i += 2 {
		token = QueryTok3[i+0]
		if token == char {
			token = QueryTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = QueryTok2[1] /* unknown char */
	}
	if QueryDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", QueryTokname(token), uint(char))
	}
	return char, token
}

func QueryParse(Querylex QueryLexer) int {
	return QueryNewParser().Parse(Querylex)
}

func (Queryrcvr *QueryParserImpl) Parse(Querylex QueryLexer) int {
	var Queryn int
	var Querylval QuerySymType
	var QueryVAL QuerySymType
	var QueryDollar []QuerySymType
	_ = QueryDollar // silence set and not used
	QueryS := make([]QuerySymType, QueryMaxDepth)

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	Querystate := 0
	Querychar := -1
	Querytoken := -1 // Querychar translated into internal numbering
	Queryrcvr.lookahead = func() int { return Querychar }
	defer func() {
		// Make sure we report no lookahead when not parsing.
		Querystate = -1
		Querychar = -1
		Querytoken = -1
	}()
	Queryp := -1
	goto Querystack

ret0:
	return 0

ret1:
	return 1

Querystack:
	/* put a state and value onto the stack */
	if QueryDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", QueryTokname(Querytoken), QueryStatname(Querystate))
	}

	Queryp++
	if Queryp >= len(QueryS) {
		nyys := make([]QuerySymType, len(QueryS)*2)
		copy(nyys, QueryS)
		QueryS = nyys
	}
	QueryS[Queryp] = QueryVAL
	QueryS[Queryp].yys = Querystate

Querynewstate:
	Queryn = QueryPact[Querystate]
	if Queryn <= QueryFlag {
		goto Querydefault /* simple state */
	}
	if Querychar < 0 {
		Querychar, Querytoken = Querylex1(Querylex, &Querylval)
	}
	Queryn += Querytoken
	if Queryn < 0 || Queryn >= QueryLast {
		goto Querydefault
	}
	Queryn = QueryAct[Queryn]
	if QueryChk[Queryn] == Querytoken { /* valid shift */
		Querychar = -1
		Querytoken = -1
		QueryVAL = Querylval
		Querystate = Queryn
		if Errflag > 0 {
			Errflag--
		}
		goto Querystack
	}

Querydefault:
	/* default state action */
	Queryn = QueryDef[Querystate]
	if Queryn == -2 {
		if Querychar < 0 {
			Querychar, Querytoken = Querylex1(Querylex, &Querylval)
		}

		/* look through exception table */
		xi := 0
		for {
			if QueryExca[xi+0] == -1 && QueryExca[xi+1] == Querystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			Queryn = QueryExca[xi+0]
			if Queryn < 0 || Queryn == Querytoken {
				break
			}
		}
		Queryn = QueryExca[xi+1]
		if Queryn < 0 {
			goto ret0
		}
	}
	if Queryn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			Querylex.Error(QueryErrorMessage(Querystate, Querytoken))
			Nerrs++
			if QueryDebug >= 1 {
				__yyfmt__.Printf("%s", QueryStatname(Querystate))
				__yyfmt__.Printf(" saw %s\n", QueryTokname(Querytoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for Queryp >= 0 {
				Queryn = QueryPact[QueryS[Queryp].yys] + QueryErrCode
				if Queryn >= 0 && Queryn < QueryLast {
					Querystate = QueryAct[Queryn] /* simulate a shift of "error" */
					if QueryChk[Querystate] == QueryErrCode {
						goto Querystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if QueryDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", QueryS[Queryp].yys)
				}
				Queryp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if QueryDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", QueryTokname(Querytoken))
			}
			if Querytoken == QueryEofCode {
				goto ret1
			}
			Querychar = -1
			Querytoken = -1
			goto Querynewstate /* try again in the same state */
		}
	}

	/* reduction by production Queryn */
	if QueryDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", Queryn, QueryStatname(Querystate))
	}

	Querynt := Queryn
	Querypt := Queryp
	_ = Querypt // guard against "declared and not used"

	Queryp -= QueryR2[Queryn]
	// Queryp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if Queryp+1 >= len(QueryS) {
		nyys := make([]QuerySymType, len(QueryS)*2)
		copy(nyys, QueryS)
		QueryS = nyys
	}
	QueryVAL = QueryS[Queryp+1]

	/* consult goto table to find next state */
	Queryn = QueryR1[Queryn]
	Queryg := QueryPgo[Queryn]
	Queryj := Queryg + QueryS[Queryp].yys + 1

	if Queryj >= QueryLast {
		Querystate = QueryAct[Queryg]
	} else {
		Querystate = QueryAct[Queryj]
		if QueryChk[Querystate] != -Queryn {
			Querystate = QueryAct[Queryg]
		}
	}
	// dummy call; replaced with literal code
	switch Querynt {

	case 1:
		QueryDollar = QueryS[Querypt-5 : Querypt+1]
		//line query.y:46
		{
			Querylex.(*QueryLex).Query.Selects = QueryDollar[2].selectTermList
			Querylex.(*QueryLex).Query.Wheres = QueryDollar[4].whereClause
		}
	case 2:
		QueryDollar = QueryS[Querypt-3 : Querypt+1]
		//line query.y:51
		{
			Querylex.(*QueryLex).Query.Selects = QueryDollar[2].selectTermList
		}
	case 3:
		QueryDollar = QueryS[Querypt-1 : Querypt+1]
		//line query.y:57
		{
			QueryVAL.selectTermList = QueryDollar[1].selectTermList
		}
	case 4:
		QueryDollar = QueryS[Querypt-1 : Querypt+1]
		//line query.y:63
		{
			QueryVAL.selectTermList = []SelectTerm{QueryDollar[1].selectTerm}
		}
	case 5:
		QueryDollar = QueryS[Querypt-3 : Querypt+1]
		//line query.y:67
		{
			QueryVAL.selectTermList = append([]SelectTerm{QueryDollar[1].selectTerm}, QueryDollar[3].selectTermList...)
		}
	case 6:
		QueryDollar = QueryS[Querypt-2 : Querypt+1]
		//line query.y:71
		{
			QueryVAL.selectTermList = []SelectTerm{{Tag: QueryDollar[2].str}}
		}
	case 7:
		QueryDollar = QueryS[Querypt-1 : Querypt+1]
		//line query.y:77
		{
			QueryVAL.selectTerm = SelectTerm{Tag: QueryDollar[1].str}
		}
	case 8:
		QueryDollar = QueryS[Querypt-1 : Querypt+1]
		//line query.y:81
		{
			QueryVAL.selectTerm = SelectTerm{Tag: QueryDollar[1].str}
		}
	case 9:
		QueryDollar = QueryS[Querypt-1 : Querypt+1]
		//line query.y:88
		{
			letter := Querylex.(*QueryLex).NextLetter()
			QueryDollar[1].whereTerm.Letter = letter
			if QueryDollar[1].whereTerm.IsPredicate {
				QueryVAL.whereClause = WrapTermInSelect(QueryDollar[1].whereTerm.SQL, QueryDollar[1].whereTerm.Letter)
			} else { // have a full select clause
				QueryVAL.whereClause = WhereClause{SQL: QueryDollar[1].whereTerm.SQL, Letter: QueryDollar[1].whereTerm.Letter}
			}
		}
	case 10:
		QueryDollar = QueryS[Querypt-2 : Querypt+1]
		//line query.y:98
		{
			letter := Querylex.(*QueryLex).NextLetter()
			QueryDollar[1].whereTerm.Letter = letter
			if QueryDollar[1].whereTerm.IsPredicate {
				QueryVAL.whereClause = WrapTermInSelectWithTime(QueryDollar[1].whereTerm.SQL, QueryDollar[1].whereTerm.Letter, QueryDollar[2].str)
			} else { // have a full select clause
				QueryVAL.whereClause = WhereClause{SQL: QueryDollar[1].whereTerm.SQL, Letter: QueryDollar[1].whereTerm.Letter}
			}
		}
	case 11:
		QueryDollar = QueryS[Querypt-3 : Querypt+1]
		//line query.y:108
		{
			letter := Querylex.(*QueryLex).NextLetter()
			QueryDollar[1].whereTerm.Letter = letter
			var firstTerm = QueryDollar[1].whereTerm.GetClause()
			sql := fmt.Sprintf(`
	select distinct uuid
	from
	%s as %s
	union
	%s`, firstTerm.SQL, firstTerm.Letter, QueryDollar[3].whereClause.SQL)
			QueryVAL.whereClause = WhereClause{SQL: sql, Letter: firstTerm.Letter}
		}
	case 12:
		QueryDollar = QueryS[Querypt-3 : Querypt+1]
		//line query.y:121
		{
			letter := Querylex.(*QueryLex).NextLetter()
			QueryDollar[1].whereTerm.Letter = letter
			var firstTerm = QueryDollar[1].whereTerm.GetClause()
			sql := fmt.Sprintf(`
	select distinct %s.uuid
	from
	%s as %s
	inner join
	(%s) as %s
	on %s.uuid = %s.uuid`, firstTerm.Letter, firstTerm.SQL, firstTerm.Letter, QueryDollar[3].whereClause.SQL, QueryDollar[3].whereClause.Letter, firstTerm.Letter, QueryDollar[3].whereClause.Letter)
			QueryVAL.whereClause = WhereClause{SQL: sql, Letter: firstTerm.Letter}
		}
	case 13:
		QueryDollar = QueryS[Querypt-2 : Querypt+1]
		//line query.y:135
		{
			sql := fmt.Sprintf(`
	select distinct data.uuid
	from
	data
	where data.uuid not in (%s)`, QueryDollar[2].whereClause.SQL)
			QueryVAL.whereClause = WhereClause{SQL: sql, Letter: QueryDollar[2].whereClause.Letter}
		}
	case 14:
		QueryDollar = QueryS[Querypt-3 : Querypt+1]
		//line query.y:147
		{
			if QueryDollar[1].str == "uuid" {
				QueryVAL.whereTerm = WhereTerm{Key: QueryDollar[1].str, Op: QueryDollar[2].str, Val: QueryDollar[3].str, SQL: fmt.Sprintf(`data.uuid LIKE %s`, QueryDollar[3].str), IsPredicate: true}
			} else {
				QueryVAL.whereTerm = WhereTerm{Key: QueryDollar[1].str, Op: QueryDollar[2].str, Val: QueryDollar[3].str, SQL: fmt.Sprintf(`data.dkey = "%s" and data.dval LIKE %s`, QueryDollar[1].str, QueryDollar[3].str), IsPredicate: true}
			}
		}
	case 15:
		QueryDollar = QueryS[Querypt-3 : Querypt+1]
		//line query.y:155
		{
			if QueryDollar[1].str == "uuid" {
				QueryVAL.whereTerm = WhereTerm{Key: QueryDollar[1].str, Op: QueryDollar[2].str, Val: QueryDollar[3].str, SQL: fmt.Sprintf(`data.uuid = %s`, QueryDollar[3].str), IsPredicate: true}
			} else {
				QueryVAL.whereTerm = WhereTerm{Key: QueryDollar[1].str, Op: QueryDollar[2].str, Val: QueryDollar[3].str, SQL: fmt.Sprintf(`data.dkey = "%s" and data.dval = %s`, QueryDollar[1].str, QueryDollar[3].str), IsPredicate: true}
			}
		}
	case 16:
		QueryDollar = QueryS[Querypt-3 : Querypt+1]
		//line query.y:163
		{
			if QueryDollar[1].str == "uuid" {
				QueryVAL.whereTerm = WhereTerm{Key: QueryDollar[1].str, Op: QueryDollar[2].str, Val: QueryDollar[3].str, SQL: fmt.Sprintf(`data.uuid != %s`, QueryDollar[3].str), IsPredicate: true}
			} else {
				QueryVAL.whereTerm = WhereTerm{Key: QueryDollar[1].str, Op: QueryDollar[2].str, Val: QueryDollar[3].str, SQL: fmt.Sprintf(`data.dkey = "%s" and data.dval != %s`, QueryDollar[1].str, QueryDollar[3].str), IsPredicate: true}
			}
		}
	case 17:
		QueryDollar = QueryS[Querypt-2 : Querypt+1]
		//line query.y:171
		{
			if QueryDollar[2].str == "uuid" {
				QueryVAL.whereTerm = WhereTerm{Key: QueryDollar[1].str, Op: QueryDollar[1].str, SQL: `data.uuid is not null`, IsPredicate: true}
			} else {
				QueryVAL.whereTerm = WhereTerm{Key: QueryDollar[2].str, Op: QueryDollar[1].str, SQL: fmt.Sprintf(`data.dkey = "%s"`, QueryDollar[2].str), IsPredicate: true}
			}
		}
	case 18:
		QueryDollar = QueryS[Querypt-3 : Querypt+1]
		//line query.y:179
		{
			QueryVAL.whereTerm = WhereTerm{SQL: fmt.Sprintf(`(%s)`, QueryDollar[2].whereClause.SQL), IsPredicate: false}
		}
	case 19:
		QueryDollar = QueryS[Querypt-2 : Querypt+1]
		//line query.y:185
		{
			//TODO: fix!
			QueryVAL.str = fmt.Sprintf(timePredicateRange, ">=", _time.Now())
		}
	case 20:
		QueryDollar = QueryS[Querypt-2 : Querypt+1]
		//line query.y:190
		{
			template := `select distinct uuid, dkey, timestamp as maxtime from data
					where timestamp <= "%s"
					order by timestamp desc`
			QueryVAL.str = fmt.Sprintf(template, QueryDollar[2].time.Format(_time.RFC3339))
		}
	case 21:
		QueryDollar = QueryS[Querypt-2 : Querypt+1]
		//line query.y:197
		{
			template := `select distinct uuid, dkey, max(timestamp) as maxtime from data
					where timestamp <= "%s"
					group by dkey, uuid order by timestamp desc`
			QueryVAL.str = fmt.Sprintf(template, QueryDollar[2].time.Format(_time.RFC3339))
		}
	case 22:
		QueryDollar = QueryS[Querypt-2 : Querypt+1]
		//line query.y:204
		{
			template := `select distinct uuid, dkey, timestamp as maxtime from data
					where timestamp >= "%s"
					order by timestamp desc`
			QueryVAL.str = fmt.Sprintf(template, QueryDollar[2].time.Format(_time.RFC3339))
		}
	case 23:
		QueryDollar = QueryS[Querypt-2 : Querypt+1]
		//line query.y:211
		{
			template := `select distinct uuid, dkey, min(timestamp) as maxtime from data
					where timestamp >= "%s"
					group by dkey, uuid order by timestamp desc`
			QueryVAL.str = fmt.Sprintf(template, QueryDollar[2].time.Format(_time.RFC3339))
		}
	case 25:
		QueryDollar = QueryS[Querypt-1 : Querypt+1]
		//line query.y:224
		{
			QueryVAL.time = QueryDollar[1].time
		}
	case 26:
		QueryDollar = QueryS[Querypt-2 : Querypt+1]
		//line query.y:228
		{
			QueryVAL.time = QueryDollar[1].time.Add(QueryDollar[2].timediff)
		}
	case 27:
		QueryDollar = QueryS[Querypt-2 : Querypt+1]
		//line query.y:234
		{
			foundtime, err := parseAbsTime(QueryDollar[1].str, QueryDollar[2].str)
			if err != nil {
				Querylex.(*QueryLex).Error(fmt.Sprintf("Could not parse time \"%v %v\" (%v)", QueryDollar[1].str, QueryDollar[2].str, err.Error()))
			}
			QueryVAL.time = foundtime
		}
	case 28:
		QueryDollar = QueryS[Querypt-1 : Querypt+1]
		//line query.y:242
		{
			num, err := strconv.ParseInt(QueryDollar[1].str, 10, 64)
			if err != nil {
				Querylex.(*QueryLex).Error(fmt.Sprintf("Could not parse integer \"%v\" (%v)", QueryDollar[1].str, err.Error()))
			}
			QueryVAL.time = _time.Unix(num, 0)
		}
	case 29:
		QueryDollar = QueryS[Querypt-1 : Querypt+1]
		//line query.y:250
		{
			found := false
			for _, format := range supported_formats {
				t, err := _time.Parse(format, QueryDollar[1].str)
				if err != nil {
					continue
				}
				QueryVAL.time = t
				found = true
				break
			}
			if !found {
				Querylex.(*QueryLex).Error(fmt.Sprintf("No time format matching \"%v\" found", QueryDollar[1].str))
			}
		}
	case 30:
		QueryDollar = QueryS[Querypt-1 : Querypt+1]
		//line query.y:266
		{
			QueryVAL.time = _time.Now()
		}
	case 31:
		QueryDollar = QueryS[Querypt-2 : Querypt+1]
		//line query.y:272
		{
			var err error
			QueryVAL.timediff, err = parseReltime(QueryDollar[1].str, QueryDollar[2].str)
			if err != nil {
				Querylex.(*QueryLex).Error(fmt.Sprintf("Error parsing relative time \"%v %v\" (%v)", QueryDollar[1].str, QueryDollar[2].str, err.Error()))
			}
		}
	case 32:
		QueryDollar = QueryS[Querypt-3 : Querypt+1]
		//line query.y:280
		{
			newDuration, err := parseReltime(QueryDollar[1].str, QueryDollar[2].str)
			if err != nil {
				Querylex.(*QueryLex).Error(fmt.Sprintf("Error parsing relative time \"%v %v\" (%v)", QueryDollar[1].str, QueryDollar[2].str, err.Error()))
			}
			QueryVAL.timediff = addDurations(newDuration, QueryDollar[3].timediff)
		}
	}
	goto Querystack /* stack new state and value */
}
