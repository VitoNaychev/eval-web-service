package interp_test

import (
	"testing"

	"github.com/VitoNaychev/eval-web-service/interp"
	"github.com/VitoNaychev/eval-web-service/testutil/assert"
)

func TestParser(t *testing.T) {
	cases := []struct {
		Name           string
		Input          []interp.Token
		ExpectedOutput []interp.Token
		ExpectedError  error
	}{
		{
			Name: "error on missing question token",
			Input: []interp.Token{
				&interp.NumberToken{"10"},
			},
			ExpectedOutput: nil,
			ExpectedError:  interp.ErrInvalidSyntax,
		},
		{
			Name: "error on missing punctuation token",
			Input: []interp.Token{
				&interp.QuestionToken{"What is"},
				&interp.NumberToken{"10"},
			},
			ExpectedOutput: nil,
			ExpectedError:  interp.ErrInvalidSyntax,
		},
		{
			Name: "removes nonsignificant tokens (one number expression)",
			Input: []interp.Token{
				&interp.QuestionToken{"What is"},
				&interp.NumberToken{"10"},
				&interp.PunctuationToken{"?"},
			},
			ExpectedOutput: []interp.Token{
				&interp.NumberToken{"10"},
			},
			ExpectedError: nil,
		},
		{
			Name: "error on statement ending on operand",
			Input: []interp.Token{
				&interp.QuestionToken{"What is"},
				&interp.NumberToken{"10"},
				&interp.OperandToken{"plus"},
				&interp.PunctuationToken{"?"},
			},
			ExpectedOutput: nil,
			ExpectedError:  interp.ErrInvalidSyntax,
		},
		{
			Name: "removes nonsignificant tokens (number, operand, number expression)",
			Input: []interp.Token{
				&interp.QuestionToken{"What is"},
				&interp.NumberToken{"10"},
				&interp.OperandToken{"plus"},
				&interp.NumberToken{"20"},
				&interp.PunctuationToken{"?"},
			},
			ExpectedOutput: []interp.Token{
				&interp.NumberToken{"10"},
				&interp.OperandToken{"plus"},
				&interp.NumberToken{"20"},
			},
			ExpectedError: nil,
		},
		{
			Name: "error on statement with two questions",
			Input: []interp.Token{
				&interp.QuestionToken{"What is"},
				&interp.NumberToken{"10"},
				&interp.OperandToken{"plus"},
				&interp.NumberToken{"20"},
				&interp.QuestionToken{"What is"},
				&interp.PunctuationToken{"?"},
			},
			ExpectedOutput: nil,
			ExpectedError:  interp.ErrInvalidSyntax,
		},
		{
			Name: "error on statement with two punctuation marks",
			Input: []interp.Token{
				&interp.QuestionToken{"What is"},
				&interp.NumberToken{"10"},
				&interp.PunctuationToken{"?"},
				&interp.OperandToken{"plus"},
				&interp.NumberToken{"20"},
				&interp.PunctuationToken{"?"},
			},
			ExpectedOutput: nil,
			ExpectedError:  interp.ErrInvalidSyntax,
		},
		{
			Name: "error on statement with operand after punctuation token",
			Input: []interp.Token{
				&interp.QuestionToken{"What is"},
				&interp.NumberToken{"10"},
				&interp.PunctuationToken{"?"},
				&interp.OperandToken{"plus"},
				&interp.NumberToken{"20"},
				&interp.PunctuationToken{"?"},
				&interp.OperandToken{"multiplied by"},
			},
			ExpectedOutput: nil,
			ExpectedError:  interp.ErrInvalidSyntax,
		},
		{
			Name: "error on statement with number after punctuation token",
			Input: []interp.Token{
				&interp.QuestionToken{"What is"},
				&interp.NumberToken{"10"},
				&interp.PunctuationToken{"?"},
				&interp.OperandToken{"plus"},
				&interp.NumberToken{"20"},
				&interp.PunctuationToken{"?"},
				&interp.NumberToken{"42"},
			},
			ExpectedOutput: nil,
			ExpectedError:  interp.ErrInvalidSyntax,
		},
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			gotTokens, gotError := interp.Parse(test.Input)

			assert.Equal(t, gotTokens, test.ExpectedOutput)
			assert.Equal(t, gotError, test.ExpectedError)
		})
	}
}
