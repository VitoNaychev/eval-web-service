package service

import (
	"errors"
)

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

	err := e.recordExpressionError(expr, MethodValidate, interpErr)
	if err != nil {
		return false, err
	}

	return false, interpErr
}

func (e *EvalService) Execute(expr string) (int, error) {
	result, interpErr := e.interp.Execute(expr)
	if interpErr == nil {
		return result, nil
	}

	err := e.recordExpressionError(expr, MethodExecute, interpErr)
	if err != nil {
		return -1, err
	}

	return -1, interpErr
}

func (e *EvalService) GetExpressionErrors() ([]ExpressionError, error) {
	exprErrors, err := e.errorRepo.GetAll()
	if err != nil {
		return nil, NewEvalServiceError(err.Error())
	}

	return exprErrors, nil
}

func (e *EvalService) recordExpressionError(expr string, method MethodType, interpErr error) error {
	errorType, err := evalServiceErrorToErrorType(interpErr)
	if err != nil {
		return err
	}

	exprError := ExpressionError{
		Expression: expr,
		Method:     method,
		Type:       errorType,
	}

	repoErr := e.errorRepo.Increment(&exprError)
	if repoErr != nil {
		return NewEvalServiceError(repoErr.Error())
	}

	return nil
}

func evalServiceErrorToErrorType(err error) (ErrorType, error) {
	switch {
	case errors.Is(err, ErrNonMathQuestion):
		return ErrorTypeNonMathQuestion, nil
	case errors.Is(err, ErrUnsupportedOperation):
		return ErrorTypeUnsupportedOperand, nil
	case errors.Is(err, ErrInvalidSyntax):
		return ErrorTypeInvalidSyntax, nil
	default:
		return ErrorType(-1), NewUnsupportedInterpreterError(err.Error())
	}
}
