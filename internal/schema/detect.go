package schema

import (
	yaml "github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
)

type GVK struct {
	Group   string
	Version string
	Kind    string
}

// DetectGVK extracts apiVersion+kind from a document body. Returns ok=false
// when the document isn't a recognizable Kubernetes manifest.
func DetectGVK(body ast.Node) (GVK, bool) {
	if body == nil {
		return GVK{}, false
	}
	var head struct {
		APIVersion string `yaml:"apiVersion"`
		Kind       string `yaml:"kind"`
	}
	if err := yaml.NodeToValue(body, &head); err != nil {
		return GVK{}, false
	}
	if head.APIVersion == "" || head.Kind == "" {
		return GVK{}, false
	}
	group, version := splitOnce(head.APIVersion, '/')
	if version == "" {
		version = group
		group = ""
	}
	return GVK{Group: group, Version: version, Kind: head.Kind}, true
}

func splitOnce(s string, sep byte) (string, string) {
	for i := 0; i < len(s); i++ {
		if s[i] == sep {
			return s[:i], s[i+1:]
		}
	}
	return s, ""
}
