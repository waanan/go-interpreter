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
	pos   int
	str   string
	cache *Token
}

func (s *Scanner) PutBack(token Token) {
	if s.cache == nil {
		s.cache = &token
	} else {
		panic("scanner cache exist wrong")
	}
}

func (s *Scanner) Feed(str string) {
	s.str = str
	s.pos = 0
}

func (s *Scanner) NextWithWhite() Token {
	if s.cache != nil {
		token := s.cache
		s.cache = nil
		return *token
	}
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
	Var  []string
	Exp1 []*Exp
	Body *Exp
}
type MinusExp struct {
	Exp1 *Exp
}

type EqualExp struct {
	Exp1 *Exp
	Exp2 *Exp
}

type ConsExp struct {
	Exp1 *Exp
	Exp2 *Exp
}

type NullExp struct {
	Exp1 *Exp
}

type CarExp struct {
	Exp1 *Exp
}

type CdrExp struct {
	Exp1 *Exp
}

type EmtyListExp struct {
}

type ListExp struct {
	ExpList []*Exp
}

type CondExp struct {
	ExpCond []*Exp
	ExpBody []*Exp
}

type UnpackExp struct {
	Vars []string
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
		str := make([]string, 0)
		strv := make([]*Exp, 0)
		for {
			t = s.Next()
			if t.Type == ReservedToken && t.GetStr() == "in" {
				break
			} else {
				if t.Type != IdentifierToken {
					panic("LetExp meet no identifier")
				}
				iden := t.GetStr()
				str = append(str, iden)
				AssertNextReservedToken(s, "LetExp", "=")
				exp1 := parse_exp(s)
				strv = append(strv, &exp1)
			}
		}
		body := parse_exp(s)
		return LetExp{str, strv, &body}
	}
	if t.Type == ReservedToken && t.GetStr() == "minus" {
		AssertNextReservedToken(s, "MinusExp", "(")
		exp1 := parse_exp(s)
		AssertNextReservedToken(s, "MinusExp", ")")
		return MinusExp{Exp1: &exp1}
	}
	if t.Type == ReservedToken && t.GetStr() == "equal?" {
		AssertNextReservedToken(s, "EqualExp", "(")
		exp1 := parse_exp(s)
		AssertNextReservedToken(s, "EqualExp", ",")
		exp2 := parse_exp(s)
		AssertNextReservedToken(s, "EqualExp", ")")
		return EqualExp{Exp1: &exp1, Exp2: &exp2}
	}
	if t.Type == ReservedToken && t.GetStr() == "null?" {
		AssertNextReservedToken(s, "NullExp", "(")
		exp1 := parse_exp(s)
		AssertNextReservedToken(s, "NullExp", ")")
		return NullExp{Exp1: &exp1}
	}
	if t.Type == ReservedToken && t.GetStr() == "cons" {
		AssertNextReservedToken(s, "ConsExp", "(")
		exp1 := parse_exp(s)
		AssertNextReservedToken(s, "ConsExp", ",")
		exp2 := parse_exp(s)
		AssertNextReservedToken(s, "ConsExp", ")")
		return ConsExp{Exp1: &exp1, Exp2: &exp2}
	}
	if t.Type == ReservedToken && t.GetStr() == "car" {
		AssertNextReservedToken(s, "CarExp", "(")
		exp1 := parse_exp(s)
		AssertNextReservedToken(s, "CarExp", ")")
		return CarExp{Exp1: &exp1}
	}
	if t.Type == ReservedToken && t.GetStr() == "cdr" {
		AssertNextReservedToken(s, "CdrExp", "(")
		exp1 := parse_exp(s)
		AssertNextReservedToken(s, "CdrExp", ")")
		return CdrExp{Exp1: &exp1}
	}
	if t.Type == ReservedToken && t.GetStr() == "emptylist" {
		return EmtyListExp{}
	}
	// ListExp := list ( Exp1 {,Exp}* ) | list()
	if t.Type == ReservedToken && t.GetStr() == "list" {
		AssertNextReservedToken(s, "ListExp", "(")
		exps := make([]*Exp, 0)
		token := s.Next()
		if token.Type == ReservedToken && token.GetStr() == ")" {
			return ListExp{ExpList: exps}
		} else {
			s.PutBack(token)
		}
		exp1 := parse_exp(s)
		exps = append(exps, &exp1)
		for {
			token := s.Next()
			if token.Type == ReservedToken && token.GetStr() == ")" {
				break
			} else if token.Type == ReservedToken && token.GetStr() == "," {
				exp := parse_exp(s)
				exps = append(exps, &exp)
			} else {
				panic("listExp syntax wrong")
			}
		}
		return ListExp{ExpList: exps}
	}
	if t.Type == ReservedToken && t.GetStr() == "cond" {
		cexp := make([]*Exp, 0)
		bexp := make([]*Exp, 0)
		for {
			token := s.Next()
			if token.Type == ReservedToken && token.GetStr() == "end" {
				break
			} else if token.Type == ReservedToken && token.GetStr() == "{" {
				exp := parse_exp(s)
				cexp = append(cexp, &exp)
				AssertNextReservedToken(s, "CondExp", ">")
				expbody := parse_exp(s)
				bexp = append(bexp, &expbody)
				AssertNextReservedToken(s, "CondExp", "}")
			} else {
				panic("CondExp syntax wrong")
			}
		}
		return CondExp{ExpCond: cexp, ExpBody: bexp}
	}
	if t.Type == ReservedToken && t.GetStr() == "unpack" {
		vexp := make([]string, 0)
		for {
			token := s.Next()
			if token.Type == ReservedToken && token.GetStr() == "=" {
				break
			} else {
				if token.Type != IdentifierToken {
					panic("LetExp meet no identifier")
				}
				vexp = append(vexp, token.GetStr())
			}
		}
		exp1 := parse_exp(s)
		token := s.Next()
		if token.Type == ReservedToken && token.GetStr() == "in" {
			exp2 := parse_exp(s)
			return UnpackExp{Vars: vexp, Exp1: &exp1, Exp2: &exp2}
		}
		panic("unpackExp syntax wrong")
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
		str += "\n" + white + "  " + lexp.Var[0]
		str += "\n" + ExpPrettyStr(lexp.Exp1[0], lvl+1)
		str += "\n" + ExpPrettyStr(lexp.Body, lvl+1)
	}
	return str
}
