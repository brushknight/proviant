package http

import (
	"encoding/json"
	"github.com/brushknight/proviant/internal/config"
	"github.com/brushknight/proviant/internal/errors"
	"github.com/brushknight/proviant/internal/i18n"
	"github.com/brushknight/proviant/internal/pkg/category"
	"github.com/brushknight/proviant/internal/pkg/list"
	"github.com/brushknight/proviant/internal/pkg/product"
	"github.com/brushknight/proviant/internal/pkg/product_category"
	"github.com/brushknight/proviant/internal/pkg/service"
	"github.com/brushknight/proviant/internal/pkg/stock"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Server struct {
	productRepo         *product.Repository
	listRepo            *list.Repository
	categoryRepo        *category.Repository
	productCategoryRepo *product_category.Repository
	stockRepo           *stock.Repository
	relationService     *service.RelationService
	router              *mux.Router
	l                   i18n.Localizer
	cfg                 config.Config
}

func (s *Server) Run(hostPort string) error {
	return http.ListenAndServe(hostPort, s.router)
}

func (s *Server) parseJSON(r *http.Request, model interface{}) error {
	return json.NewDecoder(r.Body).Decode(model)
}

func (s *Server) getLocale(r *http.Request) i18n.Locale {
	return i18n.LocaleFromString(r.Header.Get("User-Locale"))
}

func (s *Server) handleBadRequest(w http.ResponseWriter, locale i18n.Locale, error string, params ...interface{}) {
	m := i18n.NewMessage(error, params...)
	response := Response{
		Status: BadRequest,
		Error:  s.l.T(m, locale),
	}

	s.jsonResponse(w, response)
}

func (s *Server) handleError(w http.ResponseWriter, locale i18n.Locale, error errors.CustomError) {
	response := Response{
		Status: error.Code(),
		Error:  s.l.T(error.Message(), locale),
	}

	s.jsonResponse(w, response)
}

func (s *Server) jsonResponse(w http.ResponseWriter, response Response) {
	payload, err := json.Marshal(response)
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.Status)
	_, err = w.Write(payload)
	if err != nil {
		log.Println(err)
	}
	return
}

func NewServer(productRepo *product.Repository,
	listRepo *list.Repository,
	categoryRepo *category.Repository,
	productCategoryRepo *product_category.Repository,
	stockRepo *stock.Repository,
	relationService *service.RelationService,
	l i18n.Localizer,
	cfg config.Config) *Server {

	server := &Server{
		productRepo:         productRepo,
		listRepo:            listRepo,
		categoryRepo:        categoryRepo,
		productCategoryRepo: productCategoryRepo,
		stockRepo:           stockRepo,
		relationService:     relationService,
		l:                   l,
		cfg:                 cfg,
	}

	router := mux.NewRouter()

	apiV1Router := router.PathPrefix("/api/v1").Subrouter()

	// product routes
	apiV1Router.HandleFunc("/product/{id}/", server.getProduct).Methods("GET")
	apiV1Router.HandleFunc("/product/", server.getProducts).Methods("GET")
	apiV1Router.HandleFunc("/product/", server.createProduct).Methods("POST")
	apiV1Router.HandleFunc("/product/{id}/", server.updateProduct).Methods("PUT")
	apiV1Router.HandleFunc("/product/{id}/", server.deleteProduct).Methods("DELETE")
	// category routes
	apiV1Router.HandleFunc("/category/{id}/", server.getCategory).Methods("GET")
	apiV1Router.HandleFunc("/category/", server.getCategories).Methods("GET")
	apiV1Router.HandleFunc("/category/", server.createCategory).Methods("POST")
	apiV1Router.HandleFunc("/category/{id}/", server.updateCategory).Methods("PUT")
	apiV1Router.HandleFunc("/category/{id}/", server.deleteCategory).Methods("DELETE")
	// list routes
	apiV1Router.HandleFunc("/list/{id}/", server.getList).Methods("GET")
	apiV1Router.HandleFunc("/list/", server.getLists).Methods("GET")
	apiV1Router.HandleFunc("/list/", server.createList).Methods("POST")
	apiV1Router.HandleFunc("/list/{id}/", server.updateList).Methods("PUT")
	apiV1Router.HandleFunc("/list/{id}/", server.deleteList).Methods("DELETE")
	// stock routers
	apiV1Router.HandleFunc("/product/{id}/stock/", server.getStock).Methods("GET")
	apiV1Router.HandleFunc("/product/{id}/add/", server.addStock).Methods("POST")
	apiV1Router.HandleFunc("/product/{id}/consume/", server.consumeStock).Methods("POST")
	apiV1Router.HandleFunc("/product/{product_id}/stock/{id}/", server.deleteStock).Methods("DELETE")
	apiV1Router.HandleFunc("/i18n/missing/", server.getMissingTranslations).Methods("GET")

	if cfg.Mode == config.ModeWeb {
		router.PathPrefix("/static").Handler(http.FileServer(http.Dir("./public/")))

		spa := spaHandler{staticPath: "public", indexPath: "index.html"}
		router.PathPrefix("/").Handler(spa)
	}

	server.router = router

	return server
}
