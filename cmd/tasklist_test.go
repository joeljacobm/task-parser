package cmd

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestDecode(t *testing.T) {
	t.Run("decode task file, valid yaml, return no error", func(t *testing.T) {
		tl := TasksList{}
		f, err := os.Open("../testfiles/valid.yaml")
		assert.NoError(t, err)
		defer f.Close()
		assert.NoError(t, tl.decode(f))
	})
	t.Run("decode task file, invalid yaml, return error", func(t *testing.T) {
		tl := TasksList{}
		f, err := os.Open("../testfiles/invalid.yaml")
		assert.NoError(t, err)
		defer f.Close()
		assert.Error(t, tl.decode(f))
	})

}
func TestValidateAndRunTasks(t *testing.T) {
	t.Run("test create_file, path doesn't exist", func(t *testing.T) {
		tempDir := t.TempDir()
		task := Task{
			Name:        "create tmp file",
			TaskType:    "create_file",
			AbortOnFail: false,
			Args:        map[string]string{"path": tempDir + "/test.file"},
		}
		var tl TasksList
		tl.Tasks = append(tl.Tasks, task)
		err := tl.validateAndRunTasks()
		assert.NoError(t, err)
		assert.Equal(t, 0, len(tl.Unsuccessful))
		assert.Equal(t, 1, len(tl.Successful))
		assert.Equal(t, "create tmp file", tl.Successful[0])
	})

	t.Run("test create_file, path exists", func(t *testing.T) {
		f, _ := os.CreateTemp(t.TempDir(), "")
		task := Task{
			Name:        "create tmp file",
			TaskType:    "create_file",
			AbortOnFail: false,
			Args:        map[string]string{"path": f.Name()},
		}
		var tl TasksList
		tl.Tasks = append(tl.Tasks, task)
		err := tl.validateAndRunTasks()
		assert.NoError(t, err)
		assert.Equal(t, 1, len(tl.Unchanged))
		assert.Equal(t, 0, len(tl.Unsuccessful))
		assert.Equal(t, 0, len(tl.Successful))
	})

	t.Run("test put_content with content", func(t *testing.T) {
		f, _ := os.CreateTemp(t.TempDir(), "")
		task := Task{
			Name:        "put content",
			TaskType:    "put_content",
			AbortOnFail: false,
			Args:        map[string]string{"path": f.Name(), "content": "testing"},
		}

		var tl TasksList
		tl.Tasks = append(tl.Tasks, task)
		err := tl.validateAndRunTasks()
		contents, err := os.ReadFile(f.Name())
		assert.NoError(t, err)
		assert.Equal(t, 1, len(tl.Successful))
		assert.Equal(t, 0, len(tl.Unsuccessful))
		assert.Equal(t, 0, len(tl.Unchanged))
		assert.Equal(t, "testing", string(contents))
	})

	t.Run("test put_content with content and append", func(t *testing.T) {
		f, _ := os.CreateTemp(t.TempDir(), "")
		task := Task{
			Name:        "put content",
			TaskType:    "put_content",
			AbortOnFail: false,
			Args:        map[string]string{"path": f.Name(), "content": "testing", "append": "true"},
		}

		_, err := f.WriteString("initial msg ")
		assert.NoError(t, err)
		var tl TasksList
		tl.Tasks = append(tl.Tasks, task)
		err = tl.validateAndRunTasks()
		contents, err := os.ReadFile(f.Name())
		assert.NoError(t, err)
		assert.Equal(t, 1, len(tl.Successful))
		assert.Equal(t, 0, len(tl.Unsuccessful))
		assert.Equal(t, "initial msg testing", string(contents))
	})
	t.Run("test put_content with content and append set to false", func(t *testing.T) {
		f, _ := os.CreateTemp(t.TempDir(), "")
		task := Task{
			Name:        "put content",
			TaskType:    "put_content",
			AbortOnFail: false,
			Args:        map[string]string{"path": f.Name(), "content": "testing", "append": "false"},
		}

		_, err := f.WriteString("initial msg ")
		assert.NoError(t, err)
		var tl TasksList
		tl.Tasks = append(tl.Tasks, task)
		err = tl.validateAndRunTasks()
		contents, err := os.ReadFile(f.Name())
		assert.NoError(t, err)
		assert.Equal(t, 1, len(tl.Successful))
		assert.Equal(t, 0, len(tl.Unsuccessful))
		assert.Equal(t, "testing", string(contents))
	})

	t.Run("test put_content with no content", func(t *testing.T) {
		f, _ := os.CreateTemp(t.TempDir(), "")
		task := Task{
			Name:        "put content",
			TaskType:    "put_content",
			AbortOnFail: false,
			Args:        map[string]string{"path": f.Name(), "content": "testing", "append": "true"},
		}
		var tl TasksList
		tl.Tasks = append(tl.Tasks, task)
		err := tl.validateAndRunTasks()
		assert.NoError(t, err)
		assert.Equal(t, 1, len(tl.Successful))
		assert.Equal(t, 0, len(tl.Unsuccessful))
	})

}

func TestValidateAndRunMultipleTasks(t *testing.T) {
	t.Run("don't fail if the same task is executed twice", func(t *testing.T) {
		tempDir := t.TempDir()
		tl := TasksList{Tasks: []Task{
			{
				Name:        "create dir",
				TaskType:    "create_dir",
				AbortOnFail: true,
				Args:        map[string]string{"path": tempDir},
			},
			{
				Name:     "create file",
				TaskType: "create_file",
				Args:     map[string]string{"path": tempDir + "/test.file"},
			},
			{
				Name:        "create file2",
				TaskType:    "create_file",
				AbortOnFail: true,
				Args:        map[string]string{"path": tempDir + "/test.file"},
			},
			{
				Name:        "create file2",
				TaskType:    "rm_dir",
				AbortOnFail: true,
				Args:        map[string]string{"path": tempDir, "recursive": "true"},
			},
		}}
		err := tl.validateAndRunTasks()
		assert.NoError(t, err)
		assert.Equal(t, 2, len(tl.Successful))
		assert.Equal(t, 2, len(tl.Unchanged))
		assert.Equal(t, 0, len(tl.Unsuccessful))
	})

}
