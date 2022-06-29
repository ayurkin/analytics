package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	pb "analytics/internal/transport/grpc/proto"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	authGrpcAddr := os.Getenv("AUTH_ADDR")
	log.Println("authGrpcAddr: " + authGrpcAddr)

	conn, err := grpc.Dial(authGrpcAddr, grpc.WithInsecure())

	if err != nil {
		log.Println(err)
	}
	client := pb.NewAuthApiClient(conn)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
		rest, err := client.Authenticate(context.Background(), &pb.AuthRequest{Username: "evg", AccessToken: "", RefreshToken: ""})
		if err != nil {
			log.Println(err)
		}

		log.Println(rest)
	})
	http.ListenAndServe(":3000", r)
}
