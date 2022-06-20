/*
Package print 输出解析包 确认处理情况的输出函数
*/
package print

import (
	"github.com/esonhugh/compiler/grammar"
	"github.com/esonhugh/compiler/grammarLL1"
	"github.com/esonhugh/compiler/lexer"
	"github.com/esonhugh/compiler/util/transfer"
	"fmt"
	"github.com/gookit/color"
	"strings"
)

// PrintGrammar 和 Grammar 库一起使用 用于输出语法分析结果
func PrintGrammar(gram []*grammar.Production) {
	start := "E"
	matched := ""
	for i := 0; i < len(gram); i++ {
		if gram[i].Type == "kill" {
			fmt.Printf("规约：%v\n", transfer.Transfer(gram[i].Target))
			start = strings.Replace(start, gram[i].Target, "", 1)
			matched += gram[i].Target
			color.Green.Printf("%s", transfer.Transfer(matched))
			color.Red.Printf(transfer.Transfer(start + "\n"))
			i++
		} else {
			fmt.Printf("推导：%v--->%v\n", transfer.Transfer(gram[i].Origin), transfer.Transfer(gram[i].Next))
			start = strings.Replace(start, gram[i].Origin, gram[i].Next, 1)
			color.Green.Printf("%s", transfer.Transfer(matched))
			color.Red.Printf(transfer.Transfer(start + "\n"))
		}
	}
}

// PrintGrammarLL1 和 GrammarLL1 库一起使用 用于输出语法分析结果
func PrintGrammarLL1(gram []*grammarLL1.Production) {
	start := "E"
	matched := ""
	for i := 0; i < len(gram)-1; i++ {
		if gram[i].Type == "kill" {
			fmt.Printf("规约：%v\n", transfer.Transfer(gram[i].Target))
			start = strings.Replace(start, gram[i].Target, "", 1)
			matched += gram[i].Target
			color.Green.Printf("%s", transfer.Transfer(matched))
			color.Red.Printf(transfer.Transfer(start + "\n"))
		} else {
			fmt.Printf("推导：%v--->%v\n", transfer.Transfer(gram[i].Origin), transfer.Transfer(gram[i].Next))
			start = strings.Replace(start, gram[i].Origin, gram[i].Next, 1)
			color.Green.Printf("%s", transfer.Transfer(matched))
			color.Red.Printf(transfer.Transfer(start + "\n"))
		}
	}
}

// PrintLexer 和 Lexer 库一起使用 用于输出词法分析结果
func PrintToken(tokens []*lexer.Token) bool {
	err := false
	for i := 0; i < len(tokens); i++ {
		if tokens[i].Typ == lexer.ERROR {
			err = true
		}
		tokens[i].Show()
	}
	fmt.Println("\n ")
	return err
}
