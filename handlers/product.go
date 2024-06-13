package handlers

import (
	"25fd/micro-service/data"
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Product struct {
	l *log.Logger
}

func NewProduct(l *log.Logger) *Product {
	return &Product{l}
}

// func (p *Product) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
// 	if r.Method == http.MethodGet {
// 		p.getProducts(rw, r)
// 		return
// 	}

// 	if r.Method == http.MethodPost {
// 		p.addProduct(rw, r)
// 		return
// 	}

// 	if r.Method == http.MethodPut {
// 		rg := regexp.MustCompile(`/([0-9]+)`)
// 		g := rg.FindAllStringSubmatch(r.URL.Path, -1)

// 		p.l.Println(g)

// 		if len(g) != 1 || len(g[0]) != 2 {
// 			http.Error(rw, "Invalid request", http.StatusBadRequest)
// 		}
// 		id, err := strconv.Atoi(g[0][1])

// 		if err != nil {
// 			http.Error(rw, "Invalid Id", http.StatusBadRequest)
// 		}
// 		p.updateProduct(id, rw, r)
// 	}

// 	rw.WriteHeader(http.StatusMethodNotAllowed)
// }

func (p *Product) GetProducts(rw http.ResponseWriter, r *http.Request) {
	products := data.GetProducts()
	err := products.ToJSON(rw)
	if err != nil {
		http.Error(rw, "oh snap..", http.StatusInternalServerError)
		return
	}
}

func (p *Product) AddProduct(rw http.ResponseWriter, r *http.Request) {
	prod := r.Context().Value(KeyProduct{}).(*data.Product)

	p.l.Printf("Product %#v", prod)
	data.AddProduct(prod)
}

func (p *Product) UpdateProduct(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]

	//reading from the context
	prod := r.Context().Value(KeyProduct{}).(*data.Product)

	intId, _ := strconv.Atoi(id)
	errE := data.UpdatedProduct(intId, prod)
	if errE != nil {
		http.Error(rw, errE.Error(), http.StatusBadRequest)
		return
	}
}

type KeyProduct struct{}

func (p *Product) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		prod := &data.Product{}
		p.l.Println(r.Body)
		err := prod.FromJSON(r.Body)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
