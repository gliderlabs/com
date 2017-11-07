package viper

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

func newTestProvider(t *testing.T, path, config string) *Provider {
	fs := afero.NewMemMapFs()
	err := afero.WriteFile(fs, path, []byte(config), 0644)
	if err != nil {
		t.Fatal(err)
	}
	v := viper.New()
	v.SetFs(fs)
	return &Provider{v}
}

func getString(v interface{}) string {
	vv, _ := v.(string)
	return vv
}

func fatal(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func TestLoadToml(t *testing.T) {
	provider := newTestProvider(t, "/etc/test.toml", `
[Test]
foo = "foobar"
`)
	settings, err := provider.Load("test", []string{"/etc"})
	fatal(t, err)
	if got := getString(settings.Get("Test.foo")); got != "foobar" {
		t.Fatalf("got %#v; want %#v", got, "foobar")
	}
}

func TestLoadYaml(t *testing.T) {
	provider := newTestProvider(t, "/etc/test.yaml", `
Test:
  foo: "foobar"
`)
	settings, err := provider.Load("test", []string{"/etc"})
	fatal(t, err)
	if got := getString(settings.Get("Test.foo")); got != "foobar" {
		t.Fatalf("got %#v; want %#v", got, "foobar")
	}
}

func TestLoadJson(t *testing.T) {
	provider := newTestProvider(t, "/etc/test.json", `
{"Test":
  {"foo": "foobar"}
}`)
	settings, err := provider.Load("test", []string{"/etc"})
	fatal(t, err)
	if got := getString(settings.Get("Test.foo")); got != "foobar" {
		t.Fatalf("got %#v; want %#v", got, "foobar")
	}
}

// func TestLoadNotFound(t *testing.T) {
// 	provider := newTestProvider(t, "/var/test.toml", `
// [Test]
// foo = "foobar"
// `)
// 	_, err := provider.Load("test", []string{"/etc", "/tmp"})
// 	if err == nil {
// 		t.Fatal("expected error")
// 	}
// }

func TestNew(t *testing.T) {
	t.Parallel()
	provider := New()
	var m map[string]interface{}
	if err := provider.Unmarshal(&m); err != nil {
		t.Fatal(err)
	}
	if len(m) > 0 {
		t.Fatal("new config provider is not empty")
	}
}

func TestEmpty(t *testing.T) {
	t.Parallel()
	settings := New().New()
	var m map[string]interface{}
	if err := settings.Unmarshal(&m); err != nil {
		t.Fatal(err)
	}
	if len(m) > 0 {
		t.Fatal("empty config settings is not empty")
	}
}

func TestSub(t *testing.T) {
	t.Parallel()
	v := viper.New()
	v.Set("sub", map[string]interface{}{
		"key": "value",
	})
	provider := &Provider{v}
	sub := provider.Sub("sub")
	val := getString(sub.Get("key"))
	if val != "value" {
		t.Fatalf("expected 'value', got '%#v'", val)
	}
}
