package cmd

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"strings"
	"task-parser/validator"
)

type TasksList struct {
	Tasks        []Task
	Successful   []string
	Unsuccessful []string
	Unchanged    []string
}

// Run decodes file into a task,validates and runs the tasks.
func Run(filePath string) (string, error) {
	//Open the file using the path to the task file
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("cannot open the task file with error :%s", err.Error())
	}
	tl := TasksList{}
	if err = tl.decode(file); err != nil {
		return "", err
	}
	err = tl.validateAndRunTasks()
	if err != nil {
		return "", fmt.Errorf("Successful tasks: %s \nUnsuccessful tasks: %s \nUnChanged tasks: %s \n%s", strings.Join(tl.Successful, ","), strings.Join(tl.Unsuccessful, ","), strings.Join(tl.Unchanged, ","), err.Error())
	}
	return fmt.Sprintf("Successful tasks: %s \nUnsuccessful tasks: %s \nUnChanged tasks: %s", strings.Join(tl.Successful, ","), strings.Join(tl.Unsuccessful, ","), strings.Join(tl.Unchanged, ",")), nil
}

// decode decodes the file into a task
func (tl *TasksList) decode(i io.ReadWriter) error {
	// Decode the task file into a slice of tasks
	err := yaml.NewDecoder(i).Decode(&tl.Tasks)
	if err != nil {
		return fmt.Errorf("failed decoding the task file with error:%s", err.Error())
	}
	return nil
}

// validateAndRunTasks validates and calls the performTask method of the task object.
// A slice of successful, unsuccessful tasks and error are returned
// Program is aborted if a task with abortOnFail set to true fails
func (tl *TasksList) validateAndRunTasks() error {
	for _, t := range tl.Tasks {
		val := validator.GetNewValidator(&t, t.getValidatorFns())
		err := val.RunValidateFns()
		if err != nil {
			if val.Object.HandleAbortOnFail(err) {
				tl.Unsuccessful = append(tl.Unsuccessful, t.Name)
				return fmt.Errorf("Aborted due to Task: %s with error:%s", t.Name, err)
			}
			tl.Unsuccessful = append(tl.Unsuccessful, t.Name)
			continue
		}
		if !t.IsExecutionRequired() {
			tl.Unchanged = append(tl.Unchanged, t.Name)
			continue
		}
		err = t.performTask()
		if err != nil {
			if val.Object.HandleAbortOnFail(err) {
				tl.Unsuccessful = append(tl.Unsuccessful, t.Name)
				return fmt.Errorf("Aborted due to Task: %s with error:%s", t.Name, err)
			}
			tl.Unsuccessful = append(tl.Unsuccessful, t.Name)
			continue
		}
		tl.Successful = append(tl.Successful, t.Name)
	}

	return nil
}
