// Copyright 2024 BlueFunda, Inc.
// SPDX-License-Identifier: Apache-2.0

// Package renderers provides pluggable output renderers for release summary data.
// Each renderer transforms domain.WeeklySummary or domain.BatchSummary into a
// specific output format (Markdown, JSON, etc.).
//
// Registering a new renderer:
//
//	func init() { Register("my-format", myRenderer{}) }
package renderers

import (
	"sort"

	"github.com/release-foundry/internal/domain"
)

// Renderer transforms release data into a specific output format.
type Renderer interface {
	// Single renders a single-repository summary.
	Single(summary domain.WeeklySummary) string
	// Batch renders a multi-repository summary.
	Batch(batch domain.BatchSummary) string
	// FileExtension returns the output file extension without leading dot (e.g. "md").
	FileExtension() string
}

var registry = map[string]Renderer{}

// Register adds a renderer under the given name. Panics on duplicate names.
func Register(name string, r Renderer) {
	if _, dup := registry[name]; dup {
		panic("renderers: duplicate registration for " + name)
	}
	registry[name] = r
}

// Get returns the named renderer. ok is false if the name is not registered.
func Get(name string) (Renderer, bool) {
	r, ok := registry[name]
	return r, ok
}

// Names returns the sorted list of registered renderer names.
func Names() []string {
	names := make([]string, 0, len(registry))
	for k := range registry {
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}
