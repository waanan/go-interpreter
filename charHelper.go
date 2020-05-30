package main

var ReserverdStr = []string{
	"zero?", "if", "then", "else", "let", "in", "minus", "equal?",
	"cons", "car", "cdr", "null?", "emptylist", "list", "cond", "end",
	"unpack",
}

func IsWhiteSpace(c byte) bool {
	if c == ' ' || c == '\t' || c == '\n' {
		return true
	}
	return false
}

func IsNum(c byte) bool {
	if c >= '0' && c <= '9' {
		return true
	}
	return false
}

func IsOneWordReserved(c byte) bool {
	if c == '-' || c == '(' || c == ')' ||
		c == ',' || c == '=' || c == '>' || c == '{' || c == '}' {
		return true
	}
	return false
}

func IsMultiWordReserved(str string) bool {
	for _, v := range ReserverdStr {
		if v == str {
			return true
		}
	}
	return false
}

func IsAlpha(c byte) bool {
	if (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || c == '?' {
		return true
	}
	return false
}

func IsNoteSymbol(c byte) bool {
	if (c == ';') {
		return true
	}
	return false
}
