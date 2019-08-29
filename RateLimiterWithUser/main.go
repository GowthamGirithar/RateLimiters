package main

import (
	"RateLimiterWithUser/handlers"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	handlers.HandlerFunction(mux)
	http.ListenAndServe(":8080", mux)
	
}
