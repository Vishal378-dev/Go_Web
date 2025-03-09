package Middleware

import (
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

func Authorize(usercollection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("authorized user values ->", r.Context().Value("authorizeduser"))
		if r.Method == "GET" {
			fmt.Fprintf(w, "Hello From Middleware")
			return
		}
	}
}
