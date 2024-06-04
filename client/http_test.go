package client_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/VitoNaychev/eval-web-service/client"
	"github.com/VitoNaychev/eval-web-service/testutil/assert"
)

type StubHttpClient struct {
	err error

	spyURL         string
	spyContentType string
	spyData        io.Reader

	code     int
	response interface{}
}

func (s *StubHttpClient) Post(url string, contentType string, data io.Reader) (*http.Response, error) {
	s.spyURL = url
	s.spyContentType = contentType
	s.spyData = data

	if s.err != nil {
		return nil, s.err
	}

	response := &http.Response{}

	body := bytes.NewBuffer([]byte{})
	json.NewEncoder(body).Encode(s.response)

	response.StatusCode = s.code
	response.Body = io.NopCloser(body)

	return response, nil

}

func (s *StubHttpClient) Get(url string) (*http.Response, error) {
	s.spyURL = url

	if s.err != nil {
		return nil, s.err
	}

	response := &http.Response{}

	body := bytes.NewBuffer([]byte{})
	json.NewEncoder(body).Encode(s.response)

	response.StatusCode = s.code
	response.Body = io.NopCloser(body)

	return response, nil
}

func TestEvaluate(t *testing.T) {
	t.Run("evaluates expression", func(t *testing.T) {
		url := "example-url.com"

		expression := "What is 5 plus 10?"
		wantExpressionRequest := client.ExpressionRequest{
			Expression: expression,
		}

		wantResult := 15
		expressionResponse := client.EvaluateResponse{
			Result: wantResult,
		}

		httpClient := &StubHttpClient{
			code:     http.StatusOK,
			response: expressionResponse,
		}
		exprClient := client.NewExpressionHTTPClient(httpClient, url)

		gotResult, _ := exprClient.Evaluate(expression)

		assert.Equal(t, gotResult, wantResult)

		assert.Equal(t, httpClient.spyURL, url+client.EvaluateURL)

		var gotExpressionRequest client.ExpressionRequest
		json.NewDecoder(httpClient.spyData).Decode(&gotExpressionRequest)

		assert.Equal(t, gotExpressionRequest, wantExpressionRequest)
	})

	t.Run("parses and returns error on Status Bad Request", func(t *testing.T) {
		expression := "What is 5 plus 10?"

		wantError := client.ErrNonMathQuestion
		errorResponse := client.ErrorResponse{
			Error: wantError.Error(),
		}

		httpClient := &StubHttpClient{
			code:     http.StatusBadRequest,
			response: errorResponse,
		}
		exprClient := client.NewExpressionHTTPClient(httpClient, "/evaluate")

		_, gotError := exprClient.Evaluate(expression)

		assert.ErrorType[*client.ClientError](t, gotError)
		assert.Equal(t, gotError, wantError)
	})

	t.Run("wraps error message in ClientError on Internal Server Error", func(t *testing.T) {
		expression := "What is 5 plus 10?"

		wantError := errors.New("test error")
		errorResponse := client.ErrorResponse{
			Error: wantError.Error(),
		}

		httpClient := &StubHttpClient{
			code:     http.StatusInternalServerError,
			response: errorResponse,
		}
		exprClient := client.NewExpressionHTTPClient(httpClient, "/evaluate")

		_, gotError := exprClient.Evaluate(expression)

		assert.ErrorType[*client.ClientError](t, gotError)
		assert.Equal(t, gotError.Error(), wantError.Error())
	})
}

func TestValidate(t *testing.T) {
	t.Run("validates expression", func(t *testing.T) {
		url := "example-url.com"

		expression := "What is 5 plus 10?"
		wantExpressionRequest := client.ExpressionRequest{
			Expression: expression,
		}

		wantIsValid := true
		expressionResponse := client.ValidateResponse{
			Valid: wantIsValid,
		}

		httpClient := &StubHttpClient{
			code:     http.StatusOK,
			response: expressionResponse,
		}
		exprClient := client.NewExpressionHTTPClient(httpClient, url)

		gotIsValid, _ := exprClient.Validate(expression)

		assert.Equal(t, gotIsValid, wantIsValid)

		assert.Equal(t, httpClient.spyURL, url+client.ValidateURL)

		var gotExpressionRequest client.ExpressionRequest
		json.NewDecoder(httpClient.spyData).Decode(&gotExpressionRequest)

		assert.Equal(t, gotExpressionRequest, wantExpressionRequest)
	})

	t.Run("parses and returns error on Status Bad Request", func(t *testing.T) {
		expression := "What is 5 plus 10?"

		wantError := client.ErrNonMathQuestion
		errorResponse := client.ErrorResponse{
			Error: wantError.Error(),
		}

		httpClient := &StubHttpClient{
			code:     http.StatusBadRequest,
			response: errorResponse,
		}
		exprClient := client.NewExpressionHTTPClient(httpClient, "/evaluate")

		_, gotError := exprClient.Validate(expression)

		assert.ErrorType[*client.ClientError](t, gotError)
		assert.Equal(t, gotError, wantError)
	})

	t.Run("wraps error message in ClientError on Internal Server Error", func(t *testing.T) {
		expression := "What is 5 plus 10?"

		wantError := errors.New("test error")
		errorResponse := client.ErrorResponse{
			Error: wantError.Error(),
		}

		httpClient := &StubHttpClient{
			code:     http.StatusInternalServerError,
			response: errorResponse,
		}
		exprClient := client.NewExpressionHTTPClient(httpClient, "/evaluate")

		_, gotError := exprClient.Validate(expression)

		assert.ErrorType[*client.ClientError](t, gotError)
		assert.Equal(t, gotError.Error(), wantError.Error())
	})
}

func TestGetExpressionErrors(t *testing.T) {

}
