package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/redis.v5"
)


// Prefixes for Redis data
const STATS_PREFIX string = "stats:"
const TIMER_PREFIX string = "time:"
const CACHE_PREFIX string = "cache:"

// Cache TTL
const CACHE_TTL_SEC int64 = 30


// Global variables
var (
	router *mux.Router
	routes []string
	redisDB *redis.Client
)


func getEndpoints(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
	t, err := route.GetPathTemplate()
	if err != nil {
		return err
	}
	routes = append(routes, t)
	return nil
}


func main() {

	// Initiate redis connection for statistics and cache
	redisDB = initRedis()

	// Setup API routes
	router = mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index).Name("root")
	router.HandleFunc("/agencylist", AgencyList).Name("agencylist")
	router.HandleFunc("/routelist/{agency}", RouteList).Name("routelist")
	router.HandleFunc("/routeconfig/{agency}/{route}", RouteConfig).Name("routeconfig")
	router.HandleFunc("/predictions/{agency}/{route}/{stop}", Predictions).Name("predictions")
	router.HandleFunc("/schedule/{agency}/{route}", Schedule).Name("schedule")
	router.HandleFunc("/offroutelist/{agency}/{day}/{time}", OffRouteList).Name("offroute")
	router.HandleFunc("/stats", Stats).Name("stats")
	router.Walk(getEndpoints)

	log.Fatal(http.ListenAndServe(":8080", router))
}

