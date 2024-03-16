package main

import "regexp"

type errorWrapper struct {
	wrapped error
	message string
}

func wrapError(out string, err error) error {
	return errorWrapper{err, out}
}

func (e errorWrapper) Error() string {
	return e.message
}

func (e errorWrapper) Unwrap() error {
	return e.wrapped
}

func prepend[T any](item T, slice []T) []T {
	return append([]T{item}, slice...)
}

var flagRgx = regexp.MustCompile(`^-{0,2}([a-z\-]+)=?["']?([^"']+)?["']?$`)

func IsFlag(s, flagName string) bool {
	if found := flagRgx.FindAllStringSubmatch(s, 1); len(found) > 0 {
		return found[0][1] == flagName
	}
	return false
}

func ParseFlag(s string) (name, value string, valid bool) {
	if found := flagRgx.FindAllStringSubmatch(s, 1); len(found) > 0 {
		valid = true
		name = found[0][1]
		if len(found[0]) > 1 {
			value = found[0][2]
		}
	}
	return
}
