package utils

// NoErr panics if the given error is not nil.
// Use this for errors that should never occur in normal operation
// (e.g., during application initialization or configuration loading).
//
// Example:
//
//	cfg, err := config.Load()
//	utils.NoErr(err) // panics if err != nil
func NoErr(err error) {
	if err != nil {
		panic(err)
	}
}

// NoErrWithMsg panics with a custom message if the given error is not nil.
// Useful when you want to provide additional context about where the error occurred.
//
// Example:
//
//	utils.NoErrWithMsg(err, "failed to connect to database")
func NoErrWithMsg(err error, msg string) {
	if err != nil {
		panic(msg + ": " + err.Error())
	}
}
