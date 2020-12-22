package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gen0cide/laforge/ent"
	"github.com/gen0cide/laforge/graphql/graph"
	pb "github.com/gen0cide/laforge/grpc/proto"
	"github.com/gen0cide/laforge/grpc/server"
	"github.com/gen0cide/laforge/grpc/server/static"
	"github.com/go-chi/chi"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const defaultPort = "80"

func main() {

	pgHost, ok := os.LookupEnv("PG_HOST")
	client := &ent.Client{}

	if !ok {
		client = ent.PGOpen("postgresql://laforger:laforge@127.0.0.1/laforge")
	} else {
		client = ent.PGOpen(pgHost)
	}

	ctx := context.Background()
	defer client.Close()

	// Run the auto migration tool.
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	lis, err := net.Listen("tcp", server.Port)
	
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := handler.NewDefaultServer(graph.NewSchema(client))

	router := chi.NewRouter()

	// Add CORS middleware around every request
	// See https://github.com/rs/cors for full option listing
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080", "http://localhost:4200", "http://localhost:80"},
		AllowCredentials: true,
		Debug:            false,
	}).Handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	go http.ListenAndServe(":"+port, router)

	// secure server
	certPem,certerr := static.ReadFile(server.CertFile)
	if certerr != nil {
        fmt.Println("File reading error", certerr)
        return 
	}
	keyPem,keyerr := static.ReadFile(server.KeyFile)
	if keyerr != nil {
        fmt.Println("File reading error", keyerr)
        return 
	}

	cert, tlserr := tls.X509KeyPair(certPem, keyPem)
	if tlserr != nil {
        fmt.Println("File reading error", tlserr)
        return 
	}

	creds := credentials.NewServerTLSFromCert(&cert)
	s := grpc.NewServer(grpc.Creds(creds))


	log.Printf("Starting Laforge Server on port " + server.Port)
	
	pb.RegisterLaforgeServer(s, &server.Server{
		Client: client,
		UnimplementedLaforgeServer: pb.UnimplementedLaforgeServer{},
	})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
