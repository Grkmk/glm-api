package handlers

import (
	"net/http"

	"github.com/grkmk/glm-api/data"
)

// swagger:route POST /products products addProduct
// Create new product
//
// responses:
// 200: productResponse
// 422: errorValidation
// 501: errorResponse

// AddProduct handles POST requests to add new products
func (p *Products) AddProduct(responseWriter http.ResponseWriter, request *http.Request) {
	// fetch the product form the context
	prod := request.Context().Value(KeyProduct{}).(*data.Product)

	data.AddProduct(prod)
}
