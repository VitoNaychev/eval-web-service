package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type Client interface {
	Post(string, string, io.Reader) (*http.Response, error)
	Get(string) (*http.Response, error)
}

type ExpressionHTTPClient struct {
	client Client
	url    string
}

func NewExpressionHTTPClient(client Client, url string) *ExpressionHTTPClient {
	return &ExpressionHTTPClient{
		client: client,
		url:    url,
	}
}

func (e *ExpressionHTTPClient) Evaluate(expr string) (int, error) {
	expressionRequest := ExpressionRequest{
		Expression: expr,
	}

	body := bytes.NewBuffer([]byte{})
	json.NewEncoder(body).Encode(expressionRequest)

	response, _ := e.client.Post(e.url+EvaluateURL, "application/json", body)

	if response.StatusCode != 200 {
		return -1, handleServerError(response)
	}

	var evaluateResponse EvaluateResponse
	json.NewDecoder(response.Body).Decode(&evaluateResponse)

	return evaluateResponse.Result, nil
}

func (e *ExpressionHTTPClient) Validate(expr string) (bool, error) {
	expressionRequest := ExpressionRequest{
		Expression: expr,
	}

	body := bytes.NewBuffer([]byte{})
	json.NewEncoder(body).Encode(expressionRequest)

	response, _ := e.client.Post(e.url+ValidateURL, "application/json", body)

	if response.StatusCode != 200 {
		return false, handleServerError(response)
	}

	var validateResponse ValidateResponse
	json.NewDecoder(response.Body).Decode(&validateResponse)

	return validateResponse.Valid, nil
}

func (e *ExpressionHTTPClient) GetExpressionErrors() ([]ExpressionError, error) {
	response, _ := e.client.Post(e.url+ExpressionErrorsURL, "application/json", nil)

	if response.StatusCode != 200 {
		var errorResponse ErrorResponse
		json.NewDecoder(response.Body).Decode(&errorResponse)

		return nil, NewClientError(errorResponse.Error)
	}

	var expressionErrorsResponse []ExpressionErrorResponse
	json.NewDecoder(response.Body).Decode(&expressionErrorsResponse)

	var expressionErrors []ExpressionError
	for _, expressionErrorResponse := range expressionErrorsResponse {
		expressionErrors = append(expressionErrors,
			expressionErrorResponseToExpressionError(expressionErrorResponse))
	}

	return expressionErrors, nil
}

func expressionErrorResponseToExpressionError(r ExpressionErrorResponse) ExpressionError {
	return ExpressionError{
		Expression: r.Expression,
		Method:     r.Endpoint,
		Frequency:  r.Frequency,
		Type:       r.Type,
	}
}

func handleServerError(response *http.Response) error {
	var errorResponse ErrorResponse
	json.NewDecoder(response.Body).Decode(&errorResponse)

	if response.StatusCode == http.StatusBadRequest {
		return errorMessageToClientError(errorResponse.Error)
	}
	return NewClientError(errorResponse.Error)
}

func errorMessageToClientError(msg string) error {
	switch msg {
	case NonMathQuestionMessage:
		return ErrNonMathQuestion
	case UnsupportedOperationMessage:
		return ErrUnsupportedOperation
	case InvalidSyntaxMessasge:
		return ErrInvalidSyntax
	default:
		return NewClientError("unknown error response")
	}
}
