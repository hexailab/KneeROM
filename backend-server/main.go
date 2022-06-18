package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GlobalAppVariables struct {
	MongoClient *mongo.Client
	MongoContext context.Context
}

const MongoClientDatabaseTitle = "range_db"

var appVariables GlobalAppVariables

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// allow cross domain AJAX requests
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		next.ServeHTTP(w,r)
	})
}

func main() {
	var err error
	var cancel context.CancelFunc

	if err = godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("Invalid .env file.")
	}

	appVariables.MongoContext, cancel = context.WithCancel(context.Background())
	defer cancel()

	if appVariables.MongoClient, err = mongo.Connect(appVariables.MongoContext, options.Client().ApplyURI(uri)); err != nil {
		panic(err)
	}
	defer func() {
		if err := appVariables.MongoClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	text, err := CreateNewSessionID(AccountDescriptiveTitleTypePatient)
	fmt.Println(text.GetAccountType())


	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: RootFieldObject}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery), Mutation: RootMutationObject}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	h := handler.New(&handler.Config{
		Schema: &schema,
		Pretty: true,
		GraphiQL: true,
	})

	http.Handle("/graphql", CorsMiddleware(h))
	http.ListenAndServe(":8080", nil)
}