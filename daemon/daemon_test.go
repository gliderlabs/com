package daemon

import (
	"testing"

	"github.com/gliderlabs/com"
)

func TestEmptyRun(t *testing.T) {
	registry := &com.Registry{}
	err := Run(registry, "test") // normally blocks
	if err != ErrNoServices {
		t.Fatalf("got %#v; want %#v", err, ErrNoServices)
	}
}

func TestInitializers(t *testing.T) {
	t.Skip("TODO")
}

func TestTerminators(t *testing.T) {
	t.Skip("TODO")
}

func TestServices(t *testing.T) {
	t.Skip("TODO")
}
