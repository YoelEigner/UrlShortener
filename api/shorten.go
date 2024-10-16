package handler

import (
	"context"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

func HandleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}

	originalURL := r.FormValue("url")

	if originalURL == "" {
		http.Error(w, "URL missing", http.StatusBadRequest)
	}

	shortKey := generateShortKey()

	saveShortentedUrl(originalURL, shortKey)

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	host := r.Host
	shortURL := fmt.Sprintf("%s://%s/short/%s", scheme, host, shortKey)

	data := struct {
		OriginalURL string
		ShortUrl    string
	}{
		OriginalURL: originalURL,
		ShortUrl:    shortURL,
	}
	tmpl := template.Must(template.ParseFiles("../html/redirect.html"))
	tmpl.Execute(w, data)
}

func generateShortKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const shortKeyLength = 6

	shortKey := make([]byte, shortKeyLength)

	for i := range shortKey {
		shortKey[i] = charset[rand.Intn((len((charset))))]
	}
	return string(shortKey)
}

func saveShortentedUrl(originalUrl string, shortKey string) {
	collection := client.Database(dbName).Collection(collectionName)

	docs := bson.D{
		{Key: "originalUrl", Value: originalUrl},
		{Key: "shortUrl", Value: shortKey},
	}
	_, err := collection.InsertOne(context.TODO(), docs)
	if err != nil {
		panic(err)
	}
}
