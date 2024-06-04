package cli

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/VitoNaychev/eval-web-service/client"
)

type ExpressionClient interface {
	Evaluate(string) (int, error)
	Validate(string) (bool, error)
	GetExpressionErrors() ([]client.ExpressionError, error)
}

type CLI struct {
	client ExpressionClient
	in     *bufio.Scanner
	out    *bufio.Writer
}

func NewCLI(client ExpressionClient, in io.Reader, out io.Writer) *CLI {
	return &CLI{
		client: client,
		in:     bufio.NewScanner(in),
		out:    bufio.NewWriter(out),
	}
}

func (c *CLI) Run(ctx context.Context) {
	inputCh := make(chan string)
	go readInput(c.in, inputCh)

	for {
		writePromptSymbol(c.out, PromptSymbol)

		select {
		case <-ctx.Done():
			return
		case cmd := <-inputCh:
			if cmd == `\e` {
				return
			}

			output, err := c.executeCommand(cmd)
			if err != nil {
				c.out.WriteString("error: " + err.Error())
			} else {
				c.out.WriteString(output)
			}

			c.out.Flush()
		}
	}
}

func writePromptSymbol(out *bufio.Writer, prompt string) {
	out.WriteString(prompt + " ")
	out.Flush()
}

func readInput(in *bufio.Scanner, inputCh chan string) {
	for in.Scan() {
		input := in.Text()
		inputCh <- input
	}
}

func (c *CLI) executeCommand(cmd string) (string, error) {
	var output string

	if len(cmd) == 0 {
		return "", nil
	}

	switch {
	case strings.HasPrefix(cmd, EvaluatePrefix):
		expr := strings.TrimPrefix(cmd, EvaluatePrefix)
		result, _ := c.client.Evaluate(expr)

		output = fmt.Sprintln(result)
	case strings.HasPrefix(cmd, ValidatePrefix):
		expr := strings.TrimPrefix(cmd, ValidatePrefix)
		isValid, _ := c.client.Validate(expr)

		output = fmt.Sprintln(isValid)
	case strings.HasPrefix(cmd, ExpressionErrorsPrefix):
		exprErrors, _ := c.client.GetExpressionErrors()
		for _, exprError := range exprErrors {
			output += formatExpressionError(exprError)
		}
	default:
		output = "unknown command\n"
	}

	return output, nil
}

func formatExpressionError(e client.ExpressionError) string {
	return fmt.Sprintf("\t\"%s\"; on %s; %d times; %s\n",
		e.Expression, e.Method, e.Frequency, e.Type)
}
