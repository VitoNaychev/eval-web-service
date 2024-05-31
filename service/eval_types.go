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
	MethodExec
)

type ExpressionError struct {
	Expression string
	Method     MethodType
	Frequency  int
	Type       ErrorType
}
