package config_test

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/gliderlabs/com"
	"github.com/gliderlabs/com/config"
	"github.com/gliderlabs/com/config/viper"
	"github.com/spf13/afero"
	viperlib "github.com/spf13/viper"
)

var (
	initErr = errors.New("initerr is true")
)

type TestComponent struct {
	SettingFoo string `mapstructure:"foo"`
	SettingBaz string `mapstructure:"baz"`
}

func (c *TestComponent) InitializeConfig(cfg config.Settings) error {
	if cfg.GetBool("initerr") {
		return initErr
	}
	var keys map[string]interface{}
	cfg.Unmarshal(&keys)
	return cfg.Unmarshal(c)
}

func newTestProvider(t *testing.T, path, config string) config.Provider {
	fs := afero.NewMemMapFs()
	err := afero.WriteFile(fs, path, []byte(config), 0644)
	if err != nil {
		t.Fatal(err)
	}
	v := viperlib.New()
	v.SetFs(fs)
	return &viper.Provider{v}
}

func fatal(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func TestLoadToml(t *testing.T) {
	registry := &com.Registry{}
	obj := &TestComponent{}
	registry.Register(&com.Object{Value: obj})
	provider := newTestProvider(t, "/etc/test.toml", `
[TestComponent]
foo = "foobar"
`)
	err := config.Load(registry, provider, "test", []string{"/etc"})
	fatal(t, err)
	if obj.SettingFoo != "foobar" {
		t.Fatalf("got %#v; want %#v", obj.SettingFoo, "foobar")
	}
}

func TestLoadYaml(t *testing.T) {
	registry := &com.Registry{}
	obj := &TestComponent{}
	registry.Register(&com.Object{Value: obj})
	provider := newTestProvider(t, "/etc/test.yaml", `
TestComponent:
  foo: "foobar"
`)
	err := config.Load(registry, provider, "test", []string{"/etc"})
	fatal(t, err)
	if obj.SettingFoo != "foobar" {
		t.Fatalf("got %#v; want %#v", obj.SettingFoo, "foobar")
	}
}

func TestLoadJson(t *testing.T) {
	registry := &com.Registry{}
	obj := &TestComponent{}
	registry.Register(&com.Object{Value: obj})
	provider := newTestProvider(t, "/etc/test.json", `
{"TestComponent":
  {"foo": "foobar"}
}`)
	err := config.Load(registry, provider, "test", []string{"/etc"})
	fatal(t, err)
	if obj.SettingFoo != "foobar" {
		t.Fatalf("got %#v; want %#v", obj.SettingFoo, "foobar")
	}
}

func TestLoadFromWorkingDir(t *testing.T) {
	registry := &com.Registry{}
	wd, err := os.Getwd()
	fatal(t, err)
	// virtual FS will use real cwd since no way fake cwd resolution later
	provider := newTestProvider(t, fmt.Sprintf("%s/test.toml", wd), `
[TestComponent]
foo = "foobar"
`)
	err = config.Load(registry, provider, "test", []string{"."})
	fatal(t, err)
}

func TestLoadNotFound(t *testing.T) {
	registry := &com.Registry{}
	provider := newTestProvider(t, "/var/test.toml", `
[TestComponent]
foo = "foobar"
`)
	err := config.Load(registry, provider, "test", []string{"/etc", "/tmp"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestLoadInvalidFile(t *testing.T) {
	registry := &com.Registry{}
	provider := newTestProvider(t, "/etc/test.toml", `
#!/usr/bin/python
print "Hello world"
`)
	err := config.Load(registry, provider, "test", []string{"/etc"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestInitializerError(t *testing.T) {
	registry := &com.Registry{}
	obj := &TestComponent{}
	registry.Register(&com.Object{Value: obj})
	provider := newTestProvider(t, "/etc/test.toml", `
[TestComponent]
initerr = true
`)
	err := config.Load(registry, provider, "test", []string{"/etc"})
	if err != initErr {
		t.Fatalf("got %#v; want %#v", err, initErr)
	}
}

type stringer struct {
	s string
}

func (s stringer) String() string {
	return s.s
}

func TestConfigField(t *testing.T) {
	var c struct {
		Stringer fmt.Stringer `com:"config"`
	}
	registry := &com.Registry{}
	registry.Register(&com.Object{Value: &c, Name: "Component"})
	registry.Register(&com.Object{Value: &stringer{"Foo"}, Name: "Fooer"})
	provider := newTestProvider(t, "/etc/test.toml", `
[Component]
Stringer = "Fooer"
`)
	err := config.Load(registry, provider, "test", []string{"/etc"})
	fatal(t, err)
	if got := c.Stringer.String(); got != "Foo" {
		t.Fatalf("got %#v; want %#v", got, "Foo")
	}
}

func TestConfigFieldNoObject(t *testing.T) {
	var c struct {
		Stringer fmt.Stringer `com:"config"`
	}
	registry := &com.Registry{}
	registry.Register(&com.Object{Value: &c, Name: "Component"})
	provider := newTestProvider(t, "/etc/test.toml", `
[Component]
Stringer = "Fooer"
`)
	err := config.Load(registry, provider, "test", []string{"/etc"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestDisabled(t *testing.T) {
	registry := &com.Registry{}
	obj := &com.Object{Value: &TestComponent{}}
	registry.Register(obj)
	provider := newTestProvider(t, "/etc/test.toml", `
[TestComponent]
foo = "foobar"

[disabled]
TestComponent = true
IgnoreMe = true
`)
	err := config.Load(registry, provider, "test", []string{"/etc"})
	fatal(t, err)
	if got := obj.Enabled; got != false {
		t.Fatalf("got %#v; want %#v", got, false)
	}
}

func TestEnvOverride(t *testing.T) {
	registry := &com.Registry{}
	obj := &TestComponent{}
	registry.Register(&com.Object{Value: obj})
	provider := newTestProvider(t, "/etc/test.toml", `
[TestComponent]
foo = "foobar"
`)
	os.Setenv("TESTCOMPONENT_FOO", "bazqux")
	err := config.Load(registry, provider, "test", []string{"/etc"})
	fatal(t, err)
	var keys map[string]interface{}
	provider.Unmarshal(&keys)
	if obj.SettingFoo != "bazqux" {
		t.Fatalf("got %#v; want %#v", obj.SettingFoo, "bazqux")
	}
}
