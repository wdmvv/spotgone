package spoton

import "fmt"

type ErrBadAuth struct {
	code int
}

func (e *ErrBadAuth) Error() string {
	return fmt.Sprintf("bad auth token response, got code %d, 200 expected", e.code)
}
