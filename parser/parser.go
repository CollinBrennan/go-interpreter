package parser

import (
	"fmt"
	"interpreter/ast"
	"interpreter/lexer"
	"interpreter/token"
	"strconv"
)

type prefixParseFunction func() ast.Expression

type infixParseFunction func(ast.Expression) ast.Expression

type Parser struct {
	lexer                *lexer.Lexer
	currentToken         token.Token
	peekToken            token.Token
	errors               []string
	prefixParseFunctions map[token.TokenType]prefixParseFunction
	infixParseFunctions  map[token.TokenType]infixParseFunction
}

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -x or !x
	CALL        // function(x)
)

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
}

func New(lex *lexer.Lexer) *Parser {
	parser := &Parser{
		lexer:  lex,
		errors: []string{},
	}
	parser.nextToken()
	parser.nextToken()

	parser.prefixParseFunctions = make(map[token.TokenType]prefixParseFunction)
	parser.registerPrefix(token.IDENTIFIER, parser.parseIdentifier)
	parser.registerPrefix(token.INT, parser.parseIntegerLiteral)
	parser.registerPrefix(token.STRING, parser.parseStringLiteral)
	parser.registerPrefix(token.BANG, parser.parsePrefixExpression)
	parser.registerPrefix(token.MINUS, parser.parsePrefixExpression)
	parser.registerPrefix(token.LPAREN, parser.parseGroupedExpression)
	parser.registerPrefix(token.IF, parser.parseIfExpression)
	parser.registerPrefix(token.FUNCTION, parser.parseFunctionLiteral)
	parser.registerPrefix(token.TRUE, parser.parseBoolean)
	parser.registerPrefix(token.FALSE, parser.parseBoolean)

	parser.infixParseFunctions = make(map[token.TokenType]infixParseFunction)
	parser.registerInfix(token.PLUS, parser.parseInfixExpression)
	parser.registerInfix(token.MINUS, parser.parseInfixExpression)
	parser.registerInfix(token.SLASH, parser.parseInfixExpression)
	parser.registerInfix(token.ASTERISK, parser.parseInfixExpression)
	parser.registerInfix(token.EQ, parser.parseInfixExpression)
	parser.registerInfix(token.NOT_EQ, parser.parseInfixExpression)
	parser.registerInfix(token.LT, parser.parseInfixExpression)
	parser.registerInfix(token.GT, parser.parseInfixExpression)
	parser.registerInfix(token.LPAREN, parser.parseCallExpression)

	return parser
}

func (parser *Parser) registerPrefix(tokenType token.TokenType, function prefixParseFunction) {
	parser.prefixParseFunctions[tokenType] = function
}

func (parser *Parser) registerInfix(tokenType token.TokenType, function infixParseFunction) {
	parser.infixParseFunctions[tokenType] = function
}

func (parser *Parser) nextToken() {
	parser.currentToken = parser.peekToken
	parser.peekToken = parser.lexer.NextToken()
}

func (parser *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for parser.currentToken.Type != token.EOF {
		statement := parser.parseStatement()
		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}
		parser.nextToken()
	}

	return program
}

func (parser *Parser) parseStatement() ast.Statement {
	switch parser.currentToken.Type {
	case token.LET:
		return parser.parseLetStatement()
	case token.RETURN:
		return parser.parseReturnStatement()
	default:
		return parser.parseExpressionStatement()
	}
}

func (parser *Parser) parseLetStatement() *ast.LetStatement {
	statement := &ast.LetStatement{Token: parser.currentToken}

	if !parser.expectPeek(token.IDENTIFIER) {
		return nil
	}

	statement.Name = &ast.Identifier{Token: parser.currentToken, Value: parser.currentToken.Literal}

	if !parser.expectPeek(token.ASSIGN) {
		return nil
	}

	parser.nextToken()
	statement.Value = parser.parseExpression(LOWEST)

	if parser.peekTokenIs(token.SEMICOLON) {
		parser.nextToken()
	}

	return statement
}

func (parser *Parser) parseReturnStatement() *ast.ReturnStatement {
	statement := &ast.ReturnStatement{Token: parser.currentToken}

	parser.nextToken()
	statement.Value = parser.parseExpression(LOWEST)

	if parser.peekTokenIs(token.SEMICOLON) {
		parser.nextToken()
	}

	return statement
}

func (parser *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	statement := &ast.ExpressionStatement{Token: parser.currentToken}
	statement.Expression = parser.parseExpression(LOWEST)

	if parser.peekTokenIs(token.SEMICOLON) {
		parser.nextToken()
	}

	return statement
}

func (parser *Parser) parseExpression(precedence int) ast.Expression {
	prefixParseFunction := parser.prefixParseFunctions[parser.currentToken.Type]
	if prefixParseFunction == nil {
		parser.noPrefixParseFunctionError(parser.currentToken.Type)
		return nil
	}
	leftExpression := prefixParseFunction()

	for !parser.peekTokenIs(token.SEMICOLON) && precedence < parser.peekPrecedence() {
		infixParseFunction := parser.infixParseFunctions[parser.peekToken.Type]
		if infixParseFunction == nil {
			return leftExpression
		}

		parser.nextToken()
		leftExpression = infixParseFunction(leftExpression)
	}

	return leftExpression
}

func (parser *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    parser.currentToken,
		Operator: parser.currentToken.Literal,
	}

	parser.nextToken()
	expression.Operand = parser.parseExpression(PREFIX)

	return expression
}

func (parser *Parser) parseInfixExpression(leftOperand ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:       parser.currentToken,
		Operator:    parser.currentToken.Literal,
		LeftOperand: leftOperand,
	}

	precedence := parser.currentPrecedence()
	parser.nextToken()
	expression.RightOperand = parser.parseExpression(precedence)

	return expression
}

func (parser *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	expression := &ast.CallExpression{Token: parser.currentToken, Function: function}
	expression.Arguments = parser.parseCallArguments()
	return expression
}

func (parser *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	if parser.peekTokenIs(token.RPAREN) {
		parser.nextToken()
		return args
	}

	parser.nextToken()
	args = append(args, parser.parseExpression(LOWEST))

	for parser.peekTokenIs(token.COMMA) {
		parser.nextToken()
		parser.nextToken()
		args = append(args, parser.parseExpression(LOWEST))
	}

	if !parser.expectPeek(token.RPAREN) {
		return nil
	}

	return args
}

func (parser *Parser) parseGroupedExpression() ast.Expression {
	parser.nextToken()
	expression := parser.parseExpression(LOWEST)

	if !parser.expectPeek(token.RPAREN) {
		return nil
	}

	return expression
}

func (parser *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: parser.currentToken, Value: parser.currentToken.Literal}
}

func (parser *Parser) parseIntegerLiteral() ast.Expression {
	literal := &ast.IntegerLiteral{Token: parser.currentToken}

	value, err := strconv.ParseInt(parser.currentToken.Literal, 0, 64)
	if err != nil {
		message := fmt.Sprintf("Could not parse %q as integer", parser.currentToken.Literal)
		parser.errors = append(parser.errors, message)
		return nil
	}

	literal.Value = value
	return literal
}

func (parser *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{
		Token: parser.currentToken,
		Value: parser.currentToken.Literal,
	}
}

func (parser *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: parser.currentToken}

	if !parser.expectPeek(token.LPAREN) {
		return nil
	}

	expression.Condition = parser.parseGroupedExpression()

	if !parser.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Consequence = parser.parseBlockStatement()

	if parser.peekTokenIs(token.ELSE) {
		parser.nextToken()

		if !parser.expectPeek(token.LBRACE) {
			return nil
		}

		expression.Alternative = parser.parseBlockStatement()
	}

	return expression
}

func (parser *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: parser.currentToken}
	block.Statements = []ast.Statement{}

	parser.nextToken()

	for !parser.currentTokenIs(token.RBRACE) && !parser.currentTokenIs(token.EOF) {
		statement := parser.parseStatement()
		if statement != nil {
			block.Statements = append(block.Statements, statement)
		}

		parser.nextToken()
	}

	return block
}

func (parser *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: parser.currentToken, Value: parser.currentTokenIs(token.TRUE)}
}

func (parser *Parser) parseFunctionLiteral() ast.Expression {
	expression := &ast.FunctionLiteral{Token: parser.currentToken}

	if !parser.expectPeek(token.LPAREN) {
		return nil
	}

	expression.Parameters = parser.parseFunctionParameters()

	if !parser.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Body = parser.parseBlockStatement()

	return expression
}

func (parser *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if parser.peekTokenIs(token.RPAREN) {
		parser.nextToken()
		return identifiers
	}

	parser.nextToken()

	identifier := &ast.Identifier{Token: parser.currentToken, Value: parser.currentToken.Literal}
	identifiers = append(identifiers, identifier)

	for parser.peekTokenIs(token.COMMA) {
		parser.nextToken()
		parser.nextToken()
		identifier := &ast.Identifier{Token: parser.currentToken, Value: parser.currentToken.Literal}
		identifiers = append(identifiers, identifier)
	}

	if !parser.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
}

func (parser *Parser) currentTokenIs(tokenType token.TokenType) bool {
	return parser.currentToken.Type == tokenType
}

func (parser *Parser) peekTokenIs(tokenType token.TokenType) bool {
	return parser.peekToken.Type == tokenType
}

func (parser *Parser) expectPeek(tokenType token.TokenType) bool {
	if parser.peekTokenIs(tokenType) {
		parser.nextToken()
		return true
	} else {
		parser.peekError(tokenType)
		return false
	}
}

func (parser *Parser) peekPrecedence() int {
	if precedence, ok := precedences[parser.peekToken.Type]; ok {
		return precedence
	}
	return LOWEST
}

func (parser *Parser) currentPrecedence() int {
	if precedence, ok := precedences[parser.currentToken.Type]; ok {
		return precedence
	}
	return LOWEST
}

func (parser *Parser) Errors() []string {
	return parser.errors
}

func (parser *Parser) peekError(expectedToken token.TokenType) {
	message := fmt.Sprintf("Expected %s, got %s instead.", expectedToken, parser.peekToken.Type)
	parser.errors = append(parser.errors, message)
}

func (parser *Parser) noPrefixParseFunctionError(tokenType token.TokenType) {
	message := fmt.Sprintf("No prefix parse function for %s found.", tokenType)
	parser.errors = append(parser.errors, message)
}
