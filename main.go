package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbName = "shorturls"
var collectionName = "shorturls"
var client *mongo.Client

func initMongoClient() {
	godotenv.Load()
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	db_pass := os.Getenv("DB_PASS")
	cluster_name := os.Getenv("CLUSTER_NAME")
	opts := options.Client().ApplyURI("mongodb+srv://yoel:" + db_pass + "@" + cluster_name + "/?retryWrites=true&w=majority&appName=Cluster0").SetServerAPIOptions(serverAPI)
	var err error
	client, err = mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	if err := client.Database(dbName).RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		panic(err)
	}
	fmt.Println("successfully connected to MongoDB!")
}

func main() {
	http.HandleFunc("/", HandleHome)
	http.HandleFunc("/shorten", HandleShorten)
	http.HandleFunc("/short/", HandleRedirect)

	initMongoClient()

	fmt.Println("URL Shortener is running on :8080")
	http.ListenAndServe("0.0.0.0:8080", nil)
}
