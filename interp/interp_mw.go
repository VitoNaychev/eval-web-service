package interp

import (
	"errors"

	"github.com/VitoNaychev/eval-web-service/service"
)

type LexFunc func(string) ([]Token, error)
type ParseFunc func([]Token) ([]Token, error)
type InterpFunc func([]Token) int

type InterpMW struct {
	lex    LexFunc
	parse  ParseFunc
	interp InterpFunc
}

func NewInterpMW(lex LexFunc, parse ParseFunc, interp InterpFunc) *InterpMW {
	return &InterpMW{
		lex:    lex,
		parse:  parse,
		interp: interp,
	}
}

func (i *InterpMW) Validate(input string) (bool, error) {
	tokens, err := i.lex(input)
	if err != nil {
		return false, interpErrorToServiceError(err)
	}

	_, err = i.parse(tokens)
	if err != nil {
		return false, interpErrorToServiceError(err)
	}

	return true, nil
}

func (i *InterpMW) Evaluate(input string) (int, error) {
	tokens, err := i.lex(input)
	if err != nil {
		return -1, interpErrorToServiceError(err)
	}

	significantTokens, err := i.parse(tokens)
	if err != nil {
		return -1, interpErrorToServiceError(err)
	}

	return i.interp(significantTokens), nil
}

func interpErrorToServiceError(err error) error {
	switch {
	case errors.Is(err, ErrNonMathQuestion):
		return service.ErrNonMathQuestion
	case errors.Is(err, ErrUnsupportedOperation):
		return service.ErrUnsupportedOperation
	case errors.Is(err, ErrInvalidSyntax):
		return service.ErrInvalidSyntax
	default:
		return err
	}
}
