/*
Package transfer 产生式的替换 指 E' 等双字符 换为单字符

*/
package transfer

import "strings"

func Transfer(str string) string {
	//str = strings.ReplaceAll(str, " ", "")
	str = strings.ReplaceAll(str, "G", "E'")
	str = strings.ReplaceAll(str, "S", "T'")
	str = strings.ReplaceAll(str, "&", "ε")
	return str
}
