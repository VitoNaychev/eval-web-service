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
	err     error
}

func (s *StubInterpreter) Validate(q string) (bool, error) {
	return s.isValid, s.err
}

func (s *StubInterpreter) Exec(q string) (int, error) {
	return 0, s.err
}

type StubErrorRepository struct {
	spyExprError service.ExpressionError
	err          error
}

func (s *StubErrorRepository) Increment(exprError *service.ExpressionError) error {
	s.spyExprError = *exprError
	return s.err
}

func TestEvalService(t *testing.T) {
	t.Run("returns true on valid expression", func(t *testing.T) {
		expression := "What is 5?"
		wantValid := true

		interp := &StubInterpreter{
			isValid: wantValid,
		}
		repo := &StubErrorRepository{}
		evalSvc := service.NewEvalService(interp, repo)

		gotValid, err := evalSvc.Validate(expression)
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
		evalSvc := service.NewEvalService(interp, repo)

		gotValid, gotErr := evalSvc.Validate(expression)
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
		evalSvc := service.NewEvalService(interp, repo)

		_, _ = evalSvc.Validate(expression)

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
		evalSvc := service.NewEvalService(interp, repo)

		_, gotErr := evalSvc.Validate(expression)

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
		evalSvc := service.NewEvalService(interp, repo)

		_, gotErr := evalSvc.Validate(expression)

		assert.ErrorType[*service.EvalServiceError](t, gotErr)
		assert.Equal(t, gotErr.Error(), repoErrMessage)
	})
}
