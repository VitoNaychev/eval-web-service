package handler

import "errors"

var (
	ErrUnknownMethod          = errors.New("unknwon method type")
	ErrUnknownExpressionError = errors.New("unknwon expression error type")
)

const (
	NonMathQuesionType     = "non-math question"
	UnsupportedOperandType = "unknown operand"
	InvalidSyntaxType      = "invalid syntax"
)

type ErrorResponse struct {
	Error string `json:"message"`
}

type ExpressionErrorResponse struct {
	Expression string `json:"expression"`
	Endpoint   string `json:"endpoint"`
	Frequency  int    `json:"frequency"`
	Type       string `json:"type"`
}

type ValidateResponse struct {
	Valid  bool   `json:"valid"`
	Reason string `json:"reason,omitempty"`
}

type EvaluateResponse struct {
	Result int `json:"result"`
}

type ExpressionRequest struct {
	Expression string `json:"expression"`
}
