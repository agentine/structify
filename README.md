# structify

[![CI](https://github.com/agentine/structify/actions/workflows/ci.yml/badge.svg)](https://github.com/agentine/structify/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/agentine/structify.svg)](https://pkg.go.dev/github.com/agentine/structify)

A modern Go library for struct utilities — convert structs to maps, extract field names and values, inspect and mutate fields. Drop-in replacement for [github.com/fatih/structs](https://github.com/fatih/structs).

## Install

```bash
go get github.com/agentine/structify
```

## Usage

```go
package main

import (
	"fmt"
	"github.com/agentine/structify"
)

type Server struct {
	Name    string `json:"name"`
	ID      int    `json:"id"`
	Enabled bool   `json:"enabled"`
}

func main() {
	s := Server{Name: "web", ID: 1, Enabled: true}

	// Convert struct to map
	m := structify.Map(s)
	fmt.Println(m) // map[Name:web ID:1 Enabled:true]

	// Get field names
	names := structify.Names(s)
	fmt.Println(names) // [Name ID Enabled]

	// Get field values
	values := structify.Values(s)
	fmt.Println(values) // [web 1 true]

	// Use custom tag for map keys
	st := structify.New(s, structify.WithTag("json"))
	m = st.Map()
	fmt.Println(m) // map[name:web id:1 enabled:true]
}
```

## API

### Core Functions

| Function | Description |
|----------|-------------|
| `New(s any, opts ...Option) *Struct` | Wrap a struct value |
| `Map(s any) map[string]any` | Convert struct to map |
| `Values(s any) []any` | Extract field values |
| `Names(s any) []string` | Extract field names |
| `Fields(s any) []*Field` | Get all fields |
| `IsZero(s any) bool` | Check if all fields are zero |
| `HasZero(s any) bool` | Check if any field is zero |
| `IsStruct(s any) bool` | Check if value is a struct |
| `Name(s any) string` | Get struct type name |

### Field Methods

| Method | Description |
|--------|-------------|
| `Name() string` | Field name |
| `Value() any` | Field value |
| `Kind() reflect.Kind` | Field kind |
| `Tag(key string) string` | Get struct tag value |
| `IsZero() bool` | Check if zero value |
| `IsExported() bool` | Check if exported |
| `IsEmbedded() bool` | Check if embedded |
| `Set(val any) error` | Set field value |
| `Fields() []*Field` | Nested struct fields |

### Options

| Option | Description |
|--------|-------------|
| `WithTag(tag string)` | Use custom struct tag for map keys (default: field name) |

## Migration from fatih/structs

structify is a drop-in replacement for `fatih/structs`. Key differences:

1. **`any` instead of `interface{}`** — Modern Go style throughout.
2. **Error returns** — `Struct.Field(name)` returns `(*Field, error)` instead of panicking.
3. **Configurable tags** — `New(s, WithTag("json"))` to use any struct tag for map keys.
4. **Nil pointer safety** — Graceful handling of nil embedded pointers.

### Quick migration

```diff
- import "github.com/fatih/structs"
+ import "github.com/agentine/structify"

- m := structs.Map(s)
+ m := structify.Map(s)

- f := structs.Fields(s)
+ f := structify.Fields(s)
```

The compatibility layer provides `FillMap` and other functions matching the original API signatures.

## License

MIT
