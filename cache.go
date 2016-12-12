package main

import (
	"log"
	"time"

	"gopkg.in/redis.v5"
)


func writeToCache(request string, payload []byte, ttlSec int64) {
	expiration := time.Duration(ttlSec) * time.Second
	err := redisDB.Set("cache:" + request, payload, expiration).Err()
	if err != nil {
		log.Printf("Not able to cache data to Redis : %s", err)
	}
}


func getFromCache(request string) ([]byte, error) {
	val, err := redisDB.Get("cache:" + request).Bytes()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		log.Printf("Not able to get data from cache : %s", err)
		return nil, err
	}
	return val, nil
}


func incEndpointCounter(e string) {
	err := redisDB.Incr(STATS_PREFIX + e).Err()
	if err != nil {
		log.Printf("Not able to increase counter for %s : %s",e, err)
	}
}


func recordRequestTime(r string, timer int64) {
	err := redisDB.Set(TIMER_PREFIX + r, timer, -1).Err()
	if err != nil {
		log.Printf("Not able to increase counter for %s : %s", r, err)
	}
}

