package data

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	protos "github.com/grkmk/glm-currency/protos/currency"
	"github.com/hashicorp/go-hclog"
)

// Product defines the sturcture for an API product
// swagger:model
type Product struct {
	// the id for this user
	//
	// required: true
	// min: 1
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"gt=0"`
	SKU         string  `json:"sku" validate:"required,sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

func (p *Product) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

type Products []*Product

type ProductsDB struct {
	currency protos.CurrencyClient
	log      hclog.Logger
	rates    map[string]float64
	client   protos.Currency_SubscribeRatesClient
}

func NewProductsDB(c protos.CurrencyClient, l hclog.Logger) *ProductsDB {
	pb := &ProductsDB{c, l, make(map[string]float64), nil}

	go pb.handleUpdates()

	return pb
}

func (p *ProductsDB) handleUpdates() {
	sub, err := p.currency.SubscribeRates(context.Background())
	if err != nil {
		p.log.Error("Unable to subscribe for rates", "error", err)
		return
	}

	p.client = sub

	for {
		rr, err := sub.Recv()
		p.log.Info("Received updated rate from server", "dest", rr.GetDestination().String())

		if err != nil {
			p.log.Error("Error reveiving message", "error", err)
			return
		}

		p.rates[rr.Destination.String()] = rr.Rate
	}
}

func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *Product) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *ProductsDB) GetProducts(currency string) (Products, error) {
	if currency == "" {
		return productList, nil
	}

	rate, err := p.getRate(currency)
	if err != nil {
		p.log.Error("Unable to get rate", "currency", currency, "error", err)
		return nil, err
	}

	pr := Products{}
	for _, p := range productList {
		np := *p
		np.Price = np.Price * rate
		pr = append(pr, &np)
	}

	return pr, nil
}

func (p *ProductsDB) GetProduct(id int, currency string) (*Product, error) {
	product, _, err := findProduct(id)

	if err != nil {
		return nil, ErrProductNotFound
	}

	if currency == "" {
		return product, nil
	}

	rate, err := p.getRate(currency)
	if err != nil {
		p.log.Error("Unable to get rate", "currency", currency, "error", err)
		return nil, err
	}

	np := *product
	np.Price = np.Price * rate

	return &np, nil
}

func AddProduct(p *Product) {
	p.ID = getNextId()
	productList = append(productList, p)
}

func RemoveProduct(id int, p *Product) error {
	_, pos, err := findProduct(id)
	if err != nil {
		return err
	}

	productList = append(productList[:pos], productList[pos:]...)

	return nil
}

func UpdateProduct(id int, p *Product) error {
	_, pos, err := findProduct(id)
	if err != nil {
		return err
	}

	p.ID = id
	productList[pos] = p

	return nil
}

var ErrProductNotFound = fmt.Errorf("Product not found")

func findProduct(id int) (*Product, int, error) {
	for i, p := range productList {
		if p.ID == id {
			return p, i, nil
		}
	}

	return nil, -1, ErrProductNotFound
}

func getNextId() int {
	listOfProducts := productList[len(productList)-1]
	return listOfProducts.ID + 1
}

var productList = []*Product{
	{
		ID:          1,
		Name:        "Latte",
		Description: "Espresso with milk",
		Price:       2.45,
		SKU:         "asdf1234",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	{
		ID:          2,
		Name:        "Espresso",
		Description: "Short without milk",
		Price:       1.99,
		SKU:         "asdf1235",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}

func (p *ProductsDB) getRate(destination string) (float64, error) {
	if r, ok := p.rates[destination]; ok {
		return r, nil
	}

	rr := &protos.RateRequest{
		Base:        protos.Currencies(protos.Currencies_value["EUR"]),
		Destination: protos.Currencies(protos.Currencies_value[destination]),
	}

	resp, err := p.currency.GetRate(context.Background(), rr)
	p.rates[destination] = resp.Rate

	p.client.Send(rr)

	return resp.Rate, err
}
