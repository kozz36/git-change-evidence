// Package gotreesitterprobe characterizes a pinned experimental grammar only.
// It is not a census extractor or a Python-validity oracle.
package gotreesitterprobe

import (
	"bytes"
	"fmt"

	ts "github.com/odvcencio/gotreesitter"
	"github.com/odvcencio/gotreesitter/grammars"
)

type span struct{ Start, End int }

func parse(source []byte) (*ts.Tree, *ts.Language, error) {
	lang := grammars.PythonLanguage()
	if lang == nil {
		return nil, nil, fmt.Errorf("Python grammar unavailable")
	}
	tree, err := ts.NewParser(lang).ParseStrict(source)
	if err != nil {
		return nil, nil, err
	}
	if tree == nil || tree.RootNode() == nil || tree.ParseStoppedEarly() {
		return nil, nil, fmt.Errorf("incomplete parse")
	}
	return tree, lang, nil
}

func parseValid(source []byte) (bool, error) {
	tree, _, err := parse(source)
	if err != nil {
		return false, err
	}
	return !tree.RootNode().HasErrorOrMissing(), nil
}

func directCalls(source []byte) ([]span, error) {
	tree, lang, err := parse(source)
	if err != nil {
		return nil, err
	}
	root := tree.RootNode()
	if root.HasErrorOrMissing() {
		return nil, fmt.Errorf("grammar error or missing node")
	}
	var matches []span
	var walk func(*ts.Node)
	walk = func(n *ts.Node) {
		if n == nil {
			return
		}
		if n.Type(lang) == "call" {
			fn := n.ChildByFieldName("function", lang)
			if fn != nil && fn.Type(lang) == "attribute" {
				recv := fn.ChildByFieldName("object", lang)
				attr := fn.ChildByFieldName("attribute", lang)
				if recv != nil && attr != nil && recv.Type(lang) == "identifier" && attr.Type(lang) == "identifier" &&
					bytes.Equal(source[recv.StartByte():recv.EndByte()], []byte("os")) &&
					bytes.Equal(source[attr.StartByte():attr.EndByte()], []byte("open")) {
					matches = append(matches, span{int(n.StartByte()), int(n.EndByte())})
				}
			}
		}
		for i := 0; i < n.NamedChildCount(); i++ {
			walk(n.NamedChild(i))
		}
	}
	walk(root)
	return matches, nil
}
