package plugins

import (
	"testing"

	"github.com/gliderlabs/com/objects"
	"github.com/spf13/afero"
)

func mockLoadPlugin(reg *objects.Registry, filepath string) error {
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
	reg := &objects.Registry{}
	err := Load(reg, "test", []string{})
	if err != nil {
		t.Fatal(err)
	}
}
