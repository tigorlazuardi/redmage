package errs

import "errors"

func FindCodeOrDefault(err error, def int) int {
	for unwrap := errors.Unwrap(err); unwrap != nil; err = unwrap {
		if coder, ok := err.(interface{ GetCode() int }); ok {
			code := coder.GetCode()
			if code != 0 {
				def = code
			}
		}
	}
	return def
}
