package infrastructure

import "fmt"

type ErrUnknownJob struct {
	Opcode uint8
}

func (e *ErrUnknownJob) Error() string {
	return fmt.Sprintf("this job with opcode %d is unknown", e.Opcode)
}
