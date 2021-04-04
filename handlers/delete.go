package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/grkmk/glm-api/data"
)

// swagger:route DELETE /products/{id} products deleteProduct
// Deletes a product from the store
// responses:
// 	201: noContent

// DeleteProduct removes the product from the data store
func (p *Products) DeleteProduct(responseWriter http.ResponseWriter, request *http.Request) {
	// this will always convert because of the router
	vars := mux.Vars(request)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(responseWriter, "Unable to convert id", http.StatusBadRequest)
		return
	}

	prod := request.Context().Value(KeyProduct{}).(*data.Product)

	err = data.RemoveProduct(id, prod)
	if err == data.ErrProductNotFound {
		http.Error(responseWriter, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(responseWriter, "Something went wrong", http.StatusInternalServerError)
		return
	}
}
