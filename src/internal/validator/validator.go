package validator

type Validator struct{}
type Condition struct {
	Ok        bool
	FieldName string
	Message   string
}

func (v Validator) Validate(conditions ...Condition) map[string]string {
	errors := make(map[string]string)
	for _, c := range conditions {
		if !c.Ok {
			errors[c.FieldName] = c.Message
		}
	}
	return errors
}
