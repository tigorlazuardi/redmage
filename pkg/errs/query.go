package errs

import "errors"

func FindCodeOrDefault(err error, def int) int {
	unwrap := errors.Unwrap(err)
	for unwrap != nil {
		if coder, ok := err.(interface{ GetCode() int }); ok {
			code := coder.GetCode()
			if code != 0 {
				return code
			}
		}
		unwrap = errors.Unwrap(unwrap)
	}

	return def
}

func FindMessage(err error) string {
	unwrap := errors.Unwrap(err)
	for unwrap != nil {
		if messager, ok := err.(interface{ GetMessage() string }); ok {
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
	return code, FindMessage(err)
}
