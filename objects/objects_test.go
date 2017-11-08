package objects

import (
	"testing"
)

type Stringer interface {
	String() string
}

type Foo struct {
	msg string
}

func (f *Foo) String() string {
	return f.msg
}

func TestRegister(t *testing.T) {
	r := &Registry{}
	v := &Foo{t.Name()}
	if err := r.Register(&Object{Value: v}); err != nil {
		t.Fatal(err)
	}
	found := false
	for _, obj := range r.Objects() {
		if obj.Value == v {
			found = true
		}
	}
	if !found {
		t.Fatal("registered object not found")
	}
}

func TestLookup(t *testing.T) {
	r := &Registry{}
	v := &Foo{t.Name()}
	if err := r.Register(&Object{Name: t.Name(), Value: v}); err != nil {
		t.Fatal(err)
	}
	obj, _ := r.Lookup(t.Name())
	if obj == nil || obj.Value != v {
		t.Fatal("lookup return nil or wrong object")
	}
}

func TestStructSingleton(t *testing.T) {
	r := &Registry{}
	v1 := &Foo{t.Name()}
	var v2 struct {
		A *Foo `com:"singleton"`
	}
	if err := r.Register(&Object{Value: v1}, &Object{Value: &v2, Name: "v2"}); err != nil {
		t.Fatal(err)
	}
	if v2.A != v1 {
		t.Fatal("field not set to registered singleton after register")
	}
}

func TestInterfaceSingleton(t *testing.T) {
	r := &Registry{}
	v1 := &Foo{t.Name()}
	var v2 struct {
		A Stringer `com:"singleton"`
	}
	if err := r.Register(&Object{Value: v1}, &Object{Value: &v2, Name: "v2"}); err != nil {
		t.Fatal(err)
	}
	if v2.A != v1 {
		t.Fatal("field not set to registered singleton after register")
	}
}

func TestExtpoints(t *testing.T) {
	r := &Registry{}
	ext1 := &Foo{"ext1"}
	ext2 := &Foo{"ext2"}
	var v struct {
		A []Stringer `com:"extpoint"`
	}
	if err := r.Register(&Object{Value: &v, Name: "v"}, &Object{Value: ext1}, &Object{Value: ext2}); err != nil {
		t.Fatal(err)
	}
	if len(v.A) != 2 {
		t.Fatal("field not set to registered extensions after register")
	}
}

func TestSkipAssignedSingletons(t *testing.T) {
	r := &Registry{}
	v1 := &Foo{"registered"}
	var v2 struct {
		A *Foo `com:"singleton"`
	}
	v2.A = &Foo{"notregistered"}
	if err := r.Register(&Object{Value: v1}, &Object{Value: &v2, Name: "v2"}); err != nil {
		t.Fatal(err)
	}
	if v2.A.String() != "notregistered" {
		t.Fatal("field was assigned even though it was already set")
	}
}

func TestAssigningNonexistantFieldsNoop(t *testing.T) {
	var v struct {
		A *Foo `com:"singleton"`
	}
	obj := &Object{Value: &v}
	if obj.Assign("B", &Object{Value: &Foo{}}) != false {
		t.Fatal("assign allowed for non-existent field")
	}
}

func TestAssigningExtpointFieldsNoop(t *testing.T) {
	var v struct {
		A []Stringer `com:"extpoint"`
	}
	obj := &Object{Value: &v}
	if obj.Assign("A", &Object{Value: &Foo{}}) != false {
		t.Fatal("assign allowed for extpoint field")
	}
}

func TestConfigFieldLeftUnassigned(t *testing.T) {
	r := &Registry{}
	v1 := &Foo{t.Name()}
	var v2 struct {
		A Stringer `com:"config"`
	}
	if err := r.Register(&Object{Value: v1}, &Object{Value: &v2, Name: "v2"}); err != nil {
		t.Fatal(err)
	}
	if v2.A != nil {
		t.Fatal("config field not left unassigned")
	}
}
