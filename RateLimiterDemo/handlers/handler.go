package handlers

import "net/http"

import "golang.org/x/time/rate"

var limiter = rate.NewLimiter(2, 5) // r tokens will be put in bucket every second and 5 is max size of bucket

func HandlerFunction( mux *http.ServeMux) {
	mux.Handle("/health" , RateLimiter(HealthHandler))
}

func HealthHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("Success"))
}

func RateLimiter(next http.HandlerFunc) http.Handler {
     return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				if limiter.Allow(){
					next.ServeHTTP(writer, request)
				}else{
					// send as too many requests
					http.Error(writer, http.StatusText(429), http.StatusTooManyRequests)
				}

			return
	 })
}
