package main

func Run(str string) Val {
	p := ScanAndParse(str)
	return ValueOfProgram(p)
}

func ValueOfProgram(program *Program) Val {
	return ValueOfExp(TransOfExp(program.Exp1, EmptySEnv()), EmptyEnv())
}

func ValueOfExp(exp *Exp, env *Env) Val {
	switch (*exp).(type) {
	case ConstExp:
		return Val{Type: NumVal, Num: (*exp).(ConstExp).Num}
	case DiffExp:
		dexp := (*exp).(DiffExp)
		v1 := ValueOfExp(dexp.Exp1, env)
		v2 := ValueOfExp(dexp.Exp2, env)
		return Val{Type: NumVal, Num: v1.GetNum() - v2.GetNum()}
	case ZeroExp:
		zexp := (*exp).(ZeroExp)
		v1 := ValueOfExp(zexp.Exp1, env)
		b := v1.GetNum() == 0
		return Val{Type: BoolVal, Bool: b}
	case IfExp:
		iexp := (*exp).(IfExp)
		bv := ValueOfExp(iexp.Exp1, env)
		if bv.GetBool() {
			return ValueOfExp(iexp.Exp2, env)
		} else {
			return ValueOfExp(iexp.Exp3, env)
		}
	case NsIdentifyExp:
		return ApplyEnv((*exp).(NsIdentifyExp).Depth, env)
	case NsLetExp:
		lexp := (*exp).(NsLetExp)
		return ValueOfExp(lexp.Body,
			ExtendEnv(
				ValueOfExp(lexp.Exp1, env),
				env))
	//case LetrecExp:
	//	lrecexp := (*exp).(LetrecExp)
	//	return ValueOfExp(lrecexp.LetRecBody, ExtendEnvRec(
	//		lrecexp.PName, lrecexp.BVar, lrecexp.PBody, env))
	case NsLetRecExp:
		//lrecexp := (*exp).(LetrecExp)
		////计算最终env
		//for i := 0; i < len(lrecexp.PName); i++ {
		//	valt := Val{Type:ProcVal,var_s:lrecexp.BVar[i],body:lrecexp.PBody[i]}
		//	env = ExtendEnv(lrecexp.PName[i],valt,env)
		//}
		//tempEnv := env
		//for i := 0; i < len(lrecexp.PName); i++ {
		//	tempEnv.val.env = env
		//	tempEnv = tempEnv.next
		//}
		//return ValueOfExp(lrecexp.LetRecBody, env)
		nlexp := (*exp).(NsLetRecExp)
		valp := Val{Type: ProcVal, var_s: "", Body: nlexp.PBody, Env: env}
		env1 := ExtendEnv(valp, env)
		env1.Val.Env = env1
		return ValueOfExp(nlexp.LetRecBody, env1)
	case NsProcExp:
		pexp := (*exp).(NsProcExp)
		return Val{Type: ProcVal, var_s: "", Body: pexp.Exp1, Env: env}
	case CallExp:
		cexp := (*exp).(CallExp)
		proc := ValueOfExp(cexp.Exp1, env)
		arg := ValueOfExp(cexp.Exp2, env)
		return CallProc(proc, arg)
	}
	panic("Unkown Exp!")
}

func CallProc(val1 Val, val2 Val) Val {
	if val1.Type != ProcVal {
		panic("Call proc rator not Procedure!")
	}
	_, body, oldEnv := val1.GetProc()
	newEnv := ExtendEnv(val2, oldEnv)
	return ValueOfExp(body, newEnv)
}

func TransOfExp(exp *Exp, env *SEnv) *Exp {
	switch (*exp).(type) {
	case ConstExp:
		return exp
	case DiffExp:
		dexp := (*exp).(DiffExp)
		nexp := Exp(DiffExp{TransOfExp(dexp.Exp1, env), TransOfExp(dexp.Exp2, env)})
		return &nexp
	case ZeroExp:
		zexp := (*exp).(ZeroExp)
		nexp := Exp(ZeroExp{Exp1: TransOfExp(zexp.Exp1, env)})
		return &nexp
	case IfExp:
		iexp := (*exp).(IfExp)
		nexp := Exp(IfExp{TransOfExp(iexp.Exp1, env),
			TransOfExp(iexp.Exp2, env),
			TransOfExp(iexp.Exp3, env)})
		return &nexp
	case IdentifyExp:
		nexp := Exp(NsIdentifyExp{Depth: ApplySEnv((*exp).(IdentifyExp).Var, env)})
		return &nexp
	case LetExp:
		lexp := (*exp).(LetExp)
		nexp := Exp(NsLetExp{TransOfExp(lexp.Exp1, env),
			TransOfExp(lexp.Body,
				ExtendSEnv(lexp.Var, env))})
		return &nexp
	case ProcExp:
		pexp := (*exp).(ProcExp)
		nexp := Exp(
			NsProcExp{TransOfExp(
				pexp.Exp1,
				ExtendSEnv(pexp.Var, env),
			)})
		return &nexp
	case LetrecExp:
		letrecExp := (*exp).(LetrecExp)
		env1 := ExtendSEnv(letrecExp.PName, env)
		env2 := ExtendSEnv(letrecExp.BVar, env1)
		nexp := Exp(NsLetRecExp{TransOfExp(letrecExp.PBody, env2),
			TransOfExp(letrecExp.LetRecBody, env1),
		})
		return &nexp
	case CallExp:
		cexp := (*exp).(CallExp)
		nexp := Exp(CallExp{TransOfExp(cexp.Exp1, env), TransOfExp(cexp.Exp2, env)})
		return &nexp
	}
	panic("Unkown Exp!")
}
