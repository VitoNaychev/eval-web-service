package handler

import (
	"encoding/json"
	"net/http"
)

type ExpressionService interface {
	Evaluate(string) (int, error)
}

type ExpressionHandler struct {
	service ExpressionService
}

func NewExpressionHandler(service ExpressionService) *ExpressionHandler {
	return &ExpressionHandler{
		service: service,
	}
}

func (e *ExpressionHandler) Evaluate(w http.ResponseWriter, r *http.Request) {
	var evalRequest EvaluateRequest
	json.NewDecoder(r.Body).Decode(&evalRequest)

	result, err := e.service.Evaluate(evalRequest.Expression)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err)
		return
	}

	evalResponse := EvaluateResponse{
		Result: result,
	}
	json.NewEncoder(w).Encode(evalResponse)
}

func writeJSONError(w http.ResponseWriter, statusCode int, err error) {
	errorResponse := ErrorResponse{
		Error: err.Error(),
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorResponse)
}
