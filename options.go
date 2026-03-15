package structify

// Option configures the behavior of a Struct.
type Option func(*config)

type config struct {
	tagName string
}

func defaultConfig() *config {
	return &config{
		tagName: "structify",
	}
}

// WithTag sets the struct tag name used for map keys and tag options.
// If a field has the specified tag, its value is used as the map key.
// Default tag is "structify".
func WithTag(tag string) Option {
	return func(c *config) {
		c.tagName = tag
	}
}
