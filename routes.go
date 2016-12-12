package main

import (
	"log"
	"fmt"
	"net/http"
	"encoding/json"
	"sort"
	"time"
	"strings"

	"github.com/gorilla/mux"
	"github.com/Pryz/nextbus"
)


// Index enpoint. Will display all available API routes
func Index(w http.ResponseWriter, r *http.Request) {
	incEndpointCounter("root")
	data, _ := json.MarshalIndent(routes, "", "	")
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}


// Data structure to display the API statistics
type StatsData struct {
	RedisPing string
	Counters map[string]int64
	Timers	map[string]int64
}


// Statistics endpoint
func Stats(w http.ResponseWriter, r *http.Request) {
	var stats StatsData
	stats.Counters = make(map[string]int64)
	stats.Timers = make(map[string]int64)

	stats.RedisPing, _ = redisDB.Ping().Result()

	// Retrieve all timers
	keys, _ := redisDB.Keys(STATS_PREFIX + "*").Result()
	for _, k := range keys {
		c, err := redisDB.Get(k).Int64()
		if err != nil {
		}
		stats.Counters[strings.Split(k, ":")[1]] = c
	}

	// Retrieve all timers
	keys, _ = redisDB.Keys(TIMER_PREFIX + "*").Result()
	for _, k := range keys {
		t, err := redisDB.Get(k).Int64()
		if err != nil {
		}
		stats.Timers[strings.Split(k, ":")[1]] = t
	}

	data, _ := json.MarshalIndent(stats, "", "	")
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}


func AgencyList(w http.ResponseWriter, r *http.Request) {
	var agencies []nextbus.Agency

	// Try to retrieve result from cache
	data, _ :=  getFromCache(r.URL.String())
	if data != nil {
		w.Write(data)
		return
	}

	timerStart := time.Now()
	agencies, err := nextbus.GetAgencyList()
	if err != nil {
		log.Printf("Error during NextBus query : %s", err)
	}
	timerEnd := time.Since(timerStart)

	incEndpointCounter(mux.CurrentRoute(r).GetName())
	recordRequestTime(string(r.URL.String()), int64(timerEnd))

	data, _ = json.Marshal(agencies)
	// Write to cache before output
	writeToCache(r.URL.String(), data, CACHE_TTL_SEC)
	// Output json
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}


func RouteList(w http.ResponseWriter, r *http.Request) {
	var routes []nextbus.Route
	vars := mux.Vars(r)

	// Try to retrieve result from cache
	data, _ :=  getFromCache(r.URL.String())
	if data != nil {
		w.Write(data)
		return
	}

	timerStart := time.Now()
	routes, err := nextbus.GetRouteList(vars["agency"])
	if err != nil {
		fmt.Fprintln(w, err)
	}
	timerEnd := time.Since(timerStart)

	incEndpointCounter(mux.CurrentRoute(r).GetName())
	recordRequestTime(string(r.URL.String()), int64(timerEnd))

	data, _ = json.Marshal(routes)
	// Write to cache before output
	writeToCache(r.URL.String(), data, CACHE_TTL_SEC)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}


func RouteConfig(w http.ResponseWriter, r *http.Request) {
	var routeConfigs []nextbus.RouteConfig
	vars := mux.Vars(r)

	// Try to retrieve result from cache
	data, _ :=  getFromCache(r.URL.String())
	if data != nil {
		w.Write(data)
		return
	}

	timerStart := time.Now()
	routeConfigs, err := nextbus.GetRouteConfig(vars["agency"], vars["route"])
	if err != nil {
		fmt.Fprintln(w, err)
	}
	timerEnd := time.Since(timerStart)

	incEndpointCounter(mux.CurrentRoute(r).GetName())
	recordRequestTime(string(r.URL.String()), int64(timerEnd))

	data, _ = json.Marshal(routeConfigs)
	// Write to cache before output
	writeToCache(r.URL.String(), data, CACHE_TTL_SEC)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}


func Predictions(w http.ResponseWriter, r *http.Request) {
	var predictions []nextbus.PredictionData
	vars := mux.Vars(r)

	// Try to retrieve result from cache
	data, _ :=  getFromCache(r.URL.String())
	if data != nil {
		w.Write(data)
		return
	}

	timerStart := time.Now()
	predictions, err := nextbus.GetPredictions(vars["agency"], vars["route"], vars["stop"])
	if err != nil {
		fmt.Fprintln(w, err)
	}
	timerEnd := time.Since(timerStart)

	incEndpointCounter(mux.CurrentRoute(r).GetName())
	recordRequestTime(string(r.URL.String()), int64(timerEnd))

	data, _ = json.Marshal(predictions)
	// Write to cache before output
	writeToCache(r.URL.String(), data, CACHE_TTL_SEC)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}


func Schedule(w http.ResponseWriter, r *http.Request) {
	var routes []nextbus.ScheduleRoute
	vars := mux.Vars(r)

	// Try to retrieve result from cache
	data, _ :=  getFromCache(r.URL.String())
	if data != nil {
		w.Write(data)
		return
	}

	timerStart := time.Now()
	routes, err := nextbus.GetSchedule(vars["agency"], vars["route"])
	if err != nil {
		fmt.Fprintln(w, err)
	}
	timerEnd := time.Since(timerStart)

	incEndpointCounter(mux.CurrentRoute(r).GetName())
	recordRequestTime(string(r.URL.String()), int64(timerEnd))

	data, _ = json.Marshal(routes)
	// Write to cache before output
	writeToCache(r.URL.String(), data, CACHE_TTL_SEC)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}


func OffRouteList(w http.ResponseWriter, r *http.Request) {
	var scheduledRoutes, curRoutes []nextbus.ScheduleRoute
	var closedRoutes map[string]string

	// Try to retrieve result from cache
	data, _ :=  getFromCache(r.URL.String())
	if data != nil {
		w.Write(data)
		return
	}

	log.Println("No cache found")

	vars := mux.Vars(r)
	agency := vars["agency"]
	reqEpoch := stringToEpoch(vars["time"])
	svcClass := dayToServiceClass(vars["day"])

	timerStart := time.Now()
	rList, err := nextbus.GetRouteList(agency)
	if err != nil {
		log.Println(err)
	}
	for _, route := range rList {
		rsList, err := nextbus.GetSchedule(agency, route.Tag)
		if err != nil {
			continue
		}
		scheduledRoutes = append(scheduledRoutes, rsList...)
	}
	timerEnd := time.Since(timerStart)

	incEndpointCounter(mux.CurrentRoute(r).GetName())
	recordRequestTime(string(r.URL.String()), int64(timerEnd))

	// retrieve current route depending of "day"
	for _, route := range scheduledRoutes {
		if route.ServiceClass == svcClass {
			curRoutes = append(curRoutes, route)
		}
	}

	closedRoutes = make(map[string]string)
	for _, route := range curRoutes {
		var times []int
		for _, tr := range route.Tr {
			for _, stop := range tr.StopList {
				if stop.EpochTime != 1 {
					times = append(times, stop.EpochTime)
				}
			}
		}
		// Sort and get begin/end times
		sort.Ints(times)
		begin := times[0]
		end := times[len(times) - 1]
		if begin > reqEpoch || reqEpoch > end {
			closedRoutes[route.RouteTag] = route.RouteTitle
		}
	}

	w.Header().Set("Content-Type", "application/json")
	data, _ = json.MarshalIndent(closedRoutes, "", "	")

	// Write to cache before output
	writeToCache(r.URL.String(), data, CACHE_TTL_SEC)
	// Output data
	w.Write(data)
}

