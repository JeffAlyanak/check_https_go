package check

// Result holds information about a completed check
type Result struct {
	URL          string // URL for request
	ReturnCode   int    // Code to return to OS after check
	Status       int    // HTTP status code
	Value        string // Result text value
	VerboseValue string // Additional, optional information
	Error        error  // Error during check
}
