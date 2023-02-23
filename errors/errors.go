package errors

import (
	"regexp"
)

// NotFoundError: returns true if the queried resource does not exist and the query returns a related error for that.
// If a resource doesn't exist, the query gives below error:
// Error: graphql: Resource not found ..."
func NotFoundError(err error) bool {
	notFoundErr := "(?i)Resource not found"
	expectedErr := regexp.MustCompile(notFoundErr)
	return expectedErr.Match([]byte(err.Error()))
}
