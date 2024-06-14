package msgpackfix

type MultiError struct {
	Errors []error
}

func (m MultiError) Error() string {
	var s string
	for i := range m.Errors {
		s += m.Errors[i].Error() + "\n"
	}
	return s[:len(s)-1]
}

func appendError(curErr *error, newErr error) {
	if *curErr == nil {
		*curErr = newErr
		return
	}

	multiErr, ok := (*curErr).(MultiError)
	if ok {
		multiErr.Errors = append(multiErr.Errors, newErr)
		*curErr = &multiErr
		return
	}

	*curErr = MultiError{
		[]error{*curErr,
			newErr,
		},
	}
}
