package handler

type ErrorResponse struct {
	Error string `json:"message"`
}

type EvaluateResponse struct {
	Result int `json:"result"`
}

type EvaluateRequest struct {
	Expression string `json:"expression"`
}
