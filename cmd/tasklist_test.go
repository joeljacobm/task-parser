package cmd

import (
	"fmt"
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
		tempDir := os.TempDir()
		fmt.Println(tempDir)
		defer func() {
			os.RemoveAll(tempDir)
		}()
		task := Task{
			Name:        "create tmp file",
			TaskType:    "create_file",
			AbortOnFail: false,
			Args:        map[string]string{"path": tempDir + "/tmp"},
		}
		var tl TasksList
		tl.Tasks = append(tl.Tasks, task)
		success, unsuccess, err := tl.validateAndRunTasks()
		assert.NoError(t, err)
		assert.Equal(t, 0, len(unsuccess))
		assert.Equal(t, 1, len(success))
		assert.Equal(t, "create tmp file", success[0])
	})

	t.Run("test create_file, path exists", func(t *testing.T) {
		tempDir := os.TempDir()
		f, _ := os.CreateTemp(tempDir, "test.file")
		fmt.Println(tempDir)
		defer func() {
			os.RemoveAll(tempDir)
		}()
		task := Task{
			Name:        "create tmp file",
			TaskType:    "create_file",
			AbortOnFail: false,
			Args:        map[string]string{"path": f.Name()},
		}
		var tl TasksList
		tl.Tasks = append(tl.Tasks, task)
		success, unsuccess, err := tl.validateAndRunTasks()
		assert.NoError(t, err)
		assert.Equal(t, 1, len(unsuccess))
		assert.Equal(t, 0, len(success))
		assert.Equal(t, "create tmp file", unsuccess[0])
	})

	t.Run("test put_content with content", func(t *testing.T) {
		tempDir := os.TempDir()
		f, _ := os.CreateTemp(tempDir, "test.file")
		defer func() {
			os.RemoveAll(tempDir)
		}()
		task := Task{
			Name:        "put content",
			TaskType:    "put_content",
			AbortOnFail: false,
			Args:        map[string]string{"path": f.Name(), "content": "testing"},
		}

		var tl TasksList
		tl.Tasks = append(tl.Tasks, task)
		success, unsuccess, err := tl.validateAndRunTasks()
		contents, err := os.ReadFile(f.Name())
		assert.NoError(t, err)
		assert.Equal(t, 1, len(success))
		assert.Equal(t, 0, len(unsuccess))
		assert.Equal(t, "testing", string(contents))
	})

	t.Run("test put_content with content and append", func(t *testing.T) {
		tempDir := os.TempDir()
		f, _ := os.CreateTemp(tempDir, "test.file")
		defer func() {
			os.RemoveAll(tempDir)
		}()
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
		success, unsuccess, err := tl.validateAndRunTasks()
		contents, err := os.ReadFile(f.Name())
		assert.NoError(t, err)
		assert.Equal(t, 1, len(success))
		assert.Equal(t, 0, len(unsuccess))
		assert.Equal(t, "initial msg testing", string(contents))
	})
	t.Run("test put_content with content and append set to false", func(t *testing.T) {
		tempDir := os.TempDir()
		f, _ := os.CreateTemp(tempDir, "test.file")
		defer func() {
			os.RemoveAll(tempDir)
		}()
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
		success, unsuccess, err := tl.validateAndRunTasks()
		contents, err := os.ReadFile(f.Name())
		assert.NoError(t, err)
		assert.Equal(t, 1, len(success))
		assert.Equal(t, 0, len(unsuccess))
		assert.Equal(t, "testing", string(contents))
	})

	t.Run("test put_content with no content", func(t *testing.T) {
		tempDir := os.TempDir()
		f, _ := os.CreateTemp(tempDir, "test.file")
		defer func() {
			os.RemoveAll(tempDir)
		}()
		task := Task{
			Name:        "put content",
			TaskType:    "put_content",
			AbortOnFail: false,
			Args:        map[string]string{"path": f.Name(), "content": "testing", "append": "true"},
		}
		var tl TasksList
		tl.Tasks = append(tl.Tasks, task)
		success, unsuccess, err := tl.validateAndRunTasks()
		assert.NoError(t, err)
		assert.Equal(t, 1, len(success))
		assert.Equal(t, 0, len(unsuccess))
	})

}

func TestValidateAndRunTasksAborts(t *testing.T) {
	t.Run("test abort set to true", func(t *testing.T) {
		tempDir := os.TempDir()
		defer func() {
			os.RemoveAll(tempDir)
		}()
		task1 := Task{
			Name:        "create tmp file",
			TaskType:    "create_file",
			AbortOnFail: false,
			Args:        map[string]string{"path": tempDir + "/tmp1"},
		}
		task2 := Task{
			Name:        "create tmp file2",
			TaskType:    "create_file",
			AbortOnFail: true,
			Args:        map[string]string{"path": tempDir + "/tmp1"},
		}
		task3 := Task{
			Name:        "create tmp file3",
			TaskType:    "create_file",
			AbortOnFail: false,
			Args:        map[string]string{"path": tempDir + "/tmp3"},
		}
		var tl TasksList
		tl.Tasks = append(tl.Tasks, task1, task2, task3)
		fmt.Println(tl.Tasks)
		success, unsuccess, err := tl.validateAndRunTasks()
		assert.Error(t, err)
		assert.EqualError(t, err, "Aborted due to Task: create tmp file2 with error:path already exists")
		assert.Equal(t, 1, len(unsuccess))
		assert.Equal(t, 1, len(success))
		assert.Equal(t, "create tmp file", success[0])
		assert.Equal(t, "create tmp file2", unsuccess[0])
	})

	t.Run("test abort set to true", func(t *testing.T) {
		tempDir := os.TempDir()
		defer func() {
			os.RemoveAll(tempDir)
		}()
		task1 := Task{
			Name:        "create tmp file",
			TaskType:    "create_file",
			AbortOnFail: false,
			Args:        map[string]string{"path": tempDir + "/tmp1"},
		}
		task2 := Task{
			Name:        "create tmp file2",
			TaskType:    "create_file",
			AbortOnFail: false,
			Args:        map[string]string{"path": tempDir + "/tmp1"},
		}
		task3 := Task{
			Name:        "create tmp file3",
			TaskType:    "create_file",
			AbortOnFail: false,
			Args:        map[string]string{"path": tempDir + "/tmp3"},
		}
		var tl TasksList
		tl.Tasks = append(tl.Tasks, task1, task2, task3)
		fmt.Println(tl.Tasks)
		success, unsuccess, err := tl.validateAndRunTasks()
		assert.NoError(t, err)
		assert.Equal(t, 1, len(unsuccess))
		assert.Equal(t, 2, len(success))
		assert.Equal(t, "create tmp file", success[0])
		assert.Equal(t, "create tmp file3", success[1])
		assert.Equal(t, "create tmp file2", unsuccess[0])
	})
}
