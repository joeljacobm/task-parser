package cmd

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestIsExecutionRequired(t *testing.T) {
	t.Run("test create_dir, path doesn't exist, return true", func(t *testing.T) {
		task := Task{
			TaskType:    "create_dir",
			Args:        map[string]string{"path": "some/path"},
			AbortOnFail: true,
		}
		assert.Equal(t, true, task.IsExecutionRequired())
	})

	t.Run("test create_dir, path exists, return false", func(t *testing.T) {
		tempDir := t.TempDir()
		task := Task{
			TaskType:    "create_dir",
			Args:        map[string]string{"path": tempDir},
			AbortOnFail: true,
		}
		assert.Equal(t, false, task.IsExecutionRequired())
	})

	t.Run("test create_file, path exists, return false", func(t *testing.T) {
		tempDir := t.TempDir()
		f, _ := os.CreateTemp(tempDir, "test.file")
		task := Task{
			TaskType:    "create_file",
			Args:        map[string]string{"path": f.Name()},
			AbortOnFail: true,
		}
		assert.Equal(t, false, task.IsExecutionRequired())
	})

	t.Run("test create_file, path doesn't exist, return true", func(t *testing.T) {
		task := Task{
			TaskType:    "create_file",
			Args:        map[string]string{"path": "some/path"},
			AbortOnFail: true,
		}
		assert.Equal(t, true, task.IsExecutionRequired())
	})
	t.Run("test rm_dir, path exists, return true", func(t *testing.T) {
		tempDir := t.TempDir()
		task := Task{
			TaskType:    "rm_dir",
			Args:        map[string]string{"path": tempDir},
			AbortOnFail: true,
		}
		assert.Equal(t, true, task.IsExecutionRequired())

	})
	t.Run("test rm_dir, path doesnt exist, return false", func(t *testing.T) {
		task := Task{
			TaskType:    "rm_dir",
			Args:        map[string]string{"path": "some/path"},
			AbortOnFail: true,
		}
		assert.Equal(t, false, task.IsExecutionRequired())

	})
}
func TestValidateArguments(t *testing.T) {
	t.Run("test rm_dir, dir contains files, recursive arg not set, return error", func(t *testing.T) {
		tempDir := t.TempDir()
		_, _ = os.CreateTemp(tempDir, "test.file")
		task := Task{
			TaskType:    "rm_dir",
			Args:        map[string]string{"path": tempDir},
			AbortOnFail: true,
		}
		err := task.ValidateArguments()
		assert.Error(t, err)
		assert.EqualError(t, err, "Directory contains entries, recursive delete is required")
	})
	t.Run("test rm_dir, dir contains files, recursive arg set, return error", func(t *testing.T) {
		tempDir := t.TempDir()
		_, _ = os.CreateTemp(tempDir, "test.file")
		task := Task{
			TaskType:    "rm_dir",
			Args:        map[string]string{"path": tempDir, "recursive": "true"},
			AbortOnFail: true,
		}
		assert.NoError(t, task.ValidateArguments())
	})
}

func TestValidateTaskType(t *testing.T) {
	t.Run("test create_dir, valid type, return no error", func(t *testing.T) {
		task := Task{
			TaskType: "create_dir",
		}
		err := task.ValidateTaskType()
		assert.NoError(t, err)
	})
	t.Run("test create_directory, invalid type, return error", func(t *testing.T) {
		task := Task{
			TaskType: "create_directory",
		}
		err := task.ValidateTaskType()
		assert.Error(t, err)
	})
	t.Run("empty type, return type cannot be empty error", func(t *testing.T) {
		task := Task{
			TaskType: "",
		}
		assert.Equal(t, "type cannot be empty", task.ValidateTaskType().Error())
	})
}
