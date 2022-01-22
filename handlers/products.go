package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/zhou-en/go_mservice/product-api/data"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

//func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
//	if r.Method == http.MethodGet {
//		p.getProducts(rw, r)
//		return
//	}
//
//	// update
//	if r.Method == http.MethodPost {
//		p.addProduct(rw, r)
//		return
//	}
//
//	//
//	if r.Method == http.MethodPut {
//		p.l.Println("PUT: ", r.URL.Path)
//		// expect id in uri
//		reg := regexp.MustCompile(`/([0-9]+)`)
//		g := reg.FindAllStringSubmatch(r.URL.Path, -1)
//		if len(g) != 1 {
//			http.Error(rw, "Invalid URI", http.StatusBadRequest)
//			return
//		}
//		if len(g[0]) != 2 {
//			http.Error(rw, "Invalid URI", http.StatusBadRequest)
//			return
//		}
//
//		idString := g[0][1]
//		id, err := strconv.Atoi(idString)
//		if err != nil {
//			http.Error(rw, "Invalid URI", http.StatusBadRequest)
//			return
//		}
//		p.updateProducts(id, rw, r)
//		return
//	}
//
//	// catch the rest
//	rw.WriteHeader(http.StatusMethodNotAllowed)
//}

func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	lp := data.GetProducts()
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal JSON", http.StatusInternalServerError)
	}
}

func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Products")
	prod := r.Context().Value(KeyProduct{}).(data.Product)
	p.l.Printf("Product: %#v", prod)
	data.AddProduct(&prod)
}

func (p *Products) UpdateProducts(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
		return
	}
	p.l.Println("Handle PUT Product", id)
	prod := r.Context().Value(KeyProduct{}).(data.Product)

	err = data.UpdateProduct(id, &prod)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}
}

type KeyProduct struct {
}

func (p Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := data.Product{}

		err := prod.FromJSON(r.Body)
		if err != nil {
			http.Error(rw, "Unable to unmarshal product", http.StatusBadRequest)
			return
		}

		err = prod.Validate()
		if err != nil {
			p.l.Println("[ERROR] validating product", err)
			http.Error(
				rw,
				fmt.Sprintf("Error validating product: %s", err),
				http.StatusBadRequest,
			)
			return
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		req := r.WithContext(ctx)
		next.ServeHTTP(rw, req)
	})
}
