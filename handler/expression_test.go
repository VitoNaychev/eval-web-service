package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VitoNaychev/eval-web-service/handler"
	"github.com/VitoNaychev/eval-web-service/testutil/assert"
)

type StubExpressionService struct {
	result int
	err    error
}

func (s *StubExpressionService) Evaluate(expression string) (int, error) {
	return s.result, s.err
}

func TestEvaluate(t *testing.T) {
	t.Run("evaluates expression and returns EvaluateResponse", func(t *testing.T) {
		expression := "What is 5 plus 3?"
		result := 8

		evalRequest := handler.EvaluateRequest{
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

		var gotResponse handler.EvaluateResponse
		json.NewDecoder(response.Body).Decode(&gotResponse)

		assert.Equal(t, gotResponse, wantResponse)
	})

	t.Run("reutrns Status Bad Request and ErrorResponse on invalid expression", func(t *testing.T) {
		expression := "What is 5 plus 3?"
		errorMessage := "handler error"

		evalRequest := handler.EvaluateRequest{
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
