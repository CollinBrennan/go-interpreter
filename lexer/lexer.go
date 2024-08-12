package lexer

import "interpreter/token"

type Lexer struct {
	input        string
	position     int
	nextPosition int
	char         byte
}

func New(input string) *Lexer {
	lexer := &Lexer{
		input:        input,
		position:     0,
		nextPosition: 0,
		char:         0,
	}
	lexer.readChar()
	return lexer
}

func (lexer *Lexer) NextToken() token.Token {
	var nextToken token.Token

	lexer.skipWhitespace()

	switch lexer.char {
	case '(':
		nextToken = newToken(token.LPAREN, lexer.char)
	case ')':
		nextToken = newToken(token.RPAREN, lexer.char)
	case '{':
		nextToken = newToken(token.LBRACE, lexer.char)
	case '}':
		nextToken = newToken(token.RBRACE, lexer.char)
	case '=':
		if lexer.peekChar() == '=' {
			firstChar := lexer.char
			lexer.readChar()
			nextToken = token.Token{
				Type:    token.EQ,
				Literal: string(firstChar) + string(lexer.char),
			}
		} else {
			nextToken = newToken(token.ASSIGN, lexer.char)
		}
	case '+':
		nextToken = newToken(token.PLUS, lexer.char)
	case '-':
		nextToken = newToken(token.MINUS, lexer.char)
	case '!':
		if lexer.peekChar() == '=' {
			firstChar := lexer.char
			lexer.readChar()
			nextToken = token.Token{
				Type:    token.NOT_EQ,
				Literal: string(firstChar) + string(lexer.char),
			}
		} else {
			nextToken = newToken(token.BANG, lexer.char)
		}
	case '/':
		nextToken = newToken(token.SLASH, lexer.char)
	case '*':
		nextToken = newToken(token.ASTERISK, lexer.char)
	case '<':
		if lexer.peekChar() == '>' {
			firstChar := lexer.char
			lexer.readChar()
			nextToken = token.Token{
				Type:    token.LTGT,
				Literal: string(firstChar) + string(lexer.char),
			}
		} else {
			nextToken = newToken(token.LT, lexer.char)
		}
	case '>':
		nextToken = newToken(token.GT, lexer.char)
	case ';':
		nextToken = newToken(token.SEMICOLON, lexer.char)
	case ',':
		nextToken = newToken(token.COMMA, lexer.char)
	case '"':
		nextToken.Type = token.STRING
		nextToken.Literal = lexer.readString()
	case 0:
		nextToken.Literal = ""
		nextToken.Type = token.EOF
	default:
		if isLetter(lexer.char) {
			literal := lexer.readIdentifier()
			nextToken.Literal = literal
			nextToken.Type = token.LookupIdentifier(literal)
			return nextToken
		} else if isDigit(lexer.char) {
			nextToken.Literal = lexer.readNumber()
			nextToken.Type = token.INT
			return nextToken
		} else {
			nextToken = newToken(token.ILLEGAL, lexer.char)
		}
	}

	lexer.readChar()
	return nextToken
}

func (lexer *Lexer) readString() string {
	pos := lexer.position + 1
	for {
		lexer.readChar()
		if lexer.char == '"' || lexer.char == 0 {
			break
		}
	}
	return lexer.input[pos:lexer.position]
}

func (lexer *Lexer) readIdentifier() string {
	startPosition := lexer.position
	for isLetter(lexer.char) {
		lexer.readChar()
	}
	word := lexer.input[startPosition:lexer.position]
	return word
}

func (lexer *Lexer) readNumber() string {
	startPosition := lexer.position
	for isDigit(lexer.char) {
		lexer.readChar()
	}
	number := lexer.input[startPosition:lexer.position]
	return number
}

func isLetter(char byte) bool {
	return char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z' || char == '_'
}

func isDigit(char byte) bool {
	return char >= '0' && char <= '9'
}

func (lexer *Lexer) skipWhitespace() {
	for lexer.char == ' ' || lexer.char == '\t' || lexer.char == '\n' || lexer.char == '\r' {
		lexer.readChar()
	}
}

func newToken(tokenType token.TokenType, char byte) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: string(char),
	}
}

func (lexer *Lexer) readChar() {
	if lexer.nextPosition >= len(lexer.input) {
		lexer.char = 0
	} else {
		lexer.char = lexer.input[lexer.nextPosition]
	}
	lexer.position = lexer.nextPosition
	lexer.nextPosition += 1
}

func (lexer *Lexer) peekChar() byte {
	if lexer.nextPosition >= len(lexer.input) {
		return 0
	} else {
		return lexer.input[lexer.nextPosition]
	}
}
