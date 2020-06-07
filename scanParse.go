package main

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

type TokenType int32

const (
	EofToken TokenType = iota
	NumberToken
	IdentifierToken
	WhiteSpaceToken
	ReservedToken
	NoteToken
	UnkownToken
)

func TokenTypeToStr(tt TokenType) string {
	switch tt {
	case EofToken:
		return "EofToken"
	case NumberToken:
		return "NumberToken"
	case IdentifierToken:
		return "IdentifierToken"
	case WhiteSpaceToken:
		return "WhiteSpaceToken"
	case ReservedToken:
		return "ReservedToken"
	case NoteToken:
		return "NoteToken"
	}
	return "UnkownToken"
}

type Token struct {
	Type TokenType
	num  int
	str  string
}

func (t *Token) GetNum() (int, error) {
	if t.Type != NumberToken {
		return 0, errors.New("get num from token fail！")
	}
	return t.num, nil
}

func (t *Token) GetStr() string {
	return t.str
}

func (t *Token) GetPrettyStr() string {
	if t.Type == NumberToken {
		return fmt.Sprintf("%s:%d", TokenTypeToStr(t.Type), t.num)
	} else {
		return fmt.Sprintf("%s:%s", TokenTypeToStr(t.Type), t.str)
	}
}

type Scanner struct {
	pos int
	str string
}

func (s *Scanner) Feed(str string) {
	s.str = str
	s.pos = 0
}

func (s *Scanner) NextWithWhite() Token {
	if s.pos >= len(s.str) {
		return Token{Type: EofToken}
	}

	if IsWhiteSpace(s.str[s.pos]) {
		s.pos++
		for s.pos < len(s.str) && IsWhiteSpace(s.str[s.pos]) {
			s.pos++
		}
		return Token{Type: WhiteSpaceToken}
	}

	if IsNum(s.str[s.pos]) {
		var num int
		num = int(s.str[s.pos] - '0')
		s.pos++
		for s.pos < len(s.str) && IsNum(s.str[s.pos]) {
			num = num*10 + int(s.str[s.pos]-'0')
			s.pos++
		}
		return Token{Type: NumberToken, num: num}
	}

	if IsOneWordReserved(s.str[s.pos]) {
		s.pos++
		return Token{Type: ReservedToken, str: string(s.str[s.pos-1])}
	}

	if IsAlpha(s.str[s.pos]) {
		strBytes := make([]byte, 0)
		strBytes = append(strBytes, s.str[s.pos])
		s.pos++
		for s.pos < len(s.str) && IsAlpha(s.str[s.pos]) {
			strBytes = append(strBytes, s.str[s.pos])
			s.pos++
		}
		str := string(strBytes)
		if IsMultiWordReserved(str) {
			return Token{Type: ReservedToken, str: str}
		} else {
			return Token{Type: IdentifierToken, str: str}
		}
	}

	if IsNoteSymbol(s.str[s.pos]) {
		s.pos++
		note := make([]byte, 0)
		for s.pos < len(s.str) && s.str[s.pos] != '\n' {
			note = append(note, s.str[s.pos])
			s.pos++
		}
		s.pos++
		return Token{Type: NoteToken, str: string(note)}
	}

	s.pos++
	return Token{Type: UnkownToken, str: string(s.str[s.pos-1])}
}

func (s *Scanner) Next() Token {
	t := s.NextWithWhite()
	for t.Type == WhiteSpaceToken || t.Type == NoteToken {
		t = s.NextWithWhite()
	}
	return t
}

func DebugScanner(str string) {
	var s Scanner
	s.Feed(str)
	// TODO: 判断err
	t := s.NextWithWhite()
	for t.Type != EofToken {
		fmt.Println(t.GetPrettyStr())
		t = s.NextWithWhite()
	}
}

type Exp interface {
}

func GetExpTypeName(e Exp) string {
	return reflect.TypeOf(e).Name()
}

type ConstExp struct {
	Num int
}
type DiffExp struct {
	Exp1 *Exp
	Exp2 *Exp
}
type ZeroExp struct {
	Exp1 *Exp
}
type IfExp struct {
	Exp1 *Exp
	Exp2 *Exp
	Exp3 *Exp
}
type IdentifyExp struct {
	Var string
}
type LetExp struct {
	Var  string
	Exp1 *Exp
	Body *Exp
}

type LetrecExp struct {
	PName string
	BVar  string
	PBody *Exp
	LetRecBody *Exp
}

type ProcExp struct {
	Var  string
	Exp1 *Exp
}

type CallExp struct {
	Exp1 *Exp
	Exp2 *Exp
}

type Program struct {
	Exp1 *Exp
}

func parse(s *Scanner) Program {
	exp := parse_exp(s)
	t := s.Next()
	if t.Type != EofToken {
		panic("Progrom not meet EOF!")
	}
	return Program{Exp1: &exp}
}

func ScanAndParse(str string) *Program {
	s := Scanner{}
	s.Feed(str)
	p := parse(&s)
	return &p
}

// 从Scanner中读下一个token
// token必须是ReservedToken，并且token.GetStr() == str
// 否则panic
func AssertNextReservedToken(s *Scanner, exp_name string, str string) {
	t := s.Next()
	if t.Type != ReservedToken || t.GetStr() != str {
		panic(exp_name + " parse No " + str)
	}
}

func parse_exp(s *Scanner) Exp {
	t := s.Next()
	if t.Type == NumberToken {
		return ConstExp{Num: t.num}
	}
	if t.Type == ReservedToken && t.GetStr() == "-" {
		AssertNextReservedToken(s, "DiffExp", "(")
		exp1 := parse_exp(s)
		AssertNextReservedToken(s, "DiffExp", ",")
		exp2 := parse_exp(s)
		AssertNextReservedToken(s, "DiffExp", ")")
		return DiffExp{Exp1: &exp1, Exp2: &exp2}
	}
	if t.Type == ReservedToken && t.GetStr() == "zero?" {
		AssertNextReservedToken(s, "ZeroExp", "(")
		exp1 := parse_exp(s)
		AssertNextReservedToken(s, "ZeroExp", ")")
		return ZeroExp{Exp1: &exp1}
	}
	if t.Type == ReservedToken && t.GetStr() == "if" {
		exp1 := parse_exp(s)
		AssertNextReservedToken(s, "IfExp", "then")
		exp2 := parse_exp(s)
		AssertNextReservedToken(s, "IfExp", "else")
		exp3 := parse_exp(s)
		return IfExp{Exp1: &exp1, Exp2: &exp2, Exp3: &exp3}
	}
	if t.Type == ReservedToken && t.GetStr() == "let" {
		t = s.Next()
		if t.Type != IdentifierToken {
			panic("LetExp meet no identifier")
		}
		iden := t.GetStr()
		AssertNextReservedToken(s, "LetExp", "=")
		exp1 := parse_exp(s)
		AssertNextReservedToken(s, "LetExp", "in")
		body := parse_exp(s)
		return LetExp{iden, &exp1, &body}
	}
	if t.Type == ReservedToken && t.GetStr() == "letrec" {
		t = s.Next()
		if t.Type != IdentifierToken {
			panic("LetExp meet no identifier")
		}
		p_name := t.GetStr()
		AssertNextReservedToken(s, "LetrecExp", "(")
		t = s.Next()
		if t.Type != IdentifierToken {
			panic("LetExp meet no identifier")
		}
		b_var := t.GetStr()
		AssertNextReservedToken(s, "LetrecExp", ")")
		AssertNextReservedToken(s, "LetrecExp", "=")
		p_body := parse_exp(s)
		AssertNextReservedToken(s, "LetrecExp", "in")
		letrec_body := parse_exp(s)
		return LetrecExp{PName:p_name, BVar:b_var,PBody:&p_body, LetRecBody:&letrec_body}
	}
	if t.Type == ReservedToken && t.GetStr() == "proc" {
		AssertNextReservedToken(s, "ProcExp", "(")
		t = s.Next()
		if (t.Type != IdentifierToken) {
			panic("Proc Not Get Identifier!")
		}
		AssertNextReservedToken(s, "ProcExp", ")")
		exp1 := parse_exp(s)
		return ProcExp{Var:t.GetStr(),Exp1:&exp1}
	}
	if t.Type == ReservedToken && t.GetStr() == "(" {
		exp1 := parse_exp(s)
		exp2 := parse_exp(s)
		AssertNextReservedToken(s, "CallExp", ")")
		return CallExp{Exp1:&exp1, Exp2:&exp2}
	}
	if t.Type == IdentifierToken {
		return IdentifyExp{Var: t.GetStr()}
	}
	panic("Parser Meet No Exp Avail!")
}

func ExpPrettyStr(exp *Exp, lvl int) string {
	var white string
	for i := 0; i < lvl; i++ {
		white = white + "  "
	}
	str := white + GetExpTypeName(*exp) + ":"
	switch (*exp).(type) {
	case ConstExp:
		str += strconv.Itoa((*exp).(ConstExp).Num)
	case DiffExp:
		dexp := (*exp).(DiffExp)
		str += "\n" + ExpPrettyStr(dexp.Exp1, lvl+1)
		str += "\n" + ExpPrettyStr(dexp.Exp2, lvl+1)
	case ZeroExp:
		zexp := (*exp).(ZeroExp)
		str += "\n" + ExpPrettyStr(zexp.Exp1, lvl+1)
	case IfExp:
		iexp := (*exp).(IfExp)
		str += "\n" + ExpPrettyStr(iexp.Exp1, lvl+1)
		str += "\n" + ExpPrettyStr(iexp.Exp2, lvl+1)
		str += "\n" + ExpPrettyStr(iexp.Exp3, lvl+1)
	case IdentifyExp:
		str += (*exp).(IdentifyExp).Var
	case LetExp:
		lexp := (*exp).(LetExp)
		str += "\n" + white + "  " + lexp.Var
		str += "\n" + ExpPrettyStr(lexp.Exp1, lvl+1)
		str += "\n" + ExpPrettyStr(lexp.Body, lvl+1)
	case LetrecExp:
		lrecexp := (*exp).(LetrecExp)
		str += "\n" + white + "  " + lrecexp.PName  + ":" + lrecexp.BVar
		str += "\n" + ExpPrettyStr(lrecexp.PBody, lvl+1)
		str += "\n" + ExpPrettyStr(lrecexp.LetRecBody, lvl+1)
	case ProcExp:
		pexp := (*exp).(ProcExp)
		str += " " + pexp.Var
		str += "\n" + ExpPrettyStr(pexp.Exp1, lvl + 1)
	case CallExp:
		cexp:= (*exp).(CallExp)
		str += "\n" + ExpPrettyStr(cexp.Exp1, lvl+1)
		str += "\n" + ExpPrettyStr(cexp.Exp2, lvl+1)
	}
	return str
}
