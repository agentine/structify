package structify

// FillMap fills the given map with the struct's fields.
// This provides compatibility with fatih/structs.FillMap.
func FillMap(s any, out map[string]any) {
	m := Map(s)
	for k, v := range m {
		out[k] = v
	}
}

// FillMapFrom fills the given map using the specified Struct wrapper.
// This allows using custom options (e.g., WithTag) with FillMap behavior.
func (st *Struct) FillMap(out map[string]any) {
	m := st.Map()
	for k, v := range m {
		out[k] = v
	}
}
