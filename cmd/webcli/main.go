package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/VitoNaychev/eval-web-service/cli"
	"github.com/VitoNaychev/eval-web-service/client"
)

func main() {
	exprHTTPClient := client.NewExpressionHTTPClient(http.DefaultClient, "http://localhost:8080")

	cli := cli.NewCLI(exprHTTPClient, os.Stdin, os.Stdout)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Println("Math Expression Evaluator - Web Client 1.0")
	go cli.Run(ctx)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	<-sigCh
}
