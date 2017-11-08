package config_test

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/gliderlabs/com/config"
	"github.com/gliderlabs/com/config/viper"
	"github.com/gliderlabs/com/objects"
	"github.com/spf13/afero"
	viperlib "github.com/spf13/viper"
)

var (
	initErr = errors.New("initerr is true")
)

type TestComponent struct {
	Foo string
}

func (c *TestComponent) InitializeConfig(cfg config.Settings) error {
	if getBool(cfg.Get("initerr")) {
		return initErr
	}
	var keys map[string]interface{}
	cfg.Unmarshal(&keys)
	return cfg.Unmarshal(c)
}

func getBool(v interface{}) bool {
	vv, _ := v.(bool)
	return vv
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
	reg := &objects.Registry{}
	obj := &TestComponent{}
	reg.Register(&objects.Object{Value: obj})
	provider := newTestProvider(t, "/etc/test.toml", `
[TestComponent]
foo = "foobar"
`)
	err := config.Load(reg, provider, "test", []string{"/etc"})
	fatal(t, err)
	if obj.Foo != "foobar" {
		t.Fatalf("got %#v; want %#v", obj.Foo, "foobar")
	}
}

func TestLoadYaml(t *testing.T) {
	reg := &objects.Registry{}
	obj := &TestComponent{}
	reg.Register(&objects.Object{Value: obj})
	provider := newTestProvider(t, "/etc/test.yaml", `
TestComponent:
  foo: "foobar"
`)
	err := config.Load(reg, provider, "test", []string{"/etc"})
	fatal(t, err)
	if obj.Foo != "foobar" {
		t.Fatalf("got %#v; want %#v", obj.Foo, "foobar")
	}
}

func TestLoadJson(t *testing.T) {
	reg := &objects.Registry{}
	obj := &TestComponent{}
	reg.Register(&objects.Object{Value: obj})
	provider := newTestProvider(t, "/etc/test.json", `
{"TestComponent":
  {"foo": "foobar"}
}`)
	err := config.Load(reg, provider, "test", []string{"/etc"})
	fatal(t, err)
	if obj.Foo != "foobar" {
		t.Fatalf("got %#v; want %#v", obj.Foo, "foobar")
	}
}

func TestLoadFromWorkingDir(t *testing.T) {
	reg := &objects.Registry{}
	wd, err := os.Getwd()
	fatal(t, err)
	// virtual FS will use real cwd since no way fake cwd resolution later
	provider := newTestProvider(t, fmt.Sprintf("%s/test.toml", wd), `
[TestComponent]
foo = "foobar"
`)
	err = config.Load(reg, provider, "test", []string{"."})
	fatal(t, err)
}

func TestLoadInvalidFile(t *testing.T) {
	reg := &objects.Registry{}
	provider := newTestProvider(t, "/etc/test.toml", `
#!/usr/bin/python
print "Hello world"
`)
	err := config.Load(reg, provider, "test", []string{"/etc"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestInitializerError(t *testing.T) {
	reg := &objects.Registry{}
	obj := &TestComponent{}
	reg.Register(&objects.Object{Value: obj})
	provider := newTestProvider(t, "/etc/test.toml", `
[TestComponent]
initerr = true
`)
	err := config.Load(reg, provider, "test", []string{"/etc"})
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
	reg := &objects.Registry{}
	reg.Register(&objects.Object{Value: &c, Name: "Component"})
	reg.Register(&objects.Object{Value: &stringer{"Foo"}, Name: "Fooer"})
	provider := newTestProvider(t, "/etc/test.toml", `
[Component]
Stringer = "Fooer"
`)
	err := config.Load(reg, provider, "test", []string{"/etc"})
	fatal(t, err)
	if got := c.Stringer.String(); got != "Foo" {
		t.Fatalf("got %#v; want %#v", got, "Foo")
	}
}

func TestConfigFieldNoObject(t *testing.T) {
	var c struct {
		Stringer fmt.Stringer `com:"config"`
	}
	reg := &objects.Registry{}
	reg.Register(&objects.Object{Value: &c, Name: "Component"})
	provider := newTestProvider(t, "/etc/test.toml", `
[Component]
Stringer = "Fooer"
`)
	err := config.Load(reg, provider, "test", []string{"/etc"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestDisabled(t *testing.T) {
	reg := &objects.Registry{}
	obj := &objects.Object{Value: &TestComponent{}}
	reg.Register(obj)
	provider := newTestProvider(t, "/etc/test.toml", `
[TestComponent]
foo = "foobar"

[disabled]
TestComponent = true
IgnoreMe = true
`)
	err := config.Load(reg, provider, "test", []string{"/etc"})
	fatal(t, err)
	if got := obj.Enabled; got != false {
		t.Fatalf("got %#v; want %#v", got, false)
	}
}

func TestEnvOverride(t *testing.T) {
	reg := &objects.Registry{}
	obj := &TestComponent{}
	reg.Register(&objects.Object{Value: obj})
	provider := newTestProvider(t, "/etc/test.toml", `
[TestComponent]
foo = "foobar"
`)
	os.Setenv("TESTCOMPONENT_FOO", "bazqux")
	err := config.Load(reg, provider, "test", []string{"/etc"})
	fatal(t, err)
	var keys map[string]interface{}
	provider.Unmarshal(&keys)
	if obj.Foo != "bazqux" {
		t.Fatalf("got %#v; want %#v", obj.Foo, "bazqux")
	}
}
