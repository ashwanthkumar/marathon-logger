package mesos

import (
	"encoding/json"
	"time"

	"github.com/parnurzeal/gorequest"
)

type SlaveState struct {
	Attributes map[string]string `json:"attributes"`
	Frameworks []Framework       `json:"frameworks"`
}

func (s *SlaveState) FindExecutor(taskId string) *Executor {
	for _, framework := range s.Frameworks {
		for _, executor := range framework.Executors {
			if executor.Source == taskId {
				return &executor
			}
		}
	}

	return nil
}

type Framework struct {
	Checkpoint bool       `json:"checkpoint"`
	Executors  []Executor `json:"executors"`
}

type Executor struct {
	Container string `json:"container"`
	Directory string `json:"directory"`
	Id        string `json:"id"`
	Source    string `json:"source"`
}

func (m *MesosClient) SlaveState(slaveHost string) (SlaveState, error) {
	request := gorequest.New()
	response, body, errs := request.
		Get(slaveHost).
		Timeout(10 * time.Minute).
		End()

	var slaveState SlaveState
	if response != nil {
		if response.StatusCode == 200 && body != "" {
			err := json.Unmarshal([]byte(body), &slaveState)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}
	return slaveState, combineErrors(errs)
}
