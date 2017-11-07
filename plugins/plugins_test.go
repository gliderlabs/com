package plugins

import (
	"testing"

	"github.com/gliderlabs/com/registry"
	"github.com/spf13/afero"
)

func mockLoadPlugin(reg *registry.Registry, filepath string) error {
	return nil
}

func setupMocks() {
	pluginLoader = mockLoadPlugin
	fs = afero.NewMemMapFs()
}

func reset() {
	pluginLoader = loadPlugin
	fs = afero.NewOsFs()
}

func TestLoadTODO(t *testing.T) {
	setupMocks()
	defer reset()
	reg := &registry.Registry{}
	err := Load(reg, "test", []string{})
	if err != nil {
		t.Fatal(err)
	}
}
