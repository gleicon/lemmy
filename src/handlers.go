package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/fiorix/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/thoas/stats"
)

type ApiRouter struct {
	router *mux.Router
	config *configFile
	rc     *redis.Client
	db     *gorm.DB
}

func NewApiRouter(config *configFile, rc *redis.Client, db *gorm.DB) *ApiRouter {
	ar := ApiRouter{config: config, rc: rc, db: db}
	ar.router = mux.NewRouter()
	ar.router.HandleFunc("/health", ar.HealthHandler).Methods("GET")
	ar.router.HandleFunc("/hello", ar.HelloHandler).Methods("GET")
	ar.router.HandleFunc("/docroot", ar.DocumentRootHandler).Methods("GET")
	ar.router.HandleFunc("/httpclient", ar.HttpClientHandler).Methods("GET")

	n := negroni.New(negroni.NewRecovery(), negroni.NewStatic(http.Dir(config.DocumentRoot)), NewCommonLogger(os.Stdout))

	middleware := stats.New()
	ar.router.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		stats := middleware.Data()

		b, _ := json.Marshal(stats)

		w.Write(b)
	})
	n.Use(middleware)
	n.UseHandler(ar.router)
	n.Run(config.HTTP.Addr)
	return &ar
}

func (ar *ApiRouter) GetRouter() *mux.Router {
	return ar.router
}

func (ar *ApiRouter) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "{\"status\": \"OK\"}")
}

func (ar *ApiRouter) HttpClientHandler(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{Timeout: time.Duration(500 * time.Millisecond)}

	resp, err := client.Get("http://google.com")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s", body)
}

func (ar *ApiRouter) DocumentRootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, ar.config.DocumentRoot)
}

func (ar *ApiRouter) HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "hello")
}
