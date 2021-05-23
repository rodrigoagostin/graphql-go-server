package repository

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/rodrigagostin/graphql-server/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DATABASE   = "graphql"
	COLLECTION = "videos"
)

type VideoRepository interface {
	Save(video *model.Video)
	FindAll() []*model.Video
}

type database struct {
	client *mongo.Client
}

func New() VideoRepository {
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	// mongodb+sev://username:password@host:port
	MONGODB := os.Getenv("MONGODB")

	clientOptions := options.Client().ApplyURI(MONGODB)

	clientOptions = clientOptions.SetMaxPoolSize(50)

	dbClient, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatalf("Error to connect to mongodb %v", err)
	}

	fmt.Println("Connected to database")

	return &database{
		client: dbClient,
	}
}

func (db *database) Save(video *model.Video) {
	collection := db.client.Database(DATABASE).Collection(COLLECTION)
	_, err := collection.InsertOne(context.TODO(), video)
	if err != nil {
		log.Fatalf("Error to insert video on database %v", err)
	}
}

func (db *database) FindAll() []*model.Video {
	collection := db.client.Database(DATABASE).Collection(COLLECTION)
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatalf("Error to find videos %v", err)
	}
	defer cursor.Close(context.TODO())
	var result []*model.Video
	for cursor.Next(context.TODO()) {
		var v *model.Video
		if err := cursor.Decode(&v); err != nil {
			log.Fatalf("Error to decode document %v", err)
		}
		result = append(result, v)
	}
	return result
}
