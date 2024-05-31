package service

type ErrorType int

const (
	ErrorTypeNonMathQuestion ErrorType = iota
	ErrorTypeUnsupportedOperand
	ErrorTypeInvalidSyntax
)

type MethodType int

const (
	MethodValidate = iota
	MethodExecute
)

type ExpressionError struct {
	Expression string
	Method     MethodType
	Frequency  int
	Type       ErrorType
}
