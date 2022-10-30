package validator

// The validator package can be used to validate tasks in a generic way.
// A task should implement validatableObject inroder to use the package.
type validatableObject interface {
	ValidateTaskType() error
	ValidateArguments() error
	HandleAbortOnFail(err error) bool
}

type Validator struct {
	Object      validatableObject
	ValidateFns []func() error
}

func GetNewValidator(obj validatableObject, fns []func() error) *Validator {
	return &Validator{Object: obj, ValidateFns: fns}
}

func (v *Validator) RunValidateFns() error {
	for _, val := range v.ValidateFns {
		if err := val(); err != nil {
			return err
		}
	}
	return nil
}
