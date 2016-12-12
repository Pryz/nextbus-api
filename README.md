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

## Todo

There is a lot of room for improvements here :
* Write unit tests
* Improve usability with Flags and/or a configuration file
  * Be able to change the Cache TTL
  * Be able to change the Redis port
* Improve "stats" endpoint by adding more informations about the Cache 
