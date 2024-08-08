package ast

import (
	"bytes"
	"interpreter/token"
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

// program (root node)
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

// prefix expression
type PrefixExpression struct {
	Token    token.Token
	Operator string
	Operand  Expression
}

func (prefixExpression PrefixExpression) expressionNode()      {}
func (prefixExpression PrefixExpression) TokenLiteral() string { return prefixExpression.Token.Literal }
func (prefixExpression PrefixExpression) String() string {
	var buffer bytes.Buffer

	buffer.WriteString("(")
	buffer.WriteString(prefixExpression.Operator)
	buffer.WriteString(prefixExpression.Operand.String())
	buffer.WriteString(")")

	return buffer.String()
}

// infix expression
type InfixExpression struct {
	Token        token.Token
	Operator     string
	LeftOperand  Expression
	RightOperand Expression
}

func (infixExpression InfixExpression) expressionNode()      {}
func (infixExpression InfixExpression) TokenLiteral() string { return infixExpression.Token.Literal }
func (infixExpression InfixExpression) String() string {
	var buffer bytes.Buffer

	buffer.WriteString("(")
	buffer.WriteString(infixExpression.LeftOperand.String())
	buffer.WriteString(" " + infixExpression.Operator + " ")
	buffer.WriteString(infixExpression.RightOperand.String())
	buffer.WriteString(")")

	return buffer.String()
}

// Boolean
type Boolean struct {
	Token token.Token
	Value bool
}

func (boolean *Boolean) expressionNode()      {}
func (boolean *Boolean) TokenLiteral() string { return boolean.Token.Literal }
func (boolean *Boolean) String() string       { return boolean.Token.Literal }

// If
type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var buffer bytes.Buffer

	buffer.WriteString("if")
	buffer.WriteString(ie.Condition.String())
	buffer.WriteString(" ")
	buffer.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		buffer.WriteString("else ")
		buffer.WriteString(ie.Alternative.String())
	}

	return buffer.String()
}

// Block statement
type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var buffer bytes.Buffer

	for _, statement := range bs.Statements {
		buffer.WriteString(statement.String())
	}

	return buffer.String()
}

// Function literal
type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var buffer bytes.Buffer

	params := []string{}
	for _, param := range fl.Parameters {
		params = append(params, param.String())
	}

	buffer.WriteString(fl.TokenLiteral())
	buffer.WriteString("(")
	buffer.WriteString(strings.Join(params, ", "))
	buffer.WriteString(") ")
	buffer.WriteString(fl.Body.String())

	return buffer.String()
}

// Call expression
type CallExpression struct {
	Token     token.Token
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var buffer bytes.Buffer

	args := []string{}
	for _, arg := range ce.Arguments {
		args = append(args, arg.String())
	}

	buffer.WriteString(ce.Function.String())
	buffer.WriteString("(")
	buffer.WriteString(strings.Join(args, ", "))
	buffer.WriteString(")")

	return buffer.String()
}
