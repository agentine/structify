package structify

import (
	"errors"
	"reflect"
	"testing"
)

// --- Test struct types ---

type testStruct struct {
	Name    string `json:"name"`
	ID      int    `json:"id"`
	Enabled bool   `json:"enabled"`
}

type taggedStruct struct {
	Name    string `structify:"name"`
	Skipped string `structify:"-"`
	Empty   string `structify:"empty,omitempty"`
	Nested  inner  `structify:"nested,omitnested"`
	AsStr   int    `structify:"as_str,string"`
}

type inner struct {
	Value string
}

type nestedStruct struct {
	Name  string
	Inner inner
}

type ptrNestedStruct struct {
	Name  string
	Inner *inner
}

type embeddedStruct struct {
	inner
	Extra string
}

type unexportedStruct struct {
	Exported   string
	unexported string //nolint:unused
}

type emptyStruct struct{}

type allZeroStruct struct {
	S string
	I int
	B bool
}

type noZeroStruct struct {
	S string
	I int
	B bool
}

type setTestStruct struct {
	Name string
	ID   int
}

// --- Struct constructor tests ---

func TestNew(t *testing.T) {
	s := testStruct{Name: "web"}
	st := New(s)
	if st.Name() != "testStruct" {
		t.Errorf("expected testStruct, got %s", st.Name())
	}
}

func TestNewPointer(t *testing.T) {
	s := &testStruct{Name: "web"}
	st := New(s)
	if st.Name() != "testStruct" {
		t.Errorf("expected testStruct, got %s", st.Name())
	}
}

func TestNewDoublePointer(t *testing.T) {
	s := &testStruct{Name: "web"}
	p := &s
	st := New(p)
	if st.Name() != "testStruct" {
		t.Errorf("expected testStruct, got %s", st.Name())
	}
}

func TestNewPanicsOnNonStruct(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for non-struct")
		}
	}()
	New("not a struct")
}

// --- Map tests ---

func TestMapBasic(t *testing.T) {
	s := testStruct{Name: "web", ID: 1, Enabled: true}
	m := Map(s)

	if m["Name"] != "web" {
		t.Errorf("expected Name=web, got %v", m["Name"])
	}
	if m["ID"] != 1 {
		t.Errorf("expected ID=1, got %v", m["ID"])
	}
	if m["Enabled"] != true {
		t.Errorf("expected Enabled=true, got %v", m["Enabled"])
	}
}

func TestMapFromPointer(t *testing.T) {
	s := &testStruct{Name: "web", ID: 1}
	m := Map(s)
	if m["Name"] != "web" {
		t.Errorf("expected Name=web, got %v", m["Name"])
	}
}

func TestMapWithTag(t *testing.T) {
	s := testStruct{Name: "web", ID: 1, Enabled: true}
	st := New(s, WithTag("json"))
	m := st.Map()

	if _, ok := m["name"]; !ok {
		t.Error("expected json tag key 'name'")
	}
	if _, ok := m["id"]; !ok {
		t.Error("expected json tag key 'id'")
	}
}

func TestMapSkipTag(t *testing.T) {
	s := taggedStruct{Name: "a", Skipped: "b"}
	m := New(s).Map()

	if _, ok := m["name"]; !ok {
		t.Error("expected 'name' key")
	}
	if _, ok := m["-"]; ok {
		t.Error("did not expect skipped field")
	}
	if _, ok := m["Skipped"]; ok {
		t.Error("did not expect Skipped field")
	}
}

func TestMapOmitempty(t *testing.T) {
	s := taggedStruct{Name: "a"}
	m := New(s).Map()

	if _, ok := m["empty"]; ok {
		t.Error("expected omitempty to skip zero-valued field")
	}
}

func TestMapOmitemptyNonZero(t *testing.T) {
	s := taggedStruct{Name: "a", Empty: "present"}
	m := New(s).Map()

	if v, ok := m["empty"]; !ok || v != "present" {
		t.Errorf("expected empty=present, got %v", v)
	}
}

func TestMapOmitnested(t *testing.T) {
	s := taggedStruct{Name: "a", Nested: inner{Value: "x"}}
	m := New(s).Map()

	// With omitnested, nested struct should not be converted to map
	if _, ok := m["nested"].(map[string]any); ok {
		t.Error("expected omitnested to prevent map conversion")
	}
}

func TestMapStringOption(t *testing.T) {
	s := taggedStruct{Name: "a", AsStr: 42}
	m := New(s).Map()

	if m["as_str"] != "42" {
		t.Errorf("expected string '42', got %v (type %T)", m["as_str"], m["as_str"])
	}
}

func TestMapNestedStruct(t *testing.T) {
	s := nestedStruct{Name: "outer", Inner: inner{Value: "inner"}}
	m := Map(s)

	nested, ok := m["Inner"].(map[string]any)
	if !ok {
		t.Fatal("expected nested struct to be converted to map")
	}
	if nested["Value"] != "inner" {
		t.Errorf("expected inner, got %v", nested["Value"])
	}
}

func TestMapNestedPointerStruct(t *testing.T) {
	inner := &inner{Value: "ptr"}
	s := ptrNestedStruct{Name: "outer", Inner: inner}
	m := Map(s)

	nested, ok := m["Inner"].(map[string]any)
	if !ok {
		t.Fatal("expected nested pointer struct to be converted to map")
	}
	if nested["Value"] != "ptr" {
		t.Errorf("expected ptr, got %v", nested["Value"])
	}
}

func TestMapNestedNilPointer(t *testing.T) {
	s := ptrNestedStruct{Name: "outer", Inner: nil}
	m := Map(s)

	// nil pointer field should remain as the nil pointer value
	if m["Inner"] != (*inner)(nil) {
		t.Errorf("expected nil *inner, got %v", m["Inner"])
	}
}

func TestMapUnexportedFields(t *testing.T) {
	s := unexportedStruct{Exported: "yes"}
	m := Map(s)

	if _, ok := m["Exported"]; !ok {
		t.Error("expected Exported key")
	}
	if _, ok := m["unexported"]; ok {
		t.Error("did not expect unexported key")
	}
}

func TestMapEmptyStruct(t *testing.T) {
	m := Map(emptyStruct{})
	if len(m) != 0 {
		t.Errorf("expected empty map, got %v", m)
	}
}

func TestMapEmbeddedStruct(t *testing.T) {
	s := embeddedStruct{inner: inner{Value: "embedded"}, Extra: "extra"}
	m := Map(s)

	if _, ok := m["Extra"]; !ok {
		t.Error("expected Extra key")
	}
	// Embedded struct appears as nested
	if _, ok := m["inner"]; !ok {
		t.Error("expected inner key for embedded struct")
	}
}

// --- Values tests ---

func TestValues(t *testing.T) {
	s := testStruct{Name: "web", ID: 1, Enabled: true}
	vals := Values(s)

	if len(vals) != 3 {
		t.Fatalf("expected 3 values, got %d", len(vals))
	}
	if vals[0] != "web" {
		t.Errorf("expected web, got %v", vals[0])
	}
	if vals[1] != 1 {
		t.Errorf("expected 1, got %v", vals[1])
	}
	if vals[2] != true {
		t.Errorf("expected true, got %v", vals[2])
	}
}

func TestValuesSkipsUnexported(t *testing.T) {
	s := unexportedStruct{Exported: "yes"}
	vals := Values(s)
	if len(vals) != 1 {
		t.Errorf("expected 1 value, got %d", len(vals))
	}
}

func TestValuesSkipTag(t *testing.T) {
	s := taggedStruct{Name: "a", Skipped: "b"}
	vals := New(s).Values()
	// Should have Name, Empty, Nested, AsStr (4 fields, Skipped is "-")
	if len(vals) != 4 {
		t.Errorf("expected 4 values (skip tagged -), got %d", len(vals))
	}
}

func TestValuesEmpty(t *testing.T) {
	vals := Values(emptyStruct{})
	if len(vals) != 0 {
		t.Errorf("expected 0 values, got %d", len(vals))
	}
}

// --- Names tests ---

func TestNames(t *testing.T) {
	s := testStruct{Name: "web", ID: 1, Enabled: true}
	names := Names(s)

	if len(names) != 3 {
		t.Fatalf("expected 3 names, got %d", len(names))
	}
	if names[0] != "Name" || names[1] != "ID" || names[2] != "Enabled" {
		t.Errorf("unexpected names: %v", names)
	}
}

func TestNamesSkipsUnexported(t *testing.T) {
	names := Names(unexportedStruct{Exported: "yes"})
	if len(names) != 1 || names[0] != "Exported" {
		t.Errorf("expected [Exported], got %v", names)
	}
}

func TestNamesEmpty(t *testing.T) {
	names := Names(emptyStruct{})
	if len(names) != 0 {
		t.Errorf("expected 0 names, got %d", len(names))
	}
}

// --- Fields tests ---

func TestFields(t *testing.T) {
	s := testStruct{Name: "web", ID: 1, Enabled: true}
	fields := Fields(s)

	if len(fields) != 3 {
		t.Fatalf("expected 3 fields, got %d", len(fields))
	}
	if fields[0].Name() != "Name" {
		t.Errorf("expected Name, got %s", fields[0].Name())
	}
}

func TestFieldsSkipTag(t *testing.T) {
	s := taggedStruct{Name: "a", Skipped: "b"}
	fields := New(s).Fields()
	// Should have 4 fields (Skipped is "-")
	if len(fields) != 4 {
		t.Errorf("expected 4 fields, got %d", len(fields))
	}
}

// --- Field method tests ---

func TestFieldByName(t *testing.T) {
	s := testStruct{Name: "web", ID: 1}
	st := New(s)

	f, err := st.Field("Name")
	if err != nil {
		t.Fatal(err)
	}
	if f.Name() != "Name" {
		t.Errorf("expected Name, got %s", f.Name())
	}
	if f.Value() != "web" {
		t.Errorf("expected web, got %v", f.Value())
	}
}

func TestFieldNotFound(t *testing.T) {
	st := New(testStruct{})
	_, err := st.Field("Nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent field")
	}
	if !errors.Is(err, ErrFieldNotFound) {
		t.Errorf("expected ErrFieldNotFound, got %v", err)
	}
}

func TestFieldKind(t *testing.T) {
	s := testStruct{Name: "web", ID: 1, Enabled: true}
	fields := Fields(s)

	if fields[0].Kind() != reflect.String {
		t.Errorf("expected String, got %v", fields[0].Kind())
	}
	if fields[1].Kind() != reflect.Int {
		t.Errorf("expected Int, got %v", fields[1].Kind())
	}
	if fields[2].Kind() != reflect.Bool {
		t.Errorf("expected Bool, got %v", fields[2].Kind())
	}
}

func TestFieldTag(t *testing.T) {
	s := testStruct{Name: "web"}
	fields := Fields(s)
	if fields[0].Tag("json") != "name" {
		t.Errorf("expected 'name', got %s", fields[0].Tag("json"))
	}
}

func TestFieldIsZero(t *testing.T) {
	s := testStruct{Name: "web"}
	fields := Fields(s)
	if fields[0].IsZero() {
		t.Error("Name should not be zero")
	}
	if !fields[1].IsZero() {
		t.Error("ID should be zero")
	}
}

func TestFieldIsExported(t *testing.T) {
	s := unexportedStruct{Exported: "yes"}
	st := New(s)
	fields := st.Fields()
	if !fields[0].IsExported() {
		t.Error("Exported should be exported")
	}
}

func TestFieldIsEmbedded(t *testing.T) {
	s := embeddedStruct{inner: inner{Value: "x"}, Extra: "e"}
	fields := Fields(s)

	foundEmbedded := false
	for _, f := range fields {
		if f.IsEmbedded() {
			foundEmbedded = true
		}
	}
	if !foundEmbedded {
		t.Error("expected to find an embedded field")
	}
}

func TestFieldSet(t *testing.T) {
	s := &setTestStruct{Name: "old", ID: 1}
	st := New(s)
	f, err := st.Field("Name")
	if err != nil {
		t.Fatal(err)
	}

	if err := f.Set("new"); err != nil {
		t.Fatal(err)
	}
	if s.Name != "new" {
		t.Errorf("expected new, got %s", s.Name)
	}
}

func TestFieldSetTypeMismatch(t *testing.T) {
	s := &setTestStruct{Name: "old"}
	st := New(s)
	f, err := st.Field("Name")
	if err != nil {
		t.Fatal(err)
	}

	err = f.Set(42)
	if err == nil {
		t.Error("expected error for type mismatch")
	}
	if !errors.Is(err, ErrTypeMismatch) {
		t.Errorf("expected ErrTypeMismatch, got %v", err)
	}
}

func TestFieldSetNotSettable(t *testing.T) {
	s := setTestStruct{Name: "old"} // not a pointer, so not settable
	st := New(s)
	f, err := st.Field("Name")
	if err != nil {
		t.Fatal(err)
	}

	err = f.Set("new")
	if err == nil {
		t.Error("expected error for non-settable field")
	}
	if !errors.Is(err, ErrNotSettable) {
		t.Errorf("expected ErrNotSettable, got %v", err)
	}
}

func TestFieldNestedFields(t *testing.T) {
	s := nestedStruct{Name: "outer", Inner: inner{Value: "x"}}
	fields := Fields(s)

	var innerField *Field
	for _, f := range fields {
		if f.Name() == "Inner" {
			innerField = f
			break
		}
	}
	if innerField == nil {
		t.Fatal("Inner field not found")
	}

	nested := innerField.Fields()
	if len(nested) != 1 {
		t.Fatalf("expected 1 nested field, got %d", len(nested))
	}
	if nested[0].Name() != "Value" {
		t.Errorf("expected Value, got %s", nested[0].Name())
	}
}

func TestFieldNestedFieldsPointer(t *testing.T) {
	s := ptrNestedStruct{Name: "outer", Inner: &inner{Value: "ptr"}}
	fields := Fields(s)

	var innerField *Field
	for _, f := range fields {
		if f.Name() == "Inner" {
			innerField = f
			break
		}
	}
	if innerField == nil {
		t.Fatal("Inner field not found")
	}

	nested := innerField.Fields()
	if len(nested) != 1 {
		t.Fatalf("expected 1 nested field, got %d", len(nested))
	}
}

func TestFieldNestedFieldsNilPointer(t *testing.T) {
	s := ptrNestedStruct{Name: "outer", Inner: nil}
	fields := Fields(s)

	var innerField *Field
	for _, f := range fields {
		if f.Name() == "Inner" {
			innerField = f
			break
		}
	}
	if innerField == nil {
		t.Fatal("Inner field not found")
	}

	nested := innerField.Fields()
	if nested != nil {
		t.Error("expected nil for nil pointer nested fields")
	}
}

func TestFieldNestedFieldsNonStruct(t *testing.T) {
	s := testStruct{Name: "web"}
	fields := Fields(s)
	nested := fields[0].Fields() // String field
	if nested != nil {
		t.Error("expected nil for non-struct field")
	}
}

// --- IsZero / HasZero tests ---

func TestIsZero(t *testing.T) {
	if !IsZero(allZeroStruct{}) {
		t.Error("expected zero struct to be zero")
	}
	if IsZero(testStruct{Name: "x"}) {
		t.Error("expected non-zero struct to not be zero")
	}
}

func TestIsZeroEmpty(t *testing.T) {
	if !IsZero(emptyStruct{}) {
		t.Error("expected empty struct to be zero")
	}
}

func TestHasZero(t *testing.T) {
	if !HasZero(testStruct{Name: "x"}) {
		t.Error("expected struct with some zero fields to have zero")
	}
	s := noZeroStruct{S: "x", I: 1, B: true}
	if HasZero(s) {
		t.Error("expected fully populated struct to not have zero")
	}
}

func TestHasZeroAllZero(t *testing.T) {
	if !HasZero(allZeroStruct{}) {
		t.Error("expected all-zero struct to have zero")
	}
}

func TestHasZeroEmpty(t *testing.T) {
	if HasZero(emptyStruct{}) {
		t.Error("expected empty struct to not have zero (no fields to check)")
	}
}

// --- IsStruct tests ---

func TestIsStruct(t *testing.T) {
	if !IsStruct(testStruct{}) {
		t.Error("expected true for struct")
	}
	if !IsStruct(&testStruct{}) {
		t.Error("expected true for pointer to struct")
	}
	if IsStruct("string") {
		t.Error("expected false for string")
	}
	if IsStruct(42) {
		t.Error("expected false for int")
	}
	if IsStruct(nil) {
		t.Error("expected false for nil")
	}
}

func TestIsStructNilPointer(t *testing.T) {
	var p *testStruct
	if IsStruct(p) {
		t.Error("expected false for nil pointer to struct")
	}
}

// --- Name tests ---

func TestName(t *testing.T) {
	if Name(testStruct{}) != "testStruct" {
		t.Errorf("expected testStruct, got %s", Name(testStruct{}))
	}
	if Name(&testStruct{}) != "testStruct" {
		t.Errorf("expected testStruct from pointer, got %s", Name(&testStruct{}))
	}
}

// --- Tag parsing tests ---

func TestParseTagEmpty(t *testing.T) {
	opts := parseTag("")
	if opts.name != "" || opts.skip || opts.omitempty || opts.omitnested || opts.asString {
		t.Error("expected all defaults for empty tag")
	}
}

func TestParseTagSkip(t *testing.T) {
	opts := parseTag("-")
	if !opts.skip {
		t.Error("expected skip=true for '-' tag")
	}
}

func TestParseTagName(t *testing.T) {
	opts := parseTag("my_field")
	if opts.name != "my_field" {
		t.Errorf("expected my_field, got %s", opts.name)
	}
}

func TestParseTagOptions(t *testing.T) {
	opts := parseTag("name,omitempty,omitnested,string")
	if opts.name != "name" {
		t.Errorf("expected name, got %s", opts.name)
	}
	if !opts.omitempty {
		t.Error("expected omitempty=true")
	}
	if !opts.omitnested {
		t.Error("expected omitnested=true")
	}
	if !opts.asString {
		t.Error("expected asString=true")
	}
}

func TestParseTagUnknownOption(t *testing.T) {
	opts := parseTag("name,unknown")
	if opts.name != "name" {
		t.Errorf("expected name, got %s", opts.name)
	}
	// Unknown options should be silently ignored
	if opts.omitempty || opts.omitnested || opts.asString {
		t.Error("expected unknown option to be ignored")
	}
}

// --- FillMap compat tests ---

func TestFillMap(t *testing.T) {
	s := testStruct{Name: "web", ID: 1, Enabled: true}
	out := make(map[string]any)
	FillMap(s, out)

	if out["Name"] != "web" {
		t.Errorf("expected Name=web, got %v", out["Name"])
	}
	if out["ID"] != 1 {
		t.Errorf("expected ID=1, got %v", out["ID"])
	}
}

func TestFillMapExistingKeys(t *testing.T) {
	s := testStruct{Name: "web", ID: 1}
	out := map[string]any{"Extra": "value"}
	FillMap(s, out)

	if out["Extra"] != "value" {
		t.Error("expected existing key to be preserved")
	}
	if out["Name"] != "web" {
		t.Error("expected struct fields to be added")
	}
}

func TestStructFillMap(t *testing.T) {
	s := testStruct{Name: "web", ID: 1}
	st := New(s, WithTag("json"))
	out := make(map[string]any)
	st.FillMap(out)

	if out["name"] != "web" {
		t.Errorf("expected name=web, got %v", out["name"])
	}
}

// --- Struct method tests ---

func TestStructMap(t *testing.T) {
	st := New(testStruct{Name: "web", ID: 1})
	m := st.Map()
	if m["Name"] != "web" {
		t.Errorf("expected web, got %v", m["Name"])
	}
}

func TestStructValues(t *testing.T) {
	st := New(testStruct{Name: "web", ID: 1, Enabled: true})
	vals := st.Values()
	if len(vals) != 3 {
		t.Fatalf("expected 3 values, got %d", len(vals))
	}
}

func TestStructNames(t *testing.T) {
	st := New(testStruct{Name: "web"})
	names := st.Names()
	if len(names) != 3 {
		t.Fatalf("expected 3 names, got %d", len(names))
	}
}

func TestStructFields(t *testing.T) {
	st := New(testStruct{Name: "web"})
	fields := st.Fields()
	if len(fields) != 3 {
		t.Fatalf("expected 3 fields, got %d", len(fields))
	}
}

func TestStructIsZero(t *testing.T) {
	st := New(testStruct{})
	if !st.IsZero() {
		t.Error("expected zero")
	}
}

func TestStructHasZero(t *testing.T) {
	st := New(testStruct{Name: "x"})
	if !st.HasZero() {
		t.Error("expected has zero")
	}
}

func TestStructName(t *testing.T) {
	st := New(testStruct{})
	if st.Name() != "testStruct" {
		t.Errorf("expected testStruct, got %s", st.Name())
	}
}

// --- WithTag option test ---

func TestWithTag(t *testing.T) {
	s := testStruct{Name: "web"}
	st := New(s, WithTag("json"))
	m := st.Map()
	if _, ok := m["name"]; !ok {
		t.Error("expected json tag key")
	}
}

// --- Edge cases ---

func TestFieldSetUnexported(t *testing.T) {
	s := &unexportedStruct{Exported: "yes"}
	st := New(s)

	// Fields() returns all non-skipped fields including unexported (like fatih/structs)
	fields := st.Fields()
	if len(fields) != 2 {
		t.Errorf("expected 2 fields (exported + unexported), got %d", len(fields))
	}

	f, err := st.Field("Exported")
	if err != nil {
		t.Fatal(err)
	}
	if err := f.Set("new"); err != nil {
		t.Fatal(err)
	}
	if s.Exported != "new" {
		t.Errorf("expected new, got %s", s.Exported)
	}
}

func TestFieldSetInt(t *testing.T) {
	s := &setTestStruct{Name: "old", ID: 1}
	st := New(s)
	f, err := st.Field("ID")
	if err != nil {
		t.Fatal(err)
	}

	if err := f.Set(42); err != nil {
		t.Fatal(err)
	}
	if s.ID != 42 {
		t.Errorf("expected 42, got %d", s.ID)
	}
}

type multiTagStruct struct {
	Field1 string `structify:"f1,omitempty" json:"field_1"`
	Field2 int    `structify:"-" json:"field_2"`
}

func TestMultipleTagKeys(t *testing.T) {
	s := multiTagStruct{Field1: "val", Field2: 10}

	// With structify tag
	m1 := New(s).Map()
	if _, ok := m1["f1"]; !ok {
		t.Error("expected f1 key with structify tag")
	}
	if _, ok := m1["Field2"]; ok {
		t.Error("Field2 should be skipped with structify tag")
	}

	// With json tag
	m2 := New(s, WithTag("json")).Map()
	if _, ok := m2["field_1"]; !ok {
		t.Error("expected field_1 key with json tag")
	}
	if _, ok := m2["field_2"]; !ok {
		t.Error("expected field_2 key with json tag")
	}
}
