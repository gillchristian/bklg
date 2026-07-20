package main

import (
	"slices"
	"testing"
)

func TestSplitArgs(t *testing.T) {
	tests := []struct {
		name     string
		argv     []string
		path     string
		flagArgs []string
		extra    []string
	}{
		{"no args -> default path", nil, ".", nil, nil},
		{"path only", []string{"."}, ".", nil, nil},
		{"flag then path", []string{"--port", "9001", "."}, ".", []string{"--port", "9001"}, nil},
		{"path then flag", []string{".", "--port", "9001"}, ".", []string{"--port", "9001"}, nil},
		{"equals form", []string{"--port=9001", "."}, ".", []string{"--port=9001"}, nil},
		{"positional before flag keeps its value", []string{"systems/trail", "--port", "9000"}, "systems/trail", []string{"--port", "9000"}, nil},
		{"single-dash flag", []string{"-port", "9001", "."}, ".", []string{"-port", "9001"}, nil},
		{"dir flag with value", []string{"--dir", "knowledge", "."}, ".", []string{"--dir", "knowledge"}, nil},
		{"value-flag at end without value stays in flagArgs", []string{"--port"}, ".", []string{"--port"}, nil},
		{"extra positional captured (not silently dropped)", []string{"a", "b", "--port", "9000"}, "a", []string{"--port", "9000"}, []string{"b"}},
		{"multiple extra positionals", []string{"a", "b", "c"}, "a", nil, []string{"b", "c"}},
		{"boolean-ish flag with equals is not treated as taking next token", []string{"--dir=x", "repo"}, "repo", []string{"--dir=x"}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, flagArgs, extra := splitArgs(tt.argv)
			if path != tt.path {
				t.Errorf("path = %q, want %q", path, tt.path)
			}
			if !slices.Equal(flagArgs, tt.flagArgs) {
				t.Errorf("flagArgs = %v, want %v", flagArgs, tt.flagArgs)
			}
			if !slices.Equal(extra, tt.extra) {
				t.Errorf("extra = %v, want %v", extra, tt.extra)
			}
		})
	}
}
