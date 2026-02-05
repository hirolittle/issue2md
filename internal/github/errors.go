package github

import "errors"

var (
	ErrResourceNotFound = errors.New("resource not found")
	ErrAccessDenied     = errors.New("access denied: check token permissions")
	ErrRateLimited      = errors.New("rate limited: please provide a token")
)
