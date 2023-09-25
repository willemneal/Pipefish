package parser

import (
	"charm/source/ast"
	"charm/source/object"
	"charm/source/signature"
	"charm/source/token"
)

type FunctionTable map[string][]ast.Function

func (ft FunctionTable) Add(T TypeSystem, functionName string, f ast.Function) (ok bool) {
	if functions, ok := ft[functionName]; ok {
		functions, ok = AddInOrder(T, functions, f)
		ft[functionName] = functions
		return ok
	}
	ft[functionName] = []ast.Function{f}
	return true
}

func (p *Parser) ParamsFitSig(s signature.Signature, parameters []object.Object) bool {

	if len(parameters) == 0 && len(s) == 0 {
		return true
	}
	if len(parameters) > len(s) &&
		((len(s) == 0) || !(s[len(s)-1].VarType == "tuple")) {
		return false
	}
	if len(parameters) < len(s) &&
		!(len(s) == len(parameters)+1 && (s[len(parameters)].VarType == "tuple")) {
		return false
	}
	for i, param := range parameters {
		if i == len(s)-1 && (s[i].VarType == "tuple") {
			return true
		}

		if !IsObjectInType(p.TypeSystem, param, s[i].VarType) && s[i].VarType != "varname" {
			return false
		}
		if param.Type() == object.BLING_OBJ &&
			s[i].VarType == "bling" &&
			param.(*object.Bling).Value != s[i].VarName {
			return false
		}
		if i == len(parameters)-1 {
			return true
		}
	}

	if len(s) == len(parameters)+1 &&
		(s[len(parameters)].VarType == "tuple") {
		return true
	}

	return false
}

func UpdateEnvironment(sig signature.Signature, params []object.Object, env *object.Environment, tok token.Token) (*object.Environment, object.Object) {
	if len(sig) == 0 {
		return env, nil
	}
	sigPos := 0
	tupleAccumulator := []object.Object{}

	for paramPos := 0; paramPos < len(params); paramPos++ {
		if sig[sigPos].VarType == object.TUPLE_OBJ {
			if params[paramPos].Type() == object.BLING_OBJ {
				env.Set(sig[sigPos].VarName, &object.Tuple{Elements: tupleAccumulator})
				tupleAccumulator = []object.Object{}
				sigPos = sigPos + 2
				continue
			}
			if params[paramPos].Type() == object.TUPLE_OBJ && len(params[paramPos].(*object.Tuple).Elements) == 0 {
				if len(tupleAccumulator) == 0 && (paramPos == len(params)-1 || params[paramPos+1].Type() == object.BLING_OBJ) {
					env.Set(sig[sigPos].VarName, object.EMPTY_TUPLE)
					tupleAccumulator = []object.Object{}
					sigPos = sigPos + 1
				}
				continue
			}
			if paramPos == len(params)-1 {
				tupleAccumulator = append(tupleAccumulator, params[paramPos])
				env.Set(sig[sigPos].VarName, &object.Tuple{Elements: tupleAccumulator})
				break
			}
			tupleAccumulator = append(tupleAccumulator, params[paramPos])
			continue
		}

		if params[paramPos].Type() == object.TUPLE_OBJ && len(params[paramPos].(*object.Tuple).Elements) == 0 {
			continue
		}

		if sig[sigPos].VarType == "varname" || sig[sigPos].VarType == "varref" {
			obj, ok := env.Get(sig[sigPos].VarName)
			if !ok {
				return nil, newError("parse/deref/unknown", tok, sig[sigPos].VarName)
			}
			if obj.Type() != object.CODE_OBJ {
				return nil, newError("parse/deref/code", tok, sig[sigPos].VarName)
			}
			if obj.(*object.Code).Value.GetToken().Type != token.IDENT {
				return nil, newError("parse/deref/ident", tok, sig[sigPos].VarName)
			}
			env.Set(obj.(*object.Code).Value.GetToken().Literal, params[paramPos])
		} else {
			env.Set(sig[sigPos].VarName, params[paramPos])
		}
		sigPos++
	}
	return env, nil
}

// Yeah, this isn't very DRY, but it disentangles this from the Golang bits which is helpful.
func GetValueList(sig signature.Signature, params []object.Object) []object.Object {
	if len(sig) == 0 {
		return []object.Object{}
	}
	sigPos := 0
	result := []object.Object{}
	tupleAccumulator := []object.Object{}

	for paramPos := 0; paramPos < len(params); paramPos++ {
		if sig[sigPos].VarType == object.TUPLE_OBJ {
			if params[paramPos].Type() == object.BLING_OBJ {
				result = append(result, &object.Tuple{Elements: tupleAccumulator})
				result = append(result, &object.Bling{})
				tupleAccumulator = []object.Object{}
				sigPos = sigPos + 2
				continue
			}
			if params[paramPos].Type() == object.TUPLE_OBJ && len(params[paramPos].(*object.Tuple).Elements) == 0 {
				if len(tupleAccumulator) == 0 && (paramPos == len(params)-1 || params[paramPos+1].Type() == object.BLING_OBJ) {
					result = append(result, object.EMPTY_TUPLE)
					tupleAccumulator = []object.Object{}
					sigPos = sigPos + 1
				}
				continue
			}
			if paramPos == len(params)-1 {
				tupleAccumulator = append(tupleAccumulator, params[paramPos])
				result = append(result, &object.Tuple{Elements: tupleAccumulator})
				break
			}
			tupleAccumulator = append(tupleAccumulator, params[paramPos])
			continue
		}

		if params[paramPos].Type() == object.TUPLE_OBJ && len(params[paramPos].(*object.Tuple).Elements) == 0 {
			continue
		}

		result = append(result, params[paramPos])

		sigPos++
	}
	return result
}
