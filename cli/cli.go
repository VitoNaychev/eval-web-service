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
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		c.in.Scan()
		cmd := c.in.Text()
		if cmd == `\e` {
			return
		}

		output, err := c.executeCommand(cmd)
		if err != nil {
			c.out.WriteString("error: " + err.Error())
			c.out.Flush()
			continue
		}

		c.out.WriteString(output)
		c.out.Flush()
	}
}

func (c *CLI) executeCommand(cmd string) (string, error) {
	var output string

	if len(cmd) == 0 {
		return "", nil
	}

	switch {
	case strings.HasPrefix(cmd, ".. "):
		expr := strings.TrimPrefix(cmd, ".. ")
		result, _ := c.client.Evaluate(expr)

		output = fmt.Sprintln(result)
	case strings.HasPrefix(cmd, "?? "):
		expr := strings.TrimPrefix(cmd, "?? ")
		isValid, _ := c.client.Validate(expr)

		output = fmt.Sprintln(isValid)
	case strings.HasPrefix(cmd, "!!"):
		exprErrors, _ := c.client.GetExpressionErrors()

		for _, exprError := range exprErrors {
			output += formatExpressionError(exprError)
		}
	default:
		output = "unknown command"
	}

	return output, nil
}

func formatExpressionError(e client.ExpressionError) string {
	return fmt.Sprintf(`\t"%s"; on %s; %d times; %s\n`,
		e.Expression, e.Method, e.Frequency, e.Type)
}
