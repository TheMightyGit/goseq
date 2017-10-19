//line grammer.y:6
package parse

import __yyfmt__ "fmt"

//line grammer.y:6
import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"text/scanner"
)

var DualRunes = map[string]int{
	".": DOT,
	",": COMMA,

	"--": DOUBLEDASH,
	"-":  DASH,
	"=":  EQUAL,

	">>":  DOUBLEANGR,
	">":   ANGR,
	"/>":  SLASHANGR,
	"\\>": BACKSLASHANGR,
}

//line grammer.y:35
type yySymType struct {
	yys          int
	nodeList     *NodeList
	node         Node
	arrow        ArrowType
	arrowStem    ArrowStemType
	arrowHead    ArrowHeadType
	actorRef     ActorRef
	noteAlign    NoteAlignment
	dividerType  GapType
	blockSegList *BlockSegmentList
	attrList     *AttributeList
	attr         *Attribute

	sval string
}

const K_TITLE = 57346
const K_PARTICIPANT = 57347
const K_NOTE = 57348
const K_STYLE = 57349
const K_LEFT = 57350
const K_RIGHT = 57351
const K_OVER = 57352
const K_OF = 57353
const K_HORIZONTAL = 57354
const K_SPACER = 57355
const K_GAP = 57356
const K_LINE = 57357
const K_FRAME = 57358
const K_ALT = 57359
const K_ELSEALT = 57360
const K_ELSE = 57361
const K_END = 57362
const K_LOOP = 57363
const K_OPT = 57364
const K_PAR = 57365
const K_ELSEPAR = 57366
const K_CONCURRENT = 57367
const K_WHILST = 57368
const DASH = 57369
const DOUBLEDASH = 57370
const DOT = 57371
const EQUAL = 57372
const COMMA = 57373
const ANGR = 57374
const DOUBLEANGR = 57375
const BACKSLASHANGR = 57376
const SLASHANGR = 57377
const PARL = 57378
const PARR = 57379
const STRING = 57380
const MESSAGE = 57381
const IDENT = 57382

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"K_TITLE",
	"K_PARTICIPANT",
	"K_NOTE",
	"K_STYLE",
	"K_LEFT",
	"K_RIGHT",
	"K_OVER",
	"K_OF",
	"K_HORIZONTAL",
	"K_SPACER",
	"K_GAP",
	"K_LINE",
	"K_FRAME",
	"K_ALT",
	"K_ELSEALT",
	"K_ELSE",
	"K_END",
	"K_LOOP",
	"K_OPT",
	"K_PAR",
	"K_ELSEPAR",
	"K_CONCURRENT",
	"K_WHILST",
	"DASH",
	"DOUBLEDASH",
	"DOT",
	"EQUAL",
	"COMMA",
	"ANGR",
	"DOUBLEANGR",
	"BACKSLASHANGR",
	"SLASHANGR",
	"PARL",
	"PARR",
	"STRING",
	"MESSAGE",
	"IDENT",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line grammer.y:335

// Manages the lexer as well as the current diagram being parsed
type parseState struct {
	S     scanner.Scanner
	err   error
	atEof bool
	//diagram     *Diagram
	procInstrs []string
	nodeList   *NodeList
}

func newParseState(src io.Reader, filename string) *parseState {
	ps := &parseState{}
	ps.S.Init(src)
	ps.S.Position.Filename = filename
	//    ps.diagram = &Diagram{}

	return ps
}

func (ps *parseState) Lex(lval *yySymType) int {
	if ps.atEof {
		return 0
	}
	for {
		tok := ps.S.Scan()
		switch tok {
		case scanner.EOF:
			ps.atEof = true
			return 0
		case '#':
			ps.scanComment()
		case ':':
			return ps.scanMessage(lval)
		case '(':
			return PARL
		case ')':
			return PARR
		case '-', '>', '*', '=', '/', '\\', '.', ',':
			if res, isTok := ps.handleDoubleRune(tok); isTok {
				return res
			} else {
				ps.Error("Invalid token: " + scanner.TokenString(tok))
			}
		case scanner.String:
			tokVal := ps.S.TokenText()
			if res, err := strconv.Unquote(tokVal); err == nil {
				lval.sval = res
				return STRING
			} else {
				ps.Error("Invalid string: " + scanner.TokenString(tok) + ": " + err.Error())
			}
		case scanner.Ident:
			return ps.scanKeywordOrIdent(lval)
		default:
			ps.Error("Invalid token: " + scanner.TokenString(tok))
		}
	}
}

func (ps *parseState) handleDoubleRune(firstRune rune) (int, bool) {
	nextRune := ps.S.Peek()

	// Try the double rune
	if nextRune != scanner.EOF {
		tokStr := string(firstRune) + string(nextRune)
		if tok, hasTok := DualRunes[tokStr]; hasTok {
			ps.NextRune()
			return tok, true
		}
	}

	// Try the single rune
	tokStr := string(firstRune)
	if tok, hasTok := DualRunes[tokStr]; hasTok {
		return tok, true
	}

	return 0, false
}

func (ps *parseState) scanKeywordOrIdent(lval *yySymType) int {
	tokVal := ps.S.TokenText()
	switch strings.ToLower(tokVal) {
	case "title":
		return K_TITLE
	case "participant":
		return K_PARTICIPANT
	case "note":
		return K_NOTE
	case "left":
		return K_LEFT
	case "right":
		return K_RIGHT
	case "over":
		return K_OVER
	case "of":
		return K_OF
	case "spacer":
		return K_SPACER
	case "gap":
		return K_GAP
	case "frame":
		return K_FRAME
	case "line":
		return K_LINE
	case "style":
		return K_STYLE
	case "horizontal":
		return K_HORIZONTAL
	case "alt":
		return K_ALT
	case "elsealt":
		return K_ELSEALT
	case "par":
		return K_PAR
	case "elsepar":
		return K_ELSEPAR
	case "else":
		return K_ELSE
	case "end":
		return K_END
	case "loop":
		return K_LOOP
	case "opt":
		return K_OPT
	case "concurrent":
		return K_CONCURRENT
	case "whilst":
		return K_WHILST
	default:
		lval.sval = tokVal
		return IDENT
	}
}

// Scans a message.  A message is all characters up to the new line
func (ps *parseState) scanMessage(lval *yySymType) int {
	buf := new(bytes.Buffer)
	r := ps.NextRune()
	for (r != '\n') && (r != scanner.EOF) {
		if r == '\\' {
			nr := ps.NextRune()
			switch nr {
			case 'n':
				buf.WriteRune('\n')
			case '\\':
				buf.WriteRune('\\')
			default:
				ps.Error("Invalid backslash escape: \\" + string(nr))
			}
		} else {
			buf.WriteRune(r)
		}
		r = ps.NextRune()
	}

	lval.sval = strings.TrimSpace(buf.String())
	return MESSAGE
}

// Scans a comment.  This ignores all characters up to the new line.
func (ps *parseState) scanComment() {
	var buf *bytes.Buffer

	r := ps.NextRune()
	if r == '!' {
		// This starts a processor instruction
		buf = new(bytes.Buffer)
		r = ps.NextRune()
	}

	for (r != '\n') && (r != scanner.EOF) {
		if buf != nil {
			buf.WriteRune(r)
		}
		r = ps.NextRune()
	}

	if buf != nil {
		ps.procInstrs = append(ps.procInstrs, strings.TrimSpace(buf.String()))
	}
}

func (ps *parseState) NextRune() rune {
	if ps.atEof {
		return scanner.EOF
	}

	r := ps.S.Next()
	if r == scanner.EOF {
		ps.atEof = true
	}

	return r
}

func (ps *parseState) Error(err string) {
	errMsg := fmt.Sprintf("%s:%d: %s", ps.S.Position.Filename, ps.S.Position.Line, err)
	ps.err = errors.New(errMsg)
}

func Parse(reader io.Reader, filename string) (*NodeList, error) {
	ps := newParseState(reader, filename)
	yyParse(ps)

	// Add processing instructions to the start of the node list
	for i := len(ps.procInstrs) - 1; i >= 0; i-- {
		instrParts := strings.SplitN(ps.procInstrs[i], " ", 2)
		name, value := strings.TrimSpace(instrParts[0]), strings.TrimSpace(instrParts[1])
		ps.nodeList = &NodeList{&ProcessInstructionNode{name, value}, ps.nodeList}
	}

	if ps.err != nil {
		return nil, ps.err
	} else {
		return ps.nodeList, nil
	}
}

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyPrivate = 57344

const yyLast = 115

var yyAct = [...]int{

	2, 80, 73, 83, 29, 18, 15, 17, 19, 16,
	27, 28, 27, 28, 20, 75, 34, 104, 79, 21,
	32, 101, 99, 24, 23, 22, 78, 25, 98, 96,
	95, 77, 76, 67, 53, 52, 51, 50, 49, 30,
	90, 58, 26, 103, 26, 55, 64, 54, 91, 92,
	68, 69, 70, 71, 72, 33, 60, 61, 62, 63,
	37, 38, 89, 39, 84, 82, 81, 66, 100, 85,
	97, 94, 87, 86, 45, 46, 47, 48, 41, 42,
	43, 31, 57, 65, 74, 93, 56, 88, 44, 40,
	59, 36, 35, 13, 102, 12, 105, 106, 14, 107,
	108, 11, 109, 10, 9, 8, 7, 6, 110, 5,
	4, 112, 111, 3, 1,
}
var yyPact = [...]int{

	2, -1000, -1000, 2, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, 0, 15, -24, 33, 70,
	61, -1, -2, -3, -4, -5, -1000, -1000, -1000, -1000,
	-1000, 9, -1000, -1000, 9, 4, 24, -1000, -1000, -1000,
	4, 72, 56, -1000, -6, -1000, -1000, -1000, -1000, 2,
	2, 2, 2, 2, -1000, -25, -7, -1000, -8, -1000,
	-1000, -1000, -1000, -1000, -13, -1000, -1000, -1000, 47, 45,
	53, 52, 36, 3, 17, 19, -1000, -1000, -1000, 4,
	51, -9, -10, 50, -11, -17, -1000, -1000, 48, -18,
	-1000, -25, 5, -22, -1000, 2, 2, -1000, 2, 2,
	-1000, 2, -1000, -1000, -1000, -1000, 47, -1000, 45, 47,
	-1000, -1000, -1000,
}
var yyPgo = [...]int{

	0, 114, 0, 113, 110, 109, 107, 106, 105, 104,
	103, 101, 98, 95, 93, 92, 5, 91, 90, 89,
	88, 1, 3, 87, 86, 2, 47, 84, 81,
}
var yyR1 = [...]int{

	0, 1, 2, 2, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 4, 5, 28, 28, 24,
	24, 26, 25, 25, 25, 27, 6, 6, 7, 8,
	8, 16, 16, 16, 9, 9, 10, 21, 21, 21,
	11, 22, 22, 22, 13, 14, 12, 23, 23, 20,
	20, 20, 20, 19, 19, 19, 15, 17, 17, 17,
	18, 18, 18, 18,
}
var yyR2 = [...]int{

	0, 1, 0, 2, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 2, 3, 1, 1, 0,
	1, 3, 0, 1, 3, 3, 3, 4, 4, 4,
	6, 1, 1, 1, 2, 3, 5, 0, 3, 4,
	5, 0, 3, 4, 4, 4, 5, 0, 4, 1,
	1, 1, 1, 2, 2, 1, 2, 1, 1, 1,
	1, 1, 1, 1,
}
var yyChk = [...]int{

	-1000, -1, -2, -3, -4, -5, -6, -7, -8, -9,
	-10, -11, -13, -14, -12, 4, 7, 5, -16, 6,
	12, 17, 23, 22, 21, 25, 40, 8, 9, -2,
	39, -28, 5, 40, 40, -15, -17, 27, 28, 30,
	-19, 8, 9, 10, -20, 13, 14, 15, 16, 39,
	39, 39, 39, 39, -26, 36, -24, -26, -16, -18,
	32, 33, 34, 35, -16, 11, 11, 39, -2, -2,
	-2, -2, -2, -25, -27, 40, 39, 39, 39, 31,
	-21, 19, 18, -22, 19, 24, 20, 20, -23, 26,
	37, 31, 30, -16, 20, 39, 39, 20, 39, 39,
	20, 39, -25, 38, 39, -2, -2, -2, -2, -2,
	-21, -22, -21,
}
var yyDef = [...]int{

	2, -2, 1, 2, 4, 5, 6, 7, 8, 9,
	10, 11, 12, 13, 14, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 31, 32, 33, 3,
	15, 0, 17, 18, 19, 0, 0, 57, 58, 59,
	0, 0, 0, 55, 34, 49, 50, 51, 52, 2,
	2, 2, 2, 2, 16, 22, 26, 20, 0, 56,
	60, 61, 62, 63, 0, 53, 54, 35, 37, 41,
	0, 0, 47, 0, 23, 0, 27, 28, 29, 0,
	0, 0, 0, 0, 0, 0, 44, 45, 0, 0,
	21, 22, 0, 0, 36, 2, 2, 40, 2, 2,
	46, 2, 24, 25, 30, 38, 37, 42, 41, 37,
	39, 43, 48,
}
var yyTok1 = [...]int{

	1,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 36, 37, 38, 39, 40,
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
		//line grammer.y:83
		{
			yylex.(*parseState).nodeList = yyDollar[1].nodeList
		}
	case 2:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line grammer.y:90
		{
			yyVAL.nodeList = nil
		}
	case 3:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line grammer.y:94
		{
			yyVAL.nodeList = &NodeList{yyDollar[1].node, yyDollar[2].nodeList}
		}
	case 15:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line grammer.y:115
		{
			yyVAL.node = &TitleNode{yyDollar[2].sval}
		}
	case 16:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line grammer.y:122
		{
			yyVAL.node = &StyleNode{yyDollar[2].sval, yyDollar[3].attrList}
		}
	case 17:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line grammer.y:128
		{
			yyVAL.sval = "participant"
		}
	case 18:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line grammer.y:129
		{
			yyVAL.sval = yyDollar[1].sval
		}
	case 19:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line grammer.y:134
		{
			yyVAL.attrList = nil
		}
	case 20:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line grammer.y:138
		{
			yyVAL.attrList = yyDollar[1].attrList
		}
	case 21:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line grammer.y:145
		{
			yyVAL.attrList = yyDollar[2].attrList
		}
	case 22:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line grammer.y:152
		{
			yyVAL.attrList = nil
		}
	case 23:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line grammer.y:156
		{
			yyVAL.attrList = &AttributeList{yyDollar[1].attr, nil}
		}
	case 24:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line grammer.y:160
		{
			yyVAL.attrList = &AttributeList{yyDollar[1].attr, yyDollar[3].attrList}
		}
	case 25:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line grammer.y:167
		{
			yyVAL.attr = &Attribute{yyDollar[1].sval, yyDollar[3].sval}
		}
	case 26:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line grammer.y:174
		{
			yyVAL.node = &ActorNode{yyDollar[2].sval, false, "", yyDollar[3].attrList}
		}
	case 27:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line grammer.y:178
		{
			yyVAL.node = &ActorNode{yyDollar[2].sval, true, yyDollar[4].sval, yyDollar[3].attrList}
		}
	case 28:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line grammer.y:185
		{
			yyVAL.node = &ActionNode{yyDollar[1].actorRef, yyDollar[3].actorRef, yyDollar[2].arrow, yyDollar[4].sval}
		}
	case 29:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line grammer.y:192
		{
			yyVAL.node = &NoteNode{yyDollar[3].actorRef, nil, yyDollar[2].noteAlign, yyDollar[4].sval}
		}
	case 30:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line grammer.y:196
		{
			yyVAL.node = &NoteNode{yyDollar[3].actorRef, yyDollar[5].actorRef, yyDollar[2].noteAlign, yyDollar[6].sval}
		}
	case 31:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line grammer.y:203
		{
			yyVAL.actorRef = NormalActorRef(yyDollar[1].sval)
		}
	case 32:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line grammer.y:207
		{
			yyVAL.actorRef = PseudoActorRef("left")
		}
	case 33:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line grammer.y:211
		{
			yyVAL.actorRef = PseudoActorRef("right")
		}
	case 34:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line grammer.y:218
		{
			yyVAL.node = &GapNode{yyDollar[2].dividerType, ""}
		}
	case 35:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line grammer.y:222
		{
			yyVAL.node = &GapNode{yyDollar[2].dividerType, yyDollar[3].sval}
		}
	case 36:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line grammer.y:229
		{
			yyVAL.node = &BlockNode{&BlockSegmentList{&BlockSegment{ALT_SEGMENT, "", yyDollar[2].sval, yyDollar[3].nodeList}, yyDollar[4].blockSegList}}
		}
	case 37:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line grammer.y:236
		{
			yyVAL.blockSegList = nil
		}
	case 38:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line grammer.y:240
		{
			yyVAL.blockSegList = &BlockSegmentList{&BlockSegment{ALT_ELSE_SEGMENT, "", yyDollar[2].sval, yyDollar[3].nodeList}, nil}
		}
	case 39:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line grammer.y:244
		{
			yyVAL.blockSegList = &BlockSegmentList{&BlockSegment{ALT_SEGMENT, "", yyDollar[2].sval, yyDollar[3].nodeList}, yyDollar[4].blockSegList}
		}
	case 40:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line grammer.y:251
		{
			yyVAL.node = &BlockNode{&BlockSegmentList{&BlockSegment{PAR_SEGMENT, "", yyDollar[2].sval, yyDollar[3].nodeList}, yyDollar[4].blockSegList}}
		}
	case 41:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line grammer.y:258
		{
			yyVAL.blockSegList = nil
		}
	case 42:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line grammer.y:262
		{
			yyVAL.blockSegList = &BlockSegmentList{&BlockSegment{PAR_ELSE_SEGMENT, "", yyDollar[2].sval, yyDollar[3].nodeList}, nil}
		}
	case 43:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line grammer.y:266
		{
			yyVAL.blockSegList = &BlockSegmentList{&BlockSegment{PAR_SEGMENT, "", yyDollar[2].sval, yyDollar[3].nodeList}, yyDollar[4].blockSegList}
		}
	case 44:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line grammer.y:273
		{
			yyVAL.node = &BlockNode{&BlockSegmentList{&BlockSegment{OPT_SEGMENT, "", yyDollar[2].sval, yyDollar[3].nodeList}, nil}}
		}
	case 45:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line grammer.y:280
		{
			yyVAL.node = &BlockNode{&BlockSegmentList{&BlockSegment{LOOP_SEGMENT, "", yyDollar[2].sval, yyDollar[3].nodeList}, nil}}
		}
	case 46:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line grammer.y:287
		{
			yyVAL.node = &BlockNode{&BlockSegmentList{&BlockSegment{CONCURRENT_SEGMENT, "", "", yyDollar[3].nodeList}, yyDollar[4].blockSegList}}
		}
	case 47:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line grammer.y:294
		{
			yyVAL.blockSegList = nil
		}
	case 48:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line grammer.y:298
		{
			yyVAL.blockSegList = &BlockSegmentList{&BlockSegment{CONCURRENT_WHILST_SEGMENT, "", "", yyDollar[3].nodeList}, yyDollar[4].blockSegList}
		}
	case 49:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line grammer.y:304
		{
			yyVAL.dividerType = SPACER_GAP
		}
	case 50:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line grammer.y:305
		{
			yyVAL.dividerType = EMPTY_GAP
		}
	case 51:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line grammer.y:306
		{
			yyVAL.dividerType = LINE_GAP
		}
	case 52:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line grammer.y:307
		{
			yyVAL.dividerType = FRAME_GAP
		}
	case 53:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line grammer.y:311
		{
			yyVAL.noteAlign = LEFT_NOTE_ALIGNMENT
		}
	case 54:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line grammer.y:312
		{
			yyVAL.noteAlign = RIGHT_NOTE_ALIGNMENT
		}
	case 55:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line grammer.y:313
		{
			yyVAL.noteAlign = OVER_NOTE_ALIGNMENT
		}
	case 56:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line grammer.y:318
		{
			yyVAL.arrow = ArrowType{yyDollar[1].arrowStem, yyDollar[2].arrowHead}
		}
	case 57:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line grammer.y:324
		{
			yyVAL.arrowStem = SOLID_ARROW_STEM
		}
	case 58:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line grammer.y:325
		{
			yyVAL.arrowStem = DASHED_ARROW_STEM
		}
	case 59:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line grammer.y:326
		{
			yyVAL.arrowStem = THICK_ARROW_STEM
		}
	case 60:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line grammer.y:330
		{
			yyVAL.arrowHead = SOLID_ARROW_HEAD
		}
	case 61:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line grammer.y:331
		{
			yyVAL.arrowHead = OPEN_ARROW_HEAD
		}
	case 62:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line grammer.y:332
		{
			yyVAL.arrowHead = BARBED_ARROW_HEAD
		}
	case 63:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line grammer.y:333
		{
			yyVAL.arrowHead = LOWER_BARBED_ARROW_HEAD
		}
	}
	goto yystack /* stack new state and value */
}
