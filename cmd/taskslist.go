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
	Tasks []Task
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
	successfulTasks, unsuccessfulTasks, err := tl.validateAndRunTasks()
	if err != nil {
		return "", fmt.Errorf("Successfull tasks: %s \nUnsuccessful tasks: %s \n%s", strings.Join(successfulTasks, ","), strings.Join(unsuccessfulTasks, ","), err.Error())
	}
	return fmt.Sprintf("Successfull tasks: %s \nUnsuccessful tasks: %s", strings.Join(successfulTasks, ","), strings.Join(unsuccessfulTasks, ",")), nil
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
func (tl *TasksList) validateAndRunTasks() ([]string, []string, error) {
	var successfulTasks, unsuccessfulTasks []string
	for _, t := range tl.Tasks {
		val := validator.GetNewValidator(&t, t.getValidatorFns())
		err := val.RunValidateFns()
		if err != nil {
			if val.Object.HandleAbortOnFail(err) {
				unsuccessfulTasks = append(unsuccessfulTasks, t.Name)
				return successfulTasks, unsuccessfulTasks, fmt.Errorf("Aborted due to Task: %s with error:%s", t.Name, err)
			}
			unsuccessfulTasks = append(unsuccessfulTasks, t.Name)
			continue
		}
		err = t.performTask()
		if err != nil {
			if val.Object.HandleAbortOnFail(err) {
				unsuccessfulTasks = append(unsuccessfulTasks, t.Name)
				return successfulTasks, unsuccessfulTasks, fmt.Errorf("Aborted due to Task: %s with error:%s", t.Name, err)
			}
			unsuccessfulTasks = append(unsuccessfulTasks, t.Name)
			continue
		}
		successfulTasks = append(successfulTasks, t.Name)
	}

	return successfulTasks, unsuccessfulTasks, nil
}
