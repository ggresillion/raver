package api

import (
	"fmt"
	"net/http"
)

func Listen() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, World!")
	})
	return http.ListenAndServe(":8080", mux)
}
