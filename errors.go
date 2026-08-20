package address

import "errors"

// ErrNoMatch is returned by Normalize when no street reaches the minimum
// confidence threshold.
var ErrNoMatch = errors.New("no matching street found")
