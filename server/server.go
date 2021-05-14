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
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gen0cide/laforge/ent"
	"github.com/gen0cide/laforge/ent/ginfilemiddleware"
	"github.com/gen0cide/laforge/graphql/auth"
	"github.com/gen0cide/laforge/graphql/graph"
	pb "github.com/gen0cide/laforge/grpc/proto"
	"github.com/gen0cide/laforge/grpc/server"
	"github.com/gen0cide/laforge/grpc/server/static"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const defaultPort = ":8080"

// Defining the Graphql handler
func redirectToRootHandler(client *ent.Client) gin.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file
	// h := handler.NewDefaultServer(graph.NewSchema(client))

	return func(c *gin.Context) {
		// h.ServeHTTP(c.Writer, c.Request)
		c.Redirect(301, "/ui")
		c.Abort()
	}
}

// tempURLHandler Checks ENT to verify that the url results in a file
func tempURLHandler(client *ent.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		urlID := ctx.Param("url_id")
		fileInfo, err := client.GinFileMiddleware.Query().Where(
			ginfilemiddleware.And(
				ginfilemiddleware.URLIDEQ(urlID),
				ginfilemiddleware.AccessedEQ(false),
			),
		).
			Only(ctx)
		if err != nil {
			ctx.AbortWithStatus(404)
			return
		}
		ctx.File(fileInfo.FilePath)
		_, err = fileInfo.Update().SetAccessed(true).Save(ctx)
		if err != nil {
			ctx.AbortWithStatus(404)
			return
		}
		ctx.Next()
	}
}

// Defining the Graphql handler
func graphqlHandler(client *ent.Client) gin.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file
	h := handler.NewDefaultServer(graph.NewSchema(client))

	h.AddTransport(&transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// Check against your desired domains here
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	})

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/api/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func main() {

	// pgHost, ok := os.LookupEnv("PG_HOST")
	// client := &ent.Client{}

	// if !ok {
	// 	client = ent.PGOpen("postgresql://laforger:laforge@127.0.0.1/laforge")
	// } else {
	// 	client = ent.PGOpen(pgHost)
	// }
	client := ent.SQLLiteOpen("file:test.sqlite?_loc=auto&cache=shared&_fk=1")

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

	router := gin.Default()

	// Add CORS middleware around every request
	// See https://github.com/rs/cors for full option listing
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "PUT", "PATCH"},
		AllowHeaders:     []string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "origin", "Cache-Control", "X-Requested-With"},
		AllowCredentials: true,
	}))

	port, ok := os.LookupEnv("PORT")

	if !ok {
		port = defaultPort
	}
	gqlHandler := graphqlHandler(client)
	redirectHandler := redirectToRootHandler(client)
	router.GET("/", redirectHandler)
	// router.Static("/ui/", "./dist")
	router.GET("/ui/*filename", func(c *gin.Context) {
		filename := c.Param("filename")
		fmt.Println(filename)
		if filename == "/monitor" || filename == "/monitor/" {
			c.Redirect(301, "/ui/")
		} else {
			c.File("./dist/" + filename)
		}
	})
	router.Static("/assets/", "./dist/assets")
	router.GET("/playground", playgroundHandler())

	api := router.Group("/api")
	api.Use(auth.Middleware(client))
	api.POST("/local/login", auth.Login(client))

	// TODO: Remove Get Path once Testing is done
	api.GET("/local/login", auth.Login(client))
	//

	api.GET("/local/logout", auth.Logout(client))
	api.POST("/query", gqlHandler)
	api.GET("/query", gqlHandler)
	api.GET("/download/:url_id", tempURLHandler(client))
	go router.Run(port)

	// secure server
	certPem, certerr := static.ReadFile(server.CertFile)
	if certerr != nil {
		fmt.Println("File reading error", certerr)
		return
	}
	keyPem, keyerr := static.ReadFile(server.KeyFile)
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
		Client:                     client,
		UnimplementedLaforgeServer: pb.UnimplementedLaforgeServer{},
	})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
