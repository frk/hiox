package hxio

type WriteError struct {
	Err error
}

func (we WriteError) Error() string {
	return we.Err.Error()
}

type ReadError struct {
	Err error
}

func (re ReadError) Error() string {
	return re.Err.Error()
}
