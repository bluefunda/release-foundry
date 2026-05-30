package domain

import "strings"

// FilterConfig holds configurable label and title-prefix rules for PR filtering.
type FilterConfig struct {
	IncludeLabels        map[string]bool
	ExcludeLabels        map[string]bool
	TitlePrefixMap       map[string]string
	TitleExcludePrefixes map[string]bool
}

// DefaultFilterConfig returns the built-in filter rules (matching the package-level vars).
func DefaultFilterConfig() FilterConfig {
	return FilterConfig{
		IncludeLabels:        copyBoolMap(IncludeLabels),
		ExcludeLabels:        copyBoolMap(ExcludeLabels),
		TitlePrefixMap:       copyStringMap(TitlePrefixMap),
		TitleExcludePrefixes: copyBoolMap(TitleExcludePrefixes),
	}
}

// NewFilterConfig builds a FilterConfig from custom include/exclude label lists.
// Title prefix maps use the built-in defaults.
func NewFilterConfig(include, exclude []string) FilterConfig {
	fc := DefaultFilterConfig()
	if len(include) > 0 {
		fc.IncludeLabels = make(map[string]bool, len(include))
		for _, l := range include {
			fc.IncludeLabels[l] = true
		}
	}
	if len(exclude) > 0 {
		fc.ExcludeLabels = make(map[string]bool, len(exclude))
		for _, l := range exclude {
			fc.ExcludeLabels[l] = true
		}
	}
	return fc
}

// IsZero reports whether the FilterConfig has no rules configured.
func (fc FilterConfig) IsZero() bool {
	return len(fc.IncludeLabels) == 0 && len(fc.ExcludeLabels) == 0 &&
		len(fc.TitlePrefixMap) == 0 && len(fc.TitleExcludePrefixes) == 0
}

// InferTypeFromTitleWith extracts a conventional commit prefix using the given prefix map.
func InferTypeFromTitleWith(title string, prefixMap map[string]string) (string, bool) {
	prefix := strings.ToLower(title)
	for _, sep := range []string{"(", ":"} {
		if idx := strings.Index(prefix, sep); idx > 0 {
			prefix = prefix[:idx]
			break
		}
	}
	prefix = strings.TrimSpace(prefix)
	if mapped, ok := prefixMap[prefix]; ok {
		return mapped, true
	}
	return prefix, false
}

// IsExcludedByTitleWith returns true if the title's conventional prefix is in excludePrefixes.
func IsExcludedByTitleWith(title string, excludePrefixes map[string]bool) bool {
	prefix := strings.ToLower(title)
	for _, sep := range []string{"(", ":"} {
		if idx := strings.Index(prefix, sep); idx > 0 {
			prefix = prefix[:idx]
			break
		}
	}
	return excludePrefixes[strings.TrimSpace(prefix)]
}

// ClassifyPRWith determines the primary type using the given FilterConfig.
func ClassifyPRWith(labels []string, title string, fc FilterConfig) string {
	priority := []string{"feature", "fix", "performance", "security", "infrastructure"}
	set := make(map[string]bool, len(labels))
	for _, l := range labels {
		set[l] = true
	}
	for _, p := range priority {
		if set[p] {
			return p
		}
	}
	if t, ok := InferTypeFromTitleWith(title, fc.TitlePrefixMap); ok {
		return t
	}
	return "other"
}

func copyBoolMap(m map[string]bool) map[string]bool {
	c := make(map[string]bool, len(m))
	for k, v := range m {
		c[k] = v
	}
	return c
}

func copyStringMap(m map[string]string) map[string]string {
	c := make(map[string]string, len(m))
	for k, v := range m {
		c[k] = v
	}
	return c
}
