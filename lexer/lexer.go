package lexer

import (
	"bytes"
	"compiler/lexer/define"
	"compiler/util"
	"io"
	"os"
	"path/filepath"
)

const EndToken = "$"

type Lexer struct {
	*util.Stream
	endToken string
}

func FromFile(path string) []*Token {
	absPath, err := filepath.Abs(path)
	if nil != err {
		panic(err)
	}
	f, err := os.Open(absPath)
	if nil != err {
		panic(err)
	}
	defer f.Close()

	return NewLexer(f, EndToken).Analyse()
}

func Analyse(source string) []*Token {
	return NewLexer(bytes.NewBufferString(source), EndToken).Analyse()
}

func NewLexer(r io.Reader, et string) *Lexer {
	s := util.NewStream(r, EndToken)
	return &Lexer{Stream: s, endToken: et}
}

func (l *Lexer) Analyse() []*Token {
	tokens := make([]*Token, 0)
	for l.HasNext() {
		c := l.Next()
		if c == EndToken {
			break
		}
		lookahead := l.Peek()

		if c == " " || c == "\n" || c == "\t" {
			continue
		}

		// 分析注释
		if c == "#" {
			l.PutBack(c)
			tokens = append(tokens, l.MakeComment())
			continue
		}

		// 抓括号
		if define.IsBracket(c) {
			tokens = append(tokens, NewToken(BRACKET, c))
			continue
		}

		// 是个字符串
		if define.IsStringWrap(c) {
			l.PutBack(c)
			tokens = append(tokens, l.MakeString())
			continue
		}

		if define.IsLetter(c) {
			l.PutBack(c)
			tokens = append(tokens, l.MakeVarOrKeyword())
			continue
		}
		if define.IsNumber(c) {
			l.PutBack(c)
			tokens = append(tokens, l.MakeNumber())
			continue
		}

		// 遇到操作符，可能是数字
		if (c == "+" || c == "-" || c == ".") && define.IsNumber(lookahead) {
			var lastToken *Token = nil
			if len(tokens) > 0 {
				lastToken = tokens[len(tokens)-1]
			}

			if lastToken == nil || !lastToken.IsValue() || lastToken.IsOperator() {
				l.PutBack(c)
				tokens = append(tokens, l.MakeNumber())
				continue
			}
		}

		if define.IsOperator(c) {
			l.PutBack(c)
			tokens = append(tokens, l.MakeOp())
			continue
		}
		l.PutBack(c)
		tokens = append(tokens, l.MakeErr())

	}

	return tokens
}

func (l *Lexer) MakeComment() *Token {
	s := ""
	for l.HasNext() {
		c := l.Next()
		if c == "\n" {
			break
		} else {
			s += c
		}
	}
	return NewToken(COMMENT, s)
}

func (l *Lexer) MakeString() *Token {
	s := ""
	state := 0
	for l.HasNext() {
		c := l.Next()
		switch state {
		case 0:
			if c == `'` {
				state = 1
			} else {
				state = 2
			}
			s += c
		case 1:
			if `'` == c {
				return NewToken(STRING, s+c)
			} else {
				s += c
			}
		case 2:
			if `"` == c {
				return NewToken(STRING, s+c)
			} else {
				s += c
			}
		}
	}

	panic("make string failed")
}

func (l *Lexer) MakeVarOrKeyword() *Token {
	s := ""
	for l.HasNext() {
		lookahead := l.Peek()
		if define.IsLiteral(lookahead) {
			s += lookahead
		} else if lookahead == "." {
			return l.MakeErr(s)
		} else {
			break
		}
		l.Next()
	}

	if define.IsKeyTypes(s) {
		return NewToken(TYPE, s)
	}

	if define.IsKeyword(s) {
		return NewToken(KEYWORD, s)
	}

	if s == "true" || s == "false" {
		return NewToken(BOOLEAN, s)
	}

	return NewToken(VARIABLE, s)
}

func (l *Lexer) MakeOp() *Token {
	state := 0

	for l.HasNext() {
		lookahead := l.Next()
		switch state {
		case 0:
			switch lookahead {
			case "+":
				state = 1
			case "-":
				state = 2
			case "*":
				state = 3
			case `/`:
				state = 4
			case `>`:
				state = 5
			case `<`:
				state = 6
			case `=`:
				state = 7
			case `!`:
				state = 8
			case `&`:
				state = 9
			case `|`:
				state = 10
			case `^`:
				state = 11
			case `%`:
				state = 12
			case `:`:
				state = 13
			case ",":
				return NewToken(OPERATOR, ",")
			case ";":
				return NewToken(OPERATOR, ";")
			}
		case 1:
			switch lookahead {
			case `+`:
				return NewToken(OPERATOR, "++")
			case `=`:
				return NewToken(OPERATOR, "+=")
			default:
				l.PutBack(lookahead)
				return NewToken(OPERATOR, "+")
			}
		case 2:
			switch lookahead {
			case `-`:
				return NewToken(OPERATOR, "--")
			case `=`:
				return NewToken(OPERATOR, "-=")
			default:
				l.PutBack(lookahead)
				return NewToken(OPERATOR, "-")
			}
		case 3:
			switch lookahead {
			case `=`:
				return NewToken(OPERATOR, "*=")
			default:
				l.PutBack(lookahead)
				return NewToken(OPERATOR, "*")
			}
		case 4:
			switch lookahead {
			case `=`:
				return NewToken(OPERATOR, "/=")
			default:
				l.PutBack(lookahead)
				return NewToken(OPERATOR, "/")
			}
		case 5:
			switch lookahead {
			case `=`:
				return NewToken(OPERATOR, ">=")
			case `>`:
				return NewToken(OPERATOR, ">>")
			default:
				l.PutBack(lookahead)
				return NewToken(OPERATOR, ">")
			}
		case 6:
			switch lookahead {
			case `=`:
				return NewToken(OPERATOR, "<=")
			case `<`:
				return NewToken(OPERATOR, "<<")
			case `>`:
				return NewToken(OPERATOR, "<>")
			default:
				l.PutBack(lookahead)
				return NewToken(OPERATOR, "<")
			}
		case 7:
			switch lookahead {
			case `=`:
				return NewToken(OPERATOR, "==")
			default:
				l.PutBack(lookahead)
				return NewToken(OPERATOR, "=")
			}
		case 8:
			switch lookahead {
			case `=`:
				return NewToken(OPERATOR, "!=")
			default:
				l.PutBack(lookahead)
				return NewToken(OPERATOR, "!")
			}
		case 9:
			switch lookahead {
			case `&`:
				return NewToken(OPERATOR, "&&")
			case `=`:
				return NewToken(OPERATOR, "&=")
			default:
				l.PutBack(lookahead)
				return NewToken(OPERATOR, "&")
			}
		case 10:
			switch lookahead {
			case `|`:
				return NewToken(OPERATOR, "||")
			case `=`:
				return NewToken(OPERATOR, "|=")
			default:
				l.PutBack(lookahead)
				return NewToken(OPERATOR, "|")
			}
		case 11:
			switch lookahead {
			case `^`:
				return NewToken(OPERATOR, "^^")
			case `=`:
				return NewToken(OPERATOR, "^=")
			default:
				l.PutBack(lookahead)
				return NewToken(OPERATOR, "^")
			}
		case 12:
			switch lookahead {
			case `=`:
				return NewToken(OPERATOR, "%=")
			default:
				l.PutBack(lookahead)
				return NewToken(OPERATOR, "%")
			}
		case 13:
			switch lookahead {
			case "=":
				return NewToken(OPERATOR, ":=")
			default:
				panic("makeOp failed")
			}
		}
	}
	panic("makeOp failed")
}

func (l *Lexer) MakeErr(prefix ...string) *Token {
	s := ""
	row := l.GetLine()
	column := l.GetColumn()
	if len(prefix) > 0 {
		for _, pre := range prefix {
			s += pre
			column -= len(pre)
		}
	}
	for l.HasNext() {
		c := l.Next()
		if define.IsOperator(c) || define.IsBracket(c) {
			l.PutBack(c)
			break
		} else {
			s += c
		}
	}
	return NewTokenWithLocation(ERROR, s, row, column)
}

func (l *Lexer) MakeNumber() *Token {
	state := 0
	s := ""
	for l.HasNext() {
		lookahead := l.Peek()
		switch state {
		case 0:
			if "0" == lookahead {
				state = 1
			} else if define.IsNumber(lookahead) {
				state = 2
			} else if `+` == lookahead || `-` == lookahead {
				state = 3
			} else if lookahead == `.` {
				state = 5
			}
		case 1:
			if lookahead == "0" {
				state = 1
			} else if define.IsNumber(lookahead) {
				state = 2
			} else if lookahead == "." {
				state = 4
			} else {
				return NewToken(INTEGER, s)
			}
		case 2:
			if define.IsNumber(lookahead) {
				state = 2
			} else if lookahead == "." {
				state = 4
			} else if define.IsLetter(lookahead) {
				state = -2
			} else {
				return NewToken(INTEGER, s)
			}
		case 3:
			if define.IsNumber(lookahead) {
				state = 2
			} else if lookahead == "." {
				state = 5
			} else {
				state = -2
			}
		case 4:
			if "." == lookahead {
				state = -2
			} else if define.IsNumber(lookahead) {
				state = 6
			} else if define.IsLetter(lookahead) || define.IsBracket(lookahead) || define.IsOperator(lookahead) {
				state = -2
			} else {
				return NewToken(FLOAT, s)
			}
		case 5:
			if define.IsNumber(lookahead) {
				state = 6
			} else {
				state = -2
			}
		case 6:
			if define.IsNumber(lookahead) {
				state = 6
			} else if "." == lookahead {
				state = -2
			} else {
				return NewToken(FLOAT, s)
			}
		case -1:
			return NewToken(ERROR, s)
		case -2:
			return l.MakeErr(s)
			//for l.HasNext() {
			//	c := l.Next()
			//	if util.IsOperator(c) || util.IsBracket(c) {
			//		l.PutBack(c)
			//		break
			//	} else {
			//		s += c
			//	}
			//}
			//return NewToken(ERROR, s)
		}
		if state >= 0 {
			l.Next()
			s += lookahead
		}
	}
	return nil
}
