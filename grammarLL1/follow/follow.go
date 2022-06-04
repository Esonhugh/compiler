package follow

import (
	"github.com/esonhugh/compiler/grammarLL1/first"
	"github.com/esonhugh/compiler/grammarLL1/rule"
	"github.com/esonhugh/compiler/grammarLL1/util"
	"github.com/esonhugh/compiler/util/transfer"
	"fmt"
	"strings"
)

// FollowSet follow集合定义
type FollowSet map[string]map[string]struct{}

// GetFollowSet 获取 follow 集合
// 与 FirstSet 类似 
func GetFollowSet(rule *rule.Rule, start string, firstSet first.FirstSet) FollowSet {
	followSet := make(FollowSet)
	if len(firstSet) == 0 {
		return followSet
	}
	for key := range firstSet {
		followSet[key] = make(map[string]struct{})
	}

	followSet[start]["#"] = struct{}{}

	var changed bool
	for {
		changed = false

		for left, right := range rule.Rules {
			for i := 0; i < len(right); i++ {
				// 对每一个字符及进行遍历
				for index, char := range right[i] { //char B
					if util.IsTerminal(byte(char)) {
						continue
					}
					offset := 1
					for {
						if index+offset == len(right[i]) { // 到末尾了
							// A->bB
							if removeEmptyAndMergeSet(followSet[string(char)], followSet[left]) != 0 {
								changed = true
							}
							break
						} else { // 未到末尾
							if util.IsTerminal(right[i][index+offset]) { // A->Bb
								if mergeSet(followSet[string(char)], map[string]struct{}{string(right[i][index+offset]): {}}) != 0 {
									changed = true
								}
								break
							} else { // A-> BC
								if rule.HaveEmptySet(string(right[i][index+offset])) {
									if removeEmptyAndMergeSet(followSet[string(char)], firstSet[string(right[i][index+offset])]) != 0 {
										changed = true
									}
									offset++
									continue
								} else {
									if removeEmptyAndMergeSet(followSet[string(char)], firstSet[string(right[i][index+offset])]) != 0 {
										changed = true
									}
									break
								}
							}
						}
					}
				}
			}
		}

		if !changed {
			break
		}
	}
	return followSet
}

func removeEmptyAndMergeSet(a map[string]struct{}, b map[string]struct{}) int {
	//delete(b, "&")
	count := 0
	for key, value := range b {
		if key == "&" {
			continue
		}
		if _, ok := a[key]; !ok {
			count++
		}
		a[key] = value
	}
	return count
}

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

func (f FollowSet) String() string {
	var build strings.Builder
	for key, value := range f {
		build.WriteString(fmt.Sprintf("FOLLOW(%s) = { ", transfer.Transfer(key)))
		for item := range value {
			build.WriteString(fmt.Sprintf("%s ", transfer.Transfer(item)))
		}
		build.WriteString("}\n")
	}
	return build.String()
}
