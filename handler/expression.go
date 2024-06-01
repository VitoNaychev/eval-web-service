package handler

import (
	"encoding/json"
	"net/http"

	"github.com/VitoNaychev/eval-web-service/service"
)

type ExpressionService interface {
	Evaluate(string) (int, error)
	Validate(string) (bool, error)
	GetExpressionErrors() ([]service.ExpressionError, error)
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

func (e *ExpressionHandler) GetExpressionErrors(w http.ResponseWriter, r *http.Request) {
	exprErrors, _ := e.service.GetExpressionErrors()

	var exprErrorsResponse []ExpressionErrorResponse
	for _, exprError := range exprErrors {
		exprErrorResponse, err := exprErrorToExprErrorResponse(exprError)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err)
		}
		exprErrorsResponse = append(exprErrorsResponse, exprErrorResponse)
	}

	json.NewEncoder(w).Encode(exprErrorsResponse)
}

func exprErrorToExprErrorResponse(e service.ExpressionError) (ExpressionErrorResponse, error) {
	endpoint, err := serviceMethodToEndpoint(e.Method)
	if err != nil {
		return ExpressionErrorResponse{}, err
	}

	errType, err := serviceTypeToHandlerType(e.Type)
	if err != nil {
		return ExpressionErrorResponse{}, err
	}

	return ExpressionErrorResponse{
		Expression: e.Expression,
		Endpoint:   endpoint,
		Frequency:  e.Frequency,
		Type:       errType,
	}, nil
}

func serviceMethodToEndpoint(m service.MethodType) (string, error) {
	switch m {
	case service.MethodExecute:
		return EvaluateEndpoint, nil
	case service.MethodValidate:
		return ValidateEndpoint, nil
	default:
		return "", ErrUnknownMethod
	}
}

func serviceTypeToHandlerType(t service.ErrorType) (string, error) {
	switch t {
	case service.ErrorTypeNonMathQuestion:
		return NonMathQuesionType, nil
	case service.ErrorTypeUnsupportedOperand:
		return UnsupportedOperandType, nil
	case service.ErrorTypeInvalidSyntax:
		return InvalidSyntaxType, nil
	default:
		return "", ErrUnknownExpressionError
	}
}

func writeJSONError(w http.ResponseWriter, statusCode int, err error) {
	errorResponse := ErrorResponse{
		Error: err.Error(),
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorResponse)
}
