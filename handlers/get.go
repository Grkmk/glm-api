package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// swagger:route GET /products products listProducts
// Returns a list of products
// responses:
// 	200: productsResponse

// GetProducts returns the products from the data store
func (p *Products) GetProducts(responseWriter http.ResponseWriter, request *http.Request) {
	p.l.Debug("Get all records")
	responseWriter.Header().Add("Content-type", "application/json")

	cur := request.URL.Query().Get("currency")

	listOfProducts, err := p.productDB.GetProducts(cur)
	if err != nil {
		http.Error(responseWriter, "Unable to get products", http.StatusInternalServerError)
	}

	// serialize the list to JSON
	err = listOfProducts.ToJSON(responseWriter)
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

	cur := request.URL.Query().Get("currency")

	p.l.Debug("Get record", "id", id)
	responseWriter.Header().Add("Content-type", "application/json")

	product, err := p.productDB.GetProduct(id, cur)
	if err != nil {
		http.Error(responseWriter, "Unable to find product for id", http.StatusBadRequest)
		return
	}

	err = product.ToJSON(responseWriter)
	if err != nil {
		http.Error(responseWriter, "Unable to marshal json", http.StatusInternalServerError)
	}
}
