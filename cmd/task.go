package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
)

type Task struct {
	Name        string            `yaml:"name,omitempty"`
	TaskType    string            `yaml:"type"`
	AbortOnFail bool              `yaml:"abortOnFail,omitempty"`
	Args        map[string]string `yaml:"args,omitempty"`
}

var supportedArguments = map[string]bool{
	"path":      true,
	"content":   true,
	"append":    true,
	"recursive": true,
}

// performTask() performs the given operation specified in the task entry.
func (t *Task) performTask() error {
	path := t.Args["path"]
	switch t.TaskType {
	case "create_dir":
		if err := os.Mkdir(path, os.ModePerm); err != nil {
			return err
		}
	case "create_file":
		f, err := os.OpenFile(path, os.O_CREATE, 0644)
		if err != nil {
			return err
		}
		defer f.Close()
	case "put_content":
		var f *os.File
		var err error
		appendArg, ok := t.getArgument("append")
		if ok && appendArg == "true" {
			f, err = os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return err
			}
		} else {
			f, err = os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
			if err != nil {
				return err
			}
		}
		if content := t.getContent(); content != "" {
			_, err := f.WriteString(content)
			if err != nil {
				return err
			}
		}
		defer f.Close()
	case "rm_dir":
		recurse, ok := t.getArgument("recursive")
		if ok && recurse == "true" {
			if err := os.RemoveAll(path); err != nil {
				return err
			}
			return nil
		}
		if err := os.Remove(path); err != nil {
			return err
		}
	case "rm_file":
		if err := os.Remove(path); err != nil {
			return err
		}

	}
	return nil
}

func (t *Task) getContent() string {
	if content, ok := t.Args["content"]; ok {
		return content
	}
	return ""
}

func (t *Task) getValidatorFns() []func() error {
	return []func() error{t.ValidateTaskType, t.ValidateArguments}
}

func (t *Task) ValidateTaskType() error {

	switch t.TaskType {
	case "create_dir", "create_file", "rm_file", "rm_dir", "put_content":
	case "":
		// type is a required field so it cannot be empty
		return errors.New("type cannot be empty")
	default:
		return errors.New("invalid type")
	}
	return nil
}

func (t *Task) ValidateArguments() error {
	for arg, _ := range t.Args {
		if !supportedArguments[arg] {
			return errors.New("unsupported argument")
		}
		path, ok := t.getArgument("path")
		if !ok {
			return errors.New("path cannot be empty")
		}
		_, err := os.Stat(path)
		switch t.TaskType {
		case "create_dir", "create_file":
			if err == nil {
				return errors.New("path already exists")
			}
		case "put_content":
			if os.IsNotExist(err) {
				return errors.New("path doesn't exist")
			}
		case "rm_dir":
			if os.IsNotExist(err) {
				return errors.New("path doesn't exist")
			}
			entry, err := os.ReadDir(path)
			if err != nil {
				return err
			}
			if len(entry) > 0 {
				_, ok := t.getArgument("recursive")
				if !ok {
					return errors.New("Directory contains entries, recursive delete is required")
				}
			}
		case "rm_file":
			if os.IsNotExist(err) {
				return errors.New("path doesn't exist")
			}
		}
	}
	return nil
}

func (t *Task) HandleAbortOnFail(err error) bool {
	if t.AbortOnFail {
		return true
	}
	log.Println(fmt.Sprintf("Task: %s failed with error:%s", t.Name, err.Error()))
	return false
}

func (t *Task) getArgument(arg string) (string, bool) {
	if argValue, ok := t.Args[arg]; ok {
		return argValue, true
	}
	return "", false
}
