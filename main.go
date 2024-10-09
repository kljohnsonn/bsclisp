package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"unicode"
)

type Token int

const (
	EOF = iota
	NUMBER
	IDENT // ;
	ILLEGAL

	// Infix ops
	ADD // +
	SUB // -
	MUL // *
	DIV // /

	ASSIGN // =
	DEFINE
	SEMICOLON
	LPAREN
	RPAREN
	LESSTHAN
	GREATERTHAN
)

type Position struct {
	line   int
	column int
}

type Lexer struct {
	pos    Position
	reader *bufio.Reader
}

func (l *Lexer) Lex() (Position, Token, string) {
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return l.pos, EOF, ""
			}

			panic(err)
		}

		l.pos.column++

		switch r {
		case '\n':
			l.resetPos()

		case ';':
			return l.pos, SEMICOLON, ";"
		case '+':
			return l.pos, ADD, "+"
		case '-':
			return l.pos, SUB, "-"
		case '*':
			return l.pos, MUL, "*"
		case '/':
			return l.pos, DIV, "/"
		case '=':
			return l.pos, ASSIGN, "="

		case '(':
			return l.pos, LPAREN, "("

		case ')':
			return l.pos, RPAREN, ")"

		case '<':
			return l.pos, LESSTHAN, "<"

		case '>':
			return l.pos, GREATERTHAN, ">"

		default:
			if unicode.IsSpace(r) {
				continue
			} else if unicode.IsDigit(r) {
				startPos := l.pos

				l.Backup()

				num := l.lexInt()

				return startPos, NUMBER, num
			} else if unicode.IsLetter(r) {
				startPos := l.pos

				l.Backup()

				ident := l.lexIdent()

				return startPos, IDENT, ident

			} else {
				return l.pos, ILLEGAL, string(r)
			}
		}
	}
}

func (l *Lexer) resetPos() {
	l.pos.line++
	l.pos.column = 0
}

var tokens = []string{
	EOF:    "EOF",
	IDENT:  "IDENT",
	NUMBER: "NUMBER",

	// Infix ops
	ADD: "+",
	SUB: "-",
	MUL: "*",
	DIV: "/",

	ASSIGN:      "=",
	DEFINE:      "DEFINE",
	SEMICOLON:   ";",
	LPAREN:      "LPAREN",
	RPAREN:      "RPAREN",
	LESSTHAN:    "<",
	GREATERTHAN: ">",
}

func (t Token) String() string {
	return tokens[t]
}

func (l *Lexer) lexInt() string {
	var number string

	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return number
			}

			panic(err)
		}

		l.pos.column++

		if unicode.IsDigit(r) {
			number = number + string(r)
		} else {
			l.Backup()

			return number
		}
	}
}

func (l *Lexer) lexIdent() string {
	var ident string

	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return ident
			}

			panic(err)
		}

		l.pos.column++

		if unicode.IsLetter(r) {
			ident = ident + string(r)
		} else {
			l.Backup()
			break
		}
	}

	switch ident {
	case "define":
		return ident
	}

	return ident
}

func (l *Lexer) Backup() {
	if err := l.reader.UnreadRune(); err != nil {
		panic(err)
	}

	l.pos.column--
}

func NewLexer(reader io.Reader) *Lexer {
	return &Lexer{
		pos:    Position{line: 1, column: 0},
		reader: bufio.NewReader(reader),
	}
}

func main() {
	file, err := os.Open("input.lisp")

	if err != nil {
		panic(err)
	}

	lexer := NewLexer(file)

	for {
		pos, tok, lit := lexer.Lex()

		if tok == EOF {
			break
		}

		fmt.Printf("%d:%d\t%s\t%s\n", pos.line, pos.column, tok, lit)
	}
}
