package api

import (
	"fmt"
	"log"
	"net/http"
)

func root(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprintf(w, "Welcome to DiscordSoundBoard API root !")
}

func Listen() {
	http.HandleFunc("/", root)
	log.Println("Listening on 8080")
	http.ListenAndServe(":8080", nil)
}
