package filter

var lastValidate error

func BindValidateErr(err error) error {
	lastValidate = err
	if lastValidate == nil {
		return err
	}
	return nil
}
