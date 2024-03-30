package response_errors

import "errors"

var InvalidUserIDErr = errors.New("user ID must be more than 0")
var InvalidUsernameErr = errors.New("username is required")
