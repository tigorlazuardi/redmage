package errs

import (
	"errors"
)

func FindCodeOrDefault(err error, def int) int {
	if coder, ok := err.(interface{ GetCode() int }); ok {
		code := coder.GetCode()
		if code != 0 {
			return code
		}
	}

	for unwrap := errors.Unwrap(err); unwrap != nil; unwrap = errors.Unwrap(unwrap) {
		if coder, ok := unwrap.(interface{ GetCode() int }); ok {
			code := coder.GetCode()
			if code != 0 {
				return code
			}
		}
	}

	return def
}

func FindMessage(err error) string {
	if messager, ok := err.(interface{ GetMessage() string }); ok {
		message := messager.GetMessage()
		if message != "" {
			return message
		}
	}

	for unwrap := errors.Unwrap(err); unwrap != nil; unwrap = errors.Unwrap(unwrap) {
		if messager, ok := unwrap.(interface{ GetMessage() string }); ok {
			message := messager.GetMessage()
			if message != "" {
				return message
			}
		}
	}

	return err.Error()
}

func HTTPMessage(err error) (code int, message string) {
	code = FindCodeOrDefault(err, 500)
	if code >= 500 {
		return code, err.Error()
	}
	message = FindMessage(err)
	return code, message
}
