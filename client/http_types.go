package client

const (
	EvaluateURL         = "/evaluate"
	ValidateURL         = "/validate"
	ExpressionErrorsURL = "/errors"
)

const (
	NonMathQuestionMessage      = "non-math question"
	UnsupportedOperationMessage = "unsupported operation"
	InvalidSyntaxMessasge       = "invalid syntax"
)

type ErrorResponse struct {
	Error string `json:"message"`
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

type ExpressionErrorResponse struct {
	Expression string `json:"expression"`
	Method     string `json:"method"`
	Frequency  int    `json:"frequency"`
	Type       string `json:"type"`
}
