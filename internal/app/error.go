package app

type StatusError struct {
	Code int
	Err  error
}

func (err StatusError) Error() string {
	return err.Err.Error()
}

func (err StatusError) Status() int {
	return err.Code
}
