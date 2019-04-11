package emoji

type log struct {
	Errors    bool
	Downloads bool
}

// Verbose is an accessor to two configurable elements that have default values.
//
// - `Verbose.Errors` determines if errors should be printed to stdout.
//	 Default: false
// - `Verbose.Downloads` determines if files downloaded should be printed to
// stdout.
//	 Default: false
var Verbose = log{Errors: false, Downloads: false}
