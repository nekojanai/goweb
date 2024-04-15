package utils

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
)

func LogReq(req *http.Request) {
	data, err := httputil.DumpRequest(req, false)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%v\n", data)
}

func ErrorHandler(w http.ResponseWriter, r *http.Request, status int, err error) {
	if err != nil {
		w.WriteHeader(status)
		switch status {
		case 404:
			fmt.Fprintf(w, "404")
		}
	}
}
