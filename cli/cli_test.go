package cli_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/VitoNaychev/eval-web-service/cli"
	"github.com/VitoNaychev/eval-web-service/client"
	"github.com/VitoNaychev/eval-web-service/testutil/assert"
)

type StubExpressionClient struct {
	result     int
	isValid    bool
	exprErrors []client.ExpressionError

	spyEvaluateExpr string
	spyValidateExpr string
}

func (s *StubExpressionClient) Evaluate(expr string) (int, error) {
	s.spyEvaluateExpr = expr
	return s.result, nil
}

func (s *StubExpressionClient) Validate(expr string) (bool, error) {
	s.spyValidateExpr = expr
	return s.isValid, nil
}

func (s *StubExpressionClient) GetExpressionErrors() ([]client.ExpressionError, error) {
	return s.exprErrors, nil
}

func TestCLI(t *testing.T) {
	exitCmd := "\\e\n"

	t.Run("returns on signal from ctx.Done()", func(t *testing.T) {
		in := strings.NewReader("")
		out := bytes.NewBuffer([]byte{})

		exprClient := &StubExpressionClient{}

		exprCli := cli.NewCLI(exprClient, in, out)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		returnedChan := make(chan struct{})
		go func() {
			exprCli.Run(ctx)
			returnedChan <- struct{}{}
		}()

		deadlineChan := make(chan struct{})
		go func() {
			time.Sleep(2 * time.Millisecond)
			deadlineChan <- struct{}{}
		}()

		select {
		case <-deadlineChan:
			t.Errorf("didn't exit on signal from ctx.Done()")
		case <-returnedChan:
		}
	})

	t.Run("returns on exit command", func(t *testing.T) {
		in := strings.NewReader(exitCmd)
		out := bytes.NewBuffer([]byte{})

		exprClient := &StubExpressionClient{}

		exprCli := cli.NewCLI(exprClient, in, out)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		exprCli.Run(ctx)

		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			t.Errorf("didn't exit CLI on command")
		}
	})

	t.Run("evaluates a math expression", func(t *testing.T) {
		expr := "What is 5 plus 10?"
		cmd := fmt.Sprint(".. ", expr)
		wantResult := 15
		in := strings.NewReader(cmd + "\n" + exitCmd)
		out := bytes.NewBuffer([]byte{})

		exprClient := &StubExpressionClient{
			result: wantResult,
		}

		exprCli := cli.NewCLI(exprClient, in, out)

		exprCli.Run(context.Background())

		assert.Equal(t, exprClient.spyEvaluateExpr, expr)
		assert.Equal(t, out.String(), fmt.Sprintln(wantResult))
	})

	t.Run("validates a math expresion", func(t *testing.T) {
		expr := "What is 5 plus 10?"
		cmd := fmt.Sprint("?? ", expr)
		wantIsValid := true
		in := strings.NewReader(cmd + "\n" + exitCmd)
		out := bytes.NewBuffer([]byte{})

		exprClient := &StubExpressionClient{
			isValid: wantIsValid,
		}

		exprCli := cli.NewCLI(exprClient, in, out)

		exprCli.Run(context.Background())

		assert.Equal(t, exprClient.spyValidateExpr, expr)
		assert.Equal(t, out.String(), fmt.Sprintln(wantIsValid))
	})

	t.Run("returns expression errors", func(t *testing.T) {
		exprError := client.ExpressionError{
			Expression: "Who is the president of the US?",
			Method:     "Validate()",
			Frequency:  3,
			Type:       "non-math question",
		}
		wantErrOutput := fmt.Sprintf(`\t"%s"; on %s; %d times; %s\n`,
			exprError.Expression, exprError.Method, exprError.Frequency, exprError.Type)

		cmd := "!!"
		in := strings.NewReader(cmd + "\n" + exitCmd)
		out := bytes.NewBuffer([]byte{})

		exprClient := &StubExpressionClient{
			exprErrors: []client.ExpressionError{exprError},
		}

		exprCli := cli.NewCLI(exprClient, in, out)

		exprCli.Run(context.Background())

		assert.Equal(t, out.String(), wantErrOutput)
	})

	t.Run("skips empty input", func(t *testing.T) {
		in := strings.NewReader("\n" + exitCmd)
		out := bytes.NewBuffer([]byte{})

		exprClient := &StubExpressionClient{}

		exprCli := cli.NewCLI(exprClient, in, out)

		exprCli.Run(context.Background())

		assert.Equal(t, out.Available(), 0)
	})

	t.Run("prints \"unknown command\" on unknown command", func(t *testing.T) {
		cmd := "## What is 5?"
		in := strings.NewReader(cmd + "\n" + exitCmd)
		out := bytes.NewBuffer([]byte{})

		exprClient := &StubExpressionClient{}

		exprCli := cli.NewCLI(exprClient, in, out)

		exprCli.Run(context.Background())

		assert.Equal(t, out.String(), "unknown command")
	})
}
