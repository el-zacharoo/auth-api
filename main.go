package main

import (
	"net/http"

	"github.com/el-zacharoo/auth/handler"
	"github.com/rs/cors"
)

const port = "localhost:8082"

func main() {
	h := &handler.SignInServer{}

	c := setCORS()
	mux := http.NewServeMux()
	mux.HandleFunc("/user", h.GetUser)
	mux.HandleFunc("/login", h.SignIn)
	hndlr := c.Handler(mux)

	http.ListenAndServe(
		port,
		hndlr,
	)
}

func setCORS() *cors.Cors {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"Access-Control-Allow-Origin", "Content-Type"},
		AllowedMethods:   []string{"POST"},
		AllowCredentials: true,
	})
	return c
}
