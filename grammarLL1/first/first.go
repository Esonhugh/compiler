package first

import (
	"github.com/esonhugh/compiler/grammarLL1/rule"
	"github.com/esonhugh/compiler/grammarLL1/util"
	"github.com/esonhugh/compiler/util/transfer"
	"fmt"
	"strings"
)

// FirstSet 存储 FIRST 集
// 一个二维度的表格 存储 FIRST 集
type FirstSet map[string]map[string]struct{}

// GetFirstSet 根据规则构建 FIRST 集
func GetFirstSet(rules *rule.Rule) FirstSet {
	firstSet := make(FirstSet)
	var changed bool
	for {
		changed = false
		for key, r := range rules.Rules {
			// key 左值 r推导值
			if firstSet[key] == nil {
				firstSet[key] = make(map[string]struct{})
			}
			for _, v := range r {
				// 遍历产生式
				// 第一个是终结符,直接将终结符加进first集
				if util.IsTerminal(v[0]) {
					// 直接合并
					if mergeSet(firstSet[key], map[string]struct{}{string(v[0]): {}}) != 0 {
						changed = true
					}
					continue
				} else {
					// 第一个是非终结符 去空 合并
					if removeEmptyAndMergeSet(firstSet[key], firstSet[string(v[0])]) != 0 {
						changed = true
					}
				}
			}
		}
		if !changed {
			break
		}
	}
	return firstSet
}

// removeEmptyAndMergeSet
// 去掉空终结符 合并 a 和 b 的 map
// 返回合并变化的个数
func removeEmptyAndMergeSet(a map[string]struct{}, b map[string]struct{}) int {
	flag := false
	// b 集合去空
	if _, flag = b["&"]; flag {
		delete(b, "&")
	}
	count := 0
	for key, value := range b {
		if _, ok := a[key]; !ok {
			count++
		}
		a[key] = value
	}
	if flag {
		b["&"] = struct{}{}
	}
	return count
}

// mergeSet 合并 a 和 b 的 集合
// 返回合并变化的个数
func mergeSet(a map[string]struct{}, b map[string]struct{}) int {
	count := 0
	for key, value := range b {
		if _, ok := a[key]; !ok {
			count++
		}
		a[key] = value
	}
	return count
}

// String 转换为字符串 使得 First 集合可以打印出来
func (f FirstSet) String() string {
	var build strings.Builder
	for key, value := range f {
		build.WriteString(fmt.Sprintf("FIRST(%s) = { ", transfer.Transfer(key)))
		for item := range value {
			build.WriteString(fmt.Sprintf("%s ", transfer.Transfer(item)))
		}
		build.WriteString("}\n")
	}
	return build.String()
}

// haveEmpty 检查 FIRST 集合是否包含空
func (f FirstSet) haveEmpty(first string) bool {
	_, ok := f[first]["&"]
	return ok
}

// IsInFirstSet 检查是否在 FIRST 集合中
func (f FirstSet) IsInFirstSet(first string, target string) bool {
	for key := range f[first] {
		if key == target {
			return true
		}
	}
	return false
}
