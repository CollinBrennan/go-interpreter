package ast

import (
	"bytes"
	"interpreter/token"
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

func (program *Program) TokenLiteral() string {
	if len(program.Statements) > 0 {
		return program.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (program *Program) String() string {
	var out bytes.Buffer
	for _, statement := range program.Statements {
		out.WriteString(statement.String())
	}
	return out.String()
}

// let
type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (LetStatement *LetStatement) statementNode()       {}
func (letStatement *LetStatement) TokenLiteral() string { return letStatement.Token.Literal }
func (letStatement *LetStatement) String() string {
	var buffer bytes.Buffer

	buffer.WriteString(letStatement.TokenLiteral() + " ")
	buffer.WriteString(letStatement.Name.String())
	buffer.WriteString(" = ")

	if letStatement.Value != nil {
		buffer.WriteString(letStatement.Value.String())
	}

	buffer.WriteString(";")
	return buffer.String()
}

// identifier
type Identifier struct {
	Token token.Token
	Value string
}

func (identifier *Identifier) expressionNode()      {}
func (identifier *Identifier) TokenLiteral() string { return identifier.Token.Literal }
func (identifier *Identifier) String() string       { return identifier.Value }

// return
type ReturnStatement struct {
	Token token.Token
	Value Expression
}

func (returnStatement *ReturnStatement) statementNode()       {}
func (returnStatement *ReturnStatement) TokenLiteral() string { return returnStatement.Token.Literal }
func (returnStatement *ReturnStatement) String() string {
	var buffer bytes.Buffer

	buffer.WriteString(returnStatement.TokenLiteral() + " ")

	if returnStatement.Value != nil {
		buffer.WriteString(returnStatement.Value.String())
	}

	buffer.WriteString(";")
	return buffer.String()
}

// expression
type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (expressionStatement *ExpressionStatement) statementNode() {}
func (expressionStatement *ExpressionStatement) TokenLiteral() string {
	return expressionStatement.Token.Literal
}
func (expressionStatement *ExpressionStatement) String() string {
	if expressionStatement.Expression != nil {
		return expressionStatement.Expression.String()
	}
	return ""
}

// int
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (integerLiteral IntegerLiteral) expressionNode()      {}
func (integerLiteral IntegerLiteral) TokenLiteral() string { return integerLiteral.Token.Literal }
func (integerLiteral IntegerLiteral) String() string       { return integerLiteral.Token.Literal }
