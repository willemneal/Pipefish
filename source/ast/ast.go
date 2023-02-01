package ast

import (
	"bytes"
	"reflect"

	"charm/source/signature"
	"charm/source/token"
)

// The base Node interface
type Node interface {
	GetToken() token.Token
	String() string
}

// Nodes in alphabetical order. Other structures and functions are in a separate section at the bottom.

type ApplicationExpression struct {
	Left  Node
	Right Node
	Token token.Token
}

func (ae *ApplicationExpression) GetToken() token.Token { return ae.Token }
func (ae *ApplicationExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ae.Left.String())
	out.WriteString(" (")
	out.WriteString(ae.Right.String())
	out.WriteString("))")

	return out.String()
}

type AssignmentExpression struct {
	Left  Node
	Right Node
	Token token.Token
}

func (ae *AssignmentExpression) GetToken() token.Token { return ae.Token }
func (ae *AssignmentExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ae.Left.String())
	out.WriteString(" = ")
	out.WriteString(ae.Right.String())
	out.WriteString(")")

	return out.String()
}

type Bling struct {
	Token token.Token
	Value string
}

func (bl *Bling) GetToken() token.Token { return bl.Token }
func (bl *Bling) String() string        { return bl.Value }

type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (b *BooleanLiteral) GetToken() token.Token { return b.Token }
func (b *BooleanLiteral) String() string        { return b.Token.Literal }

type BuiltInExpression struct {
	Name string
}

func (bi *BuiltInExpression) GetToken() token.Token {
	return token.Token{Type: token.BUILTIN, Literal: bi.Name, Line: -1}
}
func (bi *BuiltInExpression) String() string       { return "builtin \"" + bi.Name + "\"" }

type CodeLiteral struct {
	Token token.Token
	Right Node
}

func (cl *CodeLiteral) GetToken() token.Token { return cl.Token }
func (cl *CodeLiteral) String() string {
	return "code (" + cl.Right.String() + ")"
}

type EmptyTuple struct {
	Token token.Token
	Value string
}

func (et *EmptyTuple) GetToken() token.Token { return et.Token }
func (et *EmptyTuple) String() string        { return "()" }

type ExecExpression struct {
	Token token.Token
	Left  Node
	Right Node
}
func (ee *ExecExpression) GetToken() token.Token { return ee.Token }
func (ee *ExecExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ee.Left.String())
	out.WriteString(" exec ")
	out.WriteString(ee.Right.String())
	out.WriteString(")")

	return out.String()
}

type Expression struct {
	Token token.Token // the first token of the expression
	Node  Node
}

func (es *Expression) GetToken() token.Token { return es.Token }
func (es *Expression) String() string {
	return es.Node.String()
}

type FloatLiteral struct {
	Token token.Token
	Value float64
}

func (fl *FloatLiteral) GetToken() token.Token { return fl.Token }
func (fl *FloatLiteral) String() string        { return fl.Token.Literal }

type FuncExpression struct {
	Token token.Token
	Function
}

func (fe *FuncExpression) GetToken() token.Token { return fe.Token }
func (fe *FuncExpression) String() string {
	result := "func " + fe.Sig.String() + " : " + fe.Body.String()
	if fe.Given != nil {
		result = "(" + result + ")" + " given : " + "(" + fe.Given.String() + ")"
	}

	return result
}

type GolangExpression struct {
	Token       token.Token
	ObjectCode  func(args ...any) any
	Raw         []bool
	Sig         signature.Signature
	ReturnTypes signature.Signature
}

func (ge *GolangExpression) GetToken() token.Token { return ge.Token }
func (ge *GolangExpression) String() string        { return ge.Token.Literal }

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) GetToken() token.Token { return i.Token }
func (i *Identifier) String() string        { return i.Value }

type IndexExpression struct {
	Token token.Token // The [ token
	Left  Node
	Index Node
}

func (ie *IndexExpression) GetToken() token.Token { return ie.Token }
func (ie *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")

	return out.String()
}

type InfixExpression struct {
	Token    token.Token
	Left     Node
	Operator string
	Right    Node
	Args     []Node
}

func (ie *InfixExpression) GetToken() token.Token { return ie.Token }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

type IntegerLiteral struct {
	Token token.Token
	Value int
}

func (il *IntegerLiteral) GetToken() token.Token { return il.Token }
func (il *IntegerLiteral) String() string        { return il.Token.Literal }

type LazyInfixExpression struct {
	Token    token.Token
	Left     Node
	Operator string
	Right    Node
}

func (ie *LazyInfixExpression) GetToken() token.Token { return ie.Token }
func (ie *LazyInfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

type ListExpression struct {
	Token token.Token // The [ token
	List  Node
}

func (le *ListExpression) GetToken() token.Token { return le.Token }
func (le *ListExpression) String() string {
	var out bytes.Buffer

	out.WriteString("[")
	if le.List == nil {
		out.WriteString("<nil list! error! bad!>")
	}
	out.WriteString(le.List.String() + " ")
	out.WriteString("]")

	return out.String()
}

type LogExpression struct {
	Token    token.Token
	Args     []Node
	Code     Node
	LogType  LogType
}

func (le *LogExpression) GetToken() token.Token { return le.Token }
func (le *LogExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(le.Code.String())
	out.WriteString(") \\ ")
	for i, arg := range(le.Args) {
		out.WriteString(arg.String())
		if i + 1 < len(le.Args) {
			out.WriteString(", ")
		}
	}
	return out.String()
}

type Nothing struct {
	Token token.Token
}

func (ne *Nothing) GetToken() token.Token { return ne.Token }
func (ne *Nothing) String() string {return ""}

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Args     []Node
}

func (pe *PrefixExpression) GetToken() token.Token { return pe.Token }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(" ")
	for i, v := range(pe.Args) {
		out.WriteString(v.String())
		if i < (len(pe.Args) - 1) && ! (reflect.TypeOf(v) == reflect.TypeOf(&Bling{})) && 
			! (reflect.TypeOf(pe.Args[i + 1]) == reflect.TypeOf(&Bling{})) {
				out.WriteString(",")
		} 
		if i < (len(pe.Args) - 1) {
			out.WriteString(" ")
		}
	}
	out.WriteString(")")

	return out.String()
}

type SetExpression struct {
	Token token.Token // The [ token
	Set   Node
}

func (se *SetExpression) GetToken() token.Token { return se.Token }
func (se *SetExpression) String() string {
	var out bytes.Buffer

	out.WriteString("{")
	out.WriteString(se.Set.String() + " ")
	out.WriteString("}")

	return out.String()
}

type StreamingExpression struct {
	Token    token.Token
	Left     Node
	Operator string
	Right    Node
}

func (se *StreamingExpression) GetToken() token.Token { return se.Token }
func (se *StreamingExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(se.Left.String())
	out.WriteString(" " + se.Operator + " ")
	out.WriteString(se.Right.String())
	out.WriteString(")")

	return out.String()
}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) GetToken() token.Token { return sl.Token }
func (sl *StringLiteral) String() string        { return "\"" + sl.Token.Literal + "\"" }

type StructExpression struct {
	Token token.Token
	Sig   signature.Signature
}

func (st *StructExpression) GetToken() token.Token { return st.Token }
func (st *StructExpression) String() string        { return "struct " + st.Sig.String() }

type SuffixExpression struct {
	Token    token.Token 
	Operator string
	Args     []Node
}

func (se *SuffixExpression) GetToken() token.Token { return se.Token }
func (se *SuffixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	for i, v := range(se.Args) {
		out.WriteString(v.String())
		if i < (len(se.Args) - 1) && ! (reflect.TypeOf(v) == reflect.TypeOf(&Bling{})) && 
			! (reflect.TypeOf(se.Args[i + 1]) == reflect.TypeOf(&Bling{})) {
				out.WriteString(",")
		} 
		if i < (len(se.Args) - 1) {
			out.WriteString(" ")
		}
	}
	out.WriteString(") ")
	out.WriteString(se.Operator)
	return out.String()
}

type TypeLiteral struct {
	Token token.Token
	Value string
}

func (t *TypeLiteral) GetToken() token.Token { return t.Token }
func (t *TypeLiteral) String() string        { return t.Value }

type UnfixExpression struct {
	Token    token.Token
	Operator string
}

func (uf *UnfixExpression) GetToken() token.Token { return uf.Token }
func (uf *UnfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(uf.Operator)
	out.WriteString(")")

	return out.String()
}

// And other useful stuff.

type Function = struct {
	Sig     signature.Signature
	Rets    signature.Signature
	Body    Node
	Given   Node
	Cmd     bool
	Private bool
}

type LogType int

const (
	LogUser LogType = iota
	LogReturn = 1
	LogIf = 2
	LogStart = 3
)

type FnTreeNode struct {
	Fn     *Function
	Branch []*TypeNodePair
}

type TypeNodePair struct { // This exists because we need an *ordered* collection of type-node pairs.
	TypeName string
	Node     *FnTreeNode
}

func (tree FnTreeNode) String() string {
	result := "["
	for i, v := range tree.Branch {
		result = result + v.TypeName
		if v.Node.Fn != nil {
			result = result + "func " + v.Node.Fn.Sig.String()
		} else {
			result = result + v.Node.String()
		}
		if i < len(tree.Branch)-1 {
			result = result + ", "
		}
	}
	return result + "]"
}

func (tree FnTreeNode) IndentString(indent string) string {
	result := ""
	for _, v := range tree.Branch {
		result = result + "\n" + indent + v.TypeName
		if v.Node.Fn != nil {
			result = result + "func " + v.Node.Fn.Sig.String()
		} else {
			result = result + v.Node.IndentString(indent+"    ")
		}
	}
	return result
}
