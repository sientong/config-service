package helper

type ValidationError struct {
	Msg string
}

func (e ValidationError) Error() string {
	return "validation error: " + e.Msg
}

func PanicIfError(err error) {
	if err != nil {
		panic(err)
	}
}
