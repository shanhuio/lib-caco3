package caco3

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"shanhu.io/misc/errcode"
)

// buildAction is a structure for creating the digest of the execution of a
// rule.
type buildAction struct {
	Rule     string `json:",omitempty"`
	RuleType string `json:",omitempty"`
	Deps     map[string]string
}

func makeRuleDigest(t, name string, v interface{}) (string, error) {
	buf := new(bytes.Buffer)
	fmt.Fprintln(buf, t)
	fmt.Fprintln(buf, name)
	bs, err := json.Marshal(v)
	if err != nil {
		return "", errcode.Annotate(err, "json marshal")
	}
	buf.Write(bs)
	sum := sha256.Sum256(buf.Bytes())
	return "sha256:" + hex.EncodeToString(sum[:]), nil
}
