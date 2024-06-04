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

	response, _ := e.client.Post(e.url, "application/json", body)

	if response.StatusCode != 200 {
		var errorResponse ErrorResponse
		json.NewDecoder(response.Body).Decode(&errorResponse)

		if response.StatusCode == http.StatusBadRequest {
			return -1, errorMessageToClientError(errorResponse.Error)
		}
		return -1, NewClientError(errorResponse.Error)
	}

	var evaluateResponse EvaluateResponse
	json.NewDecoder(response.Body).Decode(&evaluateResponse)

	return evaluateResponse.Result, nil
}

func (e *ExpressionHTTPClient) Validate(expr string) (bool, error) {
	return false, nil
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
