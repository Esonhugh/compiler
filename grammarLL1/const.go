/*
Package grammarLL1 一个 LL(1) 语法分析器包
*/
package grammarLL1


// Production 产生式子
type Production struct {
	Type   string
	Target string
	Origin string
	Next   string
}
