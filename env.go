package main

import (
	"fmt"
)

type ValType int32

const (
	NumVal ValType = iota
	BoolVal
	ProcVal
)

func ValTypeToStr(vt ValType) string {
	switch vt {
	case NumVal:
		return "NumVal"
	case BoolVal:
		return "BoolVal"
	case ProcVal:
		return "ProcVal"
	}
	panic("Unknow Val Type!")
}

type Val struct {
	Type ValType
	num  int
	bool bool
	// use by proc
	var_s string
	body *Exp
	env *Env
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

func (v *Val) GetProc() (var_s string, exp *Exp, env *Env) {
	if v.Type != ProcVal {
		panic(fmt.Sprintf("Get proc from Val Fail %+v", *v))
	}
	return v.var_s, v.body, v.env
}

func (v *Val) GetPrettyStr() string {
	if v.Type == NumVal {
		return fmt.Sprintf("%s:%d", ValTypeToStr(v.Type), v.num)
	} else if v.Type == BoolVal {
		return fmt.Sprintf("%s:%t", ValTypeToStr(v.Type), v.bool)
	} else {
		return fmt.Sprintf("%s:[%s]\n[Body]\n%s\n{Env}%s",
			ValTypeToStr(v.Type), v.var_s, ExpPrettyStr(v.body, 0), v.env.PrettyStr())
	}
}

type Env struct {
	key string
	val Val
	next *Env
}

func EmptyEnv() *Env {
	return nil
}

func IsEmptyEnv(env *Env) bool {
	if env == nil {
		return true
	}
	return  false
}

func ExtendEnv(key string, val Val, env *Env) *Env {
	return &Env{key:key,val:val,next:env}
}

func ApplyEnv(key string, env *Env) Val {
	if IsEmptyEnv(env) {
		panic("Not found var in Env:" + key)
	}
	if env.key == key {
		return env.val
	} else {
		return  ApplyEnv(key, env.next)
	}
}

func (e *Env) PrettyStr() string {
	str := ""
	for !IsEmptyEnv(e) {
		str += "[" + e.key + ":" + e.val.GetPrettyStr() + "]"
		e = e.next
	}
	return str
}