package util

// IsTerminal 判断是否是终结符 
// 也就是"是否是一个 大写字母 因为终结符是小写字母或其他特殊字符表示的"
func IsTerminal(a byte) bool {
	if a < 'A' || a > 'Z' {
		return true
	}
	return false
}
