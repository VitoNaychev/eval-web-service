package service

import (
	"errors"
)

type ExpressionService struct {
	interp        Interpreter
	exprErrorRepo ExprErrorRepository
}

func NewExpressionService(interp Interpreter, exprErrorRepo ExprErrorRepository) *ExpressionService {
	return &ExpressionService{
		interp:        interp,
		exprErrorRepo: exprErrorRepo,
	}
}

func (e *ExpressionService) Validate(expr string) (bool, error) {
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

func (e *ExpressionService) Evaluate(expr string) (int, error) {
	result, interpErr := e.interp.Evaluate(expr)
	if interpErr == nil {
		return result, nil
	}

	err := e.recordExpressionError(expr, MethodEvaluate, interpErr)
	if err != nil {
		return -1, err
	}

	return -1, interpErr
}

func (e *ExpressionService) GetExpressionErrors() ([]ExpressionError, error) {
	exprErrors, err := e.exprErrorRepo.GetAll()
	if err != nil {
		return nil, NewExpressionServiceError(err.Error())
	}

	return exprErrors, nil
}

func (e *ExpressionService) recordExpressionError(expr string, method MethodType, interpErr error) error {
	errorType, err := evalServiceErrorToErrorType(interpErr)
	if err != nil {
		return err
	}

	exprError := ExpressionError{
		Expression: expr,
		Method:     method,
		Type:       errorType,
	}

	repoErr := e.exprErrorRepo.Increment(&exprError)
	if repoErr != nil {
		return NewExpressionServiceError(repoErr.Error())
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
