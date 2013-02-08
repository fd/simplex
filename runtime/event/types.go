package event

type (
	Event interface {
		Event() string
	}

	ErrorEvent struct{ Err error }
)

func Error(err error) ErrorEvent {
	return ErrorEvent{err}
}

func (err ErrorEvent) Event() string {
	return err.Error()
}

func (err ErrorEvent) Error() string {
	return err.Err.Error()
}
