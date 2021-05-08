package handlers

import "errors"

var ErrEOF = errors.New("EOF")
var ErrUnexpectedEOF = errors.New("unexpected EOF, value required")
var ErrNotAllowed = errors.New("not alowed")
var ErrMoveOn = errors.New("move on")
