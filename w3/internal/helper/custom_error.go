package helper

type CustomError struct {
	Message string
	Code    int
}

func (c CustomError) Error() string {

	return c.Message
}

func (c CustomError) HTTPStatusCode() int {
	return c.Code
}
