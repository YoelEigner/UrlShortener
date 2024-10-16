package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"go.mongodb.org/mongo-driver/bson"
)

func HandleShorten(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if request.HTTPMethod != "POST" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusMethodNotAllowed,
			Body:       "Invalid request method",
		}, nil
	}

	var originalURL string
	if request.Headers["Content-Type"] == "application/x-www-form-urlencoded" {
		params, _ := url.ParseQuery(request.Body)
		originalURL = params.Get("url")
	} else if request.Headers["Content-Type"] == "application/json" {
		var body map[string]string
		json.Unmarshal([]byte(request.Body), &body)
		originalURL = body["url"]
	}

	if originalURL == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "URL missing",
		}, nil
	}

	shortKey := generateShortKey()

	err := saveShortentedUrl(originalURL, shortKey)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error saving shortened URL",
		}, nil
	}

	host := request.Headers["Host"]
	shortURL := fmt.Sprintf("https://%s/short/%s", host, shortKey)

	htmlContent := strings.Replace("./html/redirect.html", "{{.OriginalURL}}", originalURL, -1)
	htmlContent = strings.Replace(htmlContent, "{{.ShortURL}}", shortURL, -1)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type": "text/html",
		},
		Body: htmlContent,
	}, nil
}

func generateShortKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const shortKeyLength = 6

	shortKey := make([]byte, shortKeyLength)

	for i := range shortKey {
		shortKey[i] = charset[rand.Intn(len(charset))]
	}
	return string(shortKey)
}

func saveShortentedUrl(originalUrl string, shortKey string) error {
	collection := client.Database(dbName).Collection(collectionName)

	docs := bson.D{
		{Key: "originalUrl", Value: originalUrl},
		{Key: "shortUrl", Value: shortKey},
	}
	_, err := collection.InsertOne(context.TODO(), docs)
	return err
}
