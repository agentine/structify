# structify — Go Struct Utilities

**Replaces:** [github.com/fatih/structs](https://github.com/fatih/structs)
**Package:** `github.com/agentine/structify`
**Language:** Go (1.22+)
**License:** MIT

## Why

fatih/structs has 4,258 importers and 3,900 stars but was archived in October 2018 (last commit January 2017). The README states "This project is not maintained anymore." Single maintainer (fatih). No fork has gained traction (all forks have 0 stars). The library predates Go generics and uses `interface{}` throughout.

## Scope

Drop-in replacement for fatih/structs with modern Go improvements:
- All original functionality preserved (struct→map, struct→slice, field names, field inspection, zero-value checking, nested struct support)
- Generics support where beneficial (type-safe value extraction)
- `any` instead of `interface{}`
- Compatibility layer: exported aliases and adapter functions so existing `fatih/structs` users can migrate with minimal changes

## Architecture

### Package Structure

```
structify/
├── structify.go       # Core Struct type and methods
├── field.go           # Field type and methods
├── tags.go            # Tag parsing and options
├── convert.go         # Map/Values/Names conversion functions
├── options.go         # Functional options for configuration
├── compat.go          # fatih/structs compatibility aliases
├── doc.go             # Package documentation
└── *_test.go          # Tests for each file
```

### Core API

```go
// Core type wrapping a struct value
type Struct struct { ... }

// Constructor
func New(s any) *Struct

// Conversion methods
func (s *Struct) Map() map[string]any
func (s *Struct) Values() []any
func (s *Struct) Names() []string
func (s *Struct) Fields() []*Field

// Inspection methods
func (s *Struct) Field(name string) (*Field, error)
func (s *Struct) IsZero() bool
func (s *Struct) HasZero() bool
func (s *Struct) Name() string

// Top-level convenience functions
func Map(s any) map[string]any
func Values(s any) []any
func Names(s any) []string
func Fields(s any) []*Field
func IsZero(s any) bool
func HasZero(s any) bool
func IsStruct(s any) bool
func Name(s any) string

// Field type
type Field struct { ... }
func (f *Field) Name() string
func (f *Field) Value() any
func (f *Field) Kind() reflect.Kind
func (f *Field) Tag(key string) string
func (f *Field) IsZero() bool
func (f *Field) IsExported() bool
func (f *Field) IsEmbedded() bool
func (f *Field) Set(val any) error
func (f *Field) Fields() []*Field  // nested struct fields
```

### Improvements Over fatih/structs

1. **`any` everywhere** — replaces `interface{}` for modern Go style
2. **Error returns** — `Field()` returns `(*Field, error)` instead of panicking
3. **Configurable tag** — `New(s, WithTag("json"))` to use any struct tag for map keys
4. **Nested struct control** — options for flatten vs. nested map output
5. **Nil pointer safety** — graceful handling of nil embedded pointers
6. **Context-aware** — no global state, all config via functional options

### Compatibility Layer (compat.go)

Provide a migration path:
- Same function signatures as fatih/structs where possible
- `FillMap(s any, out map[string]any)` for parity
- Document migration guide in README

## Deliverables

1. Project scaffolding (go.mod, directory structure, CI config)
2. Core `Struct` type with `Map()`, `Values()`, `Names()`, `Fields()`
3. `Field` type with inspection and mutation methods
4. Tag parsing with options (`omitempty`, `-`, `omitnested`, `string`)
5. Compatibility layer and migration documentation
6. Comprehensive test suite (target >95% coverage)
