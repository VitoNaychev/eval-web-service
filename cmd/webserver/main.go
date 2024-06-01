package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/VitoNaychev/eval-web-service/handler"
	"github.com/VitoNaychev/eval-web-service/interp"
	"github.com/VitoNaychev/eval-web-service/repo"
	"github.com/VitoNaychev/eval-web-service/service"
)

func main() {
	exprErrorRepo := repo.NewInMemoryExprErrorRepository()

	exprInterp := interp.NewInterpMW(interp.Lex, interp.Parse, interp.Interpret)

	exprService := service.NewExpressionService(exprInterp, exprErrorRepo)

	exprHandler := handler.NewExpressionHandler(exprService)

	router := handler.NewRouter(exprHandler)

	fmt.Println("Evaluate service listening on :8080...")
	log.Fatal(http.ListenAndServe(":8080", router.Handler))
}
