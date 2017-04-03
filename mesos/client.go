package mesos

import (
	"errors"
	"github.com/parnurzeal/gorequest"
)

type Mesos interface {
	SlaveState(slaveHost string) (SlaveState, error)
}

type mesosClient struct {
	Request *gorequest.SuperAgent
}

func NewMesosClient() Mesos {
	client := new(mesosClient)
	client.Request = gorequest.New()
	return client
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
