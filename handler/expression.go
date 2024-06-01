package handler

import (
	"encoding/json"
	"net/http"
)

type ExpressionService interface {
	Evaluate(string) (int, error)
	Validate(string) (bool, error)
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
	var exprRequest ExpressionRequest
	json.NewDecoder(r.Body).Decode(&exprRequest)

	result, err := e.service.Evaluate(exprRequest.Expression)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err)
		return
	}

	evalResponse := EvaluateResponse{
		Result: result,
	}
	json.NewEncoder(w).Encode(evalResponse)
}

func (e *ExpressionHandler) Validate(w http.ResponseWriter, r *http.Request) {
	var exprRequest ExpressionRequest
	json.NewDecoder(r.Body).Decode(&exprRequest)

	isValid, err := e.service.Validate(exprRequest.Expression)

	var reason string
	if err != nil {
		reason = err.Error()
	}

	validateResponse := ValidateResponse{
		Valid:  isValid,
		Reason: reason,
	}
	json.NewEncoder(w).Encode(validateResponse)
}

func writeJSONError(w http.ResponseWriter, statusCode int, err error) {
	errorResponse := ErrorResponse{
		Error: err.Error(),
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorResponse)
}
