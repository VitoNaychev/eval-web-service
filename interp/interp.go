package interp

import "strconv"

func Interpret(tokens []Token) int {
	result := parseNumberToken(tokens[0].(*NumberToken))

	opPtr := 1
	numPtr := 2

	for numPtr < len(tokens) {
		numToken := tokens[numPtr].(*NumberToken)
		number := parseNumberToken(numToken)

		opToken := tokens[opPtr].(*OperandToken)
		result = executeOperandToken(opToken, result, number)

		opPtr += 2
		numPtr += 2
	}

	return result
}

func parseNumberToken(token *NumberToken) int {
	str := token.GetToken().(string)
	num, _ := strconv.Atoi(str)

	return num
}

func executeOperandToken(token *OperandToken, left, right int) int {
	var result int

	switch token.GetToken().(string) {
	case "plus":
		result = left + right
	case "minus":
		result = left - right
	case "multiplied by":
		result = left * right
	case "divided by":
		result = left / right
	}

	return result
}
