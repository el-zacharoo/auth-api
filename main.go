package main

import (
	"net/http"

	"github.com/el-zacharoo/auth/handler"
	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	pbcnn "github.com/el-zacharoo/auth/internal/gen/auth/v1/authv1connect"
)

const port = "localhost:8082"

func main() {

	svc := &handler.SignInServer{}

	c := setCORS()
	mux := http.NewServeMux()
	path, h := pbcnn.NewAuthServiceHandler(svc)
	mux.Handle(path, h)
	hndlr := c.Handler(mux)

	http.ListenAndServe(
		port,
		h2c.NewHandler(hndlr, &http2.Server{}),
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
