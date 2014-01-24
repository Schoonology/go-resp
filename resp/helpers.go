package resp

import (
	"fmt"
	"strings"
)

type Error struct {
	Type    string
	Message string
}

func (err Error) Error() string {
	return fmt.Sprintf("%v: %v", err.Type, err.Message)
}

func parseError(str string) error {
	pieces := strings.SplitN(str, " ", 2)
	return Error{pieces[0], pieces[1]}
}

type Status struct {
	Message string
}

func NewMulti(vals ...interface{}) []interface{} {
	return vals
}
