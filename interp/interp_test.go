package interp_test

import (
	"testing"

	"github.com/VitoNaychev/eval-web-service/interp"
	"github.com/VitoNaychev/eval-web-service/testutil/assert"
)

func TestInerpreter(t *testing.T) {
	cases := []struct {
		Name           string
		Input          []interp.Token
		ExpectedResult int
	}{
		{
			Name: "just a number",
			Input: []interp.Token{
				&interp.NumberToken{"8"},
			},
			ExpectedResult: 8,
		},
		{
			Name: "addition",
			Input: []interp.Token{
				&interp.NumberToken{"8"},
				&interp.OperandToken{"plus"},
				&interp.NumberToken{"3"},
			},
			ExpectedResult: 11,
		},
		{
			Name: "subtraction",
			Input: []interp.Token{
				&interp.NumberToken{"8"},
				&interp.OperandToken{"minus"},
				&interp.NumberToken{"3"},
			},
			ExpectedResult: 5,
		},
		{
			Name: "multiplciation",
			Input: []interp.Token{
				&interp.NumberToken{"5"},
				&interp.OperandToken{"multiplied by"},
				&interp.NumberToken{"7"},
			},
			ExpectedResult: 35,
		},
		{
			Name: "division",
			Input: []interp.Token{
				&interp.NumberToken{"42"},
				&interp.OperandToken{"divided by"},
				&interp.NumberToken{"6"},
			},
			ExpectedResult: 7,
		},
		{
			Name: "sequence of operations",
			Input: []interp.Token{
				&interp.NumberToken{"42"},
				&interp.OperandToken{"divided by"},
				&interp.NumberToken{"6"},
				&interp.OperandToken{"plus"},
				&interp.NumberToken{"3"},
				&interp.OperandToken{"multiplied by"},
				&interp.NumberToken{"8"},
			},
			ExpectedResult: 80,
		},
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			gotResult := interp.Interpret(test.Input)

			assert.Equal(t, gotResult, test.ExpectedResult)
		})
	}
}
