// File: main.go
package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

type Email struct {
	Email string
}

func emailSet(w http.ResponseWriter, r *http.Request) {
	var p Email

	err := decodeJSONBody(w, r, &p)
	if err != nil {
		var mr *malformedRequest
		fmt.Print(mr.msg)

		if errors.As(err, &mr) {
			http.Error(w, mr.msg, mr.status)

		} else {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		}
		return
	}

	fmt.Fprintf(w, "%v", p)
}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/email", emailSet)

	log.Println("Starting server on :4000...")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
