package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var totalRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_reuqest_total",
		Help: "Number of get requests",
	},
	[]string{"path"},
)

func prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)

		totalRequests.WithLabelValues("path").Inc()
	})
}

func init(){
	prometheus.Register(totalRequests)
}

func main(){
	router := mux.NewRouter()

	router.Path("/metrics").Handler(promhttp.Handler())
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	log.Println("server started on port 3000")
	err := http.ListenAndServe(":3000", router)
	log.Fatal(err)
}