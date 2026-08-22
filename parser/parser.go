package parser

import (
	"potatolang/ast"
	"potatolang/lexer"
	"potatolang/token"
)

type Parser struct {
	l *lexer.Lexer

	curToken token.Token
	peekToken token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l:l}

	// read two token, so curToken and peekToken is set
	p.NextToken()
	p.NextToken()

	return p
}

func (p *Parser) NextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}

	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.NextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
		case token.LET:
			return p.parseLetStatement()
		default:
			return nil
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}
	stmt.Names = []*ast.Identifier{}

	for {
		isMut := false

		if p.expectPeek(token.MUT) {
			isMut = true
		}

		if !p.expectPeek(token.IDENT) {
			return nil
		}

		// add Identifier to the slice
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal, IsMut: isMut}

		stmt.Names = append(stmt.Names, ident)

		if p.peekTokenIs(token.COMMA) {
			p.NextToken() // move curToken to ',' and loop again
			continue
		}

		if p.expectPeek(token.ASSIGN) {
			break
		}

		return nil
	}

	for !p.curTokenIs(token.SEMICOLON) {
		p.NextToken()
	}

	return stmt
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t 
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.NextToken()
		return true
	} else {
		return false
	}
}