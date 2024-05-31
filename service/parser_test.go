package service_test

import (
	"testing"

	"github.com/VitoNaychev/eval-web-service/service"
	"github.com/VitoNaychev/eval-web-service/testutil/assert"
)

func TestParser(t *testing.T) {
	cases := []struct {
		Name           string
		Input          []service.Token
		ExpectedOutput []service.Token
		ExpectedError  error
	}{
		{
			Name: "error on missing question token",
			Input: []service.Token{
				&service.NumberToken{"10"},
			},
			ExpectedOutput: nil,
			ExpectedError:  service.ErrInvalidSyntax,
		},
		{
			Name: "error on missing punctuation token",
			Input: []service.Token{
				&service.QuestionToken{"What is"},
				&service.NumberToken{"10"},
			},
			ExpectedOutput: nil,
			ExpectedError:  service.ErrInvalidSyntax,
		},
		{
			Name: "removes nonsignificant tokens (one number expression)",
			Input: []service.Token{
				&service.QuestionToken{"What is"},
				&service.NumberToken{"10"},
				&service.PunctuationToken{"?"},
			},
			ExpectedOutput: []service.Token{
				&service.NumberToken{"10"},
			},
			ExpectedError: nil,
		},
		{
			Name: "error on statement ending on operand",
			Input: []service.Token{
				&service.QuestionToken{"What is"},
				&service.NumberToken{"10"},
				&service.OperandToken{"plus"},
				&service.PunctuationToken{"?"},
			},
			ExpectedOutput: nil,
			ExpectedError:  service.ErrInvalidSyntax,
		},
		{
			Name: "removes nonsignificant tokens (number, operand, number expression)",
			Input: []service.Token{
				&service.QuestionToken{"What is"},
				&service.NumberToken{"10"},
				&service.OperandToken{"plus"},
				&service.NumberToken{"20"},
				&service.PunctuationToken{"?"},
			},
			ExpectedOutput: []service.Token{
				&service.NumberToken{"10"},
				&service.OperandToken{"plus"},
				&service.NumberToken{"20"},
			},
			ExpectedError: nil,
		},
		{
			Name: "error on statement with two questions",
			Input: []service.Token{
				&service.QuestionToken{"What is"},
				&service.NumberToken{"10"},
				&service.OperandToken{"plus"},
				&service.NumberToken{"20"},
				&service.QuestionToken{"What is"},
				&service.PunctuationToken{"?"},
			},
			ExpectedOutput: nil,
			ExpectedError:  service.ErrInvalidSyntax,
		},
		{
			Name: "error on statement with two punctuation marks",
			Input: []service.Token{
				&service.QuestionToken{"What is"},
				&service.NumberToken{"10"},
				&service.PunctuationToken{"?"},
				&service.OperandToken{"plus"},
				&service.NumberToken{"20"},
				&service.PunctuationToken{"?"},
			},
			ExpectedOutput: nil,
			ExpectedError:  service.ErrInvalidSyntax,
		},
		{
			Name: "error on statement with operand after punctuation token",
			Input: []service.Token{
				&service.QuestionToken{"What is"},
				&service.NumberToken{"10"},
				&service.PunctuationToken{"?"},
				&service.OperandToken{"plus"},
				&service.NumberToken{"20"},
				&service.PunctuationToken{"?"},
				&service.OperandToken{"multiplied by"},
			},
			ExpectedOutput: nil,
			ExpectedError:  service.ErrInvalidSyntax,
		},
		{
			Name: "error on statement with number after punctuation token",
			Input: []service.Token{
				&service.QuestionToken{"What is"},
				&service.NumberToken{"10"},
				&service.PunctuationToken{"?"},
				&service.OperandToken{"plus"},
				&service.NumberToken{"20"},
				&service.PunctuationToken{"?"},
				&service.NumberToken{"42"},
			},
			ExpectedOutput: nil,
			ExpectedError:  service.ErrInvalidSyntax,
		},
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			gotTokens, gotError := service.Parse(test.Input)

			assert.Equal(t, gotTokens, test.ExpectedOutput)
			assert.Equal(t, gotError, test.ExpectedError)
		})
	}
}
