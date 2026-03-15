package structify

import "strings"

// tagOptions represents parsed struct tag options.
type tagOptions struct {
	name       string
	skip       bool // "-"
	omitempty  bool
	omitnested bool
	asString   bool // "string"
}

// parseTag parses a struct tag value into its name and options.
// The tag format is: "name,option1,option2"
// Recognized options: omitempty, omitnested, string
// A tag value of "-" means the field should be skipped.
func parseTag(tag string) tagOptions {
	opts := tagOptions{}

	if tag == "" {
		return opts
	}

	if tag == "-" {
		opts.skip = true
		return opts
	}

	parts := strings.Split(tag, ",")
	opts.name = parts[0]

	for _, opt := range parts[1:] {
		switch opt {
		case "omitempty":
			opts.omitempty = true
		case "omitnested":
			opts.omitnested = true
		case "string":
			opts.asString = true
		}
	}

	return opts
}
