package handler

import (
	"context"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Short struct {
	ID          string `bson:"_id"`
	OriginalURL string `bson:"originalUrl"`
	ShortURL    string `bson:"shortUrl"`
}

func HandleRedirect(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	parts := strings.Split(request.Path, "/short/")

	if len(parts) < 2 || parts[1] == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Short key missing",
		}, nil
	}

	shortKey := parts[1]
	var result Short

	filter := bson.D{{Key: "shortUrl", Value: strings.TrimSpace(shortKey)}}
	collection := client.Database(dbName).Collection(collectionName)
	err := collection.FindOne(context.TODO(), filter).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusNotFound,
				Body:       "Short URL not found",
			}, nil
		} else {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       "Server error",
			}, nil
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusFound,
		Headers: map[string]string{
			"Location": result.OriginalURL,
		},
		Body: "",
	}, nil
}
