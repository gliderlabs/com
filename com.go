package com

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
)

const (
	tagSingleton = "singleton"
	tagExtpoint  = "extpoint"
	tagConfig    = "config"
)

var (
	DefaultRegistry = &Registry{}
)

// An object in the registry
type Object struct {
	Value   interface{}
	Name    string            // Optional
	Fields  map[string]*Field // Populated with the field names that were injected and their corresponding *Object.
	Enabled bool

	reflectType  reflect.Type
	reflectValue reflect.Value
	pkgPath      string
}

func (o *Object) FQN() string {
	return strings.ToLower(fmt.Sprintf("%s#%s", o.pkgPath, o.Name))
}

func (o *Object) String() string {
	if !o.Enabled {
		return fmt.Sprintf("%s[disabled]", o.FQN())
	}
	var fields []string
	for k, _ := range o.Fields {
		fields = append(fields, k)
	}
	return fmt.Sprintf("%s[%s]", o.FQN(), strings.Join(fields, " "))
}

func (o *Object) Assign(field string, obj *Object) bool {
	f, ok := o.Fields[field]
	if !ok {
		return false
	}
	// don't assign if already assigned
	if !isNilOrZero(f.reflectValue, f.reflectValue.Type()) {
		return false
	}
	// don't assign to extpoints becasue as slices they are handled differently
	if !f.Extpoint && obj.reflectType.AssignableTo(f.reflectValue.Type()) {
		f.reflectValue.Set(reflect.ValueOf(obj.Value))
		return true
	}
	return false
}

type Field struct {
	Object   *Object
	Name     string
	Config   bool
	Extpoint bool

	comTag       string
	reflectValue reflect.Value
}

type Registry struct {
	sync.Mutex
	objects  []*Object
	disabled map[string]bool
}

// Register objects. The Object documentation describes
// the impact of various fields.
func (r *Registry) Register(objects ...*Object) error {
	r.Lock()
	defer r.Unlock()
	r.initDisabled()
	for _, o := range objects {
		o.reflectType = reflect.TypeOf(o.Value)
		o.reflectValue = reflect.ValueOf(o.Value)
		if !isStructPtr(o.reflectType) {
			continue
		}
		o.Fields = make(map[string]*Field)
		for i := 0; i < o.reflectValue.Elem().NumField(); i++ {
			field := o.reflectValue.Elem().Field(i)
			fieldName := o.reflectType.Elem().Field(i).Name
			fieldTag, ok := o.reflectType.Elem().Field(i).Tag.Lookup("com")
			if ok && field.CanSet() {
				o.Fields[fieldName] = &Field{
					Name:         fieldName,
					Config:       fieldTag == tagConfig,
					Extpoint:     fieldTag == tagExtpoint,
					comTag:       fieldTag,
					reflectValue: field,
				}
			}
		}
		o.pkgPath = strings.TrimSuffix(o.reflectType.Elem().PkgPath(), "/com")
		if o.Name == "" {
			o.Name = o.reflectType.Elem().Name()
		}
		if o.Name == "" && o.reflectType.Elem().PkgPath() == "" {
			return errors.New("unable to register object without name when it has no package path")
		}
		o.Enabled = true
		if _, disabled := r.disabled[o.FQN()]; disabled {
			o.Enabled = false
		}
		r.objects = append(r.objects, o)
	}
	return r.reload()
}

func (r *Registry) Lookup(name string) (*Object, error) {
	// all matching is done case insensitive
	name = strings.ToLower(name)
	var matches []*Object
	for _, obj := range r.Objects() {
		// first match any exact FQN
		if obj.FQN() == name {
			return obj, nil
		}
		// name matches added to slice
		if strings.ToLower(obj.Name) == name {
			matches = append(matches, obj)
		}
	}
	// if only one matched name, return
	if len(matches) == 1 {
		return matches[0], nil
	}
	// if more than one, error
	if len(matches) > 1 {
		return nil, errors.New("ambiguous name for lookup")
	}
	// now attempt suffix matches
	matches = matches[:0]
	for _, obj := range r.Objects() {
		if strings.HasSuffix(strings.ToLower(obj.pkgPath), name) {
			matches = append(matches, obj)
		}
	}
	if len(matches) == 1 {
		return matches[0], nil
	}
	if len(matches) > 1 {
		return nil, errors.New("ambiguous name for lookup")
	}
	return nil, errors.New("object not found")
}

func (r *Registry) SetEnabled(fqn string, enabled bool) {
	r.Lock()
	defer r.Unlock()
	r.initDisabled()
	r.disabled[fqn] = !enabled
	for _, o := range r.objects {
		if o.FQN() == fqn {
			o.Enabled = enabled
			break
		}
	}
}

func (r *Registry) initDisabled() {
	if r.disabled == nil {
		r.disabled = make(map[string]bool)
	}
}

// Objects returns all enabled objects.
func (r *Registry) Enabled() []*Object {
	r.Lock()
	defer r.Unlock()
	var objects []*Object
	for _, o := range r.objects {
		if o.Enabled {
			objects = append(objects, o)
		}
	}
	return objects
}

// Objects returns all known objects.
func (r *Registry) Objects() []*Object {
	r.Lock()
	defer r.Unlock()
	var objects []*Object
	for _, o := range r.objects {
		objects = append(objects, o)
	}
	return objects
}

func (r *Registry) Reload() error {
	r.Lock()
	defer r.Unlock()
	return r.reload()
}

func (r *Registry) reload() error {
	for _, o := range r.objects {
		if err := r.populateSingletons(o); err != nil {
			return err
		}
		if err := r.populateExtpoints(o); err != nil {
			return err
		}
	}
	return nil
}

func (r *Registry) populateSingletons(o *Object) error {
	for k, f := range o.Fields {
		if f.Config || f.Extpoint {
			continue
		}
		for _, existing := range r.objects {
			if existing.Enabled && existing.reflectType.AssignableTo(f.reflectValue.Type()) {
				o.Assign(k, existing)
				break
			}
		}
	}
	return nil
}

func (r *Registry) populateExtpoints(o *Object) error {
	for _, f := range o.Fields {
		if !f.Extpoint {
			continue
		}
		var objects []reflect.Value
		for _, existing := range r.objects {
			if existing.Enabled && existing.reflectType.AssignableTo(f.reflectValue.Type().Elem()) {
				objects = append(objects, existing.reflectValue)
			}
		}
		f.reflectValue.Set(reflect.MakeSlice(f.reflectValue.Type(), 0, len(objects)))
		for _, obj := range objects {
			f.reflectValue.Set(reflect.Append(f.reflectValue, obj))
		}
	}
	return nil
}

func isStructPtr(t reflect.Type) bool {
	return t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct
}

func isNilOrZero(v reflect.Value, t reflect.Type) bool {
	switch v.Kind() {
	default:
		return reflect.DeepEqual(v.Interface(), reflect.Zero(t).Interface())
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
}
