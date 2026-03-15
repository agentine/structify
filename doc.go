// Package structify provides utilities for working with Go structs using reflection.
//
// It is a modern, drop-in replacement for github.com/fatih/structs with
// improvements including generics support, error returns instead of panics,
// configurable struct tags, and nil pointer safety.
//
// Basic usage:
//
//	type Server struct {
//		Name    string `json:"name"`
//		ID      int    `json:"id"`
//		Enabled bool   `json:"enabled"`
//	}
//
//	s := Server{Name: "web", ID: 1, Enabled: true}
//
//	// Convert struct to map
//	m := structify.Map(s) // map[Name:web ID:1 Enabled:true]
//
//	// Get field names
//	names := structify.Names(s) // [Name, ID, Enabled]
//
//	// Get field values
//	values := structify.Values(s) // [web, 1, true]
//
//	// Check if struct has zero values
//	structify.HasZero(s) // false
//
//	// Use custom tag for map keys
//	st := structify.New(s, structify.WithTag("json"))
//	m = st.Map() // map[name:web id:1 enabled:true]
package structify
