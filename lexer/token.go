package lexer

import (
	"fmt"
	"github.com/gookit/color"
)

type TokenType int

var variableMap map[string]int

// 初始化变量符号表
func init() {
	variableMap = make(map[string]int)
}

const (
	KEYWORD  TokenType = 1		// 关键字类型为 1
	TYPE     TokenType = 2		// 类型类型为 2
	VARIABLE TokenType = 3		// 变量类型为 3
	OPERATOR TokenType = 4		// 操作符类型为 4
	BRACKET  TokenType = 5		// 括号类型为 5
	STRING   TokenType = 6		// 字符串类型为 6
	FLOAT    TokenType = 7		// 浮点数类型为 7
	BOOLEAN  TokenType = 8		// 布尔类型为 8
	INTEGER  TokenType = 9		// 整数类型为 9
	COMMENT  TokenType = 10		// 注释类型为 10
	ERROR    TokenType = -1		// 错误类型为 -1
	END      TokenType = -2		// 结束类型为 -2
)

// 判断 Token 类型 并创建其可输出的字符串 用于提示和输出
func (tt TokenType) String() string {
	switch tt {

	case KEYWORD:
		return "keyword "
	case VARIABLE:
		return "variable"
	case TYPE:
		return "type    "
	case OPERATOR:
		return "operator"
	case BRACKET:
		return "bracket "
	case STRING:
		return "string  "
	case FLOAT:
		return "float   "
	case BOOLEAN:
		return "boolean "
	case INTEGER:
		return "integer "
	case COMMENT:
		return "comment "
	case ERROR:
		return "error   "
	case END:
		return "end     "
	}

	panic("unexpected token type")
}

// Token 结构体 标识处理每一个词法单元
type Token struct {
	Typ    TokenType
	Value  string
	ID     int
	Column int
	Row    int
}

// NewToken 创建一个词法对象
func NewToken(t TokenType, v string) *Token {
	ID := 0
	if t == 3 {
		if id, ok := variableMap[v]; !ok {
			variableMap[v] = len(variableMap) + 1
			ID = variableMap[v]
		} else {
			ID = id
		}
	}
	return &Token{Typ: t, Value: v, ID: ID}
}

// NewTokenWithLocation 创建一个带有位置信息的词法对象
func NewTokenWithLocation(t TokenType, v string, row int, column int) *Token {
	ID := 0
	if t == 3 {
		if id, ok := variableMap[v]; !ok {
			variableMap[v] = len(variableMap) + 1
			ID = variableMap[v]
		} else {
			ID = id
		}
	}
	return &Token{Typ: t, Value: v, ID: ID, Row: row, Column: column}
}

// IsVariable 判断词法对象是否是变量类型
func (t *Token) IsVariable() bool {
	return t.Typ == VARIABLE
}

// IsScalar 判断词法对象是否是标量类型
func (t *Token) IsScalar() bool {
	return t.Typ == FLOAT || t.Typ == BOOLEAN || t.Typ == INTEGER || t.Typ == STRING
}

// IsNumber 判断词法对象是否是数字类型
func (t *Token) IsNumber() bool {
	return t.Typ == INTEGER || t.Typ == FLOAT
}

// IsOperator 判断词法对象是否是操作符类型
func (t *Token) IsOperator() bool {
	return t.Typ == OPERATOR
}

// IsKeyword 判断词法对象是否是括号类型
func (t *Token) IsBracket() bool {
	return t.Typ == BRACKET
}

// IsEnd 判断词法对象是否是结束类型
func (t *Token) IsEnd() bool {
	return t.Typ == END
}

// Show 输出词法对象 格式化输出结果
func (t *Token) Show() {
	if t.Typ == ERROR {
		color.Red.Printf("【%d:%d】ERROR: You Have an Lexical error near \"%s\" at line %d column %d.\n", t.Row, t.Column, t.Value, t.Row, t.Column)
	} else if t.Typ == VARIABLE {
		color.Green.Printf("< %v, %v, id=%v >\n", t.Typ, t.Value, t.ID)
	} else {
		fmt.Printf("< %v, %v >\n", t.Typ, t.Value)
	}
}

// IsValue 判断词法对象是否是变量或者标量
func (t *Token) IsValue() bool {
	return t.IsVariable() || t.IsScalar()
}

// IsType 判断词法对象是否是类型定义类型
func (t *Token) IsType() bool {
	return t.Typ == TYPE
}
