package main

import (
	"slices"
	"testing"
)

func TestBuildToolArgs(t *testing.T) {
	args := []string{"test"}
	tools := Tools{Allow: []string{"a", "b"}, Deny: []string{"c"}}

	expected := []string{"test", "--allow-tool", "a", "--allow-tool", "b", "--deny-tool", "c"}
	actual := buildToolArgs(args, tools)

	if !slices.Equal(expected, actual) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}
