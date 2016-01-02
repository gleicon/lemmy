# Lemmy

Dynamic HTTP/HTTPS Loadbalancer and API gateway with opinionated choices.

Lemmy keeps an external table on backends that serve specific VHOSTS.

You can plug in middlewares compatible with net/http and negroni to provide oauth, throttling and filtering.

Currently lemmy uses redis to keep dynamic data and backend config. It can be easily replaced.

## Redis layout

lemmy:<vhost> is a zset with backends sorted by current connections
lemmy:deactivated:<vhost> is a zset with backends that presented errors, sorted by error number
lemmy:stats:<vhost> realtime stats on vhosts


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


## Baseline routes

	Healthcheck: curl http://localhost:8080/health | python -mjson.tool
	Stats: curl http://localhost:8080/stats | python -mjson.tool


