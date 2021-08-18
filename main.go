package main

import (
	"fmt"
	"log"
	"net/http"
	"main/tools"

	// initialize mock DB with blank identifier.
	_ "main/db"
)

func main() {
	mux := http.NewServeMux()

	toolHandlers := tools.Handlers()                 // 1

	mux.HandleFunc("/api/tools", toolHandlers.Root)  // 2
  mux.HandleFunc("/api/tools/", toolHandlers.Items)
  mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		for k, v := range r.Header {
			fmt.Println(k, v)
		}
  mux.HandleFunc("/myProblem", toolHandlers.ThrowError)  

		fmt.Fprintf(w, "Hi! check your terminal's output for more info on the request headers")
})

	fmt.Println("API running on port :3000")
	log.Fatal(http.ListenAndServe(":3000", mux))
}