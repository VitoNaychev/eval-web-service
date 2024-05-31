package interp_test

import (
	"testing"

	"github.com/VitoNaychev/eval-web-service/interp"
	"github.com/VitoNaychev/eval-web-service/testutil/assert"
)

func TestLexer(t *testing.T) {
	cases := []struct {
		Name           string
		Input          string
		ExpectedTokens []interp.Token
		ExpectedError  error
	}{
		{
			Name:  "just a question",
			Input: "What is",
			ExpectedTokens: []interp.Token{
				&interp.QuestionToken{"What is"},
			},
			ExpectedError: nil,
		},
		{
			Name:           "error on non-math question",
			Input:          "How is",
			ExpectedTokens: nil,
			ExpectedError:  interp.ErrNonMathQuestion,
		},
		{
			Name:           "error on non-math question (just a number)",
			Input:          "5",
			ExpectedTokens: nil,
			ExpectedError:  interp.ErrNonMathQuestion,
		},
		{
			Name:           "error on non-math question (just an operand)",
			Input:          "plus",
			ExpectedTokens: nil,
			ExpectedError:  interp.ErrNonMathQuestion,
		},
		{
			Name:           "error on non-math question (just punctuation)",
			Input:          "?",
			ExpectedTokens: nil,
			ExpectedError:  interp.ErrNonMathQuestion,
		},
		{
			Name:  "question with a number",
			Input: "What is 5",
			ExpectedTokens: []interp.Token{
				&interp.QuestionToken{"What is"},
				&interp.NumberToken{"5"},
			},
			ExpectedError: nil,
		},
		{
			Name:  "question with an operand",
			Input: "What is plus",
			ExpectedTokens: []interp.Token{
				&interp.QuestionToken{"What is"},
				&interp.OperandToken{"plus"},
			},
			ExpectedError: nil,
		},
		{
			Name:  "question with punctuation",
			Input: "What is ?",
			ExpectedTokens: []interp.Token{
				&interp.QuestionToken{"What is"},
				&interp.PunctuationToken{"?"},
			},
			ExpectedError: nil,
		},
		{
			Name:           "error on question with unknown operand",
			Input:          "What is cubed",
			ExpectedTokens: nil,
			ExpectedError:  interp.ErrUnsupportedOperation,
		},
		{
			Name:           "error on question with number and unknown operand",
			Input:          "What is 3 cubed",
			ExpectedTokens: nil,
			ExpectedError:  interp.ErrUnsupportedOperation,
		},
		{
			Name:           "error on question with known operand and unknown operand",
			Input:          "What is 3 cubed",
			ExpectedTokens: nil,
			ExpectedError:  interp.ErrUnsupportedOperation,
		},
		{
			Name:           "error on question with punctuation and unknown operand",
			Input:          "What is cubed?",
			ExpectedTokens: nil,
			ExpectedError:  interp.ErrUnsupportedOperation,
		},
		{
			Name:  "question with number and operand",
			Input: "What is 3 plus",
			ExpectedTokens: []interp.Token{
				&interp.QuestionToken{"What is"},
				&interp.NumberToken{"3"},
				&interp.OperandToken{"plus"},
			},
			ExpectedError: nil,
		},
		{
			Name:  "question with number and punctuation",
			Input: "What is 3?",
			ExpectedTokens: []interp.Token{
				&interp.QuestionToken{"What is"},
				&interp.NumberToken{"3"},
				&interp.PunctuationToken{"?"},
			},
			ExpectedError: nil,
		},
		{
			Name:  "question with operand and number",
			Input: "What is plus 3",
			ExpectedTokens: []interp.Token{
				&interp.QuestionToken{"What is"},
				&interp.OperandToken{"plus"},
				&interp.NumberToken{"3"},
			},
			ExpectedError: nil,
		},
		{
			Name:  "question with operand and punctuation",
			Input: "What is plus?",
			ExpectedTokens: []interp.Token{
				&interp.QuestionToken{"What is"},
				&interp.OperandToken{"plus"},
				&interp.PunctuationToken{"?"},
			},
			ExpectedError: nil,
		},
		{
			Name:  "question with punctuation and number",
			Input: "What is ? plus",
			ExpectedTokens: []interp.Token{
				&interp.QuestionToken{"What is"},
				&interp.PunctuationToken{"?"},
				&interp.OperandToken{"plus"},
			},
			ExpectedError: nil,
		},
		{
			Name:  "question with punctuation and operand",
			Input: "What is ? 3",
			ExpectedTokens: []interp.Token{
				&interp.QuestionToken{"What is"},
				&interp.PunctuationToken{"?"},
				&interp.NumberToken{"3"},
			},
			ExpectedError: nil,
		},
		{
			Name:  "question with number, operand, number",
			Input: "What is 3 plus 3",
			ExpectedTokens: []interp.Token{
				&interp.QuestionToken{"What is"},
				&interp.NumberToken{"3"},
				&interp.OperandToken{"plus"},
				&interp.NumberToken{"3"},
			},
			ExpectedError: nil,
		},
		{
			Name:  "question with number, operand, number, punctuation",
			Input: "What is 3 plus 3?",
			ExpectedTokens: []interp.Token{
				&interp.QuestionToken{"What is"},
				&interp.NumberToken{"3"},
				&interp.OperandToken{"plus"},
				&interp.NumberToken{"3"},
				&interp.PunctuationToken{"?"},
			},
			ExpectedError: nil,
		},
		{
			Name:  "question with number and operand, number, operand, punctuation",
			Input: "What is 3 plus 10 minus ?",
			ExpectedTokens: []interp.Token{
				&interp.QuestionToken{"What is"},
				&interp.NumberToken{"3"},
				&interp.OperandToken{"plus"},
				&interp.NumberToken{"10"},
				&interp.OperandToken{"minus"},
				&interp.PunctuationToken{"?"},
			},
			ExpectedError: nil,
		},
		{
			Name:  "question with number and operand, number, operand, number, punctuation",
			Input: "What is 3 plus 10 minus 5?",
			ExpectedTokens: []interp.Token{
				&interp.QuestionToken{"What is"},
				&interp.NumberToken{"3"},
				&interp.OperandToken{"plus"},
				&interp.NumberToken{"10"},
				&interp.OperandToken{"minus"},
				&interp.NumberToken{"5"},
				&interp.PunctuationToken{"?"},
			},
			ExpectedError: nil,
		},
		{
			Name:  "question with number and operand, number, operand, number, punctuation",
			Input: "What is 3plus10minus5?",
			ExpectedTokens: []interp.Token{
				&interp.QuestionToken{"What is"},
				&interp.NumberToken{"3"},
				&interp.OperandToken{"plus"},
				&interp.NumberToken{"10"},
				&interp.OperandToken{"minus"},
				&interp.NumberToken{"5"},
				&interp.PunctuationToken{"?"},
			},
			ExpectedError: nil,
		},
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			gotTokens, gotError := interp.Lex(test.Input)

			assert.Equal(t, gotTokens, test.ExpectedTokens)
			assert.Equal(t, gotError, test.ExpectedError)
		})
	}
}
