package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
)

type Backend struct {
	URL          *url.URL
	Alive        bool
	ReverseProxy *httputil.ReverseProxy
}

type BackendPool struct {
	servers []*Backend
	index   uint32
}

func NewBackendPool() *BackendPool {
	return &BackendPool{}
}

func (bp *BackendPool) AddBackend(BackendURL string) {
	url, _ := url.Parse(BackendURL)
	bp.servers = append(bp.servers, &Backend{
		URL:          url,
		Alive:        true,
		ReverseProxy: httputil.NewSingleHostReverseProxy(url),
	})

}

func (bp *BackendPool) NextBackend() *Backend {
	n := len(bp.servers)
	if n == 0 {
		return nil
	}
	next := atomic.AddUint32(&bp.index, 1) % uint32(n)
	return bp.servers[next]
}

func (bp *BackendPool) LoadBalanceHandler(w http.ResponseWriter, r *http.Request) {
	backend := bp.NextBackend()
	if backend != nil {
		backend.ReverseProxy.ServeHTTP(w, r)
		return
	}
	http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
}

func main() {
	backendPool := NewBackendPool()
	backendPool.AddBackend("http://localhost:8080")
	backendPool.AddBackend("http://localhost:8081")
	backendPool.AddBackend("http://localhost:8082")

	http.HandleFunc("/", backendPool.LoadBalanceHandler)
	fmt.Println("Running the server on 8079")
	log.Fatal(http.ListenAndServe(":8079", nil))
}
