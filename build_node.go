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

	outs     []string
	outImage string

	ruleType string
	rule     buildRule
	ruleMeta *buildRuleMeta
}

func (n *buildNode) mainOut() string {
	if len(n.outs) > 0 {
		return n.outs[0]
	}
	return ""
}
