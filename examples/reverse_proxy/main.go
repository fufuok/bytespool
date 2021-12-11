package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/fufuok/bytespool"
)

// Ref: https://github.com/fufuok/reverse-proxy
func main() {
	target, _ := url.Parse("https://cn.bing.com")
	proxy := httputil.NewSingleHostReverseProxy(target)

	proxy.BufferPool = bytespool.NewBufPool(32 * 1024)

	director := proxy.Director
	proxy.Director = func(r *http.Request) {
		director(r)
		r.Host = target.Host
	}

	http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		log.Println(req.RequestURI)
		proxy.ServeHTTP(rw, req)
	})

	log.Println("try: http://127.0.0.1:8080")
	log.Fatalln(http.ListenAndServe(":8080", nil))
}
