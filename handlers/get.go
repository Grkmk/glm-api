package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/grkmk/glm-api/data"
)

// swagger:route GET /products products listProducts
// Returns a list of products
// responses:
// 	200: productsResponse

// GetProducts returns the products from the data store
func (p *Products) GetProducts(responseWriter http.ResponseWriter, request *http.Request) {
	// fetch the procust from the data store
	listOfProducts := data.GetProducts()

	// serialize the list to JSON
	err := listOfProducts.ToJSON(responseWriter)
	if err != nil {
		http.Error(responseWriter, "Unable to marshal json", http.StatusInternalServerError)
	}
}

// swagger:route GET /products/{id} products listProduct
// Returns a product
// responses:
// 	200: productResponse

// GetProduct returns the product from the data store
func (p *Products) GetProduct(responseWriter http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(responseWriter, "Unable to convert id", http.StatusBadRequest)
		return
	}

	product, err := data.GetProduct(id)
	if err != nil {
		http.Error(responseWriter, "Unable to find product for id", http.StatusBadRequest)
		return
	}

	err = product.ToJSON(responseWriter)
	if err != nil {
		http.Error(responseWriter, "Unable to marshal json", http.StatusInternalServerError)
	}
}
