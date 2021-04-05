package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grkmk/glm-api/data"
	"github.com/hashicorp/go-hclog"
)

type Products struct {
	l         hclog.Logger
	productDB *data.ProductsDB
}

func NewProducts(l hclog.Logger, productionDB *data.ProductsDB) *Products {
	return &Products{l, productionDB}
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
