package structify

import (
	"fmt"
	"reflect"
)

// Map converts a struct to a map[string]any. Convenience wrapper.
func Map(s any) map[string]any {
	return New(s).Map()
}

// Values returns the field values of a struct as []any. Convenience wrapper.
func Values(s any) []any {
	return New(s).Values()
}

// Names returns the field names of a struct. Convenience wrapper.
func Names(s any) []string {
	return New(s).Names()
}

// Fields returns the fields of a struct. Convenience wrapper.
func Fields(s any) []*Field {
	return New(s).Fields()
}

// IsZero returns true if all exported fields of the struct are zero-valued.
func IsZero(s any) bool {
	return New(s).IsZero()
}

// HasZero returns true if any exported field of the struct is zero-valued.
func HasZero(s any) bool {
	return New(s).HasZero()
}

// IsStruct returns true if the given value is a struct or pointer to struct.
func IsStruct(s any) bool {
	v := reflect.ValueOf(s)
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return false
		}
		v = v.Elem()
	}
	return v.Kind() == reflect.Struct
}

// Name returns the struct type name. Convenience wrapper.
func Name(s any) string {
	return New(s).Name()
}

// structToMap converts a struct reflect.Value to a map[string]any.
func structToMap(val reflect.Value, tagName string) map[string]any {
	t := val.Type()
	out := make(map[string]any)

	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		fv := val.Field(i)

		// Skip unexported fields
		if !sf.IsExported() {
			continue
		}

		opts := parseTag(sf.Tag.Get(tagName))

		// Skip fields tagged with "-"
		if opts.skip {
			continue
		}

		// Determine map key
		key := sf.Name
		if opts.name != "" {
			key = opts.name
		}

		// Handle omitempty
		if opts.omitempty && fv.IsZero() {
			continue
		}

		// Handle string option — format as string
		if opts.asString {
			out[key] = fmt.Sprintf("%v", fv.Interface())
			continue
		}

		// Resolve pointers
		v := fv
		for v.Kind() == reflect.Ptr {
			if v.IsNil() {
				break
			}
			v = v.Elem()
		}

		// Handle nested structs
		if v.Kind() == reflect.Struct && !opts.omitnested {
			out[key] = structToMap(v, tagName)
			continue
		}

		out[key] = fv.Interface()
	}

	return out
}

// structValues extracts field values from a struct.
func structValues(val reflect.Value, tagName string) []any {
	t := val.Type()
	values := make([]any, 0, t.NumField())

	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		fv := val.Field(i)

		if !sf.IsExported() {
			continue
		}

		opts := parseTag(sf.Tag.Get(tagName))
		if opts.skip {
			continue
		}

		values = append(values, fv.Interface())
	}

	return values
}

// structNames extracts field names from a struct.
func structNames(val reflect.Value, tagName string) []string {
	t := val.Type()
	names := make([]string, 0, t.NumField())

	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)

		if !sf.IsExported() {
			continue
		}

		opts := parseTag(sf.Tag.Get(tagName))
		if opts.skip {
			continue
		}

		names = append(names, sf.Name)
	}

	return names
}

// isZero returns true if all exported fields are zero-valued.
func isZero(val reflect.Value, tagName string) bool {
	t := val.Type()

	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		fv := val.Field(i)

		if !sf.IsExported() {
			continue
		}

		opts := parseTag(sf.Tag.Get(tagName))
		if opts.skip {
			continue
		}

		if !fv.IsZero() {
			return false
		}
	}

	return true
}

// hasZero returns true if any exported field is zero-valued.
func hasZero(val reflect.Value, tagName string) bool {
	t := val.Type()

	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		fv := val.Field(i)

		if !sf.IsExported() {
			continue
		}

		opts := parseTag(sf.Tag.Get(tagName))
		if opts.skip {
			continue
		}

		if fv.IsZero() {
			return true
		}
	}

	return false
}
