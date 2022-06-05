package main

import (
	"github.com/esonhugh/compiler/grammar"
	"github.com/esonhugh/compiler/grammarLL1"
	"github.com/esonhugh/compiler/lexer"
	servicePrint "github.com/esonhugh/compiler/print"
	"github.com/gookit/color"
	"os"
	"strings"
)

// MakeToken 词法分析 同时输出结果
func MakeToken(code string) []*lexer.Token {
	code = strings.Replace(code, "{", " begin ", -1)
	code = strings.Replace(code, "}", " end ", -1)
	tokens := lexer.Analyse(code)
	err := servicePrint.PrintToken(tokens)
	if err {
		color.Redln("Lexer ERR, please check")
		os.Exit(-1)
	}
	return tokens
}

// Grammar 语法分析 同时输出结果 递归下降 也就是自顶向下解析
func Grammar(tokens []*lexer.Token) {
	gram, correct := grammar.Analyse(tokens)
	if !correct {
		panic("语法推导失败")
	}
	servicePrint.PrintGrammar(gram)
}

// GrammarLL1 语法分析 同时输出结果 LL1 解析
func GrammarLL1(tokens []*lexer.Token) {
	gram, correct := grammarLL1.Analyze(tokens, "E->TG\nG->ATG|&\nT->FS\nS->MFS|&\nF->(E)|i\nA->+|-\nM->*|/", "E")
	if !correct {
		panic("语法推导失败")
	}
	servicePrint.PrintGrammarLL1(gram)
}

// 实验
func main() {
	main_proxy()
}

// main_proxy main 函数代理
func main_proxy() {
	main_proxy_sysy()
	main_proxy_grammar()
	main_proxy_LL1()
}

// main_proxy_LL1 LL11 分析实验 4 and 6
func main_proxy_LL1() {
	// LL1 Grammar Parser 专题 4 6
	//tokens := MakeToken("i*(i-i)/(i+i)*((i+i-i*i/i))")
	tokens := MakeToken("i+i")
	//tokens := MakeToken("i*i**")
	//tokens := MakeToken("i+i*i(")
	//tokens := MakeToken("i+i*i/i-i)")
	//tokens := MakeToken("i+)i-i(")
	//tokens := MakeToken("(i-i)(i/i)")
	//Grammar(tokens)
	GrammarLL1(tokens)
}

// main_proxy_grammar 语法分析器 结果 递归下降专题 3
func main_proxy_grammar() {
	// 递归下降 专题3
	// tokens := MakeToken("i*(i-i)/(i+i)*((i+i-i*i/i))")
	tokens := MakeToken("i*i+i-i")
	// tokens := MakeToken("i**i++i--i")
	//tokens := MakeToken("i+i")
	//tokens := MakeToken("i*i**")
	//tokens := MakeToken("i+i*i(")
	//tokens := MakeToken("i+i*i/i-i)")
	//tokens := MakeToken("i+)i-i(")
	//tokens := MakeToken("(i-i)(i/i)")
	Grammar(tokens)
}

// main_sysy SysY 输入文件
var main_sysy = `
int a;
int main() {
a = 10;
b = 0x1fff;
c = 027;
// d = 20;
/* e= 2; */
if (a >< 0) {
return a;
} else {
return 0;
}

for(int i=1;i<=10;i++) {
b_C--; // comment
B123 = 1234567; 
a = 1;
int c=12;
}
}
`

// main_proxy_sysy SysY 语言词法分析
// 如果需要文件读入 使用 lexer.FromFile 进行读取
func main_proxy_sysy() {
	// SysY 分析 专题1
	tokens := MakeToken(main_sysy)
	servicePrint.PrintToken(tokens)
	// Grammar(tokens)
}
