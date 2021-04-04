package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/grkmk/glm-api/data"
	protos "github.com/grkmk/glm-currency/protos/currency"
)

type Products struct {
	l  *log.Logger
	cc protos.CurrencyClient
}

func NewProducts(l *log.Logger, cc protos.CurrencyClient) *Products {
	return &Products{l, cc}
}

type KeyProduct struct{}

func (p Products) ValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		prod := &data.Product{}

		err := prod.FromJSON(request.Body)
		if err != nil {
			http.Error(responseWriter, "Unable to unmarshal json", http.StatusBadRequest)
			return
		}

		err = prod.Validate()
		if err != nil {
			http.Error(responseWriter, fmt.Sprintf("Error validating product: %s", err), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(request.Context(), KeyProduct{}, prod)
		req := request.WithContext(ctx)

		next.ServeHTTP(responseWriter, req)
	})
}
