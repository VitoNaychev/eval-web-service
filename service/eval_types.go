package service

import "fmt"

type EvalServiceError struct {
	msg string
}

func NewEvalServiceError(msg string) error {
	return &EvalServiceError{
		msg: msg,
	}
}

func (e *EvalServiceError) Error() string {
	return e.msg
}

var (
	ErrNonMathQuestion      = NewEvalServiceError("non-math question")
	ErrUnsupportedOperation = NewEvalServiceError("unsupported operation")
	ErrInvalidSyntax        = NewEvalServiceError("invalid syntax")
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
	if evalServiceError, ok := target.(**EvalServiceError); ok {
		*evalServiceError = &EvalServiceError{msg: u.msg}
		return true
	}
	return false
}

type Interpreter interface {
	Validate(string) (bool, error)
	Execute(string) (int, error)
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
	MethodExecute
)

type ExpressionError struct {
	Expression string
	Method     MethodType
	Frequency  int
	Type       ErrorType
}
