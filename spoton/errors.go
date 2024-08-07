package spoton

import "fmt"

type ErrBadRequest struct {
	reqtype string
	code    int
}

func (e *ErrBadRequest) Error() string {
	return fmt.Sprintf("bad %s request, got code %d, 200 expected", e.reqtype, e.code)
}

type ErrNoAuth struct{}

func (e *ErrNoAuth) Error() string {
	return "auth token is empty - have you called SetAuth previously?"
}
