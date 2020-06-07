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
		if bv.GetBool() {
			return ValueOfExp(iexp.Exp2, env)
		} else {
			return ValueOfExp(iexp.Exp3, env)
		}
	case IdentifyExp:
		return ApplyEnv((*exp).(IdentifyExp).Var, env)
	case LetExp:
		lexp := (*exp).(LetExp)
		return ValueOfExp(lexp.Body,
			ExtendEnv(lexp.Var,
				ValueOfExp(lexp.Exp1, env),
				env))
	case ProcExp:
		pexp := (*exp).(ProcExp)
		return Val{Type: ProcVal, var_s: pexp.Var, body: pexp.Exp1, env: env}
	case CallExp:
		cexp := (*exp).(CallExp)
		proc := ValueOfExp(cexp.Exp1, env)
		args := make([]*Val, 0)
		for _, v := range cexp.Exp2 {
			arg := ValueOfExp(v, env)
			args = append(args, &arg)
		}
		return CallProc(proc, args, env)
	}
	panic("Unkown Exp!")
}

func CallProc(val1 Val, val2 []*Val, env *Env) Val {
	if val1.Type != ProcVal {
		panic("Call proc rator not Procedure!")
	}
	varS, body, _ := val1.GetProc()
	newEnv := env
	for i, v := range val2 {
		newEnv = ExtendEnv(varS[i], *v, newEnv)
	}
	return ValueOfExp(body, newEnv)
}
