package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/nishipy/microservice-go-practice/data"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	//GET
	if r.Method == http.MethodGet {
		p.getProducts(rw, r)
		return
	}

	//POST
	if r.Method == http.MethodPost {
		p.addProducts(rw, r)
		return
	}

	//PUT = update
	if r.Method == http.MethodPut {
		p.l.Println("Handle PUT Product")
		//expected the id in URI
		//https://pkg.go.dev/regexp
		reg := regexp.MustCompile(`/([0-9]+)`)
		// if user execute `curl localhost:9090/1 -X PUT -d ...`
		// then g is [[/1 1]]. Each value means:
		// - `/1`: the string matched
		// - `1`: the information for matched string
		g := reg.FindAllStringSubmatch(r.URL.Path, -1)
		if len(g) != 1 {
			p.l.Println("Invalid URI more than one id")
			http.Error(rw, "Invalid URI", http.StatusBadGateway)
			return
		}
		if len(g[0]) != 2 {
			p.l.Println("Invalid URI more than one capture group")
			http.Error(rw, "Invalid URI", http.StatusBadGateway)
			return
		}
		idString := g[0][1]
		id, err := strconv.Atoi(idString)
		if err != nil {
			p.l.Println("Invalid URI unable to convert to number", idString)
			http.Error(rw, "Invalid URI", http.StatusBadGateway)
			return
		}

		p.updateProducts(id, rw, r)
		return
	}

	//catch all
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Products) getProducts(rw http.ResponseWriter, r *http.Request) {
	lp := data.GetProducts()
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Products) addProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Product")

	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json.", http.StatusBadRequest)
	}

	data.AddProducts(prod)

}

func (p *Products) updateProducts(id int, rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle PUT Product")

	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json.", http.StatusBadRequest)
	}

	err = data.UpdateProducts(id, prod)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}
}
