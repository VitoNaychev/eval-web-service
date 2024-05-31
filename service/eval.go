package service

import (
	"errors"
	"fmt"
)

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

func (u *UnsupportedInterpreterError) As(err any) bool {
	if evalServiceError, ok := err.(**EvalServiceError); ok {
		(*evalServiceError).msg = u.msg

		return true
	}
	return false
}

type Interpreter interface {
	Validate(string) (bool, error)
	Exec(string) (int, error)
}

type ErrorRepository interface {
	Increment(*ExpressionError) error
}

type EvalService struct {
	interp    Interpreter
	errorRepo ErrorRepository
}

func NewEvalService(interp Interpreter, errorRepo ErrorRepository) *EvalService {
	return &EvalService{
		interp:    interp,
		errorRepo: errorRepo,
	}
}

func (e *EvalService) Validate(expr string) (bool, error) {
	isValid, interpErr := e.interp.Validate(expr)
	if isValid {
		return isValid, nil
	}

	errorType, err := evalServiceErrorToErrorType(interpErr)
	if err != nil {
		return isValid, err
	}

	exprError := ExpressionError{
		Expression: expr,
		Method:     MethodValidate,
		Type:       errorType,
	}

	repoErr := e.errorRepo.Increment(&exprError)
	if repoErr != nil {
		return isValid, NewEvalServiceError(repoErr.Error())
	}

	return isValid, interpErr
}

func evalServiceErrorToErrorType(err error) (ErrorType, error) {
	if errors.Is(err, ErrNonMathQuestion) {
		return ErrorTypeNonMathQuestion, nil
	} else if errors.Is(err, ErrUnsupportedOperation) {
		return ErrorTypeUnsupportedOperand, nil
	} else if errors.Is(err, ErrInvalidSyntax) {
		return ErrorTypeInvalidSyntax, nil
	}

	return ErrorType(-1), NewUnsupportedInterpreterError(err.Error())
}
