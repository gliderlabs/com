package objects

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
)

const (
	TagSingleton = "singleton"
	TagExtpoint  = "extpoint"
	TagConfig    = "config"
)

// Object represents an object and its metadata in a registry
type Object struct {
	Value   interface{}
	Name    string
	Fields  map[string]*Field
	Enabled bool
	PkgPath string

	reflectType  reflect.Type
	reflectValue reflect.Value
}

// New creates a new Object by value and name
func New(v interface{}, name string) *Object {
	return &Object{Value: v, Name: name}
}

// FQN returns a fully qualified name for a unique object in a registry
func (o *Object) FQN() string {
	return strings.ToLower(fmt.Sprintf("%s#%s", o.PkgPath, o.Name))
}

// Assign will set a named field of the object value if it has not already
// been assigned. It will not assign to fields marked as extension points.
// It will return true if the assignment is successful.
func (o *Object) Assign(field string, obj *Object) bool {
	f, ok := o.Fields[field]
	if !ok {
		return false
	}
	// don't assign if already assigned
	if !isNilOrZero(f.reflectValue, f.reflectValue.Type()) {
		return false
	}
	// don't assign to extpoints because as slices they are handled differently
	if !f.Extpoint && obj.reflectType.AssignableTo(f.reflectValue.Type()) {
		f.reflectValue.Set(reflect.ValueOf(obj.Value))
		return true
	}
	return false
}

// Field represents metadata of a field in an Object value's struct.
type Field struct {
	Object   *Object
	Name     string
	Config   bool
	Extpoint bool
	Tag      string

	reflectValue reflect.Value
}

// Registry is a container for objects.
type Registry struct {
	sync.Mutex
	objects  []*Object
	disabled map[string]bool
}

// Register adds objects to the registry.
func (r *Registry) Register(objects ...*Object) error {
	r.Lock()
	defer r.Unlock()
	r.initDisabled()
	for _, o := range objects {
		// set up type and value reflection
		o.reflectType = reflect.TypeOf(o.Value)
		o.reflectValue = reflect.ValueOf(o.Value)

		// if not a struct, ignore
		if !isStructPtr(o.reflectType) {
			continue
		}

		// collect tagged fields
		o.Fields = make(map[string]*Field)
		for i := 0; i < o.reflectValue.Elem().NumField(); i++ {
			field := o.reflectValue.Elem().Field(i)
			fieldName := o.reflectType.Elem().Field(i).Name
			fieldTag, ok := o.reflectType.Elem().Field(i).Tag.Lookup("com")
			if ok && field.CanSet() {
				o.Fields[fieldName] = &Field{
					Name:         fieldName,
					Config:       fieldTag == TagConfig,
					Extpoint:     fieldTag == TagExtpoint,
					Tag:          fieldTag,
					reflectValue: field,
				}
			}
		}

		// set normalized package path. if the package is "com" we assume it
		// contains the component for its parent, so we strip it off.
		o.PkgPath = strings.TrimSuffix(o.reflectType.Elem().PkgPath(), "/com")

		// the default name is set by the name of the struct type
		if o.Name == "" {
			o.Name = o.reflectType.Elem().Name()
		}

		// error if the object has no package path
		if o.Name == "" && o.reflectType.Elem().PkgPath() == "" {
			return errors.New("unable to register object without name when it has no package path")
		}

		// enable unless already marked as disabled
		_, disabled := r.disabled[o.FQN()]
		o.Enabled = !disabled

		// append object to registry list of objects
		r.objects = append(r.objects, o)
	}

	// re-populate registered objects
	return r.reload()
}

// Lookup will attempt to find an object in the registry...
// 1. if it matches the object FQN exactly
// 2. if it matches a single object Name
// 3. if it matches a single object by package path suffix
func (r *Registry) Lookup(name string) (*Object, error) {
	// TODO: allow to choose to ignore disabled
	// TODO: match suffix for full FQN? (pkgpath+name)

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
		if strings.HasSuffix(strings.ToLower(obj.PkgPath), name) {
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

// SetEnabled will set whether an object is enabled.
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

// Enabled returns all enabled objects.
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

// Objects returns all registered objects.
func (r *Registry) Objects() []*Object {
	r.Lock()
	defer r.Unlock()
	var objects []*Object
	for _, o := range r.objects {
		objects = append(objects, o)
	}
	return objects
}

// Reload will go over all objects in the registry and attempt to populate
// fields with com struct tags with other objects in the registry.
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
