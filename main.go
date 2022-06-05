package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	tempval, temp := r.URL.Query()["temp"]
	nodeval, node := r.URL.Query()["node"]

	if temp && node {
		json.NewEncoder(w).Encode(map[string]string{"status": "200", "temp": tempval[0], "node": nodeval[0]})
	} else {
		json.NewEncoder(w).Encode(map[string]string{"status": "500"})
	}
}

func main() {
	fmt.Println("Started backend on port 8080")
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
