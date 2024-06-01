package handler

import "net/http"

const (
	EvaluateEndpoint            = "/evaluate"
	ValidateEndpoint            = "/validate"
	GetExpressionErrorsEndpoint = "/errors"
)

type expressionHandler interface {
	Evaluate(w http.ResponseWriter, r *http.Request)
	Validate(w http.ResponseWriter, r *http.Request)
	GetExpressionErrors(w http.ResponseWriter, r *http.Request)
}

type Router struct {
	Handler http.Handler
}

func NewRouter(handler expressionHandler) *Router {
	mux := http.NewServeMux()
	mux.HandleFunc(EvaluateEndpoint, handler.Evaluate)
	mux.HandleFunc(ValidateEndpoint, handler.Validate)
	mux.HandleFunc(GetExpressionErrorsEndpoint, handler.GetExpressionErrors)

	return &Router{
		Handler: mux,
	}
}
