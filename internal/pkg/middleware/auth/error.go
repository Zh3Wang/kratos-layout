package auth

import "github.com/go-kratos/kratos/v2/errors"

var (
	errAuthFail      = errors.New(401, "AUTHENTICATION FAILED", "Missing token or token incorrect")
	errAuthTypeError = errors.New(401, "AUTHENTICATION FAILED", "token type error")
)
