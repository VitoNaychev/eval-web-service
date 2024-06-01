package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VitoNaychev/eval-web-service/handler"
	"github.com/VitoNaychev/eval-web-service/service"
	"github.com/VitoNaychev/eval-web-service/testutil/assert"
)

type StubExpressionService struct {
	result     int
	isValid    bool
	exprErrors []service.ExpressionError
	err        error
}

func (s *StubExpressionService) Evaluate(expression string) (int, error) {
	return s.result, s.err
}

func (s *StubExpressionService) Validate(expression string) (bool, error) {
	return s.isValid, s.err
}

func (s *StubExpressionService) GetExpressionErrors() ([]service.ExpressionError, error) {
	return s.exprErrors, s.err
}

func TestEvaluate(t *testing.T) {
	t.Run("evaluates expression and returns EvaluateResponse", func(t *testing.T) {
		expression := "What is 5 plus 3?"
		result := 8

		evalRequest := handler.ExpressionRequest{
			Expression: expression,
		}
		wantResponse := handler.EvaluateResponse{
			Result: result,
		}

		body := bytes.NewBuffer([]byte{})
		json.NewEncoder(body).Encode(evalRequest)

		request, _ := http.NewRequest(http.MethodGet, "/", body)
		response := httptest.NewRecorder()

		exprService := &StubExpressionService{
			result: result,
		}
		exprHandler := handler.NewExpressionHandler(exprService)

		exprHandler.Evaluate(response, request)
		assert.Equal(t, response.Code, http.StatusOK)

		var gotResponse handler.EvaluateResponse
		json.NewDecoder(response.Body).Decode(&gotResponse)

		assert.Equal(t, gotResponse, wantResponse)
	})

	t.Run("reutrns Status Bad Request and ErrorResponse on invalid expression", func(t *testing.T) {
		expression := "What is 5 plus 3?"
		errorMessage := "handler error"

		evalRequest := handler.ExpressionRequest{
			Expression: expression,
		}
		wantErrorResponse := handler.ErrorResponse{
			Error: errorMessage,
		}

		body := bytes.NewBuffer([]byte{})
		json.NewEncoder(body).Encode(evalRequest)

		request, _ := http.NewRequest(http.MethodGet, "/", body)
		response := httptest.NewRecorder()

		exprService := &StubExpressionService{
			err: errors.New(errorMessage),
		}
		exprHandler := handler.NewExpressionHandler(exprService)

		exprHandler.Evaluate(response, request)
		assert.Equal(t, response.Code, http.StatusBadRequest)

		var gotErrorResponse handler.ErrorResponse
		json.NewDecoder(response.Body).Decode(&gotErrorResponse)

		assert.Equal(t, gotErrorResponse, wantErrorResponse)
	})
}

func TestValidate(t *testing.T) {

	t.Run("sets Valid field to true on valid expression", func(t *testing.T) {
		expression := "What is 5 plus 3?"
		wantIsValid := true

		evalRequest := handler.ExpressionRequest{
			Expression: expression,
		}
		wantResponse := handler.ValidateResponse{
			Valid: true,
		}

		body := bytes.NewBuffer([]byte{})
		json.NewEncoder(body).Encode(evalRequest)

		request, _ := http.NewRequest(http.MethodGet, "/", body)
		response := httptest.NewRecorder()

		exprService := &StubExpressionService{
			isValid: wantIsValid,
		}
		exprHandler := handler.NewExpressionHandler(exprService)

		exprHandler.Validate(response, request)
		assert.Equal(t, response.Code, http.StatusOK)

		var gotResponse handler.ValidateResponse
		json.NewDecoder(response.Body).Decode(&gotResponse)

		assert.Equal(t, gotResponse, wantResponse)
	})

	t.Run("sets Valid field to false and Reason to error message on invalid expression", func(t *testing.T) {
		expression := "What is 5 plus 3?"
		errorMessage := "handler error"

		evalRequest := handler.ExpressionRequest{
			Expression: expression,
		}
		wantResponse := handler.ValidateResponse{
			Valid:  false,
			Reason: errorMessage,
		}

		body := bytes.NewBuffer([]byte{})
		json.NewEncoder(body).Encode(evalRequest)

		request, _ := http.NewRequest(http.MethodGet, "/", body)
		response := httptest.NewRecorder()

		exprService := &StubExpressionService{
			err: errors.New(errorMessage),
		}
		exprHandler := handler.NewExpressionHandler(exprService)

		exprHandler.Validate(response, request)
		assert.Equal(t, response.Code, http.StatusOK)

		var gotResponse handler.ValidateResponse
		json.NewDecoder(response.Body).Decode(&gotResponse)

		assert.Equal(t, gotResponse, wantResponse)
	})
}

func TestGetErrors(t *testing.T) {
	t.Run("returns expression errors from service", func(t *testing.T) {
		exprError := service.ExpressionError{
			Expression: "example expression",
			Method:     service.MethodValidate,
			Frequency:  3,
			Type:       service.ErrorTypeInvalidSyntax,
		}
		wantResponse := []handler.ExpressionErrorResponse{
			{
				Expression: exprError.Expression,
				Endpoint:   handler.ValidateEndpoint,
				Frequency:  exprError.Frequency,
				Type:       handler.InvalidSyntaxType,
			},
		}

		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		exprService := &StubExpressionService{
			exprErrors: []service.ExpressionError{exprError},
		}
		exprHandler := handler.NewExpressionHandler(exprService)

		exprHandler.GetExpressionErrors(response, request)
		assert.Equal(t, response.Code, http.StatusOK)

		var gotResponse []handler.ExpressionErrorResponse
		json.NewDecoder(response.Body).Decode(&gotResponse)

		assert.Equal(t, gotResponse, wantResponse)
	})

	t.Run("returns Internal Server Error on unknown method type from service", func(t *testing.T) {
		exprError := service.ExpressionError{
			Expression: "example expression",
			Method:     service.MethodType(-1),
			Frequency:  3,
			Type:       service.ErrorTypeInvalidSyntax,
		}
		wantResponse := handler.ErrorResponse{
			Error: handler.ErrUnknownMethod.Error(),
		}

		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		exprService := &StubExpressionService{
			exprErrors: []service.ExpressionError{exprError},
		}
		exprHandler := handler.NewExpressionHandler(exprService)

		exprHandler.GetExpressionErrors(response, request)
		assert.Equal(t, response.Code, http.StatusInternalServerError)

		var gotResponse handler.ErrorResponse
		json.NewDecoder(response.Body).Decode(&gotResponse)

		assert.Equal(t, gotResponse, wantResponse)
	})

	t.Run("returns Internal Server Error on unknown expression error type from service", func(t *testing.T) {
		exprError := service.ExpressionError{
			Expression: "example expression",
			Method:     service.MethodValidate,
			Frequency:  3,
			Type:       service.ErrorType(-1),
		}
		wantResponse := handler.ErrorResponse{
			Error: handler.ErrUnknownExpressionError.Error(),
		}

		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		exprService := &StubExpressionService{
			exprErrors: []service.ExpressionError{exprError},
		}
		exprHandler := handler.NewExpressionHandler(exprService)

		exprHandler.GetExpressionErrors(response, request)
		assert.Equal(t, response.Code, http.StatusInternalServerError)

		var gotResponse handler.ErrorResponse
		json.NewDecoder(response.Body).Decode(&gotResponse)

		assert.Equal(t, gotResponse, wantResponse)
	})
}
