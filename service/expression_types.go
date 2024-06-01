package service

import "fmt"

type ExpressionServiceError struct {
	msg string
}

func NewExpressionServiceError(msg string) error {
	return &ExpressionServiceError{
		msg: msg,
	}
}

func (e *ExpressionServiceError) Error() string {
	return e.msg
}

var (
	ErrNonMathQuestion      = NewExpressionServiceError("non-math question")
	ErrUnsupportedOperation = NewExpressionServiceError("unsupported operation")
	ErrInvalidSyntax        = NewExpressionServiceError("invalid syntax")
)

type UnsupportedInterpreterError struct {
	msg string
}

func NewUnsupportedInterpreterError(msg string) error {
	return &UnsupportedInterpreterError{
		msg: fmt.Sprintf("unsupported interpreter error: %s", msg),
	}
}

func (u *UnsupportedInterpreterError) Error() string {
	return u.msg
}

func (u *UnsupportedInterpreterError) As(target interface{}) bool {
	if evalServiceError, ok := target.(**ExpressionServiceError); ok {
		*evalServiceError = &ExpressionServiceError{msg: u.msg}
		return true
	}
	return false
}

type Interpreter interface {
	Validate(string) (bool, error)
	Evaluate(string) (int, error)
}

type ExprErrorRepository interface {
	Increment(*ExpressionError) error
	GetAll() ([]ExpressionError, error)
}

type ErrorType int

const (
	ErrorTypeNonMathQuestion ErrorType = iota
	ErrorTypeUnsupportedOperand
	ErrorTypeInvalidSyntax
)

type MethodType int

const (
	MethodValidate = iota
	MethodEvaluate
)

type ExpressionError struct {
	Expression string
	Method     MethodType
	Frequency  int
	Type       ErrorType
}
