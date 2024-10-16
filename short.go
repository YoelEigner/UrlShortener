package main

import (
	"context"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Short struct {
	ID          string `bson:"_id"`
	OriginalURL string `bson:"originalUrl"`
	ShortURL    string `bson:"shortUrl"`
}

func HandleRedirect(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/short/")

	if len(parts) < 2 || parts[1] == "" {
		http.Error(w, "Short key missing", http.StatusBadRequest)
	}

	shortKey := parts[1]
	var result Short

	filter := bson.D{{Key: "shortUrl", Value: strings.TrimSpace(shortKey)}}
	collection := client.Database(dbName).Collection(collectionName)
	err := collection.FindOne(context.TODO(), filter).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Short URL not found", http.StatusNotFound)
		} else {
			http.Error(w, "Server error", http.StatusInternalServerError)
		}
		return
	}

	http.Redirect(w, r, result.OriginalURL, http.StatusMovedPermanently)
}
