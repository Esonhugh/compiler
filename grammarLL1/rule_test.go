package grammarLL1

import (
	"github.com/esonhugh/compiler/grammarLL1/analysisTable"
	"github.com/esonhugh/compiler/grammarLL1/first"
	"github.com/esonhugh/compiler/grammarLL1/follow"
	"github.com/esonhugh/compiler/grammarLL1/rule"
	"fmt"
	"testing"
)

func TestNewRule(t *testing.T) {
	g := rule.NewRules()
	_ = g.AddRules("E->TG\nG->ATG|&\nT->FS\nS->MFS|&\nF->(E)|i\nA->+|-\nM->*|/")
	//fmt.Println(t)
	firstSet := first.GetFirstSet(g)
	fmt.Println(firstSet.String())
	followSet := follow.GetFollowSet(g, "E", firstSet)
	fmt.Println(followSet.String())
	table := analysisTable.GetAnalyzeTable(firstSet, followSet, g)
	res := table.String()
	fmt.Println(res)
}
