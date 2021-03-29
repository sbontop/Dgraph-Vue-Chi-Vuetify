package main

import (
	"context"
	"errors"
	"flag"
	"net/http"
	"strconv"
	"strings"

	// Chi Router Libraries (Api Router)
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	// Dgraph Libraries (Graph-based database)
	"encoding/json"
	"log"

	"github.com/dgraph-io/dgo/v2"
	"github.com/dgraph-io/dgo/v2/protos/api"
	"google.golang.org/grpc"
)

type CancelFunc func()

func getDgraphClient() (*dgo.Dgraph, CancelFunc) {
	conn, err := grpc.Dial("127.0.0.1:9080", grpc.WithInsecure())
	if err != nil {
		log.Fatal("While trying to dial gRPC")
	}

	dc := api.NewDgraphClient(conn)
	dg := dgo.NewDgraphClient(dc)

	return dg, func() {
		if err := conn.Close(); err != nil {
			log.Printf("Error while closing connection:%v", err)
		}
	}
}

var routes = flag.Bool("routes", false, "Generate router documentation")

func main() {

	flag.Parse()
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Route("/buyers", func(r chi.Router) {
		r.With(paginate).Get("/", ListBuyers)
	})
	r.Route("/realbuyers", func(r chi.Router) {
		r.With(paginate).Get("/", ListRealBuyers)
		r.Get("/purchaseHistory", ListPurchaseHistory)
		r.Route("/ip/{buyer_id}", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				var buyer *Buyer
				var err error
				buyer, err = dbGetBuyer(chi.URLParam(r, "buyer_id"))
				if err != nil {
					render.Render(w, r, ErrNotFound)
					return
				}
				response := dbGetBuyersByIp(buyer.Ip)
				_, _ = w.Write([]byte(response))
			})
		})
		r.Route("/recommendations/{buyer_id}", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				var buyer *Buyer
				var err error
				buyer, err = dbGetBuyer(chi.URLParam(r, "buyer_id"))
				if err != nil {
					render.Render(w, r, ErrNotFound)
					return
				}
				response := dbGetBuyersRecomByAge(strconv.Itoa(buyer.Age))
				_, _ = w.Write([]byte(response))
			})
		})
	})

	r.Route("/products", func(r chi.Router) {
		r.With(paginate).Get("/", ListProducts)
	})
	r.Route("/transactions", func(r chi.Router) {
		r.With(paginate).Get("/", ListTransactions)
	})
	http.ListenAndServe(":3333", r)
}

func ListBuyers(w http.ResponseWriter, r *http.Request) {
	if err := render.RenderList(w, r, NewBuyerListResponse(buyers)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

func ListRealBuyers(w http.ResponseWriter, r *http.Request) {
	if err := render.RenderList(w, r, NewRealBuyerListResponse(realbuyers)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

func ListPurchaseHistory(w http.ResponseWriter, r *http.Request) {
	if err := render.RenderList(w, r, NewPurchaseHistoryListResponse(purchasesHistory)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

func ListProducts(w http.ResponseWriter, r *http.Request) {
	if err := render.RenderList(w, r, NewProductListResponse(products)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

func ListTransactions(w http.ResponseWriter, r *http.Request) {
	if err := render.RenderList(w, r, NewTransactionListResponse(transactions)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

// paginate is a stub, but very possible to implement middleware logic
// to handle the request params for handling a paginated request.
func paginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// just a stub.. some ideas are to look at URL query params for something like
		// the page number, or the limit, and send a query cursor down the chain
		next.ServeHTTP(w, r)
	})
}

type BuyerRequest struct {
	*Buyer

	// User *UserPayload `json:"user,omitempty"`

	ProtectedID string `json:"id"` // override 'id' json to have more control
}

func (a *BuyerRequest) Bind(r *http.Request) error {
	// a.Article is nil if no Article fields are sent in the request. Return an
	// error to avoid a nil pointer dereference.
	if a.Buyer == nil {
		return errors.New("missing required Buyer fields.")
	}

	// a.User is nil if no Userpayload fields are sent in the request. In this app
	// this won't cause a panic, but checks in this Bind method may be required if
	// a.User or futher nested fields like a.User.Name are accessed elsewhere.

	// just a post-process after a decode..
	a.ProtectedID = ""                           // unset the protected ID
	a.Buyer.Name = strings.ToLower(a.Buyer.Name) // as an example, we down-case
	return nil
}

type BuyerResponse struct {
	*Buyer
	Elapsed int64 `json:"elapsed"`
}

type RealBuyerResponse struct {
	*Buyer
	Elapsed int64 `json:"elapsed"`
}

type PurchaseHistoryResponse struct {
	*PurchaseHistory
	Elapsed int64 `json:"elapsed"`
}

type ProductResponse struct {
	*Product
	Elapsed int64 `json:"elapsed"`
}

type TransactionResponse struct {
	*Transaction
	Elapsed int64 `json:"elapsed"`
}

func NewBuyerResponse(buyer *Buyer) *BuyerResponse {
	resp := &BuyerResponse{Buyer: buyer}
	return resp
}

func NewRealBuyerResponse(buyer *Buyer) *RealBuyerResponse {
	resp := &RealBuyerResponse{Buyer: buyer}
	return resp
}

func NewPurchaseHistoryResponse(purchaseHistory *PurchaseHistory) *PurchaseHistoryResponse {
	resp := &PurchaseHistoryResponse{PurchaseHistory: purchaseHistory}
	return resp
}

func NewProductResponse(product *Product) *ProductResponse {
	resp := &ProductResponse{Product: product}

	return resp
}

func NewTransactionResponse(transaction *Transaction) *TransactionResponse {
	resp := &TransactionResponse{Transaction: transaction}

	return resp
}

func (rd *BuyerResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	rd.Elapsed = 10
	return nil
}

func (rd *RealBuyerResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	rd.Elapsed = 10
	return nil
}

func (rd *PurchaseHistoryResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	rd.Elapsed = 10
	return nil
}

func (rd *ProductResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	rd.Elapsed = 10
	return nil
}

func (rd *TransactionResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	rd.Elapsed = 10
	return nil
}

func NewBuyerListResponse(buyers []*Buyer) []render.Renderer {
	list := []render.Renderer{}
	for _, buyer := range buyers {
		list = append(list, NewBuyerResponse(buyer))
	}
	return list
}

func NewRealBuyerListResponse(buyers []*Buyer) []render.Renderer {
	list := []render.Renderer{}
	for _, realbuyer := range realbuyers {
		list = append(list, NewRealBuyerResponse(realbuyer))
	}
	return list
}

func NewPurchaseHistoryListResponse(purchaseHistory []*PurchaseHistory) []render.Renderer {
	list := []render.Renderer{}
	for _, purchaseHistory := range purchasesHistory {
		list = append(list, NewPurchaseHistoryResponse(purchaseHistory))
	}
	return list
}

func NewProductListResponse(products []*Product) []render.Renderer {
	list := []render.Renderer{}
	for _, product := range products {
		list = append(list, NewProductResponse(product))
	}
	return list
}

func NewTransactionListResponse(transactions []*Transaction) []render.Renderer {
	list := []render.Renderer{}
	for _, transaction := range transactions {
		list = append(list, NewTransactionResponse(transaction))
	}
	return list
}

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		StatusText:     "Error rendering response.",
		ErrorText:      err.Error(),
	}
}

var ErrNotFound = &ErrResponse{HTTPStatusCode: 404, StatusText: "Resource not found."}

type Buyer struct {
	UID  string `json:"uid"`
	ID   string `json:"buyer_id"`
	Name string `json:"buyer_name"`
	Age  int    `json:"buyer_age"`
	Ip   string
}

type Buyers struct {
	Buyers []*Buyer `json:"buyers"`
}

type Product struct {
	UID   string `json:"uid"`
	ID    string `json:"product_id"`
	Name  string `json:"product_name"`
	Price int    `json:"product_price"`
}

type Products struct {
	Products []*Product `json:"products"`
}

type Transaction struct {
	UID      string     `json:"uid"`
	ID       string     `json:"transaction_id"`
	Buyers   []*Buyer   `json:"buyer"`
	Products []*Product `json:"product"`
	Ip       string     `json:"ip"`
	Device   string     `json:"device"`
}

type Transactions struct {
	Transactions []*Transaction `json:"transactions"`
}

type PurchaseHistory struct {
	Products []*Product `json:"product"`
}

type PurchasesHistory struct {
	PurchasesHistory []*PurchaseHistory `json:"purchasesHistory"`
}

// Buyers data from dgraph
var buyers = dbGetBuyers()

func dbGetBuyers() []*Buyer {
	// Initialize dgraph
	dg, cancel := getDgraphClient()
	defer cancel()
	ctx := context.Background()
	const q = `
	query myQuery() {
		buyers(func: has(buyer_id)) {
			uid
			buyer_id
			buyer_name
			buyer_age
		}
	}	
	`
	resp, err := dg.NewTxn().Query(ctx, q)
	if err != nil {
		log.Fatal(err)
	}

	var r Buyers

	if err := json.Unmarshal(resp.GetJson(), &r); err != nil {
		log.Fatal(err)
	}
	return r.Buyers
}

// Buyers that has transactions from dgraph
var realbuyers = dbGetRealBuyers()

func dbGetRealBuyers() []*Buyer {
	// Initialize dgraph
	dg, cancel := getDgraphClient()
	defer cancel()
	ctx := context.Background()
	const q = `
	query myQuery() {
		buyers(func: has(transaction_id)) {
			buyer_id
			buyer_name
			buyer_age
			ip
		}
	}	
	`
	resp, err := dg.NewTxn().Query(ctx, q)
	if err != nil {
		log.Fatal(err)
	}

	var r Buyers

	if err := json.Unmarshal(resp.GetJson(), &r); err != nil {
		log.Fatal(err)
	}
	return r.Buyers
}

// Purchases History from dgraph
var purchasesHistory = dbGetPurchasesHistory()

func dbGetPurchasesHistory() []*PurchaseHistory {
	// Initialize dgraph
	dg, cancel := getDgraphClient()
	defer cancel()
	ctx := context.Background()
	const q = `
	query myQuery() {
		purchasesHistory(func: eq(buyer_id, "ad2ba138")) {
			product {
				product_id
				product_name
				product_price
			}
		}
	}	
	`
	resp, err := dg.NewTxn().Query(ctx, q)
	if err != nil {
		log.Fatal(err)
	}

	var r PurchasesHistory

	if err := json.Unmarshal(resp.GetJson(), &r); err != nil {
		log.Fatal(err)
	}
	return r.PurchasesHistory
}

// get buyers by ip from dgraph
func dbGetBuyersByIp(ip string) []byte {
	// Initialize dgraph
	dg, cancel := getDgraphClient()
	defer cancel()
	ctx := context.Background()
	variables := map[string]string{"$id1": ip}
	const q = `
	query myQuery($id1: string) {
		getBuyerByIpAddress(func: eq(ip, $id1)) {  
			buyer_id
			buyer_name
			buyer_age
		}
	}	
	`
	resp, err := dg.NewTxn().QueryWithVars(ctx, q, variables)
	if err != nil {
		log.Fatal(err)
	}

	var r Buyers

	if err := json.Unmarshal(resp.GetJson(), &r); err != nil {
		log.Fatal(err)
	}
	// return r.Buyers
	return resp.GetJson()
}

func dbGetBuyersRecomByAge(age string) []byte {
	// Initialize dgraph
	dg, cancel := getDgraphClient()
	defer cancel()
	ctx := context.Background()
	variables := map[string]string{"$id1": age}
	const q = `
	query myQuery($id1: int) {
		productRecom(func: eq(buyer_age, $id1)) {
			product {
			  product_name
			  product_price
			}
		}
	}	
	`
	resp, err := dg.NewTxn().QueryWithVars(ctx, q, variables)
	if err != nil {
		log.Fatal(err)
	}

	var r Products

	if err := json.Unmarshal(resp.GetJson(), &r); err != nil {
		log.Fatal(err)
	}
	// return r.Products
	return resp.GetJson()
}

// Products data from dgraph
var products = dbGetProducts()

func dbGetProducts() []*Product {
	// Initialize dgraph
	dg, cancel := getDgraphClient()
	defer cancel()
	ctx := context.Background()
	const q = `
	query myQuery() {
		products(func: has(product_id)) {
			uid
			product_id
			product_name
			product_price
		}
	}	
	`
	resp, err := dg.NewTxn().Query(ctx, q)
	if err != nil {
		log.Fatal(err)
	}
	var r Products
	if err := json.Unmarshal(resp.GetJson(), &r); err != nil {
		log.Fatal(err)
	}
	return r.Products
}

// Transactions data from dgraph
var transactions = dbGetTransactions()

func dbGetTransactions() []*Transaction {
	// Initialize dgraph
	dg, cancel := getDgraphClient()
	defer cancel()
	ctx := context.Background()
	const q = `
	query myQuery() {
		transactions(func: has(transaction_id)) {
			uid
			transaction_id
			buyer {
				buyer_id
				buyer_name
				buyer_age
			}
			product {
				product_id
				product_name
				product_price
			}
			ip
			device
		}
	}	
	`
	resp, err := dg.NewTxn().Query(ctx, q)
	if err != nil {
		log.Fatal(err)
	}
	var r Transactions
	if err := json.Unmarshal(resp.GetJson(), &r); err != nil {
		log.Fatal(err)
	}
	return r.Transactions
}

func dbGetBuyer(id string) (*Buyer, error) {
	for _, a := range realbuyers {
		if a.ID == id {
			return a, nil
		}
	}
	return nil, errors.New("realbuyer not found.")
}
