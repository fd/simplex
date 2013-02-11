package event

type (
	Event interface {
		Event() string
	}

	Error interface {
		Event
		error
	}

	error_event struct{ Err error }
)

func NewError(err error) Error {
	return error_event{err}
}

func (err error_event) Event() string {
	return err.Error()
}

func (err error_event) Error() string {
	return err.Err.Error()
}
