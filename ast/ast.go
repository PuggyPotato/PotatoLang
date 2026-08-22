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
		names = append(names, name.Value)
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

func (es *ExpressionStatement) expressionNode() {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}