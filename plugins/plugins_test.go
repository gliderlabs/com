package plugins

import (
	"testing"

	"github.com/gliderlabs/com"
	"github.com/spf13/afero"
)

func mockLoadPlugin(registry *com.Registry, filepath string) error {
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
	registry := &com.Registry{}
	err := Load(registry, "test", []string{})
	if err != nil {
		t.Fatal(err)
	}
}
