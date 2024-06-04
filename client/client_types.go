package client

type ClientError struct {
	msg string
}

func NewClientError(msg string) error {
	return &ClientError{
		msg: msg,
	}
}

func (c *ClientError) Error() string {
	return c.msg
}

var (
	ErrNonMathQuestion      = NewClientError("non-math question")
	ErrUnsupportedOperation = NewClientError("unsupported operation")
	ErrInvalidSyntax        = NewClientError("invalid syntax")
)

type ExpressionError struct {
	Expression string
	Method     string
	Frequency  int
	Type       string
}
