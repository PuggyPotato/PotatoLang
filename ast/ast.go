package ast

import (
	"bytes"
	"potatolang/token"
	"strings"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

type LetStatement struct {
	Token token.Token // the token.LET token 
	Names []*Identifier
	Values []Expression
}

func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")

	var names []string
	for _, name := range ls.Names {
		if name.IsMut {
			names = append(names, "mut " + name.Value)
		} else {
			names = append(names, name.Value)
		}
	}
	out.WriteString(strings.Join(names, ", "))
	out.WriteString(" = ")

	if len(ls.Values) > 0 {
		var values []string
		for _, val := range ls.Values {
			if val != nil {
				values = append(values, val.String())
			}
		}
		out.WriteString(strings.Join(values, ", "))
	}

	out.WriteString(";")

	return out.String()
}

type Identifier struct {
	Token token.Token // the token.IDENT token 
	Value string
	IsMut bool // Determines if the identifier is mutable
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string { return i.Value }


type ReturnStatement struct {
	Token token.Token // the 'return' token 
	ReturnValues []Expression
}

func (rs *ReturnStatement) statementNode () {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")


	if len(rs.ReturnValues) > 0 {
		var values []string
		for _, value := range rs.ReturnValues {
			values = append(values, value.String())
		}
		out.WriteString(strings.Join(values, ", "))
	}

	out.WriteString(";")
	
	return out.String()
}

type ExpressionStatement struct {
	Token token.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

type NumberLiteral struct {
	Token token.Token
	Value float64
}

func (nl *NumberLiteral) expressionNode() {}
func (nl *NumberLiteral) TokenLiteral() string { return nl.Token.Literal }
func (nl *NumberLiteral) String() string { return nl.Token.Literal }

type PrefixExpression struct {
	Token token.Token // the prefix token e.g. ! - 
	Operator string
	Right Expression
}

func (pe *PrefixExpression) expressionNode() {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

type InfixExpression struct {
	Token token.Token // the infix token e.g. + - * /
	Left Expression
	Operator string
	Right Expression
}

func (ie *InfixExpression) expressionNode() {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

type Boolean struct {
	Token token.Token // the token.TRUE / token.FALSE
	Value bool
}

func (b *Boolean) expressionNode() {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string { return b.Token.Literal }

type IfExpression struct {
	Token token.Token // The 'if' token 
	Condition Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode() {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

type BlockStatement struct {
	Token token.Token // the { token 
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

type TypedIdentifier struct {
	Token token.Token // the Identifier token e.g. 'x' 
	Value string // 'a'
	Type string // number, string, bool
	IsMut bool
}

func (ti *TypedIdentifier) expressionNode() {}
func (ti *TypedIdentifier) TokenLiteral() string { return ti.Token.Literal }
func (ti *TypedIdentifier) String() string { return ti.Value + ": " + ti.Type }

type FunctionLiteral struct {
	Token token.Token // The 'func' token
 	Parameters []*TypedIdentifier
  	ReturnTypes []string // e.g. ["number", "err"]
   	Body *BlockStatement
}

func (fl *FunctionLiteral) expressionNode() {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	// Put all the params into a slice e.g. (a:number, b:number)
	var params = []string{}
	for i, p := range fl.Parameters {
		if fl.Parameters[i].IsMut {
			params = append(params, "mut " + p.String())
		} else {
			params = append(params, p.String())	
		}
	}

	// Building the first half, e.g. func something
	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")

	if len(fl.ReturnTypes) > 0 {
		out.WriteString(" -> (")
		out.WriteString(strings.Join(fl.ReturnTypes, ", "))
		out.WriteString(")")
	} else{
		// write space if no return type
		out.WriteString(" ")
	}

	out.WriteString(fl.Body.String())
	
	return out.String()
}

type CallExpression struct {
	Token token.Token // the '(' token 
	Function Expression // The identifier or functionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer
	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

type NilLiteral struct {
	Token token.Token // the nil token 
}

func (nl *NilLiteral) expressionNode() {}
func (nl *NilLiteral) TokenLiteral() string { return nl.Token.Literal }
func (nl *NilLiteral) String() string { return nl.Token.Literal }

type AssignStatement struct {
	Token token.Token // the '=' token 
	Names []*Identifier
	Values []Expression 
}

func (as *AssignStatement) statementNode() {}
func (as *AssignStatement) TokenLiteral() string { return as.Token.Literal }
func (as *AssignStatement) String() string { 
	var out bytes.Buffer

	names := []string{}

	for _, n := range as.Names {
		names = append(names, n.String())
	}
	out.WriteString(strings.Join(names, ", "))

	out.WriteString(" = ")

	values := []string{}
	for _, val := range as.Values {
		if val != nil {
			values = append(values, val.String())
		}
	}
	out.WriteString(strings.Join(values, ", "))
	out.WriteString(";")

	return out.String()
}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode() {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string { return sl.Token.Literal }
