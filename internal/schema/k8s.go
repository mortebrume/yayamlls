package schema

import "strings"

// DefaultK8sSchemaURL is the URL template applied when no
// kubernetes.schemaUrl override is configured. It targets the
// home-operations CRD mirror — core resources aren't hosted there and
// are silently skipped via the negative cache.
const DefaultK8sSchemaURL = "https://k8s-schemas.home-operations.com/" +
	"{groupSeg}{kindLower}_{versionLower}.json"

// BuildK8sURL renders a URL template against a GVK. When template is
// empty, DefaultK8sSchemaURL is used.
//
// Placeholders:
//
//	{group}         full api group, "" for core
//	{groupSeg}      "<group>/" when non-empty, "" otherwise
//	{groupFirst}    first DNS label of the group ("apps" from "apps.k8s.io")
//	{groupFirstSeg} "<groupFirst>-" when non-empty, "" otherwise
//	{kind}          kind in original case
//	{kindLower}     kind lowercased
//	{version}       version (e.g. "v2", "v1beta1")
//	{versionLower}  version lowercased
func BuildK8sURL(template, group, version, kind string) string {
	if version == "" || kind == "" {
		return ""
	}
	if template == "" {
		template = DefaultK8sSchemaURL
	}
	groupFirst, _, _ := strings.Cut(group, ".")
	groupSeg := ""
	if group != "" {
		groupSeg = group + "/"
	}
	groupFirstSeg := ""
	if groupFirst != "" {
		groupFirstSeg = groupFirst + "-"
	}
	r := strings.NewReplacer(
		"{group}", group,
		"{groupSeg}", groupSeg,
		"{groupFirst}", groupFirst,
		"{groupFirstSeg}", groupFirstSeg,
		"{kind}", kind,
		"{kindLower}", strings.ToLower(kind),
		"{version}", version,
		"{versionLower}", strings.ToLower(version),
	)
	return r.Replace(template)
}
