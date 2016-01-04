# Lemmy

Dynamic HTTP/HTTPS Loadbalancer and API gateway with opinionated choices.

Lemmy keeps an external table on backends that serve specific VHOSTS.

You can plug in middlewares compatible with net/http and negroni to provide oauth, throttling and filtering.

Currently lemmy uses redis to keep dynamic data and backend config. It can be easily replaced.

## Redis layout
	<prefix> is a string prefix set at the config file. This is to differentiate lemmy instances in case you want to share the same DB.
	<prefix>:vhosts is a set with all vhosts that an instance will answer for
	<prefix>:stats:<hostname> realtime stats for the loadbalancer at a given hostname
	<prefix>:backends:<vhost> is a scored set with all backend info. The score of active connections will serve as weight of that backend.
	<prefix>:stats:<vhost> realtime stats on vhosts

## VHost refresh
Assuming prefix is "lemmy" and the vhost list has "aceofspades" and "onparole" on it:

The *vhost* update goroutine will *SMEMBER lemmy:vhosts*, loop in the result issuing a *ZRANGEBYSCORE lemmy:backends:aceofspades -inf +inf* and *ZRANGEBYSCORE lemmy:backends:onparole* to fill up a cache in memory so the loadbalancer don't need to ping back Redis at each request. This loop is governed by the directive *vhost_refresh_rate* in the configuration file, in seconds.


This goroutine is the place to implement backend tests, server recovery and other distribution algorithms. 

## Loadbalancer algorithm

The algorithm used is a modified weighted distribution by number of connections. To avoid a connection flood to a new backend, a roundrobin distribution can be applied after N requests, set at configuration file by the directive *round_robin_after* directive. This is the number of requests sent to the same backend before a roundrobin is triggered among the remaining backends (if any) skipping the new one.

For example, in a two backend A and B setup receiving 1000 requests per second, a new backend C would attrack connections until the old ones had served their connections. To avoid that you can set *round_robin_after* to 100 and guarantee that this fill up will occur in 100 requests steps. The loadbalancer will send 100 requests to C, round robin the next requests to A and B giving time to C process the new connections, send 100 more connections to C and so on until all of them are balancing equally depending on the current connections. 

This algorithm is smarter than a pure roundrobin while allowing for a slow backend to recover and for a new backend to adapt to the workload.

Each request will increment the current connection number. In our example say we forward a connection to the VHost aceofspades to the backend 10.0.0.1: *ZINCRBY lemmy:backends:aceofspades 1 10.0.0.1*. When the connection is finished we issue *ZINCRBY lemmy:backends:aceofspades -1 10.0.0.1*.

A faulty backend is taken off schedule after *max_fails* failures. After *backend_retry_interval* seconds it will be reinserted back if not removed from the backend list. The health check is passive only for now.

## HTTP API

All operations can be done throught the HTTP API.
	
	Healthcheck: curl http://<host>:<port>/health | python -mjson.tool
	Stats: curl http://<host>:<port>/stats | python -mjson.tool
	tbd

## Depends on

  - negroni
  - golang
  - gorilla mux
  - go-redis (for redis support)
  - gorm (for your relational db needs)

## Authentication

Authentication is built in by API token

## Build

$ make dep ; make

There will be a binary at bin/




