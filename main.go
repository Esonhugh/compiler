package main

import (
	"compiler/grammar"
	"compiler/grammarLL1"
	"compiler/lexer"
	servicePrint "compiler/print"
	"github.com/gookit/color"
	"os"
)

// MakeToken 词法分析 同时输出结果
func MakeToken(code string) []*lexer.Token {
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

func main_proxy(){
	// LL1 Grammar Parser
	//tokens := MakeToken("i*(i-i)/(i+i)*((i+i-i*i/i))")
	//tokens := MakeToken("i+i")
	tokens := MakeToken("i*i**")
	//tokens := MakeToken("i+i*i(")
	//tokens := MakeToken("i+i*i/i-i)")
	//tokens := MakeToken("i+)i-i(")
	//tokens := MakeToken("(i-i)(i/i)")
	//Grammar(tokens)
	GrammarLL1(tokens)
	
	// 递归下降
	Grammar(tokens)
	
	// SysY 分析
	tokens = MakeToken("for(int i=1;i<=10;i++) begin\na_b++;#zszszszszs\nb_C--;#zszszszszs\nB123:=1234567;\na=@;\n123a=0;\na.b;\nend\n")
	servicePrint.PrintToken(tokens)

}