package caco3

import (
	"shanhu.io/text/lexing"
)

const (
	nodeSrc  = "src"
	nodeRule = "rule"
	nodeOut  = "out"
	nodeRun  = "run"
)

type buildNode struct {
	name string
	typ  string
	deps []string
	pos  *lexing.Pos

	ruleType string
	rule     buildRule
	ruleMeta *buildRuleMeta
}

func (n *buildNode) mainOut() string {
	if m := n.ruleMeta; m != nil && len(m.outs) > 0 {
		return m.outs[0]
	}
	return ""
}
