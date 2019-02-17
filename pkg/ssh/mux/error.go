package mux

import "fmt"

type (
	// ExitStatus is used for passing along exit-codes from commands
	ExitStatus struct {
		Code int
		Err  error
	}
)

func (es ExitStatus) Error() string {
	return fmt.Sprintf("exit status: %d => %v", es.Code, es.Err)
}
