package structify

import (
	"errors"
	"fmt"
	"reflect"
)

// ErrNotExported is returned when trying to set an unexported field.
var ErrNotExported = errors.New("structify: field is not exported")

// ErrNotSettable is returned when a field cannot be set.
var ErrNotSettable = errors.New("structify: field is not settable")

// ErrTypeMismatch is returned when the value type doesn't match the field type.
var ErrTypeMismatch = errors.New("structify: type mismatch")

// Field represents a single struct field.
type Field struct {
	value    reflect.Value
	field    reflect.StructField
	tagName  string
	tagOpts  tagOptions
}

// Name returns the field name.
func (f *Field) Name() string {
	return f.field.Name
}

// Value returns the field value as any.
func (f *Field) Value() any {
	return f.value.Interface()
}

// Kind returns the reflect.Kind of the field.
func (f *Field) Kind() reflect.Kind {
	return f.value.Kind()
}

// Tag returns the value of the struct tag for the given key.
func (f *Field) Tag(key string) string {
	return f.field.Tag.Get(key)
}

// IsZero returns true if the field value is the zero value for its type.
func (f *Field) IsZero() bool {
	return f.value.IsZero()
}

// IsExported returns true if the field is exported.
func (f *Field) IsExported() bool {
	return f.field.IsExported()
}

// IsEmbedded returns true if the field is an embedded (anonymous) field.
func (f *Field) IsEmbedded() bool {
	return f.field.Anonymous
}

// Set sets the field value. The field must be exported and settable.
// Returns an error if the field is unexported, not settable, or the types don't match.
func (f *Field) Set(val any) error {
	if !f.IsExported() {
		return ErrNotExported
	}
	if !f.value.CanSet() {
		return ErrNotSettable
	}

	v := reflect.ValueOf(val)
	if !v.Type().AssignableTo(f.value.Type()) {
		return fmt.Errorf("%w: cannot assign %s to %s", ErrTypeMismatch, v.Type(), f.value.Type())
	}

	f.value.Set(v)
	return nil
}

// Fields returns the nested struct fields if this field is a struct.
// Returns nil for non-struct fields. Handles nil pointer fields gracefully.
func (f *Field) Fields() []*Field {
	val := f.value
	kind := val.Kind()

	if kind == reflect.Ptr {
		if val.IsNil() {
			return nil
		}
		val = val.Elem()
		kind = val.Kind()
	}

	if kind != reflect.Struct {
		return nil
	}

	return getFields(val, f.tagName)
}

// getFields extracts fields from a struct reflect.Value.
func getFields(val reflect.Value, tagName string) []*Field {
	t := val.Type()
	fields := make([]*Field, 0, t.NumField())

	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		fv := val.Field(i)

		opts := parseTag(sf.Tag.Get(tagName))
		if opts.skip {
			continue
		}

		fields = append(fields, &Field{
			value:   fv,
			field:   sf,
			tagName: tagName,
			tagOpts: opts,
		})
	}

	return fields
}
