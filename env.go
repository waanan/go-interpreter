package main

import (
	"fmt"
)

type ValType int32

const (
	NumVal ValType = iota
	BoolVal
	PairVal
	NullVal
)

func ValTypeToStr(vt ValType) string {
	switch vt {
	case NumVal:
		return "NumVal"
	case BoolVal:
		return "BoolVal"
	case PairVal:
		return "PairVal"
	case NullVal:
		return "NullVal"
	}
	panic("Unknow Val Type!")
}

type Val struct {
	Type ValType
	num  int
	bool bool
	val1 *Val
	val2 *Val
}

func (v *Val) GetPair() (*Val, *Val) {
	if v.Type != PairVal {
		panic(fmt.Sprintf("Get pair from Val Fail %+v", *v))
	}
	return v.val1, v.val2
}

func (v *Val) GetNum() int {
	if v.Type != NumVal {
		panic(fmt.Sprintf("Get num from Val Fail %+v", *v))
	}
	return v.num
}

func (v *Val) GetBool() bool {
	if v.Type != BoolVal {
		panic(fmt.Sprintf("Get bool from Val Fail %+v", *v))
	}
	return v.bool
}

func (v *Val) GetPrettyStr() string {
	if v.Type == NumVal {
		return fmt.Sprintf("%s:%d", ValTypeToStr(v.Type), v.num)
	} else if v.Type == BoolVal {
		return fmt.Sprintf("%s:%t", ValTypeToStr(v.Type), v.bool)
	} else if v.Type == PairVal {
		return fmt.Sprintf("%s:(%s,%s)", ValTypeToStr(v.Type), v.val1.GetPrettyStr(), v.val2.GetPrettyStr())
	} else {
		return fmt.Sprintf("%s:null", ValTypeToStr(v.Type))
	}
}

type Env struct {
	key  string
	val  Val
	next *Env
}

func EmptyEnv() *Env {
	return nil
}

func IsEmptyEnv(env *Env) bool {
	if env == nil {
		return true
	}
	return false
}

func ExtendEnv(key string, val Val, env *Env) *Env {
	return &Env{key: key, val: val, next: env}
}

func ApplyEnv(key string, env *Env) Val {
	if IsEmptyEnv(env) {
		panic("Not found var in Env:" + key)
	}
	if env.key == key {
		return env.val
	} else {
		return ApplyEnv(key, env.next)
	}
}
