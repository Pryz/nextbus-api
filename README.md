# Nextbus API

This daemon provide a simple RestAPI built on top of the NextBus XML Feed.

See : http://www.nextbus.com/xmlFeedDocs/NextBusXMLFeed.pdf

It requires a Redis DB listening on port 6379

## Build it

```
make build
```

## Use it with Docker

```
docker-compose up
```

You should then be able to fetch all the APIs :

```
$ curl localhost:8080/
[
        "/",
        "/agencylist",
        "/routelist/{agency}",
        "/routeconfig/{agency}/{route}",
        "/predictions/{agency}/{route}/{stop}",
        "/schedule/{agency}/{route}",
        "/offroutelist/{agency}/{day}/{time}",
        "/stats"
]
```

## Todo

There is a lot of room for improvements here :
* Write unit tests
* Improve usability with Flags and/or a configuration file
  * Be able to change the Cache TTL
* Improve "stats" endpoint by adding more informations about the Cache 
