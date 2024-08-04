package errs

import (
	"errors"
	"net/http"

	"connectrpc.com/connect"
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

func ExtractConnectCode(err error) connect.Code {
	code := FindCodeOrDefault(err, 500)
	if code >= 500 {
		return connect.CodeInternal
	}
	switch code {
	case http.StatusNotFound:
		return connect.CodeNotFound
	case http.StatusConflict:
		return connect.CodeAlreadyExists
	case http.StatusBadRequest:
		return connect.CodeInvalidArgument
	case http.StatusForbidden:
		return connect.CodePermissionDenied
	case http.StatusFailedDependency:
		return connect.CodeUnavailable
	case http.StatusTooManyRequests:
		return connect.CodeResourceExhausted
	default:
		return connect.CodeUnknown
	}
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

func HasCode(err error, code int) bool {
	return FindCodeOrDefault(err, http.StatusInternalServerError) == code
}

func Source(err error) error {
	if err == nil {
		return nil
	}
	if _, ok := err.(Error); !ok {
		return err
	}
	for unwrap := errors.Unwrap(err); unwrap != nil; unwrap = errors.Unwrap(unwrap) {
		if _, ok := unwrap.(Error); !ok {
			return unwrap
		}
	}
	return err
}
