package mesos

import "errors"

type Mesos interface {
	SlaveState(slaveHost string) (SlaveState, error)
}

type MesosClient struct {
}

func combineErrors(errs []error) error {
	if len(errs) == 1 {
		return errs[0]
	} else if len(errs) > 1 {
		msg := "Error(s):"
		for _, err := range errs {
			msg += " " + err.Error()
		}
		return errors.New(msg)
	} else {
		return nil
	}
}
