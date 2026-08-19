package cmd

import "testing"

// depinject resolves the modules' expected keeper interfaces at container
// build time, inside NewRootCmd. A signature drift upstream, or an input
// declared as an interface nothing in the container implements, panics here
// and takes every goatd subcommand with it, version included. Neither the
// compiler nor the module tests see it: the module tests run against mocks
// generated from the same interfaces.
func TestNewRootCmdBuildsTheContainer(t *testing.T) {
	if cmd := NewRootCmd(); cmd == nil {
		t.Fatal("NewRootCmd returned nil")
	}
}
