package service_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/VitoNaychev/eval-web-service/service"
	"github.com/VitoNaychev/eval-web-service/testutil/assert"
)

type StubInterpreter struct {
	isValid bool
	result  int
	err     error
}

func (s *StubInterpreter) Validate(q string) (bool, error) {
	return s.isValid, s.err
}

func (s *StubInterpreter) Evaluate(q string) (int, error) {
	return s.result, s.err
}

func (s *StubInterpreter) Exec(q string) (int, error) {
	return 0, s.err
}

type StubErrorRepository struct {
	exprErrors []service.ExpressionError
	err        error

	spyExprError service.ExpressionError
}

func (s *StubErrorRepository) Increment(exprError *service.ExpressionError) error {
	s.spyExprError = *exprError
	return s.err
}

func (s *StubErrorRepository) GetAll() ([]service.ExpressionError, error) {
	return s.exprErrors, s.err
}

func TestValidate(t *testing.T) {
	t.Run("returns true on valid expression", func(t *testing.T) {
		expression := "What is 5?"
		wantValid := true

		interp := &StubInterpreter{
			isValid: wantValid,
		}
		repo := &StubErrorRepository{}
		exprSvc := service.NewExpressionService(interp, repo)

		gotValid, err := exprSvc.Validate(expression)
		assert.RequireNoError(t, err)

		assert.Equal(t, gotValid, wantValid)
	})

	t.Run("returns interpreter error on invalid expression", func(t *testing.T) {
		expression := "example expression"
		wantValid := false
		wantErr := service.ErrNonMathQuestion

		interp := &StubInterpreter{
			isValid: wantValid,
			err:     wantErr,
		}
		repo := &StubErrorRepository{}
		exprSvc := service.NewExpressionService(interp, repo)

		gotValid, gotErr := exprSvc.Validate(expression)
		assert.Equal(t, gotValid, wantValid)
		assert.Equal(t, gotErr, wantErr)
	})

	t.Run("persists invalid expression in repository", func(t *testing.T) {
		expression := "example expression"
		isValid := false
		err := service.ErrNonMathQuestion
		wantExprError := service.ExpressionError{
			Expression: expression,
			Method:     service.MethodValidate,
			Type:       service.ErrorTypeNonMathQuestion,
		}

		interp := &StubInterpreter{
			isValid: isValid,
			err:     err,
		}
		repo := &StubErrorRepository{}
		exprSvc := service.NewExpressionService(interp, repo)

		_, _ = exprSvc.Validate(expression)

		assert.Equal(t, repo.spyExprError, wantExprError)
	})

	t.Run("returns UnsupportedInterpreterError on unknown error from interpreter", func(t *testing.T) {
		expression := "example expression"
		isValid := false
		err := errors.New("unsupported error")
		wantErrMessage := fmt.Sprintf("unsupported interpreter error: %s", err.Error())

		interp := &StubInterpreter{
			isValid: isValid,
			err:     err,
		}
		repo := &StubErrorRepository{}
		exprSvc := service.NewExpressionService(interp, repo)

		_, gotErr := exprSvc.Validate(expression)

		assert.ErrorType[*service.UnsupportedInterpreterError](t, gotErr)
		assert.Equal(t, gotErr.Error(), wantErrMessage)
	})

	t.Run("wraps repository errors in EvalServiceError", func(t *testing.T) {
		expression := "example expression"
		isValid := false
		interpErr := service.ErrNonMathQuestion
		repoErrMessage := "repo error"

		interp := &StubInterpreter{
			isValid: isValid,
			err:     interpErr,
		}
		repo := &StubErrorRepository{
			err: errors.New(repoErrMessage),
		}
		exprSvc := service.NewExpressionService(interp, repo)

		_, gotErr := exprSvc.Validate(expression)

		assert.ErrorType[*service.ExpressionServiceError](t, gotErr)
		assert.Equal(t, gotErr.Error(), repoErrMessage)
	})
}

func TestEvaluate(t *testing.T) {
	t.Run("returns result on valid expression", func(t *testing.T) {
		expression := "What is 5?"
		wantResult := 5

		interp := &StubInterpreter{
			result: wantResult,
		}
		repo := &StubErrorRepository{}
		exprSvc := service.NewExpressionService(interp, repo)

		gotResult, err := exprSvc.Evaluate(expression)
		assert.RequireNoError(t, err)

		assert.Equal(t, gotResult, wantResult)
	})

	t.Run("returns interpreter error on invalid expression", func(t *testing.T) {
		expression := "example expression"
		wantErr := service.ErrNonMathQuestion

		interp := &StubInterpreter{
			err: wantErr,
		}
		repo := &StubErrorRepository{}
		exprSvc := service.NewExpressionService(interp, repo)

		_, gotErr := exprSvc.Evaluate(expression)
		assert.Equal(t, gotErr, wantErr)
	})

	t.Run("persists invalid expression in repository", func(t *testing.T) {
		expression := "example expression"
		err := service.ErrNonMathQuestion
		wantExprError := service.ExpressionError{
			Expression: expression,
			Method:     service.MethodEvaluate,
			Type:       service.ErrorTypeNonMathQuestion,
		}

		interp := &StubInterpreter{
			err: err,
		}
		repo := &StubErrorRepository{}
		exprSvc := service.NewExpressionService(interp, repo)

		_, _ = exprSvc.Evaluate(expression)

		assert.Equal(t, repo.spyExprError, wantExprError)
	})

	t.Run("returns UnsupportedInterpreterError on unknown error from interpreter", func(t *testing.T) {
		expression := "example expression"
		err := errors.New("unsupported error")
		wantErrMessage := fmt.Sprintf("unsupported interpreter error: %s", err.Error())

		interp := &StubInterpreter{
			err: err,
		}
		repo := &StubErrorRepository{}
		exprSvc := service.NewExpressionService(interp, repo)

		_, gotErr := exprSvc.Evaluate(expression)

		assert.ErrorType[*service.UnsupportedInterpreterError](t, gotErr)
		assert.Equal(t, gotErr.Error(), wantErrMessage)
	})

	t.Run("wraps repository errors in EvalServiceError", func(t *testing.T) {
		expression := "example expression"
		interpErr := service.ErrNonMathQuestion
		repoErrMessage := "repo error"

		interp := &StubInterpreter{
			err: interpErr,
		}
		repo := &StubErrorRepository{
			err: errors.New(repoErrMessage),
		}
		exprSvc := service.NewExpressionService(interp, repo)

		_, gotErr := exprSvc.Evaluate(expression)

		assert.ErrorType[*service.ExpressionServiceError](t, gotErr)
		assert.Equal(t, gotErr.Error(), repoErrMessage)
	})
}

func TestGetExpressionErrors(t *testing.T) {
	t.Run("returns all recorded expressions", func(t *testing.T) {
		wantExprErrors := []service.ExpressionError{
			{
				Expression: "example expression",
				Frequency:  3,
				Method:     service.MethodEvaluate,
				Type:       service.ErrorTypeNonMathQuestion,
			},
		}

		interp := &StubInterpreter{}
		repo := &StubErrorRepository{
			exprErrors: wantExprErrors,
		}
		exprSvc := service.NewExpressionService(interp, repo)

		gotExprErrors, err := exprSvc.GetExpressionErrors()
		assert.RequireNoError(t, err)

		assert.Equal(t, gotExprErrors, wantExprErrors)
	})

	t.Run("wraps repository errors in EvalServiceError", func(t *testing.T) {
		wantErrMessage := "repo error"

		interp := &StubInterpreter{}
		repo := &StubErrorRepository{
			err: errors.New(wantErrMessage),
		}
		exprSvc := service.NewExpressionService(interp, repo)

		_, gotErr := exprSvc.GetExpressionErrors()

		assert.ErrorType[*service.ExpressionServiceError](t, gotErr)
		assert.Equal(t, gotErr.Error(), wantErrMessage)
	})
}
