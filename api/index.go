package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	_ "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	dbName         = "shorturls"
	collectionName = "shorturls"
	client         *mongo.Client
	clientOnce     sync.Once
	clientError    error
)

func initMongoClient() (*mongo.Client, error) {
	clientOnce.Do(func() {
		serverAPI := options.ServerAPI(options.ServerAPIVersion1)
		db_pass := os.Getenv("DB_PASS")
		cluster_name := os.Getenv("CLUSTER_NAME")
		opts := options.Client().ApplyURI("mongodb+srv://yoel:" + db_pass + "@" + cluster_name + "/?retryWrites=true&w=majority&appName=Cluster0").SetServerAPIOptions(serverAPI)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var err error
		client, err = mongo.Connect(ctx, opts)
		if err != nil {
			clientError = fmt.Errorf("failed to connect to MongoDB: %v", err)
			return
		}

		err = client.Ping(ctx, nil)
		if err != nil {
			clientError = fmt.Errorf("failed to ping MongoDB: %v", err)
			return
		}

		fmt.Println("Successfully connected to MongoDB!")
	})

	return client, clientError
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	_, err := initMongoClient()
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Internal Server Error",
		}, nil
	}

	switch {
	case request.Path == "/":
		return HandleHome(request)
	case request.Path == "/shorten":
		return HandleShorten(request)
	case strings.HasPrefix(request.Path, "/short/"):
		return HandleRedirect(request)
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Body:       "404 Not Found",
		}, nil
	}
}

func main() {
	lambda.Start(Handler)
}
