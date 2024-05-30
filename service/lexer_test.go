package service_test

import (
	"testing"

	"github.com/VitoNaychev/eval-web-service/service"
	"github.com/VitoNaychev/eval-web-service/testutil/assert"
)

func TestLexer(t *testing.T) {
	cases := []struct {
		Name           string
		Input          string
		ExpectedTokens []service.Token
		ExpectedError  error
	}{
		{
			Name:  "just a question",
			Input: "What is",
			ExpectedTokens: []service.Token{
				&service.QuestionToken{"What is"},
			},
			ExpectedError: nil,
		},
		{
			Name:           "error on non-math question",
			Input:          "How is",
			ExpectedTokens: nil,
			ExpectedError:  service.ErrNonMathQuestion,
		},
		{
			Name:           "error on non-math question (just a number)",
			Input:          "5",
			ExpectedTokens: nil,
			ExpectedError:  service.ErrNonMathQuestion,
		},
		{
			Name:           "error on non-math question (just an operand)",
			Input:          "plus",
			ExpectedTokens: nil,
			ExpectedError:  service.ErrNonMathQuestion,
		},
		{
			Name:           "error on non-math question (just punctuation)",
			Input:          "?",
			ExpectedTokens: nil,
			ExpectedError:  service.ErrNonMathQuestion,
		},
		{
			Name:  "question with a number",
			Input: "What is 5",
			ExpectedTokens: []service.Token{
				&service.QuestionToken{"What is"},
				&service.NumberToken{"5"},
			},
			ExpectedError: nil,
		},
		{
			Name:  "question with an operand",
			Input: "What is plus",
			ExpectedTokens: []service.Token{
				&service.QuestionToken{"What is"},
				&service.OperandToken{"plus"},
			},
			ExpectedError: nil,
		},
		{
			Name:  "question with punctuation",
			Input: "What is ?",
			ExpectedTokens: []service.Token{
				&service.QuestionToken{"What is"},
				&service.PunctuationToken{"?"},
			},
			ExpectedError: nil,
		},
		{
			Name:           "error on question with unknown operand",
			Input:          "What is cubed",
			ExpectedTokens: nil,
			ExpectedError:  service.ErrUnsupportedOperation,
		},
		{
			Name:           "error on question with number and unknown operand",
			Input:          "What is 3 cubed",
			ExpectedTokens: nil,
			ExpectedError:  service.ErrUnsupportedOperation,
		},
		{
			Name:           "error on question with known operand and unknown operand",
			Input:          "What is 3 cubed",
			ExpectedTokens: nil,
			ExpectedError:  service.ErrUnsupportedOperation,
		},
		{
			Name:           "error on question with punctuation and unknown operand",
			Input:          "What is cubed?",
			ExpectedTokens: nil,
			ExpectedError:  service.ErrUnsupportedOperation,
		},
		{
			Name:  "question with number and operand",
			Input: "What is 3 plus",
			ExpectedTokens: []service.Token{
				&service.QuestionToken{"What is"},
				&service.NumberToken{"3"},
				&service.OperandToken{"plus"},
			},
			ExpectedError: nil,
		},
		{
			Name:  "question with number and punctuation",
			Input: "What is 3?",
			ExpectedTokens: []service.Token{
				&service.QuestionToken{"What is"},
				&service.NumberToken{"3"},
				&service.PunctuationToken{"?"},
			},
			ExpectedError: nil,
		},
		{
			Name:  "question with operand and number",
			Input: "What is plus 3",
			ExpectedTokens: []service.Token{
				&service.QuestionToken{"What is"},
				&service.OperandToken{"plus"},
				&service.NumberToken{"3"},
			},
			ExpectedError: nil,
		},
		{
			Name:  "question with operand and punctuation",
			Input: "What is plus?",
			ExpectedTokens: []service.Token{
				&service.QuestionToken{"What is"},
				&service.OperandToken{"plus"},
				&service.PunctuationToken{"?"},
			},
			ExpectedError: nil,
		},
		{
			Name:  "question with punctuation and number",
			Input: "What is ? plus",
			ExpectedTokens: []service.Token{
				&service.QuestionToken{"What is"},
				&service.PunctuationToken{"?"},
				&service.OperandToken{"plus"},
			},
			ExpectedError: nil,
		},
		{
			Name:  "question with punctuation and operand",
			Input: "What is ? 3",
			ExpectedTokens: []service.Token{
				&service.QuestionToken{"What is"},
				&service.PunctuationToken{"?"},
				&service.NumberToken{"3"},
			},
			ExpectedError: nil,
		},
		{
			Name:  "question with number, operand, number",
			Input: "What is 3 plus 3",
			ExpectedTokens: []service.Token{
				&service.QuestionToken{"What is"},
				&service.NumberToken{"3"},
				&service.OperandToken{"plus"},
				&service.NumberToken{"3"},
			},
			ExpectedError: nil,
		},
		{
			Name:  "question with number, operand, number, punctuation",
			Input: "What is 3 plus 3?",
			ExpectedTokens: []service.Token{
				&service.QuestionToken{"What is"},
				&service.NumberToken{"3"},
				&service.OperandToken{"plus"},
				&service.NumberToken{"3"},
				&service.PunctuationToken{"?"},
			},
			ExpectedError: nil,
		},
		{
			Name:  "question with number and operand, number, operand, punctuation",
			Input: "What is 3 plus 10 minus ?",
			ExpectedTokens: []service.Token{
				&service.QuestionToken{"What is"},
				&service.NumberToken{"3"},
				&service.OperandToken{"plus"},
				&service.NumberToken{"10"},
				&service.OperandToken{"minus"},
				&service.PunctuationToken{"?"},
			},
			ExpectedError: nil,
		},
		{
			Name:  "question with number and operand, number, operand, number, punctuation",
			Input: "What is 3 plus 10 minus 5?",
			ExpectedTokens: []service.Token{
				&service.QuestionToken{"What is"},
				&service.NumberToken{"3"},
				&service.OperandToken{"plus"},
				&service.NumberToken{"10"},
				&service.OperandToken{"minus"},
				&service.NumberToken{"5"},
				&service.PunctuationToken{"?"},
			},
			ExpectedError: nil,
		},
		{
			Name:  "question with number and operand, number, operand, number, punctuation",
			Input: "What is 3plus10minus5?",
			ExpectedTokens: []service.Token{
				&service.QuestionToken{"What is"},
				&service.NumberToken{"3"},
				&service.OperandToken{"plus"},
				&service.NumberToken{"10"},
				&service.OperandToken{"minus"},
				&service.NumberToken{"5"},
				&service.PunctuationToken{"?"},
			},
			ExpectedError: nil,
		},
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			gotTokens, gotError := service.Lex(test.Input)

			assert.Equal(t, gotTokens, test.ExpectedTokens)
			assert.Equal(t, gotError, test.ExpectedError)
		})
	}
}
