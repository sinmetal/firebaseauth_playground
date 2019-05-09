package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"firebase.google.com/go"
	"github.com/sinmetal/firebaseauth_playground/backend"
)

type FirebaseUser struct {
	UserID        string
	Email         string
	Name          string
	Picture       string
	EmailVerified bool
}

func main() {
	http.HandleFunc("/api/", apiHandler)
	http.HandleFunc("/", backend.StaticContentsHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	idToken := getIDToken(r.Header.Get("Authorization"))
	if len(idToken) < 1 {
		w.Header().Set("WWW-Authenticate", `Bearer realm=""`)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	app, err := firebase.NewApp(context.Background(), &firebase.Config{ProjectID: "metal-tile-dev1"})
	if err != nil {
		fmt.Fprintf(w, "failed firebase.NewApp:%s", err.Error())
		return
	}
	client, err := app.Auth(context.Background())
	if err != nil {
		fmt.Fprintf(w, "failed firebase.Auth:%s", err.Error())
		return
	}
	token, err := client.VerifyIDToken(r.Context(), idToken)
	if err != nil {
		fmt.Fprintf(w, "failed firebase.VerifyIDToken:%s", err.Error())
		return
	}
	user := FirebaseUser{
		UserID:        token.Claims["user_id"].(string),
		Email:         token.Claims["email"].(string),
		Name:          token.Claims["name"].(string),
		Picture:       token.Claims["picture"].(string),
		EmailVerified: token.Claims["email_verified"].(bool),
	}

	fmt.Fprintf(w, "Hello, World! %+v", user)
}

func getIDToken(headerValue string) string {
	l := strings.Split(headerValue, " ")
	if len(l) < 2 {
		return ""
	}
	return l[1]
}
