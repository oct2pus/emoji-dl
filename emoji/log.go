package emoji

type log struct {
	Errors    bool
	Downloads bool
}

// Verbose is the verbosity of the download.
//
// Verbose.Errors determines if errors should be printed to stdout.
//
// Verbose.Downloads determines if files downloaded should be printed to
// stdout.
var Verbose = log{Errors: false, Downloads: false}
