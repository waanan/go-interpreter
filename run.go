package main

func Run(str string) Val {
	p := ScanAndParse(str)
	return ValueOfProgram(p)
}

func ValueOfProgram(program *Program) Val {
	return ValueOfExp(program.Exp1, EmptyEnv())
}

func ValueOfExp(exp *Exp, env *Env) Val{
	switch (*exp).(type) {
	case ConstExp:
		return  Val{Type:NumVal,num:(*exp).(ConstExp).Num}
	case DiffExp:
		dexp := (*exp).(DiffExp)
		v1 := ValueOfExp(dexp.Exp1, env)
		v2 := ValueOfExp(dexp.Exp2, env)
		return Val{Type:NumVal, num:v1.GetNum()-v2.GetNum()}
	case ZeroExp:
		zexp := (*exp).(ZeroExp)
		v1 := ValueOfExp(zexp.Exp1, env)
		b := v1.GetNum() == 0
		return  Val{Type:BoolVal,bool:b}
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
				ValueOfExp(lexp.Exp1,env),
				env))
	}
	panic("Unkown Exp!")
}