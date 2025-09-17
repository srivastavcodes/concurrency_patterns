package main

import (
	"fmt"
	"runtime/debug"
)

type ModError struct {
	Message    string
	Inner      error
	StackTrace string

	Misc map[string]any
}

func (err ModError) Error() string {
	return err.Message
}

func wrapError(err error, messagef string, msgArgs ...any) ModError {
	return ModError{
		Message:    fmt.Sprintf(messagef, msgArgs),
		Inner:      err,
		StackTrace: string(debug.Stack()),
		Misc:       make(map[string]any),
	}
}
