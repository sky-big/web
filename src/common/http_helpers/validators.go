package http_helpers

type Validator interface {
	/**
	 	 * Resource specific validation logic
	 	 * Implementation should follow return contracts
		 * @param resource, to validate resource
		 * @return dependentErr, not nil if call dependencies error
	 	 * @return ok, return true id resource content is valid, false else
		 * @return msg, explain why resource content is not valid
	*/
	Do(resource interface{}) (dependentErr error, ok bool, msg string)
}
type CompositionalValidator struct {
	validators []Validator
}

func (t *CompositionalValidator) Add(v Validator) *CompositionalValidator {
	if t.validators == nil {
		t.validators = []Validator{}
	}
	t.validators = append(t.validators, v)
	return t
}

func (t *CompositionalValidator) Do(resource interface{}) (dependentErr error, ok bool, msg string) {
	for _, v := range t.validators {
		err, ok, msg := v.Do(resource)
		if err != nil || !ok {
			return err, ok, msg
		}
	}
	return nil, true, ""
}
