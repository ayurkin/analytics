package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	pb "analytics/internal/transport/grpc/proto"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	conn, err := grpc.Dial("auth.team3.svc.cluster.local:4000", grpc.WithInsecure())

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
