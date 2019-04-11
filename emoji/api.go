package emoji

// api represents an Emoji object returned by the MastoAPI.
type api struct {
	Shortcode       string
	StaticURL       string // doesn't seem to be used?
	URL             string
	VisibleInPicker bool // always returns false?
}
