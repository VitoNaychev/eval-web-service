package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VitoNaychev/eval-web-service/handler"
	"github.com/VitoNaychev/eval-web-service/testutil/assert"
)

type StubExpressionHandler struct {
	spyEvaluate            bool
	spyValidate            bool
	spyGetExpressionErrors bool
}

func (s *StubExpressionHandler) Evaluate(w http.ResponseWriter, r *http.Request) {
	s.spyEvaluate = true
}

func (s *StubExpressionHandler) Validate(w http.ResponseWriter, r *http.Request) {
	s.spyValidate = true
}

func (s *StubExpressionHandler) GetExpressionErrors(w http.ResponseWriter, r *http.Request) {
	s.spyGetExpressionErrors = true
}

func TestRouting(t *testing.T) {
	t.Run("routes evaluation requests for Evaluate handler", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, handler.EvaluateEndpoint, nil)
		response := httptest.NewRecorder()

		txHandler := &StubExpressionHandler{
			spyEvaluate:            false,
			spyValidate:            false,
			spyGetExpressionErrors: false,
		}
		router := handler.NewRouter(txHandler)

		router.Handler.ServeHTTP(response, request)

		assert.Equal(t, txHandler.spyEvaluate, true)
		assert.Equal(t, txHandler.spyValidate, false)
		assert.Equal(t, txHandler.spyGetExpressionErrors, false)
	})

	t.Run("routes validation requests for Validate handler", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, handler.ValidateEndpoint, nil)
		response := httptest.NewRecorder()

		txHandler := &StubExpressionHandler{
			spyEvaluate:            false,
			spyValidate:            false,
			spyGetExpressionErrors: false,
		}
		router := handler.NewRouter(txHandler)

		router.Handler.ServeHTTP(response, request)

		assert.Equal(t, txHandler.spyEvaluate, false)
		assert.Equal(t, txHandler.spyValidate, true)
		assert.Equal(t, txHandler.spyGetExpressionErrors, false)
	})

	t.Run("routes get expression errors requests for GetExpressionErrors handler", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, handler.GetExpressionErrorsEndpoint, nil)
		response := httptest.NewRecorder()

		txHandler := &StubExpressionHandler{
			spyEvaluate:            false,
			spyValidate:            false,
			spyGetExpressionErrors: false,
		}
		router := handler.NewRouter(txHandler)

		router.Handler.ServeHTTP(response, request)

		assert.Equal(t, txHandler.spyEvaluate, false)
		assert.Equal(t, txHandler.spyValidate, false)
		assert.Equal(t, txHandler.spyGetExpressionErrors, true)
	})
}
