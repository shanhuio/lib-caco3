package caco3

import (
	"shanhu.io/misc/jsonx"
	"shanhu.io/text/lexing"
)

const buildFileName = "BUILD.caco3"

func makeBuildFileNode(t string) interface{} {
	switch t {
	case ruleFileSet:
		return new(FileSet)
	case ruleBundle:
		return new(Bundle)
	}
	return nil
}

func readBuildFile(env *env, p string) ([]*buildNode, []*lexing.Error) {
	fp := env.src(p, buildFileName)
	rules, errs := jsonx.ReadSeriesFile(fp, makeBuildFileNode)
	if errs != nil {
		return nil, errs
	}

	var nodes []*buildNode

	errList := lexing.NewErrorList()

	for _, r := range rules {
		node := &buildNode{
			typ:      nodeRule,
			pos:      r.Pos,
			ruleType: r.Type,
		}

		switch v := r.V.(type) {
		case *FileSet:
			fset, err := newFileSet(env, p, v)
			if err != nil {
				errList.Add(&lexing.Error{Pos: r.Pos, Err: err})
				continue
			}
			node.rule = fset
		default:
			errList.Errorf(r.Pos, "unknown type: %q", r.Type)
			continue
		}

		if node.rule != nil {
			meta, err := node.rule.meta(env)
			if err != nil {
				errList.Errorf(r.Pos, "fail to get rule meta")
			}
			node.ruleMeta = meta
			node.name = meta.name
			node.deps = meta.deps
		}

		if node.name == p || node.name == "" {
			errList.Errorf(r.Pos, "rule has no name")
			continue
		}
		nodes = append(nodes, node)
	}

	if errs := errList.Errs(); errs != nil {
		return nil, errs
	}
	return nodes, nil
}
