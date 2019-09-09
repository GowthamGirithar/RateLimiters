package handlers

import (
	"net/http"
	"sync"
)

import "golang.org/x/time/rate"
//token bucket algorithm
var limiter = rate.NewLimiter(2, 5) // r tokens will be put in bucket every second and 5 is max size of bucket
var mtx sync.Mutex
//map[ip address][rate limiter]
var userLimiter=make(map[string]*rate.Limiter)
//add the visitor
func addVisitor( ipAddress string) {
	mtx.Lock()
	defer mtx.Unlock()
	userLimiter[ipAddress] = rate.NewLimiter(2, 5)
}
//get the visitor for rate limiter
func getVisitor(ipAddress string) *rate.Limiter{
	mtx.Lock()
	_, ok := userLimiter[ipAddress]
	mtx.Unlock()
	if !ok{
		addVisitor(ipAddress)
	}
	return userLimiter[ipAddress]
}

func HandlerFunction( mux *http.ServeMux) {
	mux.Handle("/health" , RateLimiter(HealthHandler))
}

func HealthHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("Success"))
}

func RateLimiter(next http.HandlerFunc) http.Handler {
     return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				if getVisitor(request.RemoteAddr).Allow(){
					next.ServeHTTP(writer, request)
				}else{
					// send as too many requests
					http.Error(writer, http.StatusText(429), http.StatusTooManyRequests)
				}

			return
	 })
}
