package validator

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockTask struct {
	TaskType string
}

func (mt *mockTask) ValidateTaskType() error {
	switch mt.TaskType {
	case "add":
	default:
		return errors.New("invalid type")

	}
	return nil
}

func (mt *mockTask) ValidateArguments() error {
	return nil
}

func (mt *mockTask) HandleAbortOnFail(err error) bool {
	return false
}

func TestRunValidateFns(t *testing.T) {
	t.Run("valid task type", func(t *testing.T) {
		var mt mockTask
		mt.TaskType = "add"
		assert.NoError(t, GetNewValidator(&mt, []func() error{mt.ValidateTaskType, mt.ValidateArguments}).RunValidateFns())
	})
	t.Run("invalid task type", func(t *testing.T) {
		var mt mockTask
		mt.TaskType = "remove"
		assert.EqualError(t, GetNewValidator(&mt, []func() error{mt.ValidateTaskType, mt.ValidateArguments}).RunValidateFns(), "invalid type")
	})

}
