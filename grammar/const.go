/*
Package grammar 一个简单语法分析器包 

*/
package grammar

// Procduction 产生式语法定义
type Production struct {
	Type   string
	Target string
	Origin string
	Next   string
}
