package main

import (
	"encoding/json"
	"fmt"
)

type ValType int32
type EnvType int32

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
	Num  int
	Bool bool
	// use by proc
	var_s string
	Body *Exp
	Env *Env
}

func (v *Val) GetNum() int {
	if v.Type != NumVal {
		panic(fmt.Sprintf("Get num from Val Fail %+v", *v))
	}
	return v.Num
}

func (v *Val) GetBool() bool {
	if v.Type != BoolVal {
		panic(fmt.Sprintf("Get bool from Val Fail %+v", *v))
	}
	return v.Bool
}

func (v *Val) GetProc() (var_s string, exp *Exp, env *Env) {
	if v.Type != ProcVal {
		panic(fmt.Sprintf("Get proc from Val Fail %+v", *v))
	}
	return v.var_s, v.Body, v.Env
}

func (v *Val) GetPrettyStr() string {
	b, _ :=json.MarshalIndent(v,"", "\t")
	return string(b)
}

const (
	RegularEnv EnvType = iota
	RecEnv
)


type Env struct {
	Type EnvType
	Val Val

	Bvar string
	Bexp *Exp
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

func ExtendEnv(val Val, env *Env) *Env {
	return &Env{Val:val,next:env,Type:RegularEnv}
}

func ExtendEnvRec(bvar string, bbody *Exp, env *Env) *Env {
	return &Env{Type:RecEnv,Bvar:bvar,Bexp:bbody,next:env}
}

func ApplyEnv(depth int, env *Env) Val {
	if IsEmptyEnv(env) {
		panic("Not found var in Env:" + string(depth))
	}
	if depth == 0 {
		return env.Val
	} else {
		return ApplyEnv(depth-1, env.next)
	}
}


type SEnv struct {
	key string
	next *SEnv
}
func EmptySEnv() *SEnv {
	return nil
}
func IsEmptySEnv(env *SEnv) bool {
	if env == nil {
		return true
	}
	return  false
}
func ExtendSEnv(key string, env *SEnv) *SEnv {
	return &SEnv{key:key,next:env}
}
func ApplySEnv(key string, env *SEnv) int {
	if IsEmptySEnv(env) {
		panic("Not found var in SEnv:" + key)
	}
	if key == env.key {
		return 0
	} else {
		return 1 + ApplySEnv(key, env.next)
	}
}
