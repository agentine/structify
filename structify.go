package structify

import (
	"errors"
	"fmt"
	"reflect"
)

// ErrNotStruct is returned when a non-struct value is passed.
var ErrNotStruct = errors.New("structify: value is not a struct")

// ErrFieldNotFound is returned when a field name is not found.
var ErrFieldNotFound = errors.New("structify: field not found")

// Struct wraps a struct value and provides methods for inspection and conversion.
type Struct struct {
	raw    any
	value  reflect.Value
	config *config
}

// New creates a new Struct wrapper. The input must be a struct or a pointer to a struct.
// Panics if s is not a struct or pointer to struct.
func New(s any, opts ...Option) *Struct {
	cfg := defaultConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	v := strctVal(s)
	return &Struct{
		raw:    s,
		value:  v,
		config: cfg,
	}
}

// Map converts the struct to a map[string]any.
// Map keys are determined by field names or the configured struct tag.
// Nested structs are recursively converted to maps unless the omitnested tag option is set.
// Fields tagged with "-" are skipped. Fields with omitempty are skipped if zero-valued.
func (s *Struct) Map() map[string]any {
	return structToMap(s.value, s.config.tagName)
}

// Values returns the field values as a slice.
func (s *Struct) Values() []any {
	return structValues(s.value, s.config.tagName)
}

// Names returns the exported field names.
func (s *Struct) Names() []string {
	return structNames(s.value, s.config.tagName)
}

// Fields returns all non-skipped fields.
func (s *Struct) Fields() []*Field {
	return getFields(s.value, s.config.tagName)
}

// Field returns a single field by name.
// Returns ErrFieldNotFound if the field does not exist.
func (s *Struct) Field(name string) (*Field, error) {
	t := s.value.Type()
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		if sf.Name == name {
			fv := s.value.Field(i)
			opts := parseTag(sf.Tag.Get(s.config.tagName))
			return &Field{
				value:   fv,
				field:   sf,
				tagName: s.config.tagName,
				tagOpts: opts,
			}, nil
		}
	}
	return nil, fmt.Errorf("%w: %s", ErrFieldNotFound, name)
}

// IsZero returns true if all exported fields are zero-valued.
func (s *Struct) IsZero() bool {
	return isZero(s.value, s.config.tagName)
}

// HasZero returns true if any exported field is zero-valued.
func (s *Struct) HasZero() bool {
	return hasZero(s.value, s.config.tagName)
}

// Name returns the struct type name.
func (s *Struct) Name() string {
	return s.value.Type().Name()
}

// strctVal resolves a struct value from a value or pointer.
// Panics if s is not a struct or pointer to struct.
func strctVal(s any) reflect.Value {
	v := reflect.ValueOf(s)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		panic("structify: not a struct")
	}
	return v
}
