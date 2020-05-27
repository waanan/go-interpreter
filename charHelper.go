package main

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
		c == ',' || c == '=' {
		return true
	}
	return false
}

func IsMultiWordReserved(str string) bool {
	if str == "zero?" || str == "if" || str == "then" ||
		str == "else" || str == "let" || str == "in" || str == "minus" || str == "equal?"{
		return true
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
	if(c == ';'){
		return true
	}
	return false
}
