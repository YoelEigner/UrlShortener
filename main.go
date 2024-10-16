package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbName = "shorturls"
var collectionName = "shorturls"
var client *mongo.Client

func initMongoClient() error {
	godotenv.Load()

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	db_pass := os.Getenv("DB_PASS")
	cluster_name := os.Getenv("CLUSTER_NAME")
	opts := options.Client().ApplyURI("mongodb+srv://yoel:" + db_pass + "@" + cluster_name + "/?retryWrites=true&w=majority&appName=Cluster0").SetServerAPIOptions(serverAPI)

	var err error
	client, err = mongo.Connect(context.TODO(), opts)
	if err != nil {
		return err
	}

	if err := client.Database(dbName).RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		return err
	}
	fmt.Println("Successfully connected to MongoDB!")
	return nil
}

func Handler(w http.ResponseWriter, r *http.Request) {
	if err := initMongoClient(); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	switch r.URL.Path {
	case "/":
		HandleHome(w, r)
	case "/shorten":
		HandleShorten(w, r)
	default:
		if strings.HasPrefix(r.URL.Path, "/short/") {
			HandleRedirect(w, r)
			return
		}
		http.NotFound(w, r)
	}
}

func main() {
	http.HandleFunc("/", HandleHome)
	http.HandleFunc("/shorten", HandleShorten)
	http.HandleFunc("/short/", HandleRedirect)

	initMongoClient()

	fmt.Println("URL Shortener is running on :8080")
	http.ListenAndServe(":8080", nil)
}
