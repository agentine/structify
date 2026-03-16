# Changelog

## v0.1.0 — 2026-03-15

Initial release. Drop-in replacement for [fatih/structs](https://github.com/fatih/structs) with modern Go improvements.

### Added

- Core `Struct` type with `New()` constructor and functional options (`WithTag`)
- Conversion methods: `Map()`, `Values()`, `Names()`, `Fields()`
- `Field` type with `Name`, `Value`, `Kind`, `Tag`, `IsZero`, `IsExported`, `IsEmbedded`, `Set`, nested `Fields`
- Tag parsing: `omitempty`, `-`, `omitnested`, `string` options
- Compatibility layer: `FillMap()` for fatih/structs migration
- Error returns instead of panics (`ErrFieldNotFound`, `ErrNotExported`, `ErrNotSettable`, `ErrTypeMismatch`)
- Nil pointer safety for embedded/nested structs
- Top-level convenience functions: `Map`, `Values`, `Names`, `Fields`, `IsZero`, `HasZero`, `IsStruct`, `Name`
- 96.5% test coverage with 60+ tests
- CI on Go 1.23 and 1.24 with golangci-lint
