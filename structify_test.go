package structify

import (
	"testing"
)

type testStruct struct {
	Name    string `json:"name"`
	ID      int    `json:"id"`
	Enabled bool   `json:"enabled"`
}

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

func TestMapWithTag(t *testing.T) {
	s := testStruct{Name: "web", ID: 1, Enabled: true}
	st := New(s, WithTag("json"))
	m := st.Map()

	if m["name"] != "web" {
		t.Errorf("expected name=web, got %v", m["name"])
	}
	if m["id"] != 1 {
		t.Errorf("expected id=1, got %v", m["id"])
	}
}

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

func TestValues(t *testing.T) {
	s := testStruct{Name: "web", ID: 1, Enabled: true}
	vals := Values(s)

	if len(vals) != 3 {
		t.Fatalf("expected 3 values, got %d", len(vals))
	}
	if vals[0] != "web" {
		t.Errorf("expected web, got %v", vals[0])
	}
}

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
}

func TestIsZero(t *testing.T) {
	if !IsZero(testStruct{}) {
		t.Error("expected zero struct to be zero")
	}
	if IsZero(testStruct{Name: "x"}) {
		t.Error("expected non-zero struct to not be zero")
	}
}

func TestHasZero(t *testing.T) {
	if !HasZero(testStruct{Name: "x"}) {
		t.Error("expected struct with some zero fields to have zero")
	}
	if HasZero(testStruct{Name: "x", ID: 1, Enabled: true}) {
		t.Error("expected fully populated struct to not have zero")
	}
}

func TestName(t *testing.T) {
	s := testStruct{}
	if Name(s) != "testStruct" {
		t.Errorf("expected testStruct, got %s", Name(s))
	}
}

func TestField(t *testing.T) {
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

	_, err = st.Field("Nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent field")
	}
}
