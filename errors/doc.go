// Package errors provides utilities to augment errors
//
// Architecture of a good error:
// * system information - err Error() - displayed in logs
// * user message - displayed to end user
// * status code - returned to other systems
// * error code - standard error code usable within an application
package errors
