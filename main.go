package main

import (
	"go-trace-demo/grpc_server"
	"go-trace-demo/http_server"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"time"
)

var (
	g errgroup.Group
)

func main() {
	server01 := &http.Server{
		Addr:         ":8081",
		Handler:      http_server.Router01(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	server02 := &http.Server{
		Addr:         ":8082",
		Handler:      http_server.Router02(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	server03 := &http.Server{
		Addr:         ":8083",
		Handler:      http_server.Router03(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	g.Go(func() error {

		return server01.ListenAndServe()
	})

	g.Go(func() error {
		return server02.ListenAndServe()
	})

	g.Go(func() error {
		return server03.ListenAndServe()
	})
	go grpc_server.GrpcServer()
	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
