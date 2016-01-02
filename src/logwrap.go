package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/codegangsta/negroni"
)

// CommonLogger is a middleware handler that logs the request as it goes in and the response as it goes out.
type CommonLogger struct {
	st io.Writer
}

// NewLogger returns a new Logger instance
func NewCommonLogger(st io.Writer) *CommonLogger {
	return &CommonLogger{st: st}
}

func (cl *CommonLogger) proto(r *http.Request) string {
	if r.TLS == nil {
		return "HTTP"
	} else {
		return "HTTPS"
	}
}

/*
::1 - - [14/Jul/2015:16:39:56 -0300] "GET /hello_log HTTP/1.1" 200 5 "" "curl/7.37.1"
HTTP 200 GET "/hello_log" ([::1]:49730) :: 5 bytes in 1.027115ms curl/7.37.1[negroni] Completed 200 OK in 1.175722ms
127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700] "GET /apache_pb.gif HTTP/1.0" 200 2326 "http://www.example.com/start.html" "Mozilla/4.08 [en] (Win98; I ;Nav)"

*/

func (cl *CommonLogger) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()
	next(rw, r)
	res := rw.(negroni.ResponseWriter)
	t := time.Now()
	fmt.Fprintf(cl.st, "%s - - %s \"%s %s %s\" %d %d %s %dus\n",
		r.RemoteAddr,
		t.Format("02/Jan/2006:15:04:05 -0700"),
		r.Method,
		r.URL.Path,
		r.Proto,
		res.Status(),
		res.Size(),
		r.UserAgent(),
		time.Since(start),
	)
}
