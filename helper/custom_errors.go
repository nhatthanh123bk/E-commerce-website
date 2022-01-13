package helper

import "errors"

var ErrUnAuthorization = errors.New("unAuthorization")
var ErrBadRequest = errors.New("bad request")
var ErrInternal = errors.New("internal error")
