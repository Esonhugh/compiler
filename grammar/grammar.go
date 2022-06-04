package grammar

import (
	"bytes"
	"github.com/esonhugh/compiler/lexer"
	"github.com/esonhugh/compiler/util"
	"fmt"
	"io"
	"strings"
)

// EndToken 句子终结符
const EndToken = "$"

// Grammar 对象定义了一个语法分析需要的全部对象 以及他的分析器
type Grammar struct {
	*util.Stream                     // Stream is the stream of the sentence
	productions  map[string][]string // productions is the production map use the "string"=>"string1" "string"=>"string2" ...
	endToken     string              // endToken is the end of the sentence
	tokens       *util.Queue         // tokens is the queue of the sentence
}

/* makeProductions 创建产生式的映射关系 map
由于 是单个字符匹配 所以 E' 被定义为 G
替换如下:
在文件 [utils/transfer/transfer.go] 中进行替换

E' ==> G 
T' ==> S
& ==> ε

原始的产生式如下

E→TE’
E’→ATE’|ε
T→FT’
T’→MFT’ |ε
F→(E) | i
A → + | -
M → * | /
*/
func makeProductions() map[string][]string {
	res := make(map[string][]string)
	res["E"] = []string{"TG"}
	res["G"] = []string{"ATG", "&"}
	res["T"] = []string{"FS"}
	res["S"] = []string{"MFS", "&"}
	res["F"] = []string{"(E)", "i"}
	res["A"] = []string{"+", "-"}
	res["M"] = []string{"/", "*"}
	return res
}

// Analyse 语法分析器主函数
func Analyse(raw []*lexer.Token) ([]*Production, bool) {
	// 从 E 开始 E 即为开始符号 到 EndToken 为句子结束符号
	grm := NewGrammar(raw, bytes.NewBufferString("E"), EndToken)
	productions, success := grm.Analyse(1)
	if !success {
		return nil, false
	}
	e := grm.tokens.FrontRaw()
	if e != nil {
		if e.Value.(*lexer.Token).Value == "" && e.Next() == nil {
			return productions, true
		}
		return nil, false
	}

	return productions, true
}

func NewGrammar(token []*lexer.Token, r io.Reader, et string) *Grammar {
	s := util.NewStream(r, EndToken)
	q := util.New()
	for i := 0; i < len(token); i++ {
		q.PushBack(token[i])
	}
	return &Grammar{Stream: s, endToken: et, productions: makeProductions(), tokens: q}
}



func (g *Grammar) GetNextOrigin() string {
	c := g.Stream.Next()
	lookahead := g.Stream.Peek()
	if lookahead == "'" {
		return c + g.Stream.Next()
	}
	return c
}

// isEndType judge the one word is terminal or not.
// EndType is the unspilttable token
func (g *Grammar) isEndType(r string) bool {
	dic := map[string]bool{
		"i": true,
		"+": true,
		"-": true,
		"*": true,
		"/": true,
		"(": true,
		")": true,
		"&": true,
	}
	return dic[r]
}

// Analyse the grammar
// count is the 
func (g *Grammar) Analyse(count int) (res []*Production, canKill bool) {

	// 当前推导式头部
	origin := g.GetNextOrigin()

	// 当前推导到空字符
	if origin == "&" {
		res = append(res, &Production{
			Type:   "kill",
			Target: "&",
		})
		g.tokens.PushFront(&lexer.Token{})
		return res, true
	}

	// 当前产生式首位为终结符
	if g.isEndType(origin) {
		c := g.tokens.Front()
		if c == nil {
			return res, false
		}
		currentToken := c.(*lexer.Token)
		switch origin {
		case "i":
			{
				if !currentToken.IsValue() {
					g.Stream.ClearFronts(count - 1)
					return res, false
				}
				res = append(res, &Production{
					Type:   "kill",
					Target: "i",
				})
				return res, true
			}
		case "+":
			{
				if !currentToken.IsOperator() || currentToken.Value != "+" {
					g.Stream.ClearFronts(count - 1)
					return res, false
				}
				res = append(res, &Production{
					Type:   "kill",
					Target: "+",
				})
				return res, true
			}
		case "-":
			{
				if !currentToken.IsOperator() || currentToken.Value != "-" {
					g.Stream.ClearFronts(count - 1)
					return res, false
				}
				res = append(res, &Production{
					Type:   "kill",
					Target: "-",
				})
				return res, true
			}
		case "*":
			{
				if !currentToken.IsOperator() || currentToken.Value != "*" {
					g.Stream.ClearFronts(count - 1)
					return res, false
				}
				res = append(res, &Production{
					Type:   "kill",
					Target: "*",
				})
				return res, true
			}
		case "/":
			{
				if !currentToken.IsOperator() || currentToken.Value != "/" {
					g.Stream.ClearFronts(count - 1)
					return res, false
				}
				res = append(res, &Production{
					Type:   "kill",
					Target: "/",
				})
				return res, true
			}
		case "(":
			{
				if !currentToken.IsBracket() || currentToken.Value != "(" {
					g.Stream.ClearFronts(count - 1)
					return res, false
				}
				res = append(res, &Production{
					Type:   "kill",
					Target: "(",
				})
				return res, true
			}
		case ")":
			{
				if !currentToken.IsBracket() || currentToken.Value != ")" {
					g.Stream.ClearFronts(count - 1)
					return res, false
				}
				res = append(res, &Production{
					Type:   "kill",
					Target: ")",
				})
				return res, true
			}
		}
	}
	
	// 选择一个可用的产生式子
	canUse := g.productions[origin]
	var killToken []*lexer.Token
	for i := 0; i < len(canUse); i++ {
		// 倒序插入
		for _, s := range reverse(strings.Split(canUse[i], "")) {
			g.Stream.PutBack(s)
		}
		//g.Stream.Print()
		//g.PrintToken()
		match := false
		if canUse[i] == "" {
			match = true
		}
		for j := 0; j < len(canUse[i]); j++ {
			next := len(canUse[i]) - j
			ps, kill := g.Analyse(next)
			if kill {
				res = append(res, &Production{
					Origin: origin,
					Next:   canUse[i],
					Type:   "Continue",
				})
				res = append(res, ps...)
				if j != len(canUse[i])-1 {
					p := g.tokens.Pop()
					if p != nil {
						killToken = append(killToken, p.(*lexer.Token))
					}
				}
				//g.Stream.Print()
				//g.PrintToken()
			} else {
				for _, s := range reverseAny(killToken) {
					g.tokens.PushFront(s)
				}
				break
			}
			if j == len(canUse[i])-1 {
				match = true
			}
		}
		if !match && i == len(canUse)-1 {
			break
		}
		if match {
			if canUse[i] == "" {
				return res, true
			}
			return res, true
		}

	}

	g.Stream.ClearFronts(count - 1)
	return res, false
}

// reverse 反转字符串
func reverse(s []string) []string {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

// reverseAny 
func reverseAny(s []*lexer.Token) []*lexer.Token {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}


// PrintToken 打印 token
func (g *Grammar) PrintToken() {
	e := g.tokens.FrontRaw()
	for {
		if e != nil {
			fmt.Print(e.Value.(*lexer.Token).Value)
			e = e.Next()
		} else {
			break
		}
	}
	fmt.Println("$\n ") // 结束符
}
