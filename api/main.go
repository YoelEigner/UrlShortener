package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
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

func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if err := initMongoClient(); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Internal Server Error",
		}, nil
	}

	switch req.Path {
	case "/":
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       "Welcome to the URL Shortener!",
		}, nil
	case "/shorten":
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       "Shorten URL endpoint",
		}, nil
	default:
		if strings.HasPrefix(req.Path, "/short/") {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusOK,
				Body:       "Redirecting...",
			}, nil
		}
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Body:       "404 Not Found",
		}, nil
	}
}

func main() {
	lambda.Start(Handler)
}
