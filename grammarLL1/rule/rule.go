package rule

import (
	"errors"
	"strings"
)

// 规则合集
type Rule struct {
	Rules map[string][]string
}

// 表达式 分为左右两边
type Formula struct {
	Left  string
	Right string
}

// 新建一个空的 规则
func NewRules() *Rule {
	return &Rule{Rules: make(map[string][]string)}
}

// 添加规则到规则集和中
func (r *Rule) AddRules(s string) error {
	lineRule := strings.Split(s, "\n")
	for _, t := range lineRule {
		if strings.TrimSpace(t) == "" {
			continue
		}
		c := strings.Split(t, "->")
		if len(c) != 2 || len(c[0]) != 1 {
			return errors.New("invalid arg")
		}
		right := strings.Split(strings.ReplaceAll(c[1], " ", ""), "|")
		for i := range right {
			r.Rules[c[0]] = append(r.Rules[c[0]], right[i])
		}
	}
	return nil
}

// 是否有空产生式
func (r *Rule) HaveEmptySet(first string) bool {
	for _, value := range r.Rules[first] {
		if value == "&" {
			return true
		}
	}
	return false
}

func (r *Rule) TheFirstItemIs(first, item string) string {
	for _, value := range r.Rules[first] {
		if value == item {
			return value
		}
	}
	return ""
}
 
// GetProcessMethod 获取处理方法 
// 如果遇到了
func (r *Rule) GetProcessMethod(first, end string) *Formula {
	for _, value := range r.Rules[first] {
		if value[0] == end[0] {
			return &Formula{
				Left:  first,
				Right: value,
			}
		}
	}
	return nil
}
