package handler

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
