package main

func Run(str string) Val {
	p := ScanAndParse(str)
	return ValueOfProgram(p)
}

func ValueOfProgram(program *Program) Val {
	return ValueOfExp(program.Exp1, EmptyEnv())
}

func ValueOfExp(exp *Exp, env *Env) Val {
	switch (*exp).(type) {
	case ConstExp:
		return Val{Type: NumVal, num: (*exp).(ConstExp).Num}
	case DiffExp:
		dexp := (*exp).(DiffExp)
		v1 := ValueOfExp(dexp.Exp1, env)
		v2 := ValueOfExp(dexp.Exp2, env)
		return Val{Type: NumVal, num: v1.GetNum() - v2.GetNum()}
	case ZeroExp:
		zexp := (*exp).(ZeroExp)
		v1 := ValueOfExp(zexp.Exp1, env)
		b := v1.GetNum() == 0
		return Val{Type: BoolVal, bool: b}
	case IfExp:
		iexp := (*exp).(IfExp)
		bv := ValueOfExp(iexp.Exp1, env)
		if (bv.Type == BoolVal && bv.GetBool()) || (bv.Type == NumVal && bv.GetNum() != 0) {
			return ValueOfExp(iexp.Exp2, env)
		} else {
			return ValueOfExp(iexp.Exp3, env)
		}
	case IdentifyExp:
		return ApplyEnv((*exp).(IdentifyExp).Var, env)
	case LetExp:
		lexp := (*exp).(LetExp)
		env1 := env
		for i, v := range lexp.Var {
			env1 = ExtendEnv(v,
				ValueOfExp(lexp.Exp1[i], env),
				env1)
		}
		return ValueOfExp(lexp.Body, env1)
	case MinusExp:
		mexp := (*exp).(MinusExp)
		val := ValueOfExp(mexp.Exp1, env)
		return Val{Type: NumVal, num: -val.GetNum()}
	case EqualExp:
		eexp := (*exp).(EqualExp)
		val1 := ValueOfExp(eexp.Exp1, env)
		val2 := ValueOfExp(eexp.Exp2, env)
		res := false
		if val1.GetNum() == val2.GetNum() {
			res = true
		}
		return Val{Type: BoolVal, bool: res}
	case EmtyListExp:
		return Val{Type: NullVal}
	case NullExp:
		nexp := (*exp).(NullExp)
		val1 := ValueOfExp(nexp.Exp1, env)
		res := false
		if val1.Type == NullVal {
			res = true
		}
		return Val{Type: BoolVal, bool: res}
	case CarExp:
		cexp := (*exp).(CarExp)
		val1 := ValueOfExp(cexp.Exp1, env)
		return *val1.val1
	case CdrExp:
		cexp := (*exp).(CdrExp)
		val1 := ValueOfExp(cexp.Exp1, env)
		return *val1.val2
	case ConsExp:
		cexp := (*exp).(ConsExp)
		val1 := ValueOfExp(cexp.Exp1, env)
		val2 := ValueOfExp(cexp.Exp2, env)
		return Val{val1: &val1, val2: &val2, Type: PairVal}
	case ListExp:
		explist := (*exp).(ListExp)
		exps := explist.ExpList
		if len(exps) == 0 {
			return Val{Type: NullVal}
		}
		val1 := ValueOfExp(exps[0], env)
		listTail := explist.ExpList[1:]
		listTailStru := Exp(ListExp{ExpList: listTail})
		val2 := ValueOfExp(&listTailStru, env)
		return Val{val1: &val1, val2: &val2, Type: PairVal}
	case CondExp:
		expcond := (*exp).(CondExp)
		cexp := expcond.ExpCond
		bexp := expcond.ExpBody
		if len(cexp) == 0 {
			panic("runtime error,no condition hit")
		}
		for i, v := range cexp {
			hitval := ValueOfExp(v, env)
			if hitval.GetBool() {
				return ValueOfExp(bexp[i], env)
			}
		}
		panic("runtime error,no condition hit")
	case UnpackExp:
		unpackExp := (*exp).(UnpackExp)
		vexp := unpackExp.Vars
		exp1 := unpackExp.Exp1
		exp2 := unpackExp.Exp2
		cur := ValueOfExp(exp1, env)
		for _, v := range vexp {
			if cur.Type == PairVal {
				env = ExtendEnv(v, *cur.val1, env)
			}
			cur = *cur.val2
		}
		return ValueOfExp(exp2, env)
	}

	panic("Unkown Exp!")
}
